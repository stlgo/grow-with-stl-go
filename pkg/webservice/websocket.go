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
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"

	"stl-go/grow-with-stl-go/pkg/audit"
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
	user        *string
	closing     *bool
	isAdmin     *bool
	Vhost       *string
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
var websocketFuncMap = map[string]func(*configs.WsMessage, *configs.WsMessage){
	configs.WebsocketClient: handleMessage,
}

// AppendToWebsocketFunctionMap allows us to break up the circular reference from the other packages
// It does however require them to implement an init function to append them
// TODO: maybe some form of an interface to enforce this may be necessary?
func AppendToWebsocketFunctionMap(requestType string, function func(*configs.WsMessage, *configs.WsMessage)) {
	log.Debugf("Regestering '%s' as a WebSocket Endpoint", requestType)
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
		return
	}

	remoteHost := request.RemoteAddr
	vhost := strings.Split(request.Host, ":")[0]
	session := newSession(&remoteHost, &vhost, wsConn)
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
			log.Error(err)
			session.onError()
			break
		}

		// this has to be a go routine otherwise it will block any incoming messages waiting for a command return
		go func() {
			transaction := audit.NewWSTransaction(session.requestHost, session.user, &request)
			request.IsAdmin = session.isAdmin
			// test the auth token for request validity on non auth requests
			if request.Route != nil && !strings.EqualFold(*request.Route, configs.WebsocketClient) &&
				request.Type != nil && !strings.EqualFold(*request.Type, configs.Auth) &&
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
					log.Error(err)
					session.onError()
				}
				defer session.onClose()
				transactionHelper(transaction, false)
				return
			}
			session.handleRequest(&request, transaction)

			// the user is not populated on the login transaction, this will alleviate that issue
			if transaction.User == nil && session.user != nil {
				transaction.User = session.user
			}
			transactionHelper(transaction, true)
		}()
	}
}

// common websocket close with logging
func (session *session) onClose() {
	if session.closing == nil || !*session.closing {
		session.closing = utils.BoolPointer(true)
		user := "unknown"
		if session.user != nil {
			user = *session.user
		}
		log.Infof("closing websocket for %s session %s", user, *session.sessionID)
		session.ws.Close()
		delete(sessions, *session.sessionID)
	}
}

// common websocket error handling with logging
func (session *session) onError() {
	session.onClose()
}

func (session *session) handleRequest(request *configs.WsMessage, transaction *audit.WSTransaction) {
	if request.Route != nil && request.Type != nil {
		request.Vhost = session.Vhost
		if handleMessageFunc, ok := websocketFuncMap[*request.Route]; ok {
			// reset the idle timer if it's appropriate
			go session.resetIdleTimer(request)

			if strings.EqualFold(*request.Type, configs.GetPagelet) {
				transaction.Recordable = utils.BoolPointer(false)
			}

			// do the rest
			response := &configs.WsMessage{
				Route:        request.Route,
				Type:         request.Type,
				Component:    request.Component,
				SubComponent: request.SubComponent,
			}

			handleMessageFunc(request, response)
			if err := session.webSocketSend(response); err != nil {
				log.Error(err)
				session.onError()
			}
			return
		}

		routeNotFound := fmt.Sprintf("requested route: %s is not found", *request.Route)
		log.Error(routeNotFound)
		if err := session.webSocketSend(requestErrorHelper(&routeNotFound, request)); err != nil {
			log.Error(err)
			session.onError()
		}
		return
	}
	session.onError()
}

