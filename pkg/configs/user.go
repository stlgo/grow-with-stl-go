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
	"errors"

	"stl-go/grow-with-stl-go/pkg/log"
	"stl-go/grow-with-stl-go/pkg/utils"
)

// User is our storage point for REST & WebSocket users
type User struct {
	Active         *bool           `json:"active,omitempty"`
	Admin          *bool           `json:"admin,omitempty"`
	Authentication *Authentication `json:"authentication,omitempty"`
	LastLogin      *int64          `json:"lastLogin,omitempty"`
	Location       *string         `json:"location,omitempty"`
	Vhosts         []string        `json:"vhosts,omitempty"`
}

func checkUsers() {
	if GrowSTLGo.Users == nil {
		ids := map[string]User{
			"admin": {
				Admin:    utils.BoolPointer(true),
				Location: utils.StringPointer("Saint Louis, MO. 63101"),
			},
			"user1": {
				Admin:    utils.BoolPointer(false),
				Location: utils.StringPointer("Springfield, IL. 62764"),
			},
			"user2": {
				Admin:    utils.BoolPointer(false),
				Location: utils.StringPointer("Springfield, MO. 65890"),
			},
		}

		GrowSTLGo.Users = map[string]*User{}

		for id, protoUser := range ids {
			user := User{
				Active: utils.BoolPointer(true),
				Admin:  protoUser.Admin,
				Authentication: &Authentication{
					ID: &id,
				},
				Location: protoUser.Location,
				Vhosts:   []string{"localhost", "grow-with-stlgo.localdev.org"},
			}

			if protoUser.Admin != nil && *protoUser.Admin {
				user.Vhosts = append(user.Vhosts, "grow-with-stlgo-admin.localdev.org")
			}

			if password, err := user.Authentication.GeneratePassword(); err == nil && password != nil {
				log.Warnf("Password generated for user '%s', password %s - DO NOT USE THIS FOR PRODUCTION", id, *password)

				GrowSTLGo.UsersMutex.Lock()
				GrowSTLGo.Users[id] = &user
				GrowSTLGo.UsersMutex.Unlock()
			}
		}
		rewriteConfig = true
	}
}

// ResetPassword will reset the password of a given user
func (u *User) ResetPassword(password *string) error {
	if u != nil && u.Authentication != nil && password != nil {
		backupAuth := u.Authentication
		u.Authentication.Password = password
		if err := u.Authentication.HashAuthentication(); err != nil {
			u.Authentication = backupAuth
			return err
		}
		return GrowSTLGo.persist()
	}
	return errors.New("unable to reset password: nil api user or nill password")
}

// ToggleActive will set user to enabled / disabled based on input
func (u *User) ToggleActive(active *bool) error {
	if u != nil && active != nil {
		u.Active = active
		return GrowSTLGo.persist()
	}
	return errors.New("unable to set activity: nil api user or nil active boolean")
}

// ToggleAdmin will set user to enabled / disabled based on input
func (u *User) ToggleAdmin(admin *bool) error {
	if u != nil && admin != nil {
		u.Admin = admin
		return GrowSTLGo.persist()
	}
	return errors.New("unable to set admin: nil api user or nil active boolean")
}

// Remove will remove the user
func (u *User) Remove(userID *string) error {
	if u != nil && userID != nil {
		delete(GrowSTLGo.Users, *userID)
		log.Tracef("user '%s' has been removed", *userID)
		return GrowSTLGo.persist()
	}
	return errors.New("unable to remove user")
}

// Persist will set the user
func (u *User) Persist(userID *string) error {
	if u != nil && userID != nil {
		GrowSTLGo.UsersMutex.Lock()
		GrowSTLGo.Users[*userID] = u
		GrowSTLGo.UsersMutex.Unlock()
		return GrowSTLGo.persist()
	}
	return errors.New("persist user")
}
