/*
   Copyright 2019 Dominik Madarász <zaklaus@madaraszd.net>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package core

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"io/ioutil"
	"log"

	"github.com/zaklaus/go-tiled"

	rl "github.com/zaklaus/raylib-go/raylib"
)

var (
	// CurrentSaveSystem is the primary save system to be used
	CurrentSaveSystem SaveSystem

	// CanSave specifies if we're allowed to save at this point
	CanSave Bits
)

const (
	// IsInMenu is player in menu
	IsInMenu Bits = 1 << iota

	// IsSequenceHappening is a scripted sequence happening
	IsSequenceHappening

	// IsPlayerDead is it game over yet
	IsPlayerDead

	// IsInDialogue is player currently in an active dialogue
	IsInDialogue

	// IsInChallenge is player in a danger zone which disallows saving the game
	IsInChallenge
)

// SaveSystem manages game save states
type SaveSystem struct {
	Version string      `json:"version"`
	States  []GameState `json:"gameStates"`
}

// GameState describes the serializable save state
type GameState struct {
	SaveName string
	SaveData defaultSaveData `json:"saveData"`
}

// InitSaveSystem initializes the game state system
func (s *SaveSystem) InitSaveSystem() {
	rdat, err := ioutil.ReadFile("gamesav.db")
	dat, _ := base64.StdEncoding.DecodeString(string(rdat))

	hasFailed := false

	if err == nil {
		var sav SaveSystem
		buf := bytes.NewBuffer(dat)
		dec := gob.NewDecoder(buf)
		err = dec.Decode(&sav)
		//err = jsoniter.Unmarshal(dat, &sav)

		if err != nil {
			log.Printf("Gamesav.db is broken, ignoring...\n")
			hasFailed = true
		} else {
			*s = sav
		}
	} else {
		hasFailed = true
	}

	if hasFailed {
		s.States = make([]GameState, 10)
	}

	s.Version = GameVersion
}

// SaveGame saves the game state
func (s *SaveSystem) SaveGame(slotIndex int, stateName string) bool {
	if CanSave != 0 {
		log.Printf("Cannot save the game right now! Reason: %v\n", CanSave)
		return false
	}

	state := GameState{
		SaveName: stateName,
	}

	state.SaveData = defaultSaveProvider(&state)

	s.States[slotIndex] = state

	//data, _ := jsoniter.Marshal(s)
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	enc.Encode(s)

	dat := base64.StdEncoding.EncodeToString(buf.Bytes())

	ioutil.WriteFile("gamesav.db", []byte(dat), 0755)
	return true
}

// LoadGame restores the game state
func (s *SaveSystem) LoadGame(slotIndex int) bool {
	state := &s.States[slotIndex]

	initScriptingSystem()
	defaultLoadProvider(state)
	return true
}

// ShutdownSaveSystem closes down the connection
func (s *SaveSystem) ShutdownSaveSystem() {

}

type saveData interface{}

type defaultSaveData struct {
	saveData
	CurrentMap   string           `json:"active"`
	Maps         []defaultMapData `json:"maps"`
	GameModeData []byte           `json:"gameMode"`
}

type defaultMapData struct {
	MapName     string              `json:"map"`
	Objects     []defaultObjectData `json:"objects"`
	WeatherData Weather             `json:"weather"`
}

type objectData interface{}

type defaultObjectData struct {
	Name        string `json:"objectName"`
	Type        string `json:"class"`
	Position    rl.Vector2
	Movement    rl.Vector2
	Facing      rl.Vector2
	Custom      []byte            `json:"custom"`
	Color       rl.Color          `json:"color"`
	Attenuation float32           `json:"atten"`
	Radius      float32           `json:"rad"`
	PolyLines   []*tiled.PolyLine `json:"polylines"`
}

func defaultSaveProvider(state *GameState) defaultSaveData {
	var gbuf bytes.Buffer
	genc := gob.NewEncoder(&gbuf)
	CurrentGameMode.Serialize(genc)
	save := defaultSaveData{
		CurrentMap:   CurrentMap.Name,
		Maps:         []defaultMapData{},
		GameModeData: gbuf.Bytes(),
	}

	for _, v := range Maps {
		mapData := defaultMapData{
			MapName:     v.Name,
			Objects:     []defaultObjectData{},
			WeatherData: v.Weather,
		}

		for _, b := range v.World.Objects {
			if !b.IsPersistent {
				continue
			}

			var buf bytes.Buffer
			enc := gob.NewEncoder(&buf)
			b.Serialize(b, enc)

			obj := defaultObjectData{
				Name:        b.Name,
				Type:        b.Class,
				Position:    b.Position,
				Movement:    b.Movement,
				Facing:      b.Facing,
				Color:       b.Color,
				Attenuation: b.Attenuation,
				Radius:      b.Radius,
				PolyLines:   b.PolyLines,
				Custom:      buf.Bytes(),
			}

			mapData.Objects = append(mapData.Objects, obj)
		}

		save.Maps = append(save.Maps, mapData)
	}

	return save
}

func defaultLoadProvider(state *GameState) {
	data := state.SaveData
	CanSave = 0
	FlushMaps()
	LoadMap(data.CurrentMap)

	gbuf := bytes.NewBuffer(data.GameModeData)
	gdec := gob.NewDecoder(gbuf)
	CurrentGameMode.Deserialize(gdec)

	for _, mapData := range data.Maps {
		m := LoadMap(mapData.MapName)
		m.Weather = mapData.WeatherData
		world := mapData.Objects

		for _, wo := range world {
			o, _ := m.World.FindObject(wo.Name)

			if o == nil {
				o = m.World.NewObjectPro(wo.Name, wo.Type)

				if o == nil {
					continue
				}

				m.World.AddObject(o)
			}

			o.Position = wo.Position
			o.Movement = wo.Movement
			o.Facing = wo.Facing
			o.Color = wo.Color
			o.Attenuation = wo.Attenuation
			o.Radius = wo.Radius
			o.PolyLines = wo.PolyLines

			buf := bytes.NewBuffer(wo.Custom)
			dec := gob.NewDecoder(buf)
			o.Deserialize(o, dec)
		}

		cam, _ := CurrentMap.World.FindObject("main_camera")

		if cam == nil {
			setupDefaultCamera()
		}

		m.World.InitObjects()
	}
}
