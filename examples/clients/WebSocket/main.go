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

package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"

	"stl-go/grow-with-stl-go/pkg/log"
	"stl-go/grow-with-stl-go/pkg/utils"
)

var (
	host = "localhost:10443"

	ws *websocket.Conn
	// proxy      *string
	// extraCerts *string

	user     *string
	password *string

	// WebSocket constructs
	onError    = make(chan error)
	onClose    = make(chan os.Signal, 1)
	writeMutex sync.Mutex

	// Things used in WebSocket messages
	sessionID *string
	token     *string

	rootCmd = &cobra.Command{
		Use:     "grow-with-stl-go",
		Short:   "grow-with-stl-go is a sample go application developed by stl-go for demonstration purposes, this is its REST example client",
		Run:     runAutomatedProcess,
		Version: version(),
	}
)

type wsMessage struct {
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
	Authentication *authentication `json:"authentication,omitempty"`
	Token          *string         `json:"token,omitempty"`
	RefreshToken   *string         `json:"refreshToken,omitempty"`
}

type authentication struct {
	ID       *string `json:"id,omitempty"`
	Password *string `json:"password,omitempty"`
}

func init() {
	// Add a 'version' command, in addition to the '--version' option that is auto created
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Show version",
		Long:  "Version for grow with stl-go REST client binary",
		Run: func(cmd *cobra.Command, _ []string) {
			out := cmd.OutOrStdout()
			fmt.Fprintln(out, "grow with stl-go REST client version", version())
		},
	})

	// Add the user
	u := ""
	user = &u
	rootCmd.Flags().StringVarP(
		user,
		"user",
		"u",
		"username",
		"The user used for the REST request",
	)

	// Add the user
	p := ""
	password = &p
	rootCmd.Flags().StringVarP(
		password,
		"passwd",
		"p",
		"password",
		"The password for the user specified for the REST request",
	)
}

// Version returns application version
func version() string {
	// pull the file version if it's available
	if fileVersion, err := os.ReadFile("version"); err == nil {
		return string(fileVersion)
	}
	return "dev-version"
}

func runAutomatedProcess(_ *cobra.Command, _ []string) {
	log.Info("Start automated requests")
	caCerts, err := getCertPool()
	if err != nil || caCerts == nil {
		log.Fatalf("ssl cert error %s", err)
	}

	signal.Notify(onClose, os.Interrupt)

	u := url.URL{Scheme: "wss", Host: host, Path: "/ws/v1.0.0"}
	log.Infof("connecting to %s", u.String())

	dialer := *websocket.DefaultDialer
	dialer.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: false,
		RootCAs:            caCerts,
		MinVersion:         tls.VersionTLS12,
	}

	ws, _, err = dialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer ws.Close()

	go onOpen()

	for {
		select {
		case err := <-onError:
			log.Error(err)
			onClose <- syscall.SIGINT
		case <-onClose:
			log.Info("CLosing the websocket connection ane exiting")
			ws.Close()
			return
		}
	}
}

func onOpen() {
	for {
		var message *wsMessage
		err := ws.ReadJSON(&message)
		if err != nil {
			onError <- err
			break
		}
		go onMessage(message)
	}
}

func onMessage(message *wsMessage) {
	if message != nil && message.Route != nil {
		switch *message.Route {
		case "websocketclient":
			handleClientMessages(message)
		case "seeds":
			displayJSON(message.Data)
		default:
			onError <- fmt.Errorf("unable to determine what to do with %s route", *message.Route)
		}
	}
}

func webSocketSend(message *wsMessage) error {
	writeMutex.Lock()
	defer writeMutex.Unlock()
	message.Timestamp = utils.CurrentTimeInMillis()
	message.SessionID = sessionID
	message.Token = token

	return ws.WriteJSON(message)
}

func handleClientMessages(message *wsMessage) {
	if message != nil && message.Type != nil {
		switch *message.Type {
		case "initialize":
			if message.SessionID != nil {
				sessionID = message.SessionID
				login()
				return
			}
			onError <- errors.New("no session id on initialize")
		case "auth":
			log.Infof("Authentication was %s", *message.SubComponent)
			if message.Token != nil {
				token = message.Token
				getInventory()
				return
			}
			onError <- errors.New("no toke on auth response")
		default:
			onError <- fmt.Errorf("unable to determine what to do with client type %s", *message.Type)
		}
	}
}

func login() {
	log.Info("Attempting to authenticate")

	route := "websocketclient"
	wsType := "auth"
	component := "authenticate"

	if err := webSocketSend(&wsMessage{
		Route:     &route,
		Type:      &wsType,
		Component: &component,
		Authentication: &authentication{
			ID:       user,
			Password: password,
		},
	}); err != nil {
		onError <- err
	}
}

func getInventory() {
	log.Info("Attempting to get Inventory")

	route := "seeds"
	wsType := "getInventory"
	component := "getInventory"

	if err := webSocketSend(&wsMessage{
		Route:     &route,
		Type:      &wsType,
		Component: &component,
	}); err != nil {
		onError <- err
	}
}

func displayJSON(data interface{}) {
	if data != nil {
		if jo, ok := data.(map[string]interface{}); ok {
			jsonBytes, err := json.MarshalIndent(jo, "", "\t")
			if err != nil {
				log.Fatal(err)
			}

			log.Infof("\n%s", string(jsonBytes))
		}
	}
}

func getCertPool() (*x509.CertPool, error) {
	caCert, err := os.ReadFile("etc/cert.pem")
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	return caCertPool, nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
