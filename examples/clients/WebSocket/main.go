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
	"fmt"
	"net/url"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"

	"stl-go/grow-with-stl-go/pkg/log"
)

var (
	host = "localhost:10443"

	ws *websocket.Conn
	// proxy      *string
	// extraCerts *string

	user     *string
	password *string

	rootCmd = &cobra.Command{
		Use:     "grow-with-stl-go",
		Short:   "grow-with-stl-go is a sample go application developed by stl-go for demonstration purposes, this is its REST example client",
		Run:     runAutomatedProcess,
		Version: version(),
	}
)

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

	onMessage := make(chan string)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

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

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := ws.ReadMessage()
			if err != nil {
				log.Error("read:", err)
				return
			}
			onMessage <- string(message)
		}
	}()

	for {
		select {
		case m := <-onMessage:
			log.Info(m)
		case <-interrupt:
			log.Info("interrupt")
			ws.Close()
			return
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
