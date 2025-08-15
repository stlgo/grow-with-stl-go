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
	"errors"
	"fmt"
	"strings"
	"time"

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
	Payload         *string
	QueryParamaters *string
	Start           *int64
	Elapsed         *int64
	Recordable      *bool
}

var (
	auditTables = map[string]*configs.Table{
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
	for tableName, table := range auditTables {
		if table != nil {
			if err := table.CreateTable(&tableName); err != nil {
				return err
			}
		}
	}
	return nil
}

// RecordLogin will insert a record into the users table once a user attempts a login
func RecordLogin(userID *string, protocol string, goodLogin bool) {
	db, dbErr := configs.GetSQLiteConnection()
	if dbErr != nil {
		log.Errorf("unable to record login, sqlite error: %s", dbErr)
	}
	if table, ok := auditTables["user"]; db != nil && ok && table != nil && userID != nil {
		stmt, err := db.Prepare(table.InsertSQL)
		if err != nil {
			log.Error(err)
			return
		}

		observed := time.Now().UnixMilli()

		success := 0
		if goodLogin {
			success = 1
		}

		table.WriteMutex.Lock()
		defer table.WriteMutex.Unlock()
		result, err := stmt.Exec(
			userID,
			protocol,
			success,
			observed)

		if err != nil {
			log.Error(err)
			return
		}

		rows, err := result.RowsAffected()
		if err != nil {
			log.Error(err)
			return
		}

		log.Tracef("%d rows inserted into user", rows)
		return
	}
	log.Trace("did not meet the requirements to record the login")
}

// GetLastLogins will query the database for the last recorded login of the users
func GetLastLogins() error {
	db, dbErr := configs.GetSQLiteConnection()
	if dbErr != nil {
		return fmt.Errorf("unable to get last logins, sqlite error: %s", dbErr)
	}
	if db != nil && configs.GrowSTLGo != nil && configs.GrowSTLGo.APIUsers != nil {
		rows, err := db.Query("select user, max(observed) from user where success = 1 group by user")
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var userID string
			var lastLogin *int64
			if err := rows.Scan(&userID, &lastLogin); err != nil {
				return err
			}
			configs.GrowSTLGo.APIUsersMutex.Lock()
			user, ok := configs.GrowSTLGo.APIUsers[userID]
			configs.GrowSTLGo.APIUsersMutex.Unlock()
			if ok {
				user.LastLogin = lastLogin
			}
		}
		return nil
	}
	return errors.New("the sqlite database is nil, cannot get last logins")
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

// Complete is a WSTransaction receiver function to record it to the DB if applicable
func (transaction *WSTransaction) Complete(errorMessagePresent bool) error {
	db, dbErr := configs.GetSQLiteConnection()
	if dbErr != nil {
		return fmt.Errorf("unable to record the complete transaction, sqlite error: %s", dbErr)
	}
	if transaction.Recordable != nil && *transaction.Recordable && transaction.User != nil && transaction.Start != nil && db != nil &&
		transaction.Component != nil && transaction.Type != nil && !strings.EqualFold(*transaction.Type, "keepalive") {
		if table, ok := auditTables["WebSocket"]; ok && table != nil {
			stmt, err := db.Prepare(table.InsertSQL)
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
	return nil
}

// Complete is a RESTTransaction receiver function to record it to the DB if applicable
func (transaction *RESTTransaction) Complete(httpStatusCode int) error {
	db, dbErr := configs.GetSQLiteConnection()
	if dbErr != nil {
		return fmt.Errorf("unable to record the complete transaction, sqlite error: %s", dbErr)
	}
	if transaction.Recordable != nil && *transaction.Recordable && transaction.Start != nil && db != nil {
		if table, ok := auditTables["REST"]; ok && table != nil {
			stmt, err := db.Prepare(table.InsertSQL)
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
				transaction.Payload,
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
