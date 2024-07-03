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

package webservice

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"stl-go/grow-with-stl-go/pkg/configs"
	"stl-go/grow-with-stl-go/pkg/log"
)

// getCertificates returns the cert chain in a way that the net/http server struct expects
func getCertificates() (*[]tls.Certificate, error) {
	if configs.GrowSTLGo.WebService != nil && configs.GrowSTLGo.WebService.PrivateKey != nil && configs.GrowSTLGo.WebService.PublicKey != nil {
		cert, err := tls.LoadX509KeyPair(*configs.GrowSTLGo.WebService.PublicKey, *configs.GrowSTLGo.WebService.PrivateKey)
		if err != nil {
			return nil, err
		}
		var certSlice []tls.Certificate
		certSlice = append(certSlice, cert)
		return &certSlice, nil
	}
	return nil, fmt.Errorf("unable to load certificates, check the definition in %s", *configs.ConfigFile)
}

// WebServer will run the handler functions for WebSockets
func WebServer() {
	// make sure we have a good webservice config before we proceed
	if configs.GrowSTLGo.WebService == nil || configs.GrowSTLGo.WebService.Host == nil || configs.GrowSTLGo.WebService.Port == nil {
		log.Fatalf("Invalid webservice configuration in %s, host, port or static web dir is null", *configs.ConfigFile)
	}

	webServerMux := http.NewServeMux()

	// handle WebSocket endpoints
	webServerMux.HandleFunc("/ws/v1.0.0", onOpen)

	// establish routing to static web dir if defined in the config and handle REST endpoints
	if configs.GrowSTLGo.WebService.Vhosts != nil {
		for vhost, webRoot := range configs.GrowSTLGo.WebService.Vhosts {
			if webRoot != nil {
				log.Debugf("Attempting to serve static content for %s from %s", vhost, *webRoot)
				webServerMux.HandleFunc(fmt.Sprintf("%s/", vhost), handleRESTRequest)
			}
		}
	}

	// Calculate the address and start on the host and port specified in the config
	addr := fmt.Sprintf("%s:%d", *configs.GrowSTLGo.WebService.Host, *configs.GrowSTLGo.WebService.Port)
	log.Infof("Attempting to start webservice on %s", addr)

	certs, err := getCertificates()
	if err != nil || certs == nil {
		log.Fatalf("Invalid webservice configuration in %s, ssl cert error %s", *configs.ConfigFile, err)
	}

	// configure logging & TLS for the http server
	server := &http.Server{
		Addr: addr,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: false,
			ServerName:         *configs.GrowSTLGo.WebService.Host,
			Certificates:       *certs,
			MinVersion:         tls.VersionTLS13,
		},
		Handler:      webServerMux,
		ErrorLog:     log.Logger(),
		ReadTimeout:  180 * time.Second,
		WriteTimeout: 180 * time.Second,
		IdleTimeout:  180 * time.Second,
	}

	// kick off the server, and good luck
	log.Fatal(server.ListenAndServeTLS("", ""))
}
