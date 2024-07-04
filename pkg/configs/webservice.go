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
	"errors"

	"stl-go/grow-with-stl-go/pkg/cryptography"
	"stl-go/grow-with-stl-go/pkg/log"
)

// Vhost holds the bits to define the virtual webroots we'll be using
type Vhost struct {
	Name    *string `json:"name,omitempty"`
	WebRoot *string `json:"webRoot,omitempty"`
}

// WebService is the definition for the http/webservice/webserver
type WebService struct {
	Host       *string            `json:"host,omitempty"`
	Port       *int               `json:"port,omitempty"`
	PublicKey  *string            `json:"publicKey,omitempty"`
	PrivateKey *string            `json:"privateKey,omitempty"`
	Vhosts     map[string]*string `json:"vhosts,omitempty"`
}

func checkWebService() error {
	if GrowSTLGo.WebService == nil {
		log.Info("No webservice config found, generating ssl keys, host and port info")

		etcDir, err := getEtcDir()
		if err != nil {
			return err
		}

		if etcDir != nil {
			privateKeyFile, publicKeyFile, err := cryptography.GenerateDevSSL(etcDir)
			if err != nil {
				return err
			}

			if privateKeyFile != nil && publicKeyFile != nil {
				port := 10443
				host := "localhost"
				staticWebDir := "web/grow-with-stlgo-admin"

				GrowSTLGo.WebService = &WebService{
					Host:       &host,
					Port:       &port,
					PublicKey:  publicKeyFile,
					PrivateKey: privateKeyFile,
					Vhosts:     map[string]*string{host: &staticWebDir},
				}

				err = cryptography.CheckCertValidity(publicKeyFile)
				if err != nil {
					return err
				}

				rewriteConfig = true
				return nil
			}
			return errors.New("nil private key or public key, cannot continue")
		}
		return errors.New("nil etc dir, cannot continue")
	}
	return nil
}
