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

// InitDatabase prepares the database connection
func InitDatabase() {
   CurrentSaveSystem.InitSaveSystem()
   Quests = MakeQuestManager()
}

// ShutdownDatabase disposes all database data and closes all connections
func ShutdownDatabase() {
   CurrentSaveSystem.ShutdownSaveSystem()
}
