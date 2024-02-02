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

package audit

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	_ "github.com/mutecomm/go-sqlcipher/v4" // this is required for the sqlite driver

	"stl-go/grow-with-stl-go/pkg/configs"
	"stl-go/grow-with-stl-go/pkg/log"
)

// WSTransaction will be used to record WebSocket transactions and record them to the db
type WSTransaction struct {
	Host         *string
	Component    *string
	SubComponent *string
	User         *string
	Type         *string
	Start        *int64
	Elapsed      *int64
	Recordable   *bool
}

// RESTTransaction will be used to record REST transactions and record them to the db
type RESTTransaction struct {
	Host            *string
	URI             *string
	User            *string
	Method          *string
	Payoload        *string
	QueryParamaters *string
	Elapsed         *int64
	Recordable      *bool
}

type table struct {
	CreateSQL  string
	InsertSQL  string
	Indices    []string
	WriteMutex *sync.Mutex
}

var (
	sqliteDB *sql.DB

	auditTables = map[string]*table{
		"WebSocket": {
			CreateSQL: `CREATE TABLE IF NOT EXISTS WebSocket (
				host varchar(1024) NOT NULL,
				user varchar(256) NOT NULL,
				type varchar(128) NOT NULL,
				component varcar(128) NOT NULL,
				subcomponent varchar(128) NOT NULL,
				start bigint NOT NULL,
				stop bigint NOT NULL,
				elapsed bigint NOT NULL)`,
			InsertSQL: "INSERT INTO WebSocket values(?,?,?,?,?,?,?,?)",
			Indices: []string{
				"CREATE INDEX IF NOT EXISTS wshost on WebSocket(host)",
				"CREATE INDEX IF NOT EXISTS wsuser on WebSocket(user)",
				"CREATE INDEX IF NOT EXISTS wstype on WebSocket(type)",
				"CREATE INDEX IF NOT EXISTS wscomponent on WebSocket(component)",
				"CREATE INDEX IF NOT EXISTS wssubcomponent on WebSocket(subcomponent)",
			},
		},
		"REST": {
			CreateSQL: `CREATE TABLE IF NOT EXISTS REST (
				host varchar(1024) NOT NULL,
				uri varchar(1024) NOT NULL,
				user varchar(128) NOT NULL,
				method TEXT CHECK(method IN ('GET', 'PUT', 'POST', 'PATCH', 'DELETE')) NOT NULL,
				payload TEXT,
				queryParameters text,
				httpStatusCode int NOT NULL,
				start bigint NOT NULL,
				stop bigint NOT NULL,
				elapsed bigint NOT NULL)`,
			InsertSQL: "INSERT INTO WebSocket values(?,?,?,?,?,?,?,?,?,?)",
			Indices: []string{
				"CREATE INDEX IF NOT EXISTS RESThost on REST(host)",
				"CREATE INDEX IF NOT EXISTS RESTuri on REST(uri)",
				"CREATE INDEX IF NOT EXISTS RESTuser on REST(user)",
			},
		},
		"user": {
			CreateSQL: `CREATE TABLE IF NOT EXISTS user (
				user varchar(128) NOT NULL,
				protocol TEXT CHECK(protocol IN ('WebSocket', 'REST')) NOT NULL,
				success tinyint(1) NOT NULL default 0,
				observed bigint NOT NULL)`,
			InsertSQL: "INSERT INTO user values(?,?,?,?)",
		},
	}
)

// Init will setup the default tables in the sqlite embedded database
func Init() error {
	if configs.GrowSTLGo.SQLite != nil && configs.GrowSTLGo.SQLite.FileName != nil {
		baseDir := filepath.Dir(*configs.GrowSTLGo.SQLite.FileName)
		if err := os.MkdirAll(baseDir, 0o750); err != nil {
			return (err)
		}

		// encrypted db
		if configs.GrowSTLGo.SQLite.EncryptionKey != nil {
			if key, err := configs.GrowSTLGo.SQLite.GetEncryptionKey(); err == nil && key != nil {
				dbname := fmt.Sprintf("%s?_pragma_key=x'%s'&_pragma_cipher_page_size=4096", *configs.GrowSTLGo.SQLite.FileName, *key)
				sqliteDB, err = sql.Open("sqlite3", dbname)
				if err == nil {
					if err = createTables(); err != nil {
						return (err)
					}
					log.Infof("Using encrypted aud database in %s", *configs.GrowSTLGo.SQLite.FileName)
					return nil
				}
				return (err)
			}
		}
		db, err := sql.Open("sqlite3", *configs.GrowSTLGo.SQLite.FileName)
		if err == nil {
			sqliteDB = db
			if err = createTables(); err != nil {
				return (err)
			}
			log.Infof("Using unencrypted aud database in %s", *configs.GrowSTLGo.SQLite.FileName)
			return nil
		}
		return err
	}
	return fmt.Errorf("no sutiable configuration for SQLite found in the config file %s", *configs.ConfigFile)
}

func createTables() error {
	for _, table := range auditTables {
		stmt, err := sqliteDB.Prepare(table.CreateSQL)
		if err != nil {
			return err
		}
		if _, err = stmt.Exec(); err != nil {
			return err
		}
		for _, index := range table.Indices {
			stmt, err := sqliteDB.Prepare(index)
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
