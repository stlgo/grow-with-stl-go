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
	"errors"
	"fmt"
	"strconv"

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

// Init is different than the standard init because it is called outside of the object load
func Init() {
	webservice.AppendToWebsocketFunctionMap("admin", handleMessage)
	// warm up the db connection so when we hit the page the first time it's faster
	if _, err := getUserInfo(); err != nil {
		log.Error(err)
	}
}

func handleMessage(_ *string, request *configs.WsMessage) *configs.WsMessage {
	response := configs.WsMessage{
		Route:        request.Route,
		Type:         request.Type,
		Component:    request.Component,
		SubComponent: request.SubComponent,
	}

	if request.Type != nil || (request.IsAdmin == nil || !*request.IsAdmin) {
		var err error
		switch *request.Type {
		case pageLoad:
			response.Data, err = getUserInfo()
		case generatePassword:
			response.Data = map[string]*string{
				"password": configs.GeneratePassword(),
			}
		case addUser, updateUser, updateActive, updateAdmin, removeUser:
			userAction(request, &response)
		default:
			err = fmt.Errorf("type %s not implemented", *request.Component)
		}

		if response.Data == nil {
			e := "no data round"
			log.Error(e)
			response.Error = &e
		}

		if err != nil {
			log.Error(err)
			e := err.Error()
			response.Error = &e
			response.Data = nil
		}

		return &response
	}
	err := fmt.Errorf("bad request").Error()
	response.Error = &err
	return &response
}

func updateUserActive(userID, active *string) error {
	if userID != nil && active != nil {
		b, err := strconv.ParseBool(*active)
		if err != nil {
			return err
		}
		if apiUser, ok := configs.GrowSTLGo.APIUsers[*userID]; ok {
			return apiUser.ToggleActive(&b)
		}
	}
	return errors.New("cannot update user active flag")
}

func updateUserAdmin(userID, admin *string) error {
	if userID != nil && admin != nil {
		b, err := strconv.ParseBool(*admin)
		if err != nil {
			return err
		}
		if apiUser, ok := configs.GrowSTLGo.APIUsers[*userID]; ok {
			return apiUser.ToggleAdmin(&b)
		}
	}
	return errors.New("cannot update user admin flag")
}

func userAction(request, response *configs.WsMessage) *configs.WsMessage {
	if request != nil && request.Type != nil && response != nil {
		var err error
		switch *request.Type {
		case addUser:
			err = configs.AddUser(request.Component, request.Data)
		case updateUser:
			err = configs.UpdateUser(request.Component, request.Data)
		case updateActive:
			err = updateUserActive(request.Component, request.SubComponent)
		case updateAdmin:
			err = updateUserAdmin(request.Component, request.SubComponent)
		case removeUser:
			err = configs.RemoveUser(request.Component)
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
