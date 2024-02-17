/*
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     https://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package admin

import (
	"fmt"

	"stl-go/grow-with-stl-go/pkg/configs"
	"stl-go/grow-with-stl-go/pkg/log"
	"stl-go/grow-with-stl-go/pkg/webservice"
)

const (
	addUser    = "addUser"
	updateUser = "updateUser"
	removeUser = "removeUser"
	pageLoad   = "pageLoad"
)

// Init is different than the standard init because it is called outside of the object load
func Init() {
	webservice.AppendToFunctionMap("admin", handleMessage)
}

func handleMessage(_ *string, request *configs.WsMessage) *configs.WsMessage {
	response := configs.WsMessage{
		Type:         request.Type,
		Component:    request.Component,
		SubComponent: request.SubComponent,
	}

	if request.Component != nil {
		switch *request.Component {
		case pageLoad:
			response.Data = map[string]interface{}{
				"aschiefe": map[string]interface{}{
					"lastLogin": 1708194017000,
				},
				"user": map[string]interface{}{
					"lastLogin": 1708194017000,
				},
				"admin": map[string]interface{}{
					"lastLogin": 1708194017000,
				},
			}
		case addUser:
			log.Trace(addUser)
		case updateUser:
			log.Trace(updateUser)
		case removeUser:
			log.Trace(removeUser)
		default:
			err := fmt.Sprintf("component %s not implemented", *request.Component)
			log.Error(err)
			response.Error = &err
		}
		return &response
	}
	err := fmt.Errorf("bad request").Error()
	response.Error = &err
	return &response
}
