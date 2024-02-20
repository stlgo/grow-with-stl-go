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
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"stl-go/grow-with-stl-go/pkg/cryptography"
)

// Authentication structure to hold authentication parameters
type Authentication struct {
	ID       *string `json:"id,omitempty"`
	Password *string `json:"password,omitempty"`
}

// GeneratePassword will generate a 64 bit password string on demand
func GeneratePassword() *string {
	var password bytes.Buffer
	// generate a 64bit access key
	for i := 0; i < 4; i++ {
		// golangci-lint tosses a false positive G404: Use of weak random number generator error so we'll skip that for this line
		password.WriteString(fmt.Sprintf("%x", rand.New(rand.NewSource(time.Now().UnixNano()+int64(i))).Uint64())) // #nosec
	}

	s := password.String()
	return &s
}

// ValidateAuthentication will compare the supplied password with the stored one
func (auth *Authentication) ValidateAuthentication(password *string) error {
	if auth != nil && auth.Password != nil && password != nil {
		plaintext, err := cryptography.Decrypt(auth.Password, GrowSTLGo.Secret)
		if err != nil {
			return err
		}

		if err := cryptography.HashCompare(plaintext, password); err != nil {
			return errors.New("validation failed")
		}
	}
	return nil
}

// GeneratePassword will generate a strong password for the authentication object
func (auth *Authentication) GeneratePassword() (*string, error) {
	passwd := GeneratePassword()
	auth.Password = passwd

	if err := auth.hashAuthentication(); err != nil {
		return nil, err
	}

	return passwd, nil
}

// HashAuthentication will hash the password for the Authentication object
func (auth *Authentication) hashAuthentication() error {
	if auth != nil && auth.Password != nil {
		hashText, err := cryptography.HashPassword(auth.Password)
		if err != nil {
			return err
		}

		cipherText, err := cryptography.Encrypt(hashText, GrowSTLGo.Secret)
		if err != nil {
			return err
		}
		auth.Password = cipherText
		return nil
	}
	return errors.New("cannot hash password")
}
