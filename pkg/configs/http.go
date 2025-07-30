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
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"time"

	"stl-go/grow-with-stl-go/pkg/log"
	"stl-go/grow-with-stl-go/pkg/utils"
)

var (
	proxyURL   *url.URL
	caCertPool *x509.CertPool
)

// setup the http client, add proxy and extra root ca certs if necessary
func getHTTPClient() (*http.Client, error) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS13,
		},
	}

	if GrowSTLGo.Proxy != nil && GrowSTLGo.Proxy.URL != nil {
		if proxyURL == nil {
			u, urlErr := url.Parse(*GrowSTLGo.Proxy.URL)
			if urlErr != nil {
				return nil, urlErr
			}
			proxyURL = u
		}

		if GrowSTLGo.Proxy.ExtraCACerts != nil {
			if caCertPool == nil {
				caCert, certErr := os.ReadFile(*GrowSTLGo.Proxy.ExtraCACerts)
				if certErr != nil {
					return nil, certErr
				}
				caCertPool = x509.NewCertPool()
				caCertPool.AppendCertsFromPEM(caCert)
			}

			transport = &http.Transport{
				TLSClientConfig: &tls.Config{
					MinVersion: tls.VersionTLS13,
					RootCAs:    caCertPool,
				},
				Proxy: http.ProxyURL(proxyURL),
			}
		} else {
			transport = &http.Transport{
				TLSClientConfig: &tls.Config{
					MinVersion: tls.VersionTLS13,
				},
				Proxy: http.ProxyURL(proxyURL),
			}
		}
	}

	return &http.Client{Transport: transport}, nil
}

// HTTPRequestHelper will be used by other various methods to call endpoints
func HTTPRequestHelper(r *http.Request, start time.Time) (responseText *string, httpStatusCode *int, err error) {
	requestClient, clientErr := getHTTPClient()
	if clientErr != nil {
		return nil, nil, clientErr
	}

	response, responseErr := requestClient.Do(r)
	if responseErr != nil || response == nil || response.Body == nil {
		return nil, nil, fmt.Errorf("%s returned error %s in %s", r.URL.String(), responseErr, log.FormatMilliseconds(time.Since(start).Abs().Milliseconds()))
	}

	defer response.Body.Close()
	body, bodyErr := io.ReadAll(response.Body)
	if bodyErr != nil {
		return nil, nil, fmt.Errorf("%s returned error %s in %s", r.URL.String(), bodyErr, log.FormatMilliseconds(time.Since(start).Abs().Milliseconds()))
	}

	responseText = utils.StringPointer(string(body))
	httpStatusCode = &response.StatusCode

	log.Tracef("URL %s method %s HTTP Status Code %d response %d bytes in %s",
		r.URL.String(),
		r.Method,
		*httpStatusCode,
		int64(uintptr(len(body))*reflect.TypeOf(body).Elem().Size()), log.FormatMilliseconds(time.Since(start).Abs().Milliseconds()))

	return responseText, httpStatusCode, nil
}
