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
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"stl-go/grow-with-stl-go/pkg/log"
)

func TestHTTPFunctions(t *testing.T) {
	t.Skip()
	initConfigTest()
	t.Run("Download file", func(t *testing.T) {
		err := SetGrowSTLGoConfig()
		require.NoError(t, err)

		statusCode, err := DownloadFile("https://download.geonames.org/export/zip/US.zip", http.MethodGet, "c:/temp/test.download.zip", nil)
		require.NoError(t, err)
		require.NotNil(t, statusCode)
		if err == nil && statusCode != nil {
			log.Info(*statusCode)
		}
	})
}
