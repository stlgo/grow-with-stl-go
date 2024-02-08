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
