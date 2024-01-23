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

func TestCryptographyFunctions(t *testing.T) {
	t.Run("Test generating a secret", func(t *testing.T) {
		secret, err := GenerateSecret()
		require.NoError(t, err)
		require.NotNil(t, secret)
		t.Logf("%s\n", *secret)
	})

	t.Run("Test encryption", func(t *testing.T) {
		key := "the-key-has-to-be-32-bytes-long!"
		plaintext := "passphrase"

		ciphertext, err := Encrypt(&plaintext, &key)
		require.NoError(t, err)
		require.NotNil(t, ciphertext)
		t.Logf("%s => %s\n", plaintext, *ciphertext)
	})

	t.Run("Test decryption", func(t *testing.T) {
		key := "the-key-has-to-be-32-bytes-long!"
		ciphertext := "obf::nxQ1cfkfsmXWxXeehy1wJuHiGKUZxmZpiBI7FFHciIrwzFHmenk="

		plaintext, err := Decrypt(&ciphertext, &key)
		require.NoError(t, err)
		require.NotNil(t, ciphertext)
		t.Logf("%s => %s\n", ciphertext, *plaintext)
	})

	t.Run("Test password hashing", func(t *testing.T) {
		input := "some user provided password"
		hash, err := HashPassword(&input)

		require.NoError(t, err)
		require.NotNil(t, hash)
		t.Logf("%s => %s\n", input, *hash)
	})

	t.Run("Test hash compare", func(t *testing.T) {
		input := "some user provided password"
		hash := "$2a$04$mXwoCJ5w9o31XA8YHypueuZfuvCMC7Dz2T79oof7WKjLS4zK85klu"

		err := HashCompare(&hash, &input)
		require.NoError(t, err)
		t.Logf("%s compared true to %s\n", input, hash)
	})
}
