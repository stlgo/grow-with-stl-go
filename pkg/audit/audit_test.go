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

package audit

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"stl-go/grow-with-stl-go/pkg/configs"
	"stl-go/grow-with-stl-go/pkg/log"
	"stl-go/grow-with-stl-go/pkg/utils"
)

func initConfigTest() {
	configFile := "../../etc/grow-with-stl-go.json"
	configs.ConfigFile = &configFile
	if err := configs.SetGrowSTLGoConfig(); err != nil {
		log.Fatalf("config %s", err)
	}
	if err := Init(); err != nil {
		log.Fatalf("Error starting the audit db: %s", err)
	}
}

func TestConfigFunctions(t *testing.T) {
	t.Skip()
	initConfigTest()
	t.Run("Test websocket transaction", func(t *testing.T) {
		host := "localhost"
		user := "Charlie"
		messageType := "users"
		messageComponent := "pageLoad"
		messageSubComponent := "getUsers"
		wsRequest := configs.WsMessage{
			Type:         &messageType,
			Component:    &messageComponent,
			SubComponent: &messageSubComponent,
		}
		transaction := NewWSTransaction(&host, &user, &wsRequest)
		require.NotNil(t, transaction)
		transaction.Recordable = utils.BoolPointer(true)
		err := transaction.Complete(true)
		require.NoError(t, err)
	})

	t.Run("Test REST transaction", func(t *testing.T) {
		host := "localhost"
		user := "Charlie"
		uri := "/foo/bar/glitch"
		method := "PATCH"
		transaction := NewRESTTransaction(&host, &uri, &method)
		require.NotNil(t, transaction)
		transaction.Recordable = utils.BoolPointer(true)
		transaction.User = &user
		err := transaction.Complete(http.StatusUnauthorized)
		require.NoError(t, err)
	})
}
