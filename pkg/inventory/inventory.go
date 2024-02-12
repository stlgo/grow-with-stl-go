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

package inventory

import (
	"fmt"

	"stl-go/grow-with-stl-go/pkg/configs"
	"stl-go/grow-with-stl-go/pkg/log"
	"stl-go/grow-with-stl-go/pkg/webservice"
)

const (
	inventory        = "inventory"
	getInventory     = "getInventory"
	addInventory     = "addInventory"
	reserveInventory = "reserveInventory"
	removeInventory  = "removeInventory"
)

var (
	inventoryTables = map[string]*configs.Table{
		"seeds": {
			CreateSQL: `CREATE TABLE IF NOT EXISTS seeds (
				genus varchar(512) NOT NULL,
				species varchar(512) NOT NULL,
				cultivar varchar(512) NOT NULL,
				commonName varcar(1024) NOT NULL,
				description varcar(1024) NOT NULL,
				hybrid tinyint(1) NOT NULL default 0,
				price int NOT NULL,
				perpacketcount int NOT NULL,
				packets int NOT NULL,
				image varcar(1024))`,
			InsertSQL: "INSERT INTO WebSocket values(?,?,?,?,?,?,?,?,?,?)",
			Indices: []string{
				"CREATE INDEX IF NOT EXISTS seedcn on seeds(commonName)",
			},
		}, "tools": {
			CreateSQL: `CREATE TABLE IF NOT EXISTS tools (
				name varcar(1024) NOT NULL,
				description varcar(1024) NOT NULL,
				price int NOT NULL,
				inStock int NOT NULL,
				image varcar(1024))`,
			InsertSQL: "INSERT INTO WebSocket values(?,?,?,?,?)",
			Indices: []string{
				"CREATE INDEX IF NOT EXISTS toolsName on tools(name)",
			},
		},
	}
)

// Init is called by the launch in the pkg/commands/root.go
func Init() error {
	webservice.AppendToFunctionMap(inventory, handleMessage)
	for tableName, table := range inventoryTables {
		if table != nil {
			if err := table.CreateTable(&tableName); err != nil {
				return err
			}
		}
	}
	return nil
}

func handleMessage(sessionID *string, request *configs.WsMessage) *configs.WsMessage {
	response := configs.WsMessage{
		Type:         request.Type,
		Component:    request.Component,
		SubComponent: request.SubComponent,
	}

	if request.Component != nil {
		switch *request.Component {
		case addInventory:
			log.Trace("add inventory")
		case getInventory:
			log.Trace(fmt.Sprintf("getInventory received for session %s", *sessionID))
		default:
			err := fmt.Sprintf("component %s not implemented", *request.Component)
			log.Error(err)
			response.Error = &err
		}
		return &response
	}
	err := fmt.Errorf("bad request").Error()
	response.Error = &err
	return &response
}
