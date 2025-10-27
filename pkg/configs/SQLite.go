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
	"context"
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
	UpdateSQL  string
	Indices    []string
	Defaults   map[string]string
	WriteMutex sync.Mutex
}

// SQLite is used to define the sqlite embedded database
type SQLite struct {
	FileName        *string             `json:"fileName,omitempty"`
	EncryptDatabase *bool               `json:"encryptDatabase,omitempty"`
	EncryptionKey   *string             `json:"encryptionKey,omitempty"`
	PopulatedTables map[string]struct{} `json:"populatedTables,omitempty"`
	DB              *sql.DB             `json:"-"`
}

func (c *Config) checkSQLite() error {
	if c != nil {
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
		return c.SQLite.startDB()
	}
	return errors.New("invalid config cannot check SQLite")
}

// Init will setup the default tables in the sqlite embedded database
func (sqlite *SQLite) startDB() error {
	if sqlite != nil && sqlite.FileName != nil {
		baseDir := filepath.Dir(*sqlite.FileName)
		if err := os.MkdirAll(baseDir, 0o750); err != nil {
			return err
		}

		// encrypted db
		if sqlite.EncryptionKey != nil {
			if key, err := sqlite.GetEncryptionKey(); err == nil && key != nil {
				dbname := fmt.Sprintf("%s?_pragma_key=x'%s'&_pragma_cipher_page_size=4096", *sqlite.FileName, *key)
				sqlite.DB, err = sql.Open("sqlite3", dbname)
				if err != nil {
					return err
				}
				log.Debugf("Using encrypted aud database in %s", *GrowSTLGo.SQLite.FileName)
				return nil
			}
		}
		db, sqliteErr := sql.Open("sqlite3", *sqlite.FileName)
		if sqliteErr != nil {
			return sqliteErr
		}
		sqlite.DB = db
		log.Debugf("Using unencrypted aud database in %s", *sqlite.FileName)
		return nil
	}
	return fmt.Errorf("no sutiable configuration for SQLite found in the config file %s", *ConfigFile)
}

func (sqlite *SQLite) generateEncryptionKeys() error {
	var keyBytes bytes.Buffer
	// generate 64 bit key
	for i := range 4 {
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

// GetSQLiteConnection will give access to the sqlite database
func GetSQLiteConnection() (*sql.DB, error) {
	if GrowSTLGo != nil && GrowSTLGo.SQLite != nil && GrowSTLGo.SQLite.DB != nil {
		return GrowSTLGo.SQLite.DB, nil
	}
	return nil, errors.New("sqlite object is nil")
}

// ShutdownSQLite does a clean shutdown of the sqlite database
func ShutdownSQLite() {
	if GrowSTLGo != nil && GrowSTLGo.SQLite != nil && GrowSTLGo.SQLite.DB != nil {
		log.Info("closing the SQLite database")
		if err := GrowSTLGo.SQLite.DB.Close(); err != nil {
			log.Errorf("Problems closing the SQLite database.  Error: %s", err)
		}
	}
}

// CreateTable is a helper function that will create the table & indexes in the embedded database
func (table *Table) CreateTable(tableName *string) error {
	db, dbErr := GetSQLiteConnection()
	if dbErr != nil {
		return dbErr
	}

	if db != nil && table != nil && tableName != nil {
		log.Tracef("Audit table %s was created if it didn't already exist", *tableName)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		stmt, err := db.PrepareContext(ctx, table.CreateSQL)
		if err != nil {
			return err
		}
		defer stmt.Close()
		if _, err = stmt.ExecContext(ctx); err != nil {
			return err
		}
		for _, index := range table.Indices {
			if err := createHelper(db, index); err != nil {
				return err
			}
		}

		if GrowSTLGo.SQLite.PopulatedTables == nil {
			GrowSTLGo.SQLite.PopulatedTables = make(map[string]struct{})
		}

		if _, ok := GrowSTLGo.SQLite.PopulatedTables[*tableName]; !ok && table.Defaults != nil {
			for key, sql := range table.Defaults {
				log.Tracef("Inserting default %s for table %s if it doesn't already exist", key, *tableName)
				if err := createHelper(db, sql); err != nil {
					return err
				}
			}
			GrowSTLGo.SQLite.PopulatedTables[*tableName] = struct{}{}
			if err := GrowSTLGo.persist(); err != nil {
				return err
			}
		}
		return nil
	}
	return errors.New("tables is nil cannot create tables")
}

func createHelper(db *sql.DB, index string) error {
	if db != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		stmt, err := db.PrepareContext(ctx, index)
		if err != nil {
			return fmt.Errorf("unable to prepare '%s'.  Error: %s", index, err)
		}
		defer stmt.Close()
		if _, err = stmt.ExecContext(ctx); err != nil {
			return fmt.Errorf("unable to create '%s'.  Error: %s", index, err)
		}
		return nil
	}
	return errors.New("db is nil cannot create item")
}
