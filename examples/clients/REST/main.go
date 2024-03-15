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
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
	"time"

	"stl-go/grow-with-stl-go/pkg/log"

	"github.com/spf13/cobra"
)

var (
	methods    = map[string]struct{}{http.MethodGet: {}, http.MethodPut: {}, http.MethodPost: {}, http.MethodDelete: {}}
	baseURL    = "https://localhost:10443"
	proxy      *string
	extraCerts *string
	client     *http.Client

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
	cert := "etc/cert.pem"
	extraCerts = &cert
	token, sessinID := login()
	if token != nil && sessinID != nil {
		log.Infof("%s %s", *token, *sessinID)
		headers := map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", *token),
			"sessionID":     *sessinID,
		}
		inventory := getInventory(headers)
		category := "Herb"
		commonName := "Basil"
		seedID := getHerbID(&category, &commonName, inventory)
		if seedID != nil {
			getDetail(&category, seedID, headers)
			purchaseSeed(&category, seedID, 1, headers)
			getDetail(&category, seedID, headers)
		}
	}
}

func login() (token, sessionID *string) {
	fullURL, err := url.JoinPath(baseURL, "/REST/v1.0.0/token")
	if err != nil {
		log.Fatal(err)
	}

	loginStruct := map[string]string{
		"id":       *user,
		"password": *password,
	}

	// marshall the StructJSON object to a byte array
	payloadBytes, err := json.Marshal(loginStruct)
	if err != nil {
		log.Fatal(err)
	}

	payload := string(payloadBytes)
	log.Info(payload)

	txt, httpStatusCode, err := httpRequest(fullURL, http.MethodPost, nil, &payload)
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("%s had status code of %d with the text of %s", fullURL, httpStatusCode, *txt)

	var jo map[string]string
	if err = json.Unmarshal([]byte(*txt), &jo); err != nil {
		log.Fatal(err)
	}

	tokenStr, ok := jo["token"]
	if !ok {
		log.Fatal("no token in server response")
	}

	sessionIDStr, ok := jo["sessionID"]
	if !ok {
		log.Fatal("no sessionID in server response")
	}
	return &tokenStr, &sessionIDStr
}

