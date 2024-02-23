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
	"testing"

	"github.com/stretchr/testify/require"

	"stl-go/grow-with-stl-go/pkg/configs"
	"stl-go/grow-with-stl-go/pkg/log"
)

func initConfigTest() {
	configFile := "../../etc/grow-with-stl-go.json"
	configs.ConfigFile = &configFile
	if err := configs.SetGrowSTLGoConfig(); err != nil {
		log.Fatal(err)
	}
}

func TestConfigFunctions(t *testing.T) {
	initConfigTest()
	t.Run("Test setting the config", func(t *testing.T) {
		data, err := getUserInfo()
		require.NoError(t, err)
		require.NotEmpty(t, data)

		for key, value := range data {
			log.Infof("%s %v", key, value)
		}
	})
}
