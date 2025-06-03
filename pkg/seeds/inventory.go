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

package seeds

import (
	"errors"
	"fmt"

	"github.com/google/uuid"

	"stl-go/grow-with-stl-go/pkg/configs"
	"stl-go/grow-with-stl-go/pkg/log"
)

var (
	inventoryTables = map[string]*configs.Table{
		"seeds": {
			CreateSQL: `CREATE TABLE IF NOT EXISTS seeds (
				id varchar(64) NOT NULL,
				category TEXT CHECK(category IN ('Herb', 'Onion', 'Pepper', 'Tomato')) NOT NULL,
				genus varchar(512) NOT NULL,
				species varchar(512) NOT NULL,
				cultivar varchar(512),
				commonName varcar(1024) NOT NULL,
				description varcar(2048) NOT NULL,
				hybrid tinyint(1) NOT NULL default 0,
				price real NOT NULL,
				perpacketcount int NOT NULL,
				packets int NOT NULL,
				image varcar(1024))`,
			InsertSQL: "INSERT INTO seeds values(?,?,?,?,?,?,?,?,?,?,?,?)",
			UpdateSQL: "UPDATE seeds set packets = ? where id = ?",
			Indices: []string{
				"CREATE INDEX IF NOT EXISTS seedid on seeds(id)",
				"CREATE INDEX IF NOT EXISTS seedcategory on seeds(category)",
				"CREATE INDEX IF NOT EXISTS seedcn on seeds(commonName)",
			},
			Defaults: map[string]string{
				"dill": fmt.Sprintf(`insert or ignore into seeds values ('%s', 'Herb', 'Anethum', 'graveolens', 'Ella', 'Dill Weed',
				'Ella is a dwarf dill bred for container and hydroponic growing',
				0, 2.67, 20, 100, '/images/herbs/ella_dill.jpg')`, uuid.New().String()),
				"basil": fmt.Sprintf(`insert or ignore into seeds values ('%s', 'Herb', 'Ocimum', 'basilicum', 'Genovese', 'Basil',
				'Genovese basil was first bred in the Northwest coastal port of Genoa, gateway to the Italian Riviera.',
				0, 4.28, 100, 100, '/images/herbs/genovese_basil.jpg')`, uuid.New().String()),
				"oreagno": fmt.Sprintf(`insert or ignore into seeds values ('%s', 'Herb', 'Origanum', 'vulgare', 'Greek', 'Oregano',
				'Strong oregano aroma and flavor; great for pizza and Italian cooking. Characteristic dark green leaves with white flowers.',
				0, 3.95, 50, 100, '/images/herbs/greek_oregano.jpg')`, uuid.New().String()),
				"chive": fmt.Sprintf(`insert or ignore into seeds values ('%s', 'Herb', 'Allium', 'schoenoprasum', 'Polyvert', 'Chive',
				'Suitable for growing in field or containers. Dark green leaves with very good uniformity. USDA Certified Organic.',
				0, 2.18, 100, 100, '/images/herbs/polyvert_chive.jpg')`, uuid.New().String()),

				"ailsa craig": fmt.Sprintf(`insert or ignore into seeds values ('%s', 'Onion', 'Allium', 'cepa', 'Ailsa Craig', 'Yellow',
				'Long day. Very well-known globe-shaped heirloom onion that reaches a really huge size—5 lbs is rather common',
				0, 2.67, 20, 100, '/images/onions/ailsa_craig.jpg')`, uuid.New().String()),
				"patterson": fmt.Sprintf(`insert or ignore into seeds values ('%s', 'Onion', 'Allium', 'cepa', 'Patterson', 'Yellow',
				'Patterson’ is a keeper—the longest-storing onion you can find. Straw-colored, globe-shaped bulbs with sweet, mildly pungent yellow flesh',
				0, 4.28, 100, 100, '/images/onions/patterson.jpg')`, uuid.New().String()),
				"red wing": fmt.Sprintf(`insert or ignore into seeds values ('%s', 'Onion', 'Allium', 'cepa', 'Red Wing', 'Red',
				'Uniform, large onions with deep red color. Thick skin, very hard bulbs for long storage. Consistent internal color.',
				1, 3.95, 50, 100, '/images/onions/red_wing.jpg')`, uuid.New().String()),
				"walla walla": fmt.Sprintf(`insert or ignore into seeds values ('%s', 'Onion', 'Allium', 'cepa', 'Walla Walla', 'Sweet',
				'Juicy, sweet, regional favorite. In the Northwest,  very large, flattened, ultra-mild onions',
				0, 2.18, 100, 100, '/images/onions/walla_walla.jpg')`, uuid.New().String()),

				"bell pepper": fmt.Sprintf(`insert or ignore into seeds values ('%s', 'Pepper', 'Capsicum', 'annuum', 'Ozark Giant', 'Bell',
				'Green bell peppers are bell peppers that have been harvested early. Red bell peppers have been allowed to ripen longer.',
				0, 4.58, 50, 100, '/images/peppers/green_bell.jpg')`, uuid.New().String()),
				"poblano": fmt.Sprintf(`insert or ignore into seeds values ('%s', 'Pepper', 'Capsicum', 'annuum', null, 'Poblano',
				'The poblano is a mild chili pepper originating in the state of Puebla, Mexico. Dried, it is called ancho or chile ancho',
				0, 2.96, 25, 100, '/images/peppers/poblano.jpg')`, uuid.New().String()),
				"jalapeno": fmt.Sprintf(`insert or ignore into seeds values ('%s', 'Pepper', 'Capsicum', 'annuum', 'Zapotec', 'Jalapeno',
				'This jalapeno variety from Oaxaca, Mexico which is a more flavorful, gourmet jalapeño.',
				0, 3.78, 12, 100, '/images/peppers/jalapeno.jpg')`, uuid.New().String()),
				"serrano": fmt.Sprintf(`insert or ignore into seeds values ('%s', 'Pepper', 'Capsicum', 'annuum', 'Tampiqueno', 'Serrano',
				'This serrano variety comes from the mountains of the Hidalgo and Puebla states of Mexico.  This pepper is 2-3 times hotter than jalapenos',
				0, 6.45, 15, 100, '/images/peppers/serrano.jpg')`, uuid.New().String()),

				"galahad": fmt.Sprintf(`insert or ignore into seeds values ('%s', 'Tomato', 'Solanum', 'lycopersicum', 'Galahad', 'Tomato',
				'Delicious early determinate beefsteak.',
				1, 4.58, 50, 100, '/images/tomatoes/determinate/galahad.jpg')`, uuid.New().String()),
				"plum regal": fmt.Sprintf(`insert or ignore into seeds values ('%s', 'Tomato', 'Solanum', 'lycopersicum', 'Plum Regal', 'Tomato',
				'Medium-size plants with good leaf cover produce high yields of blocky, 4 oz. plum tomatoes. Fruits have a deep red color with good flavor. Determinate.',
				1, 2.96, 25, 100, '/images/tomatoes/determinate/plum_regal.jpg')`, uuid.New().String()),
				"carbon": fmt.Sprintf(`insert or ignore into seeds values ('%s', 'Tomato', 'Solanum', 'lycopersicum', 'Carbon', 'Tomato',
				'Indeterminate heirloom. Resists cracking better than other large, black heirlooms. Blocky-round, 10-14 oz. fruit with dark olive shoulders.',
				0, 3.78, 12, 100, '/images/tomatoes/indeterminate/carbon.jpg')`, uuid.New().String()),
				"san marzano": fmt.Sprintf(`insert or ignore into seeds values ('%s', 'Tomato', 'Solanum', 'lycopersicum', 'San Marzano', 'Tomato',
				'San Marzano is considered one of the best paste tomatoes of all time, with Old World look and taste.  Indeterminate',
				0, 6.45, 15, 100, '/images/tomatoes/indeterminate/san_marzano.jpg')`, uuid.New().String()),
			},
		},
	}

	inventoryCache = make(map[string]*InventoryCategory)
)

