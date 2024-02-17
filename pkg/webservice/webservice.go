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
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"stl-go/grow-with-stl-go/pkg/configs"
	"stl-go/grow-with-stl-go/pkg/log"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

var (
	// this is a way to allow for arbitrary messages to be processed by the backend
	// the message of a specific component is shunted to that subsystem for further processing
	webServiceFunctionMap = map[string]func(w http.ResponseWriter, r *http.Request){}

	// these are things we don't want to server via REST or the UI
	webDenialPrefixes = regexp.MustCompile(`^(/.eslint.*|/pagelets.*|/package.*)`)
)

const (
	fontAwesomePrefix = "/node_modules/font-awesome/fonts/fontawesome-webfont"
)

// AppendToWebServiceFunctionMap allows us to break up the circular reference from the other packages
// It does however require them to implement an init function to append them
// TODO: maybe some form of an interface to enforce this may be necessary?
func AppendToWebServiceFunctionMap(requestType string, function func(w http.ResponseWriter, r *http.Request)) {
	log.Debugf("Regestering %s as a REST Endpoint", requestType)
	webServiceFunctionMap[requestType] = function
}

func handleRESTRequest(w http.ResponseWriter, r *http.Request) {
	uri := r.RequestURI

	if strings.EqualFold(uri, "/REST/v1.0.0/token") && strings.EqualFold(r.Method, http.MethodPost) {
		handelRESTAuthRequest(w, r)
		return
	}

	restURI := strings.TrimPrefix(uri, "/REST/v1.0.0/")
	if restFunction, ok := webServiceFunctionMap[restURI]; ok {
		id, err := handleRESTAuth(r)
		if err != nil {
			log.Error(err)
			http.Error(w, configs.BadRequestError, http.StatusBadRequest)
			return
		}
		log.Infof("User %s authenticated for %s on %s", *id, r.Method, uri)
		restFunction(w, r)
		return
	}

	serveFile(w, r)
}

// serveFile test if path and file exists, if it does send a page, else 404 or redirect
func serveFile(w http.ResponseWriter, r *http.Request) {
	uri := r.RequestURI
	// have to make a special condition for font awesome node modules
	if strings.HasPrefix(uri, fontAwesomePrefix) {
		uri = strings.Split(uri, "?")[0]
	}
	if !webDenialPrefixes.MatchString(uri) {
		path := filepath.Join(*configs.GrowSTLGo.WebService.StaticWebDir, uri)
		_, err := os.Stat(path)
		if err == nil {
			http.ServeFile(w, r, path)
			return
		}
	}

	// redirect to index.html on error
	http.Redirect(w, r, "/index.html", http.StatusFound)
}

func handelRESTAuthRequest(w http.ResponseWriter, r *http.Request) {
	defer log.FunctionTimer()()
	if r != nil {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Error(err)
			http.Error(w, configs.BadRequestError, http.StatusBadRequest)
			return
		}

		id, err := validateAPIUser(body)
		if err != nil {
			log.Error(err)
			http.Error(w, configs.UnauthorizedError, http.StatusUnauthorized)
			return
		}

		sessionID := uuid.New().String()
		token, err := createJWTToken(id, &sessionID)
		if err != nil {
			log.Error(err)
			http.Error(w, configs.IntenralServerError, http.StatusInternalServerError)
			return
		}

		b, err := json.Marshal(map[string]interface{}{"sessionID": sessionID, "token": token})
		if err != nil {
			log.Error(err)
			http.Error(w, configs.IntenralServerError, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		_, err = w.Write(b)
		if err != nil {
			log.Error(err)
			http.Error(w, configs.IntenralServerError, http.StatusInternalServerError)
			return
		}

		log.Infof("Token authentication %s successful for %s session %s", r.Method, *id, sessionID)
		return
	}
	log.Error("http request was nil")
	http.Error(w, configs.IntenralServerError, http.StatusInternalServerError)
}

func handleRESTAuth(r *http.Request) (*string, error) {
	defer log.FunctionTimer()()

	reqToken := strings.Split(r.Header.Get("Authorization"), "Bearer ")
	if len(reqToken) == 2 && reqToken[1] != "" {
		sessionID := r.Header.Get("sessionID")
		if sessionID == "" {
			return nil, errors.New("no sessionID found on request header")
		}

		token, err := jwt.Parse(reqToken[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtKey, nil
		})

		if err != nil {
			return nil, err
		}

		return validateJWTClaim(token, &sessionID, nil)
	}
	return nil, errors.New("bad token")
}
