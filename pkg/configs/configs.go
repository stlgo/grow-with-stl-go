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
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"stl-go/grow-with-stl-go/pkg/cryptography"
	"stl-go/grow-with-stl-go/pkg/log"
	"strings"
	"time"
)

var (
	// GrowSTLGo is the main config for the application
	GrowSTLGo Config
	// ConfigFile is the physical file that contains the config for the application
	ConfigFile     *string
	rewriteConfig  bool
	configReadTime int64
	watchSetup     bool
	keysToEncrypt  = regexp.MustCompile("([a-zA-Z]+sswd|[a-zA-Z]+pswd|[a-zA-Z]+assword)")
)

// constant keys used by the websocket communications
const (
	Auth         string = "auth"
	AuthComplete string = "authcomplete"
	Keepalive    string = "keepalive"
)

// Config contains the basis of the web service
type Config struct {
	Proxy  *Proxy  `json:"proxy,omitempty"`
	Secret *string `json:"secret,omitempty"`
}

// Proxy is in case we need to use a proxy for http connections this is where it goes
type Proxy struct {
	URL          *string `json:"url,omitempty"`
	ExtraCACerts *string `json:"extraCACerts,omitempty"`
}

// WebService is the definition for the webservice
type WebService struct {
	Host         *string `json:"host,omitempty"`
	Port         *int    `json:"port,omitempty"`
	PublicKey    *string `json:"publicKey,omitempty"`
	PrivateKey   *string `json:"privateKey,omitempty"`
	StaticWebDir *string `json:"staticWebDir,omitempty"`
}

// WsMessage is a request / return structure used for websockets
type WsMessage struct {
	// base components of a message
	Type         *string `json:"type,omitempty"`
	Component    *string `json:"component,omitempty"`
	SubComponent *string `json:"subComponent,omitempty"`
	SessionID    *string `json:"sessionID,omitempty"`
	Timestamp    *int64  `json:"timestamp,omitempty"`

	// additional conditional components that may or may not be involved in the request / response
	Data  interface{} `json:"data,omitempty"`
	Error *string     `json:"error,omitempty"`

	// used for auth
	Authentication *Authentication `json:"authentication,omitempty"`
}

// SetGrowSTLGoConfig sets the config for the application
func SetGrowSTLGoConfig() error {
	if ConfigFile != nil {
		// read the config file if it exists
		jsonBytes, err := os.ReadFile(*ConfigFile)
		if err != nil {
			log.Error(err)
			log.Info("No configuration found building a default configuration")
		}

		// unmarshall the config file to a hash map if it exists
		var jo map[string]interface{}
		err = json.Unmarshal(jsonBytes, &jo)
		if err != nil {
			log.Error(err)
		}

		// get the secret out of the file, if there isn't one generate it
		err = getSecret(jo)
		if err != nil {
			return err
		}

		// scan the json for keys we want to encrypt that are currently clear text
		err = scanJSON(jo)
		if err != nil {
			log.Error(err)
		}

		// map the json to the struct
		jsonBytes, err = json.Marshal(jo)
		if err != nil {
			return err
		}

		err = json.Unmarshal(jsonBytes, &GrowSTLGo)
		if err != nil {
			return err
		}

		if !watchSetup {
			watchSetup = true
			go configRecheckTimer()
		}

		return checkConfigs()
	}
	return errors.New("config file is nil")
}

// get the secret from the config
func getSecret(jo map[string]interface{}) error {
	if secretInterface, ok := jo["secret"]; ok {
		if s, ok := secretInterface.(string); ok {
			GrowSTLGo.Secret = &s
			return nil
		}
	}
	newSecret, err := cryptography.GenerateSecret()
	if err != nil {
		return err
	}
	GrowSTLGo.Secret = newSecret
	return nil
}

// scan a generic JSON object for specific keys
func scanJSON(jo any) error {
	for key, value := range jo.(map[string]interface{}) {
		switch value.(type) {
		case map[string]interface{}:
			if err := scanJSON(value); err != nil {
				return err
			}
		case []interface{}:
			if err := scanJSONArray(value); err != nil {
				return err
			}
		default:
			if keysToEncrypt.MatchString(strings.ToLower(key)) {
				if err := scanJSONHelper(jo, value, key); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// scan a generic JSON array for specific keys (utilized by scanJSON)
func scanJSONArray(anArray any) error {
	for _, value := range anArray.([]interface{}) {
		switch value.(type) {
		case map[string]interface{}:
			if err := scanJSON(value); err != nil {
				return err
			}
		case []interface{}:
			if err := scanJSONArray(value); err != nil {
				return err
			}
		}
	}
	return nil
}

func scanJSONHelper(jo any, value interface{}, key string) error {
	valueStr, ok := value.(string)
	if !ok {
		return fmt.Errorf("unable to encrypt data for key %s", key)
	}
	if !strings.HasPrefix(valueStr, cryptography.ObfuscatedPrefix) {
		cipherText, err := cryptography.Encrypt(&valueStr, GrowSTLGo.Secret)
		if err != nil {
			return fmt.Errorf("unable to encrypt data for key %s", key)
		}
		jo.(map[string]interface{})[key] = *cipherText
		rewriteConfig = true
	}
	return nil
}

func checkConfigs() error {
	configReadTime = time.Now().UnixMilli()
	return nil
}

func configRecheckTimer() {
	// move the timer to the top of the minute for execution
	time.Sleep(time.Duration(60-time.Now().Local().Second()) * time.Second)
	// execute once per minute
	for range time.NewTicker(1 * time.Minute).C {
		info, err := os.Stat(*ConfigFile)
		if err == nil {
			if !rewriteConfig && info.ModTime().UnixNano()/1000000 > configReadTime {
				log.Debugf("Update detected in %s, rechecking file", *ConfigFile)
				if err := SetGrowSTLGoConfig(); err != nil {
					log.Error(err)
				}
			}
		}
		rewriteConfig = false
	}
}
