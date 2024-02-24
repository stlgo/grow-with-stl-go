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
	"encoding/json"
	"errors"
	"fmt"

	"stl-go/grow-with-stl-go/pkg/audit"
	"stl-go/grow-with-stl-go/pkg/configs"
	"stl-go/grow-with-stl-go/pkg/log"
	"stl-go/grow-with-stl-go/pkg/webservice"
)

const (
	addUser          = "addUser"
	updateActive     = "updateActive"
	updateAdmin      = "updateAdmin"
	updateUser       = "updateUser"
	resetPassword    = "resetPassword"
	removeUser       = "removeUser"
	pageLoad         = "pageLoad"
	generatePassword = "generatePassword"
)

type checboxEnabled struct {
	Enabled *bool `json:"enabled,omitempty"`
}

// Init is different than the standard init because it is called outside of the object load
func Init() {
	webservice.AppendToFunctionMap("admin", handleMessage)
	// warm up the db connection so when we hit the page the first time it's faster
	if _, err := getUserInfo(); err != nil {
		log.Error(err)
	}
}

func handleMessage(_ *string, request *configs.WsMessage) *configs.WsMessage {
	response := configs.WsMessage{
		Type:         request.Type,
		Component:    request.Component,
		SubComponent: request.SubComponent,
	}

	if request.Component != nil || (request.IsAdmin == nil || !*request.IsAdmin) {
		switch *request.Component {
		case pageLoad:
			data, err := getUserInfo()
			if err != nil {
				log.Error(err)
				e := "Unable to retrieve user information"
				response.Error = &e
				return &response
			}
			response.Data = data
		case generatePassword:
			response.Data = map[string]*string{
				"password": configs.GeneratePassword(),
			}
		case addUser, updateUser, updateActive, updateAdmin, removeUser:
			userAction(request.Component, request, &response)
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

func updateUserActive(userID *string, data interface{}) error {
	if userID != nil {
		if bytes, err := json.Marshal(data); err == nil {
			var activeData checboxEnabled
			if err := json.Unmarshal(bytes, &activeData); err == nil && activeData.Enabled != nil {
				if apiUser, ok := configs.GrowSTLGo.APIUsers[*userID]; ok {
					return apiUser.ToggleActive(activeData.Enabled)
				}
			}
		}
	}
	return errors.New("cannot update user active flag")
}

func updateUserAdmin(userID *string, data interface{}) error {
	if userID != nil {
		if bytes, err := json.Marshal(data); err == nil {
			var activeData checboxEnabled
			if err := json.Unmarshal(bytes, &activeData); err == nil && activeData.Enabled != nil {
				if apiUser, ok := configs.GrowSTLGo.APIUsers[*userID]; ok {
					return apiUser.ToggleAdmin(activeData.Enabled)
				}
			}
		}
	}
	return errors.New("cannot update user admin flag")
}

func userAction(component *string, request, response *configs.WsMessage) *configs.WsMessage {
	if component != nil && request != nil && response != nil {
		var err error
		switch *component {
		case addUser:
			err = configs.AddUser(request.SubComponent, request.Data)

		case updateUser:
			err = configs.UpdateUser(request.SubComponent, request.Data)

		case updateActive:
			err = updateUserActive(request.SubComponent, request.Data)

		case updateAdmin:
			err = updateUserAdmin(request.SubComponent, request.Data)

		case removeUser:
			err = configs.RemoveUser(request.SubComponent)

		default:
			err = fmt.Errorf("component %s not implemented", *request.Component)
		}

		if err != nil {
			log.Error(err)
			e := err.Error()
			response.Error = &e
			return response
		}

		data, err := getUserInfo()
		if err != nil {
			log.Error(err)
			e := err.Error()
			response.Error = &e
			return response
		}
		response.Data = data
		return response
	}
	return nil
}

func getUserInfo() (map[string]interface{}, error) {
	data := make(map[string]interface{})
	lastLogins, err := audit.GetLastLogins()
	if err != nil {
		return nil, err
	}
	for userID, apiUser := range configs.GrowSTLGo.APIUsers {
		data[userID] = map[string]interface{}{
			"lastLogin": nil,
			"active":    apiUser.Active,
			"admin":     apiUser.Admin,
		}
		if lastLogin, ok := lastLogins[userID]; ok {
			data[userID] = map[string]interface{}{
				"lastLogin": lastLogin,
				"active":    apiUser.Active,
				"admin":     apiUser.Admin,
			}
		}
	}
	return data, nil
}
