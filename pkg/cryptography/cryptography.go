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
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

const (
	// keypair details
	keySize        = 4096 // 4k key
	privateKeyType = "RSA PRIVATE KEY"
	publicKeyType  = "CERTIFICATE"

	// certificate request details
	cn = "localhost" // common name
	o  = "stl-go"    // organization

	// ObfuscatedPrefix encrypted key prefix
	ObfuscatedPrefix string = "obf::"
)

// GenerateSecret returns a 32 byte AES key
func GenerateSecret() (*string, error) {
	key := make([]byte, 16)

	if _, err := rand.Read(key); err != nil {
		return nil, err
	}

	secret := fmt.Sprintf("%x", key)
	return &secret, nil
}

// HashPassword will hash passwords for storage
func HashPassword(password *string) (*string, error) {
	if password != nil {
		hashBytes, err := bcrypt.GenerateFromPassword([]byte(*password), 4)
		if err != nil {
			return nil, err
		}

		hash := string(hashBytes)
		return &hash, err
	}
	return nil, errors.New("password is nil")
}

// HashCompare will hash & compare a password to the stored hash
func HashCompare(original, request *string) error {
	if original != nil && request != nil {
		return bcrypt.CompareHashAndPassword([]byte(*original), []byte(*request))
	}
	return errors.New("nil input of original or requested strings")
}

// Encrypt takes plaintext and turns it into ciphertext
func Encrypt(plaintext, secret *string) (*string, error) {
	if plaintext != nil && secret != nil {
		c, err := aes.NewCipher([]byte(*secret))
		if err != nil {
			return nil, err
		}

		gcm, err := cipher.NewGCM(c)
		if err != nil {
			return nil, err
		}

		nonce := make([]byte, gcm.NonceSize())
		if _, err = rand.Read(nonce); err != nil {
			return nil, err
		}

		cipherText := fmt.Sprintf("%s%s", ObfuscatedPrefix, base64.URLEncoding.EncodeToString(gcm.Seal(nonce, nonce, []byte(*plaintext), nil)))
		return &cipherText, nil
	}
	return nil, errors.New("plaintext or secret is nil")
}

// Decrypt takes ciphertext and returns plaintext
func Decrypt(encodedData, secret *string) (*string, error) {
	if encodedData != nil && secret != nil {
		if !strings.Contains(*encodedData, ObfuscatedPrefix) {
			return encodedData, nil
		}

		cipherTextEncoded, err := base64.URLEncoding.DecodeString(strings.ReplaceAll(*encodedData, ObfuscatedPrefix, ""))
		if err != nil {
			return nil, err
		}

		cipherBlock, err := aes.NewCipher([]byte(*secret))
		if err != nil {
			return nil, err
		}

		aead, err := cipher.NewGCM(cipherBlock)
		if err != nil {
			return nil, err
		}

		nonceSize := aead.NonceSize()
		if len(cipherTextEncoded) < nonceSize {
			return nil, err
		}

		nonce, cipherText := cipherTextEncoded[:nonceSize], cipherTextEncoded[nonceSize:]
		plainTextEncoded, err := aead.Open(nil, nonce, cipherText, nil)
		if err != nil {
			return nil, err
		}

		plainText := string(plainTextEncoded)
		return &plainText, nil
	}
	return nil, errors.New("encodedData or secret is nil")
}
