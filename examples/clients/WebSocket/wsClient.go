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
	"strings"
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

	user     string
	password string
	wait     = false

	osInterrupt = make(chan os.Signal, 1)

	// WebSocket send message mutex
	writeMutex sync.Mutex

	// Things used in WebSocket messages
	sessionID *string
	token     *string

	rootCmd = &cobra.Command{
		Use:   "grow-with-stl-go",
		Short: "grow-with-stl-go is a sample go application developed by stl-go for demonstration purposes, this is its WebSocket example client",
		Run:   runAutomatedProcess,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				log.Info("No args entered showing default help\n\n")
				if err := cmd.Help(); err != nil {
					log.Error(err)
				}
				os.Exit(0)
			}
			return nil
		},
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
	rootCmd.Flags().StringVarP(
		&user,
		"user",
		"u",
		"username",
		"The user used for the REST request",
	)

	// Add the password
	rootCmd.Flags().StringVarP(
		&password,
		"passwd",
		"p",
		"password",
		"The password for the user specified for the REST request",
	)

	// add the wait flag
	rootCmd.Flags().BoolVarP(
		&wait,
		"wait",
		"w",
		false,
		"wait after the auto transactions",
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

func getCertPool() (*x509.CertPool, error) {
	caCert, err := os.ReadFile("etc/cert.pem")
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	return caCertPool, nil
}

func handleOsInterrupt() {
	<-osInterrupt
	log.Info("Closing active session")
	onClose()
	log.Info("Exiting the websocket client")
	os.Exit(0)
}

func runAutomatedProcess(_ *cobra.Command, _ []string) {
	log.Info("Start automated requests")
	closeWait := make(chan byte, 1)
	caCerts, err := getCertPool()
	if err != nil || caCerts == nil {
		log.Fatalf("ssl cert error %s", err)
	}

	signal.Notify(osInterrupt, os.Interrupt, syscall.SIGTERM)
	go handleOsInterrupt()

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
	<-closeWait
}

func onOpen() {
	for {
		var message *wsMessage
		err := ws.ReadJSON(&message)
		if err != nil {
			onError(err)
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
			handleSeedMessages(message)
		default:
			onError(fmt.Errorf("unable to determine what to do with %s route", *message.Route))
		}
	}
}

func onError(err error) {
	log.Error(err)
	onClose()
}

func onClose() {
	log.Info("CLosing the websocket connection and exiting")
	ws.Close()
	osInterrupt <- syscall.SIGINT
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
			onError(errors.New("no session id on initialize"))
		case "auth":
			log.Infof("Authentication was %s", *message.SubComponent)
			if message.Token != nil {
				token = message.Token
				getInventory()
				return
			}
			onError(errors.New("no token on auth response"))
		default:
			onError(fmt.Errorf("unable to determine what to do with client type %s", *message.Type))
		}
	}
}

func handleSeedMessages(message *wsMessage) {
	if message != nil && message.Type != nil {
		switch *message.Type {
		case "getInventory":
			if message.Data != nil {
				inventory := displayJSON(message)
				category := "Herb"
				commonName := "Basil"
				getSeedDetail(getSeedID(&category, &commonName, inventory))
				return
			}
			onError(errors.New("no session id on initialize"))
		case "getDetail":
			displayJSON(message)
			purchaseSeed(message.Component, message.SubComponent)
		case "purchase":
			displayJSON(message)
			if !wait {
				log.Info("Purchase completed exiting the client")
				onClose()
			}
			log.Info("Waiting for other messages from the server")
		default:
			onError(fmt.Errorf("unable to determine what to do with client type %s", *message.Type))
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
			ID:       &user,
			Password: &password,
		},
	}); err != nil {
		onError(err)
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
		onError(err)
	}
}

func getSeedDetail(seedID *string) {
	if seedID != nil {
		log.Infof("Attempting to get detail for seed ID %s", *seedID)

		route := "seeds"
		wsType := "getDetail"
		component := "Herb"

		if err := webSocketSend(&wsMessage{
			Route:        &route,
			Type:         &wsType,
			Component:    &component,
			SubComponent: seedID,
		}); err != nil {
			onError(err)
		}
		return
	}
	onError(errors.New("cannot get detail from nil seed id"))
}

func purchaseSeed(seedType, seedID *string) {
	if seedID != nil {
		log.Infof("Attempting to purchase seed ID %s", *seedID)

		route := "seeds"
		wsType := "purchase"

		if err := webSocketSend(&wsMessage{
			Route:        &route,
			Type:         &wsType,
			Component:    seedType,
			SubComponent: seedID,
			Data: &map[string]any{
				"id":       seedID,
				"quantity": 1,
			},
		}); err != nil {
			onError(err)
		}
	}
}

func getSeedID(categoryName, commonName *string, jo map[string]any) *string {
	if category, ok := jo[*categoryName]; ok {
		if items, ok := category.(map[string]interface{})["items"]; ok {
			for seedID, details := range items.(map[string]interface{}) {
				if herbName, ok := details.(map[string]interface{})["commonName"]; ok {
					if herbStr, ok := herbName.(string); ok {
						if strings.EqualFold(herbStr, *commonName) {
							return &seedID
						}
					}
				}
			}
		}
	}
	return nil
}

func displayJSON(message *wsMessage) map[string]any {
	if message != nil && message.Data != nil {
		jsonBytes, err := json.MarshalIndent(message, "", "\t")
		if err != nil {
			log.Fatal(err)
		}

		log.Infof("\n%s", string(jsonBytes))

		if jo, ok := message.Data.(map[string]any); ok {
			return jo
		}
	}
	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
