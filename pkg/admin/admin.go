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
	"slices"
	"strconv"

	"maps"

	"stl-go/grow-with-stl-go/pkg/audit"
	"stl-go/grow-with-stl-go/pkg/configs"
	"stl-go/grow-with-stl-go/pkg/locations"
	"stl-go/grow-with-stl-go/pkg/log"
	"stl-go/grow-with-stl-go/pkg/utils"
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
	if _, err := pageLoad(false); err != nil {
		return err
	}
	return nil
}

func handleMessage(request, response *configs.WsMessage) {
	if request.Type != nil && response != nil && (request.IsAdmin != nil && *request.IsAdmin) {
		var err error
		switch *request.Type {
		case pageLoadKey:
			response.Data, err = pageLoad(true)
		case generatePassword:
			response.Data = map[string]*string{
				"password": configs.GeneratePassword(),
			}
		case addUser, updateUser, updateActive, updateAdmin, removeUser:
			userAction(request, response)
		case getUserDetails:
			if request.Component != nil {
				if configs.GrowSTLGo != nil && configs.GrowSTLGo.Users != nil {
					configs.GrowSTLGo.UsersMutex.Lock()
					currentUser, ok := configs.GrowSTLGo.Users[*request.Component]
					configs.GrowSTLGo.UsersMutex.Unlock()
					if ok {
						response.Data = map[string]interface{}{
							"location": currentUser.Location,
							"vhosts":   currentUser.Vhosts,
						}
					}
				}
			}
		default:
			err = fmt.Errorf("type %s not implemented", *request.Component)
		}

		if response.Data == nil {
			e := "no data found"
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
	if userID != nil && active != nil && configs.GrowSTLGo != nil && configs.GrowSTLGo.Users != nil {
		b, err := strconv.ParseBool(*active)
		if err != nil {
			return err
		}
		configs.GrowSTLGo.UsersMutex.Lock()
		currentUser, ok := configs.GrowSTLGo.Users[*userID]
		configs.GrowSTLGo.UsersMutex.Unlock()
		if ok {
			return currentUser.ToggleActive(&b)
		}
	}
	return errors.New("cannot update user active flag")
}

func updateUserAdmin(userID, admin *string) error {
	if userID != nil && admin != nil && configs.GrowSTLGo != nil && configs.GrowSTLGo.Users != nil {
		b, err := strconv.ParseBool(*admin)
		if err != nil {
			return err
		}
		configs.GrowSTLGo.UsersMutex.Lock()
		currentUser, ok := configs.GrowSTLGo.Users[*userID]
		configs.GrowSTLGo.UsersMutex.Unlock()
		if ok {
			return currentUser.ToggleAdmin(&b)
		}
	}
	return errors.New("cannot update user admin flag")
}

func userAction(request, response *configs.WsMessage) *configs.WsMessage {
	if request != nil && request.Type != nil && response != nil {
		var err error
		switch *request.Type {
		case addUser:
			err = AddUser(request.Component, request.Data)
		case updateUser:
			err = UpdateUser(request.Component, request.Data)
		case updateActive:
			err = updateUserActive(request.Component, request.SubComponent)
		case updateAdmin:
			err = updateUserAdmin(request.Component, request.SubComponent)
		case removeUser:
			err = RemoveUser(request.Component)
		default:
			err = fmt.Errorf("component %s not implemented", *request.Component)
		}

		if err != nil {
			log.Error(err)
			e := err.Error()
			response.Error = &e
			return response
		}

		data, err := pageLoad(false)
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

func pageLoad(initialLoad bool) (map[string]interface{}, error) {
	users := make(map[string]interface{})
	err := audit.GetLastLogins()
	if err != nil {
		return nil, err
	}
	if configs.GrowSTLGo != nil && configs.GrowSTLGo.Users != nil && configs.GrowSTLGo.WebService != nil && configs.GrowSTLGo.WebService.Vhosts != nil {
		configs.GrowSTLGo.UsersMutex.Lock()
		for userID, currentUser := range configs.GrowSTLGo.Users {
			users[userID] = map[string]interface{}{
				"lastLogin": currentUser.LastLogin,
				"active":    currentUser.Active,
				"admin":     currentUser.Admin,
			}
		}
		configs.GrowSTLGo.UsersMutex.Unlock()

		data := map[string]interface{}{
			"users":   users,
			"vhosts":  slices.AppendSeq(make([]string, 0, len(configs.GrowSTLGo.WebService.Vhosts)), maps.Keys(configs.GrowSTLGo.WebService.Vhosts)),
			"version": configs.Version,
		}
		if initialLoad {
			data["zipCodes"] = locations.ZipCodeTypeahead
		}
		return data, nil
	}
	return nil, errors.New("invalid configuration cannot load page data")
}

// AddUser will add a new user to the configs
func AddUser(userID *string, data interface{}) error {
	if userID != nil && data != nil {
		if bytes, err := json.Marshal(data); err == nil {
			var user configs.User
			if err := json.Unmarshal(bytes, &user); err == nil && user.Authentication != nil && user.Authentication.ID != nil && user.Authentication.Password != nil {
				if err := user.Authentication.HashAuthentication(); err == nil {
					configs.GrowSTLGo.UsersMutex.Lock()
					_, ok := configs.GrowSTLGo.Users[*userID]
					configs.GrowSTLGo.UsersMutex.Unlock()
					if !ok {
						newUser := &configs.User{
							Active:         utils.BoolPointer(true),
							Authentication: user.Authentication,
						}

						if user.Vhosts != nil {
							for _, vhost := range user.Vhosts {
								if _, ok := configs.GrowSTLGo.WebService.Vhosts[vhost]; ok {
									newUser.Vhosts = append(newUser.Vhosts, vhost)
								}
							}
						}

						newUser.Location = user.Location
						return newUser.Persist(userID)
					}
				}
			}
		}
	}
	return errors.New("cannot add user")
}

// UpdateUser will add a new user to the configs
func UpdateUser(userID *string, data interface{}) error {
	if userID != nil && data != nil {
		if bytes, err := json.Marshal(data); err == nil {
			var user configs.User
			err := json.Unmarshal(bytes, &user)
			if err != nil {
				return fmt.Errorf("user input, cannot update user.  Error: %s", err)
			}
			configs.GrowSTLGo.UsersMutex.Lock()
			currentUser, ok := configs.GrowSTLGo.Users[*userID]
			configs.GrowSTLGo.UsersMutex.Unlock()
			if ok {
				if user.Authentication != nil && user.Authentication.Password != nil {
					if err := user.Authentication.HashAuthentication(); err != nil {
						return fmt.Errorf("bad password, cannot update user.  Error: %s", err)
					}
					currentUser.Authentication = user.Authentication
				}
				if user.Vhosts != nil {
					currentUser.Vhosts = []string{}
					for _, vhost := range user.Vhosts {
						if _, ok := configs.GrowSTLGo.WebService.Vhosts[vhost]; ok {
							currentUser.Vhosts = append(currentUser.Vhosts, vhost)
						}
					}
				}
				if user.Location != nil {
					locations.ZipCodeCacheMutex.Lock()
					_, zipOk := locations.ZipcodeLookup[*user.Location]
					locations.ZipCodeCacheMutex.Unlock()
					if zipOk {
						currentUser.Location = user.Location
					}
				}
				return currentUser.Persist(userID)
			}
		}
	}
	return errors.New("cannot update user")
}

// RemoveUser will permanently delete the user from the configs
func RemoveUser(userID *string) error {
	if userID != nil {
		configs.GrowSTLGo.UsersMutex.Lock()
		user, ok := configs.GrowSTLGo.Users[*userID]
		configs.GrowSTLGo.UsersMutex.Unlock()
		if ok {
			return user.Remove(userID)
		}
	}
	return errors.New("cannot remove user")
}
