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
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCountryFunctions(t *testing.T) {
	t.Skip()
	initConfigTest()
	t.Run("Download file", func(t *testing.T) {
		err := SetGrowSTLGoConfig()
		require.NoError(t, err)

		if GrowSTLGo.Country != nil {
			err := GrowSTLGo.Country.GetCountryData()
			require.NoError(t, err)
		}
	})
}