// InventoryCategory is a way to hold inventory categories together
type InventoryCategory struct {
	Category *string                   `json:"category,omitempty"`
	Items    map[string]*InventoryItem `json:"items,omitempty"`
}

// InventoryItem data we stored in the database
type InventoryItem struct {
	ID             *string  `json:"id,omitempty"`
	Category       *string  `json:"category,omitempty"`
	Genus          *string  `json:"genus,omitempty"`
	Species        *string  `json:"species,omitempty"`
	Cultivar       *string  `json:"cultivar,omitempty"`
	CommonName     *string  `json:"commonName,omitempty"`
	Description    *string  `json:"description,omitempty"`
	Hybrid         *bool    `json:"hybrid,omitempty"`
	Price          *float32 `json:"price,omitempty"`
	PerPacketCount *int32   `json:"perPacketCount,omitempty"`
	Packets        *int     `json:"packets,omitempty"`
	Image          *string  `json:"image,omitempty"`
}

func setupTables() error {
	for tableName, table := range inventoryTables {
		if table != nil {
			if err := table.CreateTable(&tableName); err != nil {
				return err
			}
		}
	}
	return nil
}

// getInventory will query the database for the inventory data
func getInventory() (map[string]*InventoryCategory, error) {
	if configs.SqliteDB != nil {
		if len(inventoryCache) == 0 {
			inventory := make(map[string]*InventoryCategory)
			rows, err := configs.SqliteDB.Query("select * from seeds")
			if err != nil {
				return nil, err
			}
			defer rows.Close()
			for rows.Next() {
				ii := InventoryItem{}
				if err := rows.Scan(&ii.ID, &ii.Category, &ii.Genus, &ii.Species, &ii.Cultivar, &ii.CommonName, &ii.Description, &ii.Hybrid,
					&ii.Price, &ii.PerPacketCount, &ii.Packets, &ii.Image); err != nil {
					return nil, err
				}
				if ii.ID != nil && ii.Category != nil {
					if category, ok := inventory[*ii.Category]; ok && category.Items != nil {
						category.Items[*ii.ID] = &ii
					} else {
						inventory[*ii.Category] = &InventoryCategory{
							Category: ii.Category,
							Items: map[string]*InventoryItem{
								*ii.ID: &ii,
							},
						}
					}
				}
			}
			inventoryCache = inventory
		}
		return inventoryCache, nil
	}
	return nil, errors.New("the sqlite database is nil, cannot get last logins")
}

func getDetail(categoryName, id *string) (*InventoryItem, error) {
	if categoryName != nil && id != nil {
		if category, ok := inventoryCache[*categoryName]; ok && category != nil {
			if item, ok := category.Items[*id]; ok {
				return item, nil
			}
		}
	}
	return nil, errors.New("item not found")
}

func updateItemInDB(item *InventoryItem) error {
	if table, ok := inventoryTables["seeds"]; ok && table != nil && item != nil {
		stmt, err := configs.SqliteDB.Prepare(table.UpdateSQL)
		if err != nil {
			return err
		}

		table.WriteMutex.Lock()
		defer table.WriteMutex.Unlock()
		result, err := stmt.Exec(
			item.Packets,
			item.ID)

		if err != nil {
			return err
		}

		rows, err := result.RowsAffected()
		if err != nil {
			return err
		}

		log.Tracef("%d rows updated in seeds", rows)
		return nil
	}
	return errors.New("requirements not met to update item in db")
}
