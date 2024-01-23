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

func TestAuthenticationFunctions(t *testing.T) {
	id := "Charlie"
	var testPass *string
	var auth = Authentication{
		ID: &id,
	}

	s := "the-key-has-to-be-32-bytes-long!"
	GrowSTLGo = Config{
		Secret: &s,
	}

	t.Run("Test generating a password", func(t *testing.T) {
		password, err := auth.GeneratePassword(true)
		require.NoError(t, err)
		require.NotNil(t, password)
		t.Logf("generated password %s hashed encrypted password %s\n", *password, *auth.Password)
		testPass = password
	})

	t.Run("Test validating authentication", func(t *testing.T) {
		err := auth.ValidateAuthentication(testPass)
		require.NoError(t, err)
		t.Logf("password %s validated vas validated against the stored password\n", *testPass)
	})
}
