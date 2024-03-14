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

	user     string
	password string

	version = "(dev-version)"

	rootCmd = &cobra.Command{
		Use:     "grow-with-stl-go",
		Short:   "grow-with-stl-go is a sample go application developed by stl-go for demonstration purposes",
		Run:     runAutomatedProcess,
		Version: Version(),
	}
)

func init() {
	// pull the file version if it's available
	if fileVersion, err := os.ReadFile("../../../version"); err == nil {
		version = string(fileVersion)
	}

	// Add a 'version' command, in addition to the '--version' option that is auto created
	rootCmd.AddCommand(versionCmd())

	// Add the user
	rootCmd.Flags().StringVarP(
		&user,
		"user",
		"u",
		"username",
		"The user used for the REST request",
	)

	// Add the user
	rootCmd.Flags().StringVarP(
		&password,
		"passwd",
		"p",
		"user",
		"The password for the user specified for the REST request",
	)
}

func versionCmd() *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Show version",
		Long:  "Version for grow with stl-go binary",
		Run: func(cmd *cobra.Command, _ []string) {
			out := cmd.OutOrStdout()

			fmt.Fprintln(out, "grow with stl-go version", Version())
		},
	}
	return versionCmd
}

// Version returns application version
func Version() string {
	return version
}

func runAutomatedProcess(_ *cobra.Command, _ []string) {
	log.Info("Start automated requests")
	cert := "etc/cert.pem"
	extraCerts = &cert
	fullURL, err := url.JoinPath(baseURL, "index.html")
	if err != nil {
		log.Fatal(err)
	}
	txt, httpStatusCode, err := httpRequest(fullURL, http.MethodGet, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("%s had status code of %d", fullURL, httpStatusCode)
	log.Info(*txt)
}

func httpRequest(requestedURL, method string, payload *string) (responseText *string, httpStatusCode int, err error) {
	startTime := time.Now()
	_, ok := methods[method]
	if !ok {
		return nil, 503, fmt.Errorf("invalid method requested %s", method)
	}

	var request *http.Request
	if payload != nil {
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
