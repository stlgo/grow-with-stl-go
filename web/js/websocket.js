/*
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
*/

import { Log } from './log.js';

export class WebSocketClient {
    constructor() {
        this.ws = null;
        this.timeout = null;
        this.token = null;
        this.refreshToken = null;
        this.sessionID = null;
        this.validUser = null;
        this.type = this.constructor.name.toLocaleLowerCase();

        this.functionMap = {};

        document.addEventListener('beforeunload', () => {
            this.unloadHelper();
        });

        document.addEventListener('unload', () => {
            this.unloadHelper();
        });

        this.register();
    }

    unloadHelper() {
        if (this.ws !== null && this.ws.readyState !== WebSocket.CLOSED) {
            this.ws.close();
        }
    }

    registerHandlers(type, obj) {
        this.functionMap[type] = obj;
    }

    register() {
        if (this.ws !== null) {
            this.ws.close();
            this.ws = null;
        }

        this.ws = new WebSocket(`wss://${window.location.host}/ws/v1.0.0`);

        this.registerHandlers(this.type, this);

        this.ws.onmessage = (event) => {
            this.handleMessages(event);
        };

        this.ws.onerror = (event) => {
            Log.error(`Web Socket recieved an error: ${event.code}`);
            this.wsClose(event.code);
        };

        this.ws.onopen = () => {
            this.wsOpen();
        };

        this.ws.onclose = (event) => {
            this.wsClose(event.code);
            document.dispatchEvent(new CustomEvent('WebSocketClosed'));
        };
    }

    handleMessages(message) {
        const json = JSON.parse(message.data);
        const type = json.type;
        if (Object.prototype.hasOwnProperty.call(this.functionMap, type)) {
            this.functionMap[type].handleMessage(json);
        } else {
            Log.error(`Received invalid message type ${type}`);
        }
    }

    handleMessage(json) {
        switch(json.component) {
        case 'auth':
            this.handleAuth(json);
            break;
        case 'initialize':
            this.sessionID = json.sessionID;
            document.dispatchEvent(new CustomEvent('WebSocketEstablished', {
                detail: this
            }));
            break;
        case 'keepalive':
            Log.trace(`Keepalive received: ${JSON.stringify(json)}`);
            break;
        default:
            this.authDenied();
            break;
        }
    }

    handleAuth(json) {
        switch(json.subComponent) {
        case 'approved':
            this.token = json.token;
            this.authAllowed(json);
            this.keepAlive();
            break;
        case 'refresh':
            this.refreshToken = json.refreshToken;
            break;
        case 'denied':
            Log.error('Authentication denied');
            this.authDenied();
            break;
        default:
            this.authDenied();
            break;
        }
        document.dispatchEvent(new Event('AuthComplete'));
    }

    wsOpen() {
        Log.info('WebSocket established');
    }

    wsClose(code) {
        switch (code) {
        case 1000:
            Log.info('Web Socket Closed: Normal closure: ', code);
            break;
        case 1001:
            Log.info('Web Socket Closed: An endpoint is "going away", such as a server going down or a browser having navigated away from a page:', code);
            break;
        case 1002:
            Log.info('Web Socket Closed: terminating the connection due to a protocol error: ', code);
            break;
        case 1003:
            Log.info('Web Socket Closed: terminating the connection because it has received a type of data it cannot accept: ', code);
            break;
        case 1004:
            Log.info('Web Socket Closed: Reserved. The specific meaning might be defined in the future: ', code);
            break;
        case 1005:
            Log.info('Web Socket Closed: No status code was actually present: ', code);
            break;
        case 1006:
            Log.info('Web Socket Closed: The connection was closed abnormally: ', code);
            break;
        case 1007:
            Log.info('Web Socket Closed: terminating the connection because it has received data within a message that was not consistent with the type of the message: ', code);
            break;
        case 1008:
            Log.info('Web Socket Closed: terminating the connection because it has received a message that "violates its policy": ', code);
            break;
        case 1009:
            Log.info('Web Socket Closed: terminating the connection because it has received a message that is too big for it to process: ', code);
            break;
        case 1010:
            Log.info('Web Socket Closed: client is terminating the connection because it has expected the server to negotiate one or more extension, but the server didn\'t return them in the response message of the WebSocket handshake: ', code);
            break;
        case 1011:
            Log.info('Web Socket Closed: server is terminating the connection because it encountered an unexpected condition that prevented it from fulfilling the request: ', code);
            break;
        case 1015:
            Log.info('Web Socket Closed: closed due to a failure to perform a TLS handshake (e.g., the server certificate can\'t be verified): ', code);
            break;
        default:
            Log.info('Web Socket Closed: unknown error code: ', code);
            break;
        }

        this.ws = null;
        this.token = null;
        this.refreshToken = null;
    }

    authAllowed(json) {
        Log.info(json);
    }

    authDenied() {
        this.ws.close();
        Log.error('Auth denied');
    }

    keepAlive() {
        if (this.ws !== null) {
            if (this.ws.readyState !== WebSocket.CLOSED && this.ws.readyState !== WebSocket.CONNECTING && this.token !== null) {
                // clear previous timeout
                window.clearTimeout(this.timeout);
                window.clearInterval(this.timeout);
                const json = {
                    type: this.type,
                    component: 'keepalive'
                };
                this.sendMessage(json);
            }
            this.timeout = window.setTimeout(this.keepAlive.bind(this), 60000);
        }
    }

    waitThenSendMessage(json) {
        document.addEventListener('AuthComplete', () => {
            this.sendMessage(json);
        }, {
            once: true
        });
    }

    sendMessage(json) {
        if (this.ws === null || this.ws.readyState === WebSocket.CLOSED) {
            this.register();
            this.waitThenSendMessage(json);
        } else if (this.ws.readyState === WebSocket.CONNECTING) {
            this.waitThenSendMessage(json);
        } else {
            try {
                if (this.token !== null) {
                    json.token = this.token;
                }
                if (this.refreshToken !== null) {
                    json.refreshToken = this.refreshToken;
                }
                json.sessionID = this.sessionID;
                json.timestamp = new Date().getTime();
                this.ws.send(JSON.stringify(json));
            } catch (err) {
                Log.error(err);
                window.setTimeout(this.sendMessage(json).bind(this), 250);
            }
        }
    }
}
