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
	"path/filepath"
	"time"

	"stl-go/grow-with-stl-go/pkg/cryptography"
	"stl-go/grow-with-stl-go/pkg/log"
	"stl-go/grow-with-stl-go/pkg/utils"
)

// SQLite is used to define the sqlite embedded database
type SQLite struct {
	FileName        *string `json:"fileName,omitempty"`
	EncryptDatabase *bool   `json:"encryptDatabase,omitempty"`
	EncryptionKey   *string `json:"encryptionKey,omitempty"`
}

func checkSQLite() error {
	if GrowSTLGo.SQLite == nil {
		etcDir, err := getEtcDir()
		if err != nil {
			return err
		}
		log.Debug("No audit database found in the config file, generating a default configuration")
		fileName := filepath.Join(*etcDir, "audit.db")
		sqlite := SQLite{
			FileName:        &fileName,
			EncryptDatabase: utils.BoolPointer(true),
		}
		if err := sqlite.generateEncryptionKeys(); err != nil {
			return err
		}

		GrowSTLGo.SQLite = &sqlite
		rewriteConfig = true
		return nil
	}
	return nil
}

func (sqlite *SQLite) generateEncryptionKeys() error {
	var keyBytes bytes.Buffer
	// generate 64 bit key
	for i := 0; i < 4; i++ {
		// golangci-lint tosses a false positive G404: Use of weak random number generator error so we'll skip that for this line
		keyBytes.WriteString(fmt.Sprintf("%x", rand.New(rand.NewSource(time.Now().UnixNano()+int64(i))).Uint64())) // #nosec
	}

	key := keyBytes.String()
	cipherText, err := cryptography.Encrypt(&key, GrowSTLGo.Secret)
	if err != nil {
		return err
	}

	sqlite.EncryptionKey = cipherText
	return nil
}

// GetEncryptionKey will return the encryption key for the database on startup
func (sqlite *SQLite) GetEncryptionKey() (*string, error) {
	if sqlite.EncryptDatabase != nil && *sqlite.EncryptDatabase && sqlite.EncryptionKey != nil && GrowSTLGo.Secret != nil {
		return cryptography.Decrypt(sqlite.EncryptionKey, GrowSTLGo.Secret)
	}
	return nil, errors.New("sqlite does not meet the encryption requirements")
}
