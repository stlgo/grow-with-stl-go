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
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"

	_ "github.com/mutecomm/go-sqlcipher/v4" // this is required for the sqlite driver

	"stl-go/grow-with-stl-go/pkg/cryptography"
	"stl-go/grow-with-stl-go/pkg/log"
)

var (
	// GrowSTLGo is the main config for the application
	GrowSTLGo *Config
	// ConfigFile is the physical file that contains the config for the application
	ConfigFile     *string
	etcDir         *string
	rewriteConfig  bool
	configReadTime int64
	watchSetup     bool
	writeMutex     sync.Mutex
	keysToEncrypt  = regexp.MustCompile("([a-zA-Z]+sswd|[a-zA-Z]+pswd|[a-zA-Z]+assword)")

	// SqliteDB is an embedded database
	SqliteDB *sql.DB

	// Version will be overridden by ldflags supplied in Makefile
	Version = "(dev-version)"
)

// constant keys used by the websocket communications
const (
	// UI Components
	Auth            string = "auth"
	AuthComplete    string = "authcomplete"
	GetPagelet      string = "getPagelet"
	Initialize      string = "initialize"
	Keepalive       string = "keepalive"
	UI              string = "ui"
	WebsocketClient string = "websocketclient"

	// auth subcomponents
	Approved     string = "approved"
	Authenticate string = "authenticate"
	Denied       string = "denied"
	Refresh      string = "refresh"
	Validate     string = "validate"

	// http error message "json" to return on errors as a const
	NotFoundError       = `{"error": "Not Found", "status": 404}`
	NotImplementedError = `{"error": "Not Implemented", "status": 501}`
	BadRequestError     = `{"error": "Bad Request", "status": 400}`
	InternalServerError = `{"error": "Internal Server Error", "status": 500}`
	UnauthorizedError   = `{"error": "Unauthorized", "status": 401}`
)

// Config contains the basis of the web service
type Config struct {
	APIUsers   map[string]*APIUser `json:"apiUsers,omitempty"`
	Country    *Country            `json:"country,omitempty"`
	DataDir    *string             `json:"data_dir,omitempty"`
	Proxy      *Proxy              `json:"proxy,omitempty"`
	Secret     *string             `json:"secret,omitempty"`
	SQLite     *SQLite             `json:"sqlite,omitempty"`
	WebService *WebService         `json:"webService,omitempty"`
}

func (c *Config) checkConfig() error {
	checkAPIUsers()

	for _, function := range []func() error{c.checkDataDir, c.checkCountry, checkWebService, checkSQLite, c.testRewriteConfig} {
		if err := function(); err != nil {
			log.Errorf("error calling function %s", runtime.FuncForPC(reflect.ValueOf(function).Pointer()).Name())
		}
	}

	if !watchSetup {
		watchSetup = true
		go configRecheckTimer()
	}

	configReadTime = time.Now().UnixMilli()
	return nil
}

func (c *Config) checkDataDir() error {
	if c != nil {
		if c.DataDir == nil {
			dir, err := os.MkdirTemp("", "gwstlg")
			if err != nil {
				return err
			}
			c.DataDir = &dir
		}
		return nil
	}
	return errors.New("invalid config cannot check data dir")
}

func (c *Config) checkCountry() error {
	if c != nil {
		if c.Country == nil {
			u, urlErr := url.Parse("https://download.geonames.org/export/zip/")
			if urlErr != nil {
				return urlErr
			}
			urlStr := u.String()
			country := "US"
			c.Country = &Country{
				Country: &country,
				URL:     &urlStr,
			}
			rewriteConfig = true
			return nil
		} else if c.Country.Country != nil && c.Country.URL != nil {
			_, urlErr := url.Parse(*c.Country.URL)
			if urlErr != nil {
				return urlErr
			}
			return nil
		}
	}
	return errors.New("invalid config cannot check country")
}

func (c *Config) testRewriteConfig() error {
	if c != nil {
		if rewriteConfig {
			if err := c.persist(); err != nil {
				return err
			}
		}
		return nil
	}
	return errors.New("invalid config cannot check rewrite condition")
}

func (c *Config) persist() error {
	if c != nil {
		// lets's make sure we can kick the JSON out of the config
		bytes, err := json.Marshal(GrowSTLGo)
		if err != nil {
			return err
		}

		// we do this so we can rescan for cleartext passwords
		var jo map[string]interface{}
		err = json.Unmarshal(bytes, &jo)
		if err != nil {
			return err
		}

		// scan for cleartext passwords
		err = scanJSON(jo)
		if err != nil {
			return err
		}

		return writeJSON(jo)
	}
	return errors.New("invalid config cannot persist data")
}

