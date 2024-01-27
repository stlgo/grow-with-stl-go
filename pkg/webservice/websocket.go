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
	"net/http"
	"strings"
	"sync"
	"time"

	"stl-go/grow-with-stl-go/pkg/configs"
	"stl-go/grow-with-stl-go/pkg/log"
	"stl-go/grow-with-stl-go/pkg/utils"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Session is a struct to hold information about a given session
type session struct {
	requestHost *string
	sessionID   *string
	jwt         *string
	writeMutex  sync.Mutex
	ws          *websocket.Conn
	lastUsed    *int64
	name        *string
	closing     bool
}

// sessions keeps track of open websocket sessions
var sessions = map[string]*session{}

// gorilla ws specific HTTP upgrade to WebSockets
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// start up the idle hands tester
func init() {
	go idleHandsTester()
}

// this is a way to allow for arbitrary messages to be processed by the backend
// the message of a specifc component is shunted to that subsystem for further processing
var websocketFuncMap = map[string]func(*string, *configs.WsMessage) *configs.WsMessage{
	configs.WebsocketClient: handleMessage,
}

// AppendToFunctionMap allows us to break up the circular reference from the other packages
// It does however require them to implement an init function to append them
// TODO: maybe some form of an interface to enforce this may be necessary?
func AppendToFunctionMap(requestType string, function func(*string, *configs.WsMessage) *configs.WsMessage) {
	log.Debugf("Regestering %s as a WebSocket Endpoint", requestType)
	websocketFuncMap[requestType] = function
}

// handle the origin request & upgrade to websocket
func onOpen(response http.ResponseWriter, request *http.Request) {
	// gorilla ws will give a 403 on a cross origin request, so to silence its complaints
	// upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// upgrade to websocket protocol over http
	wsConn, err := upgrader.Upgrade(response, request, nil)
	if err != nil {
		log.Errorf("Could not open websocket connection from: %s\n", request.Host)
		http.Error(response, "Could not open websocket connection", http.StatusBadRequest)
		return
	}

	remoteHost := request.Host
	session := newSession(&remoteHost, wsConn)
	log.Debugf("WebSocket session %s established with %s\n", *session.sessionID, session.ws.RemoteAddr().String())

	go session.onMessage()
}

// handle messaging to the client
func (session *session) onMessage() {
	// just in case clean up the websocket
	defer session.onClose()

	for {
		var request configs.WsMessage
		err := session.ws.ReadJSON(&request)
		if err != nil {
			session.onError(err)
			break
		}

		// this has to be a go routine otherwise it will block any incoming messages waiting for a command return
		go func() {
			// test the auth token for request validity on non auth requests
			if request.Type != nil && !strings.EqualFold(*request.Type, configs.UI) &&
				request.Component != nil && !strings.EqualFold(*request.Type, configs.Auth) &&
				request.Token != nil {
				err = session.validateWSToken(&request)
			}
			if err != nil {
				// deny the request if we get a bad token, this will force the UI to a login screen
				log.Error(err)
				e := errors.New(configs.UnauthorizedError).Error()

				ui := configs.UI
				auth := configs.Auth
				denied := configs.Denied

				response := configs.WsMessage{
					Type:         &ui,
					Component:    &auth,
					SubComponent: &denied,
					Error:        &e,
				}
				if err = session.webSocketSend(&response); err != nil {
					session.onError(err)
				}
				defer session.onClose()
				session.handleRequest(&request)

				return
			}
		}()
	}
}

// common websocket close with logging
func (session *session) onClose() {
	if !session.closing {
		session.closing = true
		n := *session.name
		if session.name == nil {
			n = "unknown"
		}
		log.Infof("closing websocket for %s session %s", n, *session.sessionID)
		session.ws.Close()
		delete(sessions, *session.sessionID)
	}
}

// common websocket error handling with logging
func (session *session) onError(err error) {
	if err != nil {
		log.Error(err)
	}
	session.onClose()
}

func (session *session) handleRequest(request *configs.WsMessage) {
	if request.Type != nil {
		if handleMessageFunc, ok := websocketFuncMap[*request.Type]; ok {
			// reset the idle timer if it's appropriate
			go session.resetIdleTimer(request)
			// do the rest
			response := handleMessageFunc(session.sessionID, request)
			if err := session.webSocketSend(response); err != nil {
				session.onError(err)
			}
			return
		}
		if err := session.webSocketSend(requestErrorHelper(fmt.Sprintf("requested type: %s is not found", *request.Type), request)); err != nil {
			session.onError(err)
		}
		return
	}
	session.onError(errors.New("invalid request type"))
}

func (session *session) resetIdleTimer(request *configs.WsMessage) {
	if !strings.EqualFold(*request.Component, configs.Keepalive) {
		session.lastUsed = utils.CurrentTimeInMillis()
	}
}

func (session *session) validateWSToken(request *configs.WsMessage) error {
	if request.Token != nil {
		tokenString := request.Token
		if request.RefreshToken != nil {
			tokenString = request.RefreshToken
		}

		token, err := jwt.Parse(*tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtKey, nil
		})

		if err != nil {
			return err
		}

		_, err = validateJWTClaim(token, session.sessionID, request)
		if err != nil {
			return err
		}

		return nil
	}
	return errors.New("invalid token")
}

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

func handleMessage(sessionID *string, request *configs.WsMessage) *configs.WsMessage {
	response := configs.WsMessage{
		Type:         request.Type,
		Component:    request.Component,
		SubComponent: request.SubComponent,
	}

	switch *request.Component {
	case configs.Keepalive:
		log.Trace(fmt.Sprintf("keepalive received for session %s", *sessionID))
	case configs.Auth:
		if err := handleWebSocketAuth(request, &response); err != nil {
			errStr := err.Error()
			response.Error = &errStr
			if sessionID != nil {
				if session, ok := sessions[*sessionID]; ok {
					if err := session.webSocketSend(&response); err != nil {
						session.onError(err)
					}
					session.onClose()
				}
			}
		}
	default:
		err := "not implemented"
		log.Error(err)
		response.Error = &err
	}
	return &response
}

func idleHandsTester() {
	time.Sleep(time.Duration(60-time.Now().Local().Second()) * time.Second)
	for range time.NewTicker(10 * time.Second).C {
		for _, session := range sessions {
			if session.lastUsed != nil {
				// 10 minute timeout
				if (time.Now().UnixMilli() - *session.lastUsed) > 120000 {
					session.onClose()
				}
				// idle abandoned connections are disconnected at 1 minute
				if (time.Now().UnixMilli()-*session.lastUsed) > 60000 && session.name == nil {
					session.onClose()
				}
			}
			// close sessions without a last used timestamp
			session.onClose()
		}
	}
}

// formats an error response in the way that we're expecting on the UI
func requestErrorHelper(err string, request *configs.WsMessage) *configs.WsMessage {
	return &configs.WsMessage{
		Type:      request.Type,
		Component: request.Component,
		Error:     &err,
	}
}

// newSession generates a new session
func newSession(requestHost *string, ws *websocket.Conn) *session {
	id := uuid.New().String()

	session := &session{
		requestHost: requestHost,
		sessionID:   &id,
		ws:          ws,
		lastUsed:    utils.CurrentTimeInMillis(),
	}

	// keep track of the session
	sessions[id] = session

	// send the init message to the client
	go session.sendInit()

	return session
}

// webSocketSend allows for the sender to be thread safe, we cannot write to the websocket at the same time
func (session *session) webSocketSend(response *configs.WsMessage) error {
	session.writeMutex.Lock()
	defer session.writeMutex.Unlock()
	response.Timestamp = utils.CurrentTimeInMillis()
	response.SessionID = session.sessionID

	return session.ws.WriteJSON(response)
}

// WebSocketSend allows of other packages to send a request for the websocket
func WebSocketSend(response *configs.WsMessage) error {
	if response.SessionID != nil {
		if session, ok := sessions[*response.SessionID]; ok {
			return session.webSocketSend(response)
		}
		return fmt.Errorf("session id %s not found", *response.SessionID)
	}
	return errors.New("no session id found in response")
}

// sendInit is generated on the onOpen event and sends the information the UI needs to startup
func (session *session) sendInit() {
	wsClient := configs.WebsocketClient
	initialize := configs.Initialize
	if err := session.webSocketSend(&configs.WsMessage{
		Type:      &wsClient,
		Component: &initialize,
	}); err != nil {
		log.Errorf("error receiving / sending init to session %s: %s\n", *session.sessionID, err)
	}
}

// Shutdown is called when the system is exiting to cleanly close all the current connections
func Shutdown() {
	for _, session := range sessions {
		session.onClose()
	}
}
