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

package weather

import "stl-go/grow-with-stl-go/pkg/configs"

// Init is different than the standard init because it is called outside of the object load
func Init() error {
	return getLocations()
}

func getLocations() error {
	// we're going to get our zip code based locations on the country here: https://download.geonames.org/export/zip/
	return configs.GrowSTLGo.Country.GetCountryData()
}
