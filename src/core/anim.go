/*
   Copyright 2018 V4 Games

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
	jsoniter "github.com/json-iterator/go"
	rl "github.com/zaklaus/raylib-go/raylib"
	"github.com/zaklaus/rurik/src/system"
)

type anim struct {
	AnimTag             string
	animStarted         bool
	PendingCurrentFrame int32
}

type animData struct {
	AnimTag      string `json:"animTag"`
	CurrentFrame int32  `json:"cframe"`
}

// NewAnim animated sprite
func (o *Object) NewAnim() {
	o.Trigger = func(o, inst *Object) {
		if o.Ase != nil {
			o.Ase.Play(o.AnimTag)

			if o.PendingCurrentFrame != -1 {
				o.Ase.CurrentFrame = o.PendingCurrentFrame
				o.PendingCurrentFrame = -1
			}
		}
		o.animStarted = true
	}

	o.Update = func(o *Object, dt float32) {
		if o.animStarted && o.Proxy == nil {
			o.Ase.Update(dt)
		}
	}

	o.Serialize = func(o *Object) string {
		if o.Ase == nil {
			return "{}"
		}

		data := animData{
			AnimTag:      o.AnimTag,
			CurrentFrame: o.Ase.CurrentFrame,
		}

		ret, _ := jsoniter.MarshalToString(&data)
		return ret
	}

	o.Deserialize = func(o *Object, data string) {
		inp := animData{}
		jsoniter.UnmarshalFromString(data, &inp)

		o.AnimTag = inp.AnimTag
		o.PendingCurrentFrame = inp.CurrentFrame
	}

	o.Finish = func(o *Object) {
		o.Texture = system.GetTexture(o.FileName + ".png")

		if o.Proxy != nil {
			o.Ase = o.Proxy.Ase
		} else {
			if o.AnimTag == "" {
				o.AnimTag = o.Meta.Properties.GetString("tag")
			}

			aseData := system.GetAnimData(o.FileName)
			o.Ase = &aseData
		}

		if o.AutoStart {
			o.Trigger(o, nil)
		}
	}

	o.GetAABB = getSpriteAABB

	o.Draw = func(o *Object) {
		if o.Ase == nil {
			return
		}

		source := getSpriteRectangle(o)
		dest := getSpriteOrigin(o)

		if DebugMode && o.DebugVisible {
			c := getSpriteAABB(o)
			rl.DrawRectangleLinesEx(c.ToFloat32(), 1, rl.Blue)
			drawTextCentered(o.Name, c.X+c.Width/2, c.Y+c.Height+2, 1, rl.White)
		}

		rl.DrawTexturePro(*o.Texture, source, dest, rl.Vector2{}, o.Rotation, SkyColor)
	}
}