func getInventory(headers map[string]string) map[string]interface{} {
	fullURL, err := url.JoinPath(baseURL, "/REST/v1.0.0/seeds/getInventory")
	if err != nil {
		log.Fatal(err)
	}

	txt, httpStatusCode, err := httpRequest(fullURL, http.MethodGet, headers, nil)
	if err != nil {
		log.Fatal(err)
	}

	var jo map[string]interface{}
	err = json.Unmarshal([]byte(*txt), &jo)
	if err != nil {
		log.Fatal(err)
	}

	jsonBytes, err := json.MarshalIndent(jo, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("%s had status code of %d with the data:\n%s", fullURL, httpStatusCode, string(jsonBytes))
	return jo
}

func getHerbID(category, commonName *string, jo map[string]any) *string {
	if herb, ok := jo[*category]; ok {
		if items, ok := herb.(map[string]interface{})["items"]; ok {
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

func getDetail(category, seedID *string, headers map[string]string) {
	fullURL, err := url.JoinPath(baseURL, fmt.Sprintf("/REST/v1.0.0/seeds/getDetail/%s/%s", *category, *seedID))
	if err != nil {
		log.Fatal(err)
	}

	txt, httpStatusCode, err := httpRequest(fullURL, http.MethodGet, headers, nil)
	if err != nil {
		log.Fatal(err)
	}

	var jo map[string]interface{}
	err = json.Unmarshal([]byte(*txt), &jo)
	if err != nil {
		log.Fatal(err)
	}

	jsonBytes, err := json.MarshalIndent(jo, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("%s had status code of %d with the data:\n%s", fullURL, httpStatusCode, string(jsonBytes))
}

func purchaseSeed(category, seedID *string, quantity int, headers map[string]string) {
	fullURL, err := url.JoinPath(baseURL, "/REST/v1.0.0/seeds/purchase")
	if err != nil {
		log.Fatal(err)
	}

	loginStruct := map[string]any{
		"id":       *seedID,
		"category": *category,
		"quantity": quantity,
	}

	// marshall the StructJSON object to a byte array
	payloadBytes, err := json.Marshal(loginStruct)
	if err != nil {
		log.Fatal(err)
	}

	payload := string(payloadBytes)
	log.Info(payload)

	txt, httpStatusCode, err := httpRequest(fullURL, http.MethodPost, headers, &payload)
	if err != nil {
		log.Fatal(err)
	}

	var jo map[string]interface{}
	err = json.Unmarshal([]byte(*txt), &jo)
	if err != nil {
		log.Fatal(err)
	}

	jsonBytes, err := json.MarshalIndent(jo, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("%s had status code of %d with the data:\n%s", fullURL, httpStatusCode, string(jsonBytes))
}

func httpRequest(requestedURL, method string, headers map[string]string, payload *string) (responseText *string, httpStatusCode int, err error) {
	startTime := time.Now()
	_, ok := methods[method]
	if !ok {
		return nil, 503, fmt.Errorf("invalid method requested %s", method)
	}

	var request *http.Request
	if payload != nil && strings.EqualFold(method, http.MethodPost) {
		request, err = http.NewRequest(method, requestedURL, strings.NewReader(*payload))
		if err != nil {
			return nil, 503, err
		}
		request.Header.Add("content-type", "application/json")
	} else {
		request, err = http.NewRequest(method, requestedURL, http.NoBody)
		if err != nil {
			return nil, 503, err
		}
	}

	if headers != nil {
		for key, value := range headers {
			request.Header.Add(key, value)
		}
	}

	c, err := getClient()
	if err != nil {
		return nil, 503, err
	}

	response, err := c.Do(request)
	if err != nil || response == nil || response.Body == nil {
		log.Errorf("%s returned error %s in %dms", requestedURL, err, time.Since(startTime).Abs().Milliseconds())
		return nil, 503, err
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Errorf("%s HTTP Code %d in %dms", requestedURL, response.StatusCode, time.Since(startTime).Abs().Milliseconds())
		return nil, 503, err
	}

	responseRawText := string(body)
	httpStatusCode = response.StatusCode

	log.Debugf("%s HTTP Code %d response %d bytes in %dms",
		requestedURL, httpStatusCode, int64(uintptr(len(body))*reflect.TypeOf(body).Elem().Size()), time.Since(startTime).Abs().Milliseconds())
	return &responseRawText, httpStatusCode, nil
}

func getClient() (*http.Client, error) {
	if client == nil {
		client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					MinVersion: tls.VersionTLS12,
				},
			},
		}

		var certPool *x509.CertPool
		var err error
		if extraCerts != nil {
			certPool, err = getCertPool()
			if err != nil {
				return nil, err
			}
		}

		if proxy != nil {
			proxyURL, err := url.Parse(*proxy)
			if err != nil {
				return nil, err
			}

			if extraCerts != nil {
				client = &http.Client{Transport: &http.Transport{
					TLSClientConfig: &tls.Config{
						MinVersion: tls.VersionTLS12,
						RootCAs:    certPool,
					},
					Proxy: http.ProxyURL(proxyURL),
				}}
			} else {
				client = &http.Client{Transport: &http.Transport{
					TLSClientConfig: &tls.Config{
						MinVersion: tls.VersionTLS12,
					},
					Proxy: http.ProxyURL(proxyURL),
				}}
			}
		} else if certPool != nil {
			client = &http.Client{Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					MinVersion: tls.VersionTLS12,
					RootCAs:    certPool,
				},
			}}
		}
	}

	return client, nil
}

func getCertPool() (*x509.CertPool, error) {
	caCert, err := os.ReadFile(*extraCerts)
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
