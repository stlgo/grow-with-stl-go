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
				description varcar(2048) NOT NULL,
				hybrid tinyint(1) NOT NULL default 0,
				price real NOT NULL,
				perpacketcount int NOT NULL,
				packets int NOT NULL,
				image varcar(1024))`,
			InsertSQL: "INSERT INTO WebSocket values(?,?,?,?,?,?,?,?,?,?)",
			Indices: []string{
				"CREATE INDEX IF NOT EXISTS seedcn on seeds(commonName)",
			},
			Defaults: map[string]string{
				"dill": `insert or ignore into seeds values ('Anethum', 'graveolens', 'Ella', 'Dill Weed',
				'Ella is a dwarf dill bred for container and hydroponic growing',
				0, 2.67, 20, 100, '/images/herbs/ella_dill.jpg')`,
				"basil": `insert or ignore into seeds values ('Ocimum', 'basilicum', 'Genovese', 'Basil',
				'Genovese basil was first bred in the Northwest coastal port of Genoa, gateway to the Italian Riviera.',
				0, 4.28, 100, 100, '/images/herbs/genovese_basil.jpg')`,
				"oreagno": `insert or ignore into seeds values ('Origanum', 'vulgare', 'Greek', 'Oregano',
				'Strong oregano aroma and flavor; great for pizza and Italian cooking. Characteristic dark green leaves with white flowers.',
				0, 3.95, 50, 100, '/images/herbs/greek_oregano.jpg')`,
				"chive": `insert or ignore into seeds values ('Allium', 'schoenoprasum', 'Polyvert', 'Chive',
				'Suitable for growing in field or containers. Dark green leaves with very good uniformity. USDA Certified Organic.',
				0, 2.18, 100, 100, '/images/herbs/polyvert_chive.jpg')`,

				"ailsa craig": `insert or ignore into seeds values ('Allium', 'cepa', 'Ailsa Craig', 'Yellow Onion',
				'Long day. Very well-known globe-shaped heirloom onion that reaches a really huge size—5 lbs is rather common',
				0, 2.67, 20, 100, '/images/onions/aisa_craig.jpg')`,
				"patterson": `insert or ignore into seeds values ('Allium', 'cepa', 'Patterson', 'Yellow Onion',
				'Patterson’ is a keeper—the longest-storing onion you can find. Straw-colored, globe-shaped bulbs with sweet, mildly pungent yellow flesh',
				0, 4.28, 100, 100, '/images/onions/patterson.jpg')`,
				"red wing": `insert or ignore into seeds values ('Allium', 'cepa', 'Red Wing', 'Red Onion',
				'Uniform, large onions with deep red color. Thick skin, very hard bulbs for long storage. Consistent internal color.',
				1, 3.95, 50, 100, '/images/onions/red_wing.jpg')`,
				"walla walla": `insert or ignore into seeds values ('Allium', 'cepa', 'Walla Walla', 'Sweet Onion',
				'Juicy, sweet, regional favorite. In the Northwest,  very large, flattened, ultra-mild onions',
				0, 2.18, 100, 100, '/images/onions/walla_walla.jpg')`,

				"bell pepper": `insert or ignore into seeds values ('Capsicum', 'annuum', 'Ozark Giant', 'Bell Pepper',
				'Green bell peppers are bell peppers that have been harvested early. Red bell peppers have been allowed to ripen longer.',
				0, 4.58, 50, 100, '/images/peppers/gree_bell.jpg')`,
				"poblano": `insert or ignore into seeds values ('Capsicum', 'annuum', 'N/A', 'Poblano',
				'The poblano is a mild chili pepper originating in the state of Puebla, Mexico. Dried, it is called ancho or chile ancho',
				0, 2.96, 25, 100, '/images/peppers/poblano.jpg')`,
				"jalapeno": `insert or ignore into seeds values ('Capsicum', 'annuum', 'Zapotec', 'Jalapeno',
				'This jalapeno variety from Oaxaca, Mexico which is a more flavorful, gourmet jalapeño.',
				0, 3.78, 12, 100, '/images/peppers/jalapeno.jpg')`,
				"serrano": `insert or ignore into seeds values ('Capsicum', 'annuum', 'Tampiqueno', 'Serrano',
				'This serrano variety comes from the mountains of the Hidalgo and Puebla states of Mexico.  This pepper is 2-3 times hotter than jalapenos',
				0, 6.45, 15, 100, '/images/peppers/serrano.jpg')`,

				"galahad": `insert or ignore into seeds values ('Solanum', 'lycopersicum', 'Galahad', 'tomato',
				'Delicious early determinate beefsteak.',
				1, 4.58, 50, 100, '/images/tomatoes/determinate/galahad.jpg')`,
				"plum regal": `insert or ignore into seeds values ('Solanum', 'lycopersicum', 'Plum Regal', 'tomato',
				'Medium-size plants with good leaf cover produce high yields of blocky, 4 oz. plum tomatoes. Fruits have a deep red color with good flavor. Determinate.',
				1, 2.96, 25, 100, '/images/tomatoes/determinate/plum_regal.jpg')`,
				"cherokee purple": `insert or ignore into seeds values ('Solanum', 'lycopersicum', 'Zapotec', 'tomato',
				'Indeterminate heirloom. Medium-large, flattened globe, 8-12 oz. fruits. Color is dusky pink with dark shoulders.',
				0, 3.78, 12, 100, '/images/tomatoes/indeterminate/cherokee_purple.jpg')`,
				"san marzano": `insert or ignore into seeds values ('Solanum', 'lycopersicum', 'San Marzano', 'tomato',
				'San Marzano is considered one of the best paste tomatoes of all time, with Old World look and taste.  Indeterminate',
				0, 6.45, 15, 100, '/images/tomatoes/indeterminate/san_marzano.jpg')`,
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
