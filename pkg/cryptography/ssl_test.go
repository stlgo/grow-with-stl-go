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

package cryptography

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSSLFunctions(t *testing.T) {
	t.Run("Test generating private key", func(t *testing.T) {
		pem, key, err := generatePrivateKey()
		require.NoError(t, err)
		require.NotNil(t, key)
		require.NotNil(t, pem)
	})

	t.Run("Test generating public key", func(t *testing.T) {
		_, privateKey, err := generatePrivateKey()
		require.NoError(t, err)

		cert, err := generatePublicKey(privateKey)
		require.NoError(t, err)
		require.NotNil(t, cert)
	})

	t.Run("Test cert validity", func(t *testing.T) {
		_, privateKey, err := generatePrivateKey()
		require.NoError(t, err)

		cert, err := generatePublicKey(privateKey)
		require.NoError(t, err)
		require.NotNil(t, cert)
	})
}
