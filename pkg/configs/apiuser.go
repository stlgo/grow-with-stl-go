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

package configs

import (
	"encoding/json"
	"errors"

	"stl-go/grow-with-stl-go/pkg/log"
	"stl-go/grow-with-stl-go/pkg/utils"
)

// APIUser is our storage point for REST & WebSocket users
type APIUser struct {
	Active         *bool           `json:"active,omitempty"`
	Authentication *Authentication `json:"authentication,omitempty"`
	LastLogin      *int64          `json:"lastLogin,omitempty"`
	Admin          *bool           `json:"admin,omitempty"`
}

func checkAPIUsers() {
	if GrowSTLGo.APIUsers == nil {
		ids := map[string]bool{
			"admin": true,
			"user":  false,
		}

		GrowSTLGo.APIUsers = map[string]*APIUser{}

		for id, isAdmin := range ids {
			localID := id
			user := APIUser{
				Active: utils.BoolPointer(isAdmin),

				Authentication: &Authentication{
					ID: &localID,
				},
			}

			if password, err := user.Authentication.GeneratePassword(); err == nil && password != nil {
				log.Warnf("Password generated for user '%s', password %s - DO NOT USE THIS FOR PRODUCTION", localID, *password)

				GrowSTLGo.APIUsers[localID] = &user
			}
		}
		rewriteConfig = true
	}
}

// AddUser will add a new user to the configs
func AddUser(userID *string, data interface{}) error {
	if userID != nil && data != nil {
		if bytes, err := json.Marshal(data); err == nil {
			var authentication Authentication
			if err := json.Unmarshal(bytes, &authentication); err == nil && authentication.ID != nil && authentication.Password != nil {
				if err := authentication.hashAuthentication(); err == nil {
					if _, ok := GrowSTLGo.APIUsers[*userID]; !ok {
						GrowSTLGo.APIUsers[*userID] = &APIUser{
							Active:         utils.BoolPointer(true),
							Authentication: &authentication,
						}
						return GrowSTLGo.persist()
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
			var authentication Authentication
			if err := json.Unmarshal(bytes, &authentication); err == nil && authentication.ID != nil && authentication.Password != nil {
				if err := authentication.hashAuthentication(); err == nil {
					if apiUser, ok := GrowSTLGo.APIUsers[*userID]; ok {
						apiUser.Authentication = &authentication
						return GrowSTLGo.persist()
					}
				}
			}
		}
	}
	return errors.New("cannot update user")
}

// RemoveUser will permanently delete the user from the configs
func RemoveUser(userID *string) error {
	if userID != nil {
		if _, ok := GrowSTLGo.APIUsers[*userID]; ok {
			delete(GrowSTLGo.APIUsers, *userID)
			return GrowSTLGo.persist()
		}
	}
	return errors.New("cannot remove user")
}

// ResetPassword will reset the password of a given user
func (apiUser *APIUser) ResetPassword(password *string) error {
	if apiUser != nil && apiUser.Authentication != nil && password != nil {
		backupAuth := apiUser.Authentication
		apiUser.Authentication.Password = password
		if err := apiUser.Authentication.hashAuthentication(); err != nil {
			apiUser.Authentication = backupAuth
			return err
		}
		return GrowSTLGo.persist()
	}
	return errors.New("unable to reset password: nil api user or nill password")
}

// ToggleActive will set user to enabled / disabled based on input
func (apiUser *APIUser) ToggleActive(active *bool) error {
	if apiUser != nil && active != nil {
		apiUser.Active = active
		log.Info(*active)
		return GrowSTLGo.persist()
	}
	return errors.New("unable to set activity: nil api user or nil active boolean")
}
