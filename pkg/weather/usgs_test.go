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

package weather

import (
	"testing"

	"github.com/stretchr/testify/require"

	"stl-go/grow-with-stl-go/pkg/configs"
	"stl-go/grow-with-stl-go/pkg/log"
)

func TestUSGSFunctions(t *testing.T) {
	t.Skip()
	configs.InitTest()

	t.Run("Test USGS Lookup", func(t *testing.T) {
		lat := -90.194509984085
		long := 38.627847981932
		resp, err := usgsSiteNumberLookup(&lat, &long)

		if err != nil {
			log.Error(err)
		}
		require.NoError(t, err)

		if resp != nil {
			log.Info(*resp)
		}
		require.NotNil(t, resp)
	})
}