func (session *session) resetIdleTimer(request *configs.WsMessage) {
	if session != nil && request != nil && request.Component != nil {
		if !strings.EqualFold(*request.Component, configs.Keepalive) {
			session.lastUsed = utils.CurrentTimeInMillis()
		}
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

// webSocketSend allows for the sender to be thread safe, we cannot write to the websocket at the same time
func (session *session) webSocketSend(response *configs.WsMessage) error {
	session.writeMutex.Lock()
	defer session.writeMutex.Unlock()
	response.Timestamp = utils.CurrentTimeInMillis()
	response.SessionID = session.sessionID

	return session.ws.WriteJSON(response)
}

// sendInit is generated on the onOpen event and sends the information the UI needs to startup
func (session *session) sendInit() {
	wsClient := configs.WebsocketClient
	initialize := configs.Initialize
	if err := session.webSocketSend(&configs.WsMessage{
		Route: &wsClient,
		Type:  &initialize,
	}); err != nil {
		log.Errorf("error receiving / sending init to session %s: %s\n", *session.sessionID, err)
	}
}

func transactionHelper(transaction *audit.WSTransaction, recordable bool) {
	if transaction != nil {
		go func() {
			if err := transaction.Complete(recordable); err != nil {
				log.Trace(err)
			}
		}()
	}
}

// The UI will either request authentication or validation, handle those situations here
func handleWebSocketAuth(request, response *configs.WsMessage) (*string, error) {
	defer log.FunctionTimer()()

	err := errors.New("not authenticated").Error()
	response.Error = &err
	denied := configs.Denied
	response.SubComponent = &denied

	if request.Component != nil && strings.EqualFold(*request.Component, configs.Authenticate) && request.SessionID != nil &&
		request.Authentication != nil && request.Authentication.ID != nil && request.Authentication.Password != nil {
		if user, ok := configs.GrowSTLGo.APIUsers[*request.Authentication.ID]; ok && user.Active != nil && *user.Active && request.Vhost != nil {
			if !slices.Contains(user.Vhosts, *request.Vhost) {
				go audit.RecordLogin(user.Authentication.ID, "WebSocket", false)
				return user.Authentication.ID, fmt.Errorf("user %s not authorized for vhost %s", *request.Authentication.ID, *request.Vhost)
			}

			if err := user.Authentication.ValidateAuthentication(request.Authentication.Password); err != nil {
				go audit.RecordLogin(user.Authentication.ID, "WebSocket", false)
				return user.Authentication.ID, fmt.Errorf("bad password attempt for user %s.  Error: %s", *request.Authentication.ID, err)
			}

			token, validTill, err := createJWTToken(request.Authentication.ID, request.SessionID)
			if err != nil || token == nil {
				go audit.RecordLogin(user.Authentication.ID, "WebSocket", false)
				return user.Authentication.ID, err
			}

			if session, ok := sessions[*request.SessionID]; ok {
				session.jwt = token
				session.isAdmin = user.Admin
			}

			approved := configs.Approved
			response.SubComponent = &approved
			response.ValidTill = validTill
			response.Token = token
			response.IsAdmin = user.Admin
			response.Error = nil

			go audit.RecordLogin(user.Authentication.ID, "WebSocket", true)
			return user.Authentication.ID, nil
		}
		return request.Authentication.ID, errors.New("user not found")
	}
	return nil, errors.New("unable to process authentication request")
}

func handleMessage(request, response *configs.WsMessage) {
	if request.Type != nil && request.SessionID != nil {
		switch *request.Type {
		case configs.GetPagelet:
			getPagelet(request, response)
		case configs.Keepalive:
			log.Trace(fmt.Sprintf("keepalive received for session %s", *request.SessionID))
		case configs.Auth:
			if session, ok := sessions[*request.SessionID]; ok {
				user, err := handleWebSocketAuth(request, response)
				session.user = user
				if err != nil {
					log.Error(err)
					if err := session.webSocketSend(response); err != nil {
						log.Error(err)
						session.onError()
					}
					session.onClose()
				}
			}
		default:
			err := fmt.Sprintf("component %s not implemented", *request.Component)
			log.Error(err)
			response.Error = &err
		}
		return
	}
	err := fmt.Errorf("bad request").Error()
	response.Error = &err
}

func getPagelet(request, response *configs.WsMessage) {
	defer log.FunctionTimer()()
	err := errors.New(configs.NotFoundError).Error()
	response.Error = &err

	if request.Component != nil {
		// ignore pages non admins shouldn't see
		if strings.EqualFold(*request.Component, "admin") && (request.IsAdmin == nil || !*request.IsAdmin) {
			log.Errorf("attempt to access the admin page by non admin user")
			return
		}
		// everyone else is free to move about the cabin
		if request.SessionID != nil {
			if session, ok := sessions[*request.SessionID]; ok && session.Vhost != nil {
				log.Info(*session.Vhost)
				if webRoot, ok := configs.GrowSTLGo.WebService.Vhosts[*session.Vhost]; ok && webRoot != nil {
					fileName := filepath.Join(*webRoot, "pagelets", fmt.Sprintf("%s.html", *request.Component))
					_, err := os.Stat(fileName)
					if err == nil {
						bytes, err := os.ReadFile(fileName)
						if err != nil {
							return
						}
						response.Data = string(bytes)
						response.Error = nil
					}
				}
			}
		}
	}
}

// NotifyAll will broadcast a user event to all websocket clients currently attached
func NotifyAll(sessionID *string, message *configs.WsMessage) {
	if message != nil {
		for _, session := range sessions {
			// don't send a notification to the user that requested the action
			if session != nil && session.sessionID != nil && sessionID != nil && strings.EqualFold(*sessionID, *session.sessionID) {
				continue
			}
			if err := session.webSocketSend(message); err != nil {
				log.Error(err)
			}
		}
	}
}

func idleHandsTester() {
	time.Sleep(time.Duration(60-time.Now().Local().Second()) * time.Second)
	for range time.NewTicker(10 * time.Second).C {
		for _, session := range sessions {
			if session != nil && session.lastUsed != nil {
				// 10 minute timeout
				if (time.Now().UnixMilli() - *session.lastUsed) > 600000 {
					session.onClose()
				}
				// idle abandoned connections are disconnected at 30 seconds
				if (time.Now().UnixMilli()-*session.lastUsed) > 30000 && session.user == nil {
					session.onClose()
				}
			}
		}
	}
}

// formats an error response in the way that we're expecting on the UI
func requestErrorHelper(err *string, request *configs.WsMessage) *configs.WsMessage {
	if err != nil {
		return &configs.WsMessage{
			Route:     request.Route,
			Type:      request.Type,
			Component: request.Component,
			Error:     err,
		}
	}
	return nil
}

// newSession generates a new session
func newSession(requestHost, vhost *string, ws *websocket.Conn) *session {
	id := uuid.New().String()
	session := &session{
		requestHost: requestHost,
		sessionID:   &id,
		ws:          ws,
		lastUsed:    utils.CurrentTimeInMillis(),
		Vhost:       vhost,
	}

	// keep track of the session
	sessions[id] = session

	// send the init message to the client
	go session.sendInit()

	return session
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

// Shutdown is called when the system is exiting to cleanly close all the current connections
func Shutdown() {
	for _, session := range sessions {
		session.onClose()
	}
}
