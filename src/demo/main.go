package main

import (
	"flag"
	"fmt"
	"path"
	"strings"

	rl "github.com/zaklaus/raylib-go/raylib"
	"github.com/zaklaus/rurik/src/core"
	"github.com/zaklaus/rurik/src/system"
)

const (
	screenW = 640
	screenH = 480

	// Apply 2x upscaling
	windowW = screenW * 2
	windowH = screenH * 2
)

var (
	dynobjCounter int
	fmapload      bool
	playMapName   string
	bloom         *bloomProg
	shadertoy     *shadertoyProg
	minimap       *minimapProg
	pulseManager  *core.Object
)

const (
	stateMenu = iota
	statePlay
)

func init() {
	rl.SetCallbackFunc(main)
}

func main() {
	dbgMode := flag.Int("debug", 1, "Enable/disable debug mode. Works only in debug builds!")
	musicVol := flag.Int("musicvol", 10, "Music volume.")
	weatherTimeScale := flag.Float64("wtimescale", 1, "Weather time scale.")
	mapName := flag.String("map", "village", "Map name to play.")
	enableProfiler := flag.Bool("profile", false, "Enable profiling.")
	forceMapLoad := flag.Bool("forceload", false, "Forces map load and skips the title screen.")
	flag.Parse()

	playMapName = *mapName
	playMapName = path.Base(playMapName)
	playMapName = strings.Split(playMapName, ".")[0]
	fmapload = *forceMapLoad

	if core.DebugMode {
		core.DebugMode = *dbgMode == 1
	}

	core.InitCore("Demo game | Rurik Engine", windowW, windowH, screenW, screenH)

	demoGame := &demoGameMode{}

	core.SetMusicVolume(float32(*musicVol) / 100)
	core.WeatherTimeScale = *weatherTimeScale

	core.Run(demoGame, *enableProfiler)
}

func initShaders() {
	bloom = newBloom()
	shadertoy = newShadertoy()
	minimap = newMinimap()
}

func updateInternals(g *demoGameMode) {
	if pulseManager == nil {
		/* w := core.CurrentMap.World
		pulseManager = w.NewObjectPro("pulse_script", "script")
		pulseManager.FileName = "pulsatingLights.js"
		w.FinalizeObject(pulseManager) */
		//pulseManager.Trigger(pulseManager, nil)
	}

	if core.DebugMode && rl.IsKeyPressed(rl.KeyF5) {
		core.FlushMaps()
		core.LoadMap(playMapName)
		core.InitMap()
		g.playState = statePlay
		return
	}

	if core.DebugMode && rl.IsKeyPressed(rl.KeyF7) {
		core.LoadNextTrack()
	}

	if core.DebugMode && rl.IsKeyPressed(rl.KeyF2) {
		core.CurrentSaveSystem.SaveGame(0, "demo")
	}

	if core.DebugMode && rl.IsKeyPressed(rl.KeyF3) {
		core.CurrentSaveSystem.LoadGame(0)
	}

	if core.DebugMode && rl.IsKeyPressed(rl.KeyF9) {
		if core.CurrentMap != nil {
			w := core.CurrentMap.World
			name := fmt.Sprintf("dynobj:%d", dynobjCounter)
			newObject := w.NewObjectPro(name, "target")
			newObject.Position = core.LocalPlayer.Position
			newObject.Color = rl.Pink
			newObject.HasLight = true
			newObject.HasSpecularLight = true
			newObject.Attenuation = 256
			newObject.Radius = 8
			w.AddObject(newObject)

			dynobjCounter++
		}
	}

	if core.DebugMode {
		wheel := rl.GetMouseWheelMove()
		if wheel != 0 {
			core.MainCamera.SetCameraZoom(core.MainCamera.Zoom + float32(wheel)*0.12)
		}
	}

	core.RenderCamera.Zoom = core.MainCamera.Zoom

	if core.WindowWasResized {
		initShaders()
	}
}

func drawBackground() {
	bgImage := system.GetTexture("bg.png")

	rows := int(system.ScreenWidth/bgImage.Width) + 1
	cols := int(system.ScreenHeight/bgImage.Height) + 1
	tileW := float32(bgImage.Width)
	tileH := float32(bgImage.Height)
	src := rl.NewRectangle(0, 0, tileW, tileH)

	for i := 0; i < cols; i++ {
		for j := 0; j < rows; j++ {
			rl.DrawTexturePro(
				*bgImage,
				src,
				rl.NewRectangle(float32(j)*tileW, float32(i)*tileH, tileW, tileH),
				rl.Vector2{},
				0,
				rl.White,
			)
		}
	}
}