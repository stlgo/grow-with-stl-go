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
	"time"

	"stl-go/grow-with-stl-go/pkg/configs"
	"stl-go/grow-with-stl-go/pkg/log"

	jwt "github.com/golang-jwt/jwt/v4"
)

// Create the JWT key used to create the signature
// TODO: use a private key for this instead of a phrase
var jwtKey = []byte("grow-with-stl-go!")

const (
	username   = "username"
	sessionID  = "sessionID"
	expiration = "exp"
)

func validateAPIUser(body []byte) (*string, error) {
	var authRequest configs.Authentication
	if err := json.Unmarshal(body, &authRequest); err != nil {
		log.Errorf("bad request body: %s", body)
		return authRequest.ID, err
	}

	if authRequest.ID == nil || authRequest.Password == nil {
		log.Errorf("bad request body: %s", body)
		return authRequest.ID, errors.New("bad request body")
	}

	apiUser, ok := configs.GrowSTLGo.APIUsers[*authRequest.ID]
	if !ok || apiUser == nil || apiUser.Active == nil || !*apiUser.Active {
		log.Errorf("User %s not found or is inactive", *authRequest.ID)
		return nil, errors.New("Unauthorized")
	}

	if err := apiUser.Authentication.ValidateAuthentication(authRequest.Password); err != nil {
		log.Errorf("Passwords did not match for request by %s", *authRequest.ID)
		return nil, errors.New("Unauthorized")
	}

	return authRequest.ID, nil
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
func createJWTToken(userid, sessionid *string) (token *string, validTill *int64, err error) {
	if userid != nil && sessionid != nil {
		// set some claims
		validTill := time.Now().Add(time.Hour * 1).UnixMilli()

		claims := make(jwt.MapClaims)
		claims[username] = *userid
		claims[sessionID] = *sessionid
		claims[expiration] = validTill

		// create the token
		jwtClaim := jwt.New(jwt.SigningMethodHS256)
		jwtClaim.Claims = claims

		// Sign and get the complete encoded token as string
		token, err := jwtClaim.SignedString(jwtKey)
		return &token, &validTill, err
	}
	return nil, nil, errors.New("nil user id of session id, cannot create JWT")
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
			log.Error(err)
			session.onError()
		}
	}
}
