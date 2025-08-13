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
	"slices"
	"strconv"

	"maps"

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
	pageLoadKey      = "pageLoad"
	generatePassword = "generatePassword"
	getUserDetails   = "getUserDetails"
)

// Init is different than the standard init because it is called outside of the object load
func Init() error {
	webservice.AppendToWebsocketFunctionMap("admin", handleMessage)
	// warm up the db connection so when we hit the page the first time it's faster
	if _, err := pageLoad(); err != nil {
		return err
	}
	return nil
}

func handleMessage(request, response *configs.WsMessage) {
	if request.Type != nil && response != nil && (request.IsAdmin != nil && *request.IsAdmin) {
		var err error
		switch *request.Type {
		case pageLoadKey:
			response.Data, err = pageLoad()
		case generatePassword:
			response.Data = map[string]*string{
				"password": configs.GeneratePassword(),
			}
		case addUser, updateUser, updateActive, updateAdmin, removeUser:
			userAction(request, response)
		case getUserDetails:
			if request.Component != nil {
				configs.GrowSTLGo.APIUsersMutex.Lock()
				apiUser, ok := configs.GrowSTLGo.APIUsers[*request.Component]
				configs.GrowSTLGo.APIUsersMutex.Unlock()
				if ok {
					response.Data = apiUser.Vhosts
				}
			}
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
		return
	}
	log.Error("bad request, nothing can be done")
}

func updateUserActive(userID, active *string) error {
	if userID != nil && active != nil {
		b, err := strconv.ParseBool(*active)
		if err != nil {
			return err
		}
		configs.GrowSTLGo.APIUsersMutex.Lock()
		apiUser, ok := configs.GrowSTLGo.APIUsers[*userID]
		configs.GrowSTLGo.APIUsersMutex.Unlock()
		if ok {
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
		configs.GrowSTLGo.APIUsersMutex.Lock()
		apiUser, ok := configs.GrowSTLGo.APIUsers[*userID]
		configs.GrowSTLGo.APIUsersMutex.Unlock()
		if ok {
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

		data, err := pageLoad()
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

func pageLoad() (map[string]interface{}, error) {
	users := make(map[string]interface{})
	err := audit.GetLastLogins()
	if err != nil {
		return nil, err
	}
	configs.GrowSTLGo.APIUsersMutex.Lock()
	for userID, apiUser := range configs.GrowSTLGo.APIUsers {
		users[userID] = map[string]interface{}{
			"lastLogin": apiUser.LastLogin,
			"active":    apiUser.Active,
			"admin":     apiUser.Admin,
		}
	}
	configs.GrowSTLGo.APIUsersMutex.Unlock()

	data := map[string]interface{}{
		"users":   users,
		"vhosts":  slices.AppendSeq(make([]string, 0, len(configs.GrowSTLGo.WebService.Vhosts)), maps.Keys(configs.GrowSTLGo.WebService.Vhosts)),
		"version": configs.Version,
	}
	return data, nil
}
