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

import "encoding/gob"

// GameMode describes main game rules and subsystems
type GameMode interface {
   Init()
   Shutdown()
   Update()
   Draw()
   DrawUI()
   DebugDraw()
   PostDraw()
   Serialize(enc *gob.Encoder)
   Deserialize(dec *gob.Decoder)
}
