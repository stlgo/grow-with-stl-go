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
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	_ "github.com/mutecomm/go-sqlcipher/v4" // this is required for the sqlite driver

	"stl-go/grow-with-stl-go/pkg/configs"
	"stl-go/grow-with-stl-go/pkg/log"
	"stl-go/grow-with-stl-go/pkg/utils"
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
	Start           *int64
	Elapsed         *int64
	Recordable      *bool
}

type table struct {
	CreateSQL  string
	InsertSQL  string
	Indices    []string
	WriteMutex sync.Mutex
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
				subcomponent varchar(128),
				success tinyint(1) NOT NULL default 0,
				start bigint NOT NULL,
				stop bigint NOT NULL,
				elapsed bigint NOT NULL)`,
			InsertSQL: "INSERT INTO WebSocket values(?,?,?,?,?,?,?,?,?)",
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
			InsertSQL: "INSERT INTO REST values(?,?,?,?,?,?,?,?,?,?)",
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
			return err
		}

		// encrypted db
		if configs.GrowSTLGo.SQLite.EncryptionKey != nil {
			if key, err := configs.GrowSTLGo.SQLite.GetEncryptionKey(); err == nil && key != nil {
				dbname := fmt.Sprintf("%s?_pragma_key=x'%s'&_pragma_cipher_page_size=4096", *configs.GrowSTLGo.SQLite.FileName, *key)
				sqliteDB, err = sql.Open("sqlite3", dbname)
				if err != nil {
					return err
				}
				if err := createTables(); err != nil {
					return err
				}
				log.Debugf("Using encrypted aud database in %s", *configs.GrowSTLGo.SQLite.FileName)
				return nil
			}
		}
		db, err := sql.Open("sqlite3", *configs.GrowSTLGo.SQLite.FileName)
		if err != nil {
			return err
		}
		sqliteDB = db
		if err := createTables(); err != nil {
			return err
		}
		log.Debugf("Using unencrypted aud database in %s", *configs.GrowSTLGo.SQLite.FileName)
		return nil
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

// NewWSTransaction creates a new WSTransaction object to record as an auditable thing
func NewWSTransaction(host, user *string, request *configs.WsMessage) *WSTransaction {
	return &WSTransaction{
		Host:         host,
		Type:         request.Type,
		Component:    request.Component,
		SubComponent: request.SubComponent,
		User:         user,
		Start:        utils.CurrentTimeInMillis(),
		Recordable:   utils.BoolPointer(true),
	}
}

// Complete is a WSTransaction receiver function to record it to the DB if applicable
func (transaction *WSTransaction) Complete(errorMessagePresent bool) error {
	if transaction.Recordable != nil && *transaction.Recordable && transaction.User != nil && transaction.Start != nil && sqliteDB != nil &&
		transaction.Component != nil && !strings.EqualFold(*transaction.Component, "keepalive") {
		if table, ok := auditTables["WebSocket"]; ok {
			stmt, err := sqliteDB.Prepare(table.InsertSQL)
			if err != nil {
				return err
			}

			start := *transaction.Start
			stop := time.Now().UnixMilli()

			success := 1
			if errorMessagePresent {
				success = 0
			}

			table.WriteMutex.Lock()
			defer table.WriteMutex.Unlock()
			result, err := stmt.Exec(
				transaction.Host,
				transaction.User,
				transaction.Type,
				transaction.Component,
				transaction.SubComponent,
				success,
				start,
				stop,
				(stop - start))

			if err != nil {
				return err
			}

			rows, err := result.RowsAffected()
			if err != nil {
				return err
			}

			log.Tracef("%d rows inserted into websocket", rows)
			return nil
		}
	}
	log.Infof("here %s", *transaction.Component)
	return nil
}

// NewRESTTransaction creates a new RESTTransaction object to record as an auditable thing
func NewRESTTransaction(host, uri, method *string) *RESTTransaction {
	return &RESTTransaction{
		Host:       host,
		URI:        uri,
		Method:     method,
		Start:      utils.CurrentTimeInMillis(),
		Recordable: utils.BoolPointer(true),
	}
}

// Complete is a RESTTransaction receiver function to record it to the DB if applicable
func (transaction *RESTTransaction) Complete(httpStatusCode int) error {
	if transaction.Recordable != nil && *transaction.Recordable && transaction.Start != nil && sqliteDB != nil {
		if table, ok := auditTables["REST"]; ok {
			stmt, err := sqliteDB.Prepare(table.InsertSQL)
			if err != nil {
				return err
			}

			start := *transaction.Start
			stop := time.Now().UnixMilli()

			table.WriteMutex.Lock()
			defer table.WriteMutex.Unlock()
			result, err := stmt.Exec(
				transaction.Host,
				transaction.URI,
				transaction.User,
				transaction.Method,
				transaction.Payoload,
				transaction.QueryParamaters,
				httpStatusCode,
				start,
				stop,
				(stop - start))

			if err != nil {
				return err
			}

			rows, err := result.RowsAffected()
			if err != nil {
				return err
			}

			log.Tracef("%d rows inserted into REST", rows)
			return nil
		}
	}
	return errors.New("transaction is not recordable")
}
