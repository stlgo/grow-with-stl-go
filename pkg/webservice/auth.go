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
	"errors"
	"fmt"
	"strings"
	"time"

	"stl-go/grow-with-stl-go/pkg/configs"
	"stl-go/grow-with-stl-go/pkg/log"

	jwt "github.com/golang-jwt/jwt/v4"
)

// Create the JWT key used to create the signature
// TODO: use a private key for this instead of a phrase
var jwtKey = []byte("airshipUI_JWT_key")

const (
	username   = "username"
	sessionID  = "sessionID"
	expiration = "exp"
)

// The UI will either request authentication or validation, handle those situations here
func handleWebSocketAuth(request, response *configs.WsMessage) error {
	defer log.FunctionTimer()()

	err := errors.New("not authenticated").Error()
	response.Error = &err
	denied := configs.Denied
	response.SubComponent = &denied

	if request.SubComponent != nil && strings.EqualFold(*request.SubComponent, configs.Authenticate) && request.Authentication != nil {
		if request.SessionID != nil {
			token, err := createJWTToken(request.Authentication.ID, request.SessionID)
			if err != nil || token == nil {
				return err
			}

			if session, ok := sessions[*request.SessionID]; ok {
				session.jwt = token
			}

			approved := configs.Approved
			response.SubComponent = &approved
			response.Token = token
			response.Error = nil

			return nil
		}
	}
	return errors.New("unable to process authentication request")
}

// validate JWT claim
func validateJWTClaim(token *jwt.Token, requestSessionID *string, request *configs.WsMessage) (*string, error) {
	if claim, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// extract the session id from the claim
		sessionID, ok := claim[sessionID].(string)
		if !ok || requestSessionID == nil || sessionID != *requestSessionID {
			return nil, fmt.Errorf("invalid JWT session id %s attempted for %s session", *requestSessionID, sessionID)
		}

		// extract the user from the claim
		if username, ok := claim[username].(string); ok {
			// test to see if we need to refresh the token
			if request != nil {
				go testForRefresh(claim, request)
			}
			return &username, nil
		}
	}
	return nil, errors.New("invalid JWT token")
}

// create a JWT (JSON Web Token)
func createJWTToken(userid, sessionid *string) (*string, error) {
	if userid != nil && sessionid != nil {
		// set some claims
		claims := make(jwt.MapClaims)
		claims[username] = *userid
		claims[sessionID] = *sessionid
		claims[expiration] = time.Now().Add(time.Hour * 1).Unix()

		// create the token
		jwtClaim := jwt.New(jwt.SigningMethodHS256)
		jwtClaim.Claims = claims

		// Sign and get the complete encoded token as string
		token, err := jwtClaim.SignedString(jwtKey)
		return &token, err
	}
	return nil, errors.New("nil user id of session id, cannot create JWT")
}

// from time to time we might want to send a refresh token to the UI.  The UI should not be in charge of requesting it
func testForRefresh(claim jwt.MapClaims, request *configs.WsMessage) {
	// for some reason the exp is stored as an float and not an int in the claim conversion
	// so we do a little dance and cast some floats to ints and everyone goes on with their lives
	if exp, ok := claim[expiration].(float64); ok {
		if int64(exp) < time.Now().Add(time.Minute*15).Unix() {
			createRefreshToken(claim, request)
		}
	}
}

// createRefreshToken will create an oauth2 refresh token based on the timeout on the UI
func createRefreshToken(claim jwt.MapClaims, request *configs.WsMessage) {
	// test to see if the session is still in existence before firing off a refresh
	if session, ok := sessions[*request.SessionID]; ok {
		// add the new expiration to the claim
		claim[expiration] = time.Now().Add(time.Hour * 1).Unix()

		// create the token
		jwtClaim := jwt.New(jwt.SigningMethodHS256)
		jwtClaim.Claims = claim

		// Sign and get the complete encoded token as string
		refreshToken, err := jwtClaim.SignedString(jwtKey)
		if err != nil {
			log.Error(err)
			return
		}

		wsClient := configs.WebsocketClient
		auth := configs.Auth
		refresh := configs.Refresh

		// test to see if the session is still in existence before firing off a message
		if err = session.webSocketSend(&configs.WsMessage{
			Type:         &wsClient,
			Component:    &auth,
			SubComponent: &refresh,
			RefreshToken: &refreshToken,
		}); err != nil {
			session.onError(err)
		}
	}
}
