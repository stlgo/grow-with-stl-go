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
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"time"

	"stl-go/grow-with-stl-go/pkg/cryptography"
	"stl-go/grow-with-stl-go/pkg/log"
	"stl-go/grow-with-stl-go/pkg/utils"
)

// Table is a struct used by various packages to create insert / update data to the embedded db
type Table struct {
	CreateSQL  string
	InsertSQL  string
	Indices    []string
	Defaults   map[string]string
	WriteMutex sync.Mutex
}

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
		log.Debug("No embedded database found in the config file, generating a default configuration")
		fileName := filepath.Join(*etcDir, "grow-with-stl-go.db")
		sqlite := SQLite{
			FileName:        &fileName,
			EncryptDatabase: utils.BoolPointer(true),
		}
		if err := sqlite.generateEncryptionKeys(); err != nil {
			return err
		}

		GrowSTLGo.SQLite = &sqlite
		rewriteConfig = true
	}
	return startDB()
}

// Init will setup the default tables in the sqlite embedded database
func startDB() error {
	if GrowSTLGo.SQLite != nil && GrowSTLGo.SQLite.FileName != nil {
		baseDir := filepath.Dir(*GrowSTLGo.SQLite.FileName)
		if err := os.MkdirAll(baseDir, 0o750); err != nil {
			return err
		}

		// encrypted db
		if GrowSTLGo.SQLite.EncryptionKey != nil {
			if key, err := GrowSTLGo.SQLite.GetEncryptionKey(); err == nil && key != nil {
				dbname := fmt.Sprintf("%s?_pragma_key=x'%s'&_pragma_cipher_page_size=4096", *GrowSTLGo.SQLite.FileName, *key)
				SqliteDB, err = sql.Open("sqlite3", dbname)
				if err != nil {
					return err
				}
				log.Debugf("Using encrypted aud database in %s", *GrowSTLGo.SQLite.FileName)
				return nil
			}
		}
		db, err := sql.Open("sqlite3", *GrowSTLGo.SQLite.FileName)
		if err != nil {
			return err
		}
		SqliteDB = db
		log.Debugf("Using unencrypted aud database in %s", *GrowSTLGo.SQLite.FileName)
		return nil
	}
	return fmt.Errorf("no sutiable configuration for SQLite found in the config file %s", *ConfigFile)
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

// CreateTable is a helper function that will create the table & indexes in the embedded database
func (table *Table) CreateTable(tableName *string) error {
	if table != nil && tableName != nil {
		log.Tracef("Audit table %s was created if it didn't already exist", *tableName)
		stmt, err := SqliteDB.Prepare(table.CreateSQL)
		if err != nil {
			return err
		}
		if _, err = stmt.Exec(); err != nil {
			return err
		}
		for _, index := range table.Indices {
			stmt, err := SqliteDB.Prepare(index)
			if err != nil {
				return err
			}
			if _, err = stmt.Exec(); err != nil {
				return err
			}
		}
		if table.Defaults != nil {
			for key, sql := range table.Defaults {
				log.Tracef("Inserting default %s for table %s if it doesn't already exist", key, *tableName)
				stmt, err := SqliteDB.Prepare(sql)
				if err != nil {
					return err
				}
				if _, err = stmt.Exec(); err != nil {
					return err
				}
			}
		}
		return nil
	}
	return errors.New("tables is nil cannot create tables")
}