// Country is set for our download of zip codes & latitude / longitudes
type Country struct {
	Country *string `json:"country,omitempty"`
	URL     *string `json:"url,omitempty"`
}

// Proxy is in case we need to use a proxy for http connections this is where it goes
type Proxy struct {
	URL          *string `json:"url,omitempty"`
	ExtraCACerts *string `json:"extraCACerts,omitempty"`
}

// WsMessage is a request / return structure used for websockets
type WsMessage struct {
	// base components of a message
	Route        *string `json:"route,omitempty"`
	Type         *string `json:"type,omitempty"`
	Component    *string `json:"component,omitempty"`
	SubComponent *string `json:"subComponent,omitempty"`
	SessionID    *string `json:"sessionID,omitempty"`
	Timestamp    *int64  `json:"timestamp,omitempty"`

	// additional conditional components that may or may not be involved in the request / response
	Data  interface{} `json:"data,omitempty"`
	Error *string     `json:"error,omitempty"`

	// used for authentication
	Authentication *Authentication `json:"authentication,omitempty"`
	Token          *string         `json:"token,omitempty"`
	ValidTill      *int64          `json:"validTill,omitempty"`
	RefreshToken   *string         `json:"refreshToken,omitempty"`
	IsAdmin        *bool           `json:"isAdmin,omitempty"`
	Vhost          *string
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

		// unmarshal the config file to a hash map if it exists
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
	}

	if GrowSTLGo != nil {
		return GrowSTLGo.checkConfig()
	}
	return errors.New("invalid config cannot continue")
}

// get the secret from the config
func getSecret(jo map[string]interface{}) error {
	if secretInterface, ok := jo["secret"]; ok {
		if s, ok := secretInterface.(string); ok {
			if GrowSTLGo == nil {
				GrowSTLGo = &Config{
					Secret: &s,
				}
			} else {
				GrowSTLGo.Secret = &s
			}

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
	if m, t := jo.(map[string]interface{}); t {
		for key, value := range m {
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
	}
	return nil
}

// scan a generic JSON array for specific keys (utilized by scanJSON)
func scanJSONArray(anArray any) error {
	if a, t := anArray.([]interface{}); t {
		for _, value := range a {
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
		if m, t := jo.(map[string]interface{}); t {
			m[key] = *cipherText
			rewriteConfig = true
		}
	}
	return nil
}

func writeJSON(jo map[string]interface{}) error {
	if ConfigFile != nil {
		// if the data has mutated then we want to write it out to disk
		writeMutex.Lock()
		defer writeMutex.Unlock()
		log.Debugf("Rewriting %s to ensure data is enciphered on disk", *ConfigFile)
		jsonBytes, err := json.MarshalIndent(jo, "", "\t")
		if err != nil {
			return err
		}

		// write to a new file
		dir := path.Dir(*ConfigFile)
		newFile := filepath.Join(dir, ".new_config.json")
		err = os.WriteFile(newFile, jsonBytes, 0o600)
		if err != nil {
			return err
		}

		// move the original off
		_, err = os.Stat(*ConfigFile)
		if err == nil {
			oldFile := filepath.Join(dir, ".old_config.json")
			err = os.Rename(*ConfigFile, oldFile)
			if err != nil {
				return err
			}
		}

		// move the new file into the proper place
		err = os.Rename(newFile, *ConfigFile)
		if err != nil {
			return err
		}

		return nil
	}
	return errors.New("config file is nil")
}

func configRecheckTimer() {
	// move the timer to the top of the minute for execution
	time.Sleep(time.Duration(60-time.Now().Local().Second()) * time.Second)
	// execute once per minute
	for range time.NewTicker(1 * time.Minute).C {
		if ConfigFile != nil {
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
}

func getEtcDir() (*string, error) {
	if etcDir == nil && ConfigFile != nil {
		// get the runtime directory
		dir, err := filepath.Abs(filepath.Dir(*ConfigFile))
		if err != nil {
			return nil, err
		}

		// create an etc dir under the root of the runtime directory
		dir, err = filepath.Abs(filepath.Join(path.Dir(dir), "etc"))
		if err != nil {
			return nil, err
		}

		// make sure the dir is there
		err = os.MkdirAll(dir, 0o750)
		if err != nil {
			return nil, err
		}

		etcDir = &dir
	}
	return etcDir, nil
}
