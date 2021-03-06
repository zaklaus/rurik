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
	"encoding/gob"
	"log"

	"github.com/zaklaus/rurik/src/system"
)

type script struct {
	WasExecuted bool
	CanRepeat   bool
	Source      string
}

type scriptData struct {
	objectData
	WasExecuted bool `json:"done"`
	CanRepeat   bool `json:"rep"`
}

// NewScript sequence script
func (o *Object) NewScript() {
	o.Trigger = func(o, inst *Object) {
		if o.FileName == "" {
			o.FileName = o.Meta.Properties.GetString("file")
		}

		data := system.GetFile("scripts/"+o.FileName, true)

		if data == nil {
			return
		}

		o.Source = string(data)

		log.Printf("Loading script %s...\n", o.FileName)

		if !o.WasExecuted || o.CanRepeat {
			updateScriptingContext()
			ScriptingContext.Set("CurrentWorld", o.world)
			ScriptingContext.Set("Self", o)
			ScriptingContext.Set("Instigator", inst)

			scriptingProfiler.StartInvocation()
			_, err := ScriptingContext.Eval(o.Source)
			scriptingProfiler.StopInvocation()

			if err != nil {
				log.Fatalf("Script error detected at '%s':%s: \n\t%s!\n", o.Name, o.FileName, err.Error())
				return
			}
		}

		o.WasExecuted = true
	}

	o.Serialize = func(o *Object, enc *gob.Encoder) {
		enc.Encode(&scriptData{
			WasExecuted: o.WasExecuted,
			CanRepeat:   o.CanRepeat,
		})
	}

	o.Deserialize = func(o *Object, dec *gob.Decoder) {
		var dat scriptData
		dec.Decode(&dat)
		o.WasExecuted = dat.WasExecuted
		o.CanRepeat = dat.CanRepeat
	}

	o.Finish = func(o *Object) {
		if o.AutoStart {
			o.Trigger(o, nil)
		}
	}
}
