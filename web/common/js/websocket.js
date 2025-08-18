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


class WebSocketClient {
    constructor(log) {
        this.ws = null;
        this.timeout = null;
        this.token = null;
        this.refreshToken = null;
        this.sessionID = null;
        this.validUser = null;
        this.route = this.constructor.name.toLocaleLowerCase();

        this.log = log;

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

        this.registerHandlers(this.route, this);

        this.ws.onmessage = (event) => {
            this.handleMessages(event);
        };

        this.ws.onerror = (event) => {
            this.log.error(`Web Socket received an error: ${event.code}`);
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
        const route = json.route;
        if (Object.prototype.hasOwnProperty.call(this.functionMap, route)) {
            this.functionMap[route].handleMessage(json);
        } else {
            this.log.error(`Received invalid message route ${route}`);
        }
    }

    handleMessage(json) {
        switch(json.type) {
        case 'auth':
            this.handleAuth(json);
            break;
        case 'getPagelet':
            if (Object.prototype.hasOwnProperty.call(json, 'error')) {
                this.log.error(json.error);
                window.history.replaceState(null, null, '/');
                window.sessionStorage.setItem('grow-with-stlgo', JSON.stringify({ timestamp: new Date().getTime(), pageType: 'home' }));
            } else {
                document.getElementById('RouterDiv').innerHTML = json.data;
                this.displayHelper([ 'RouterDiv', 'NavbarDiv' ], '');
                window.history.replaceState(null, null, `/${json.component}`);
                window.sessionStorage.setItem('grow-with-stlgo', JSON.stringify({ timestamp: new Date().getTime(), pageType: location.pathname.substring(1) }));
                document.dispatchEvent(new CustomEvent(`${json.component}`));
            }
            break;
        case 'initialize':
            this.sessionID = json.sessionID;
            document.dispatchEvent(new CustomEvent('WebSocketEstablished', {
                detail: this
            }));
            break;
        case 'keepalive':
            this.log.trace(`Keepalive received: ${JSON.stringify(json)}`);
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
            this.keepAlive();
            if (Object.prototype.hasOwnProperty.call(json, 'isAdmin') && json.isAdmin !== undefined && json.isAdmin) {
                this.displayHelper([ 'AdminNavLink' ], '');
            }
            document.dispatchEvent(new Event('AuthComplete'));
            break;
        case 'refresh':
            this.refreshToken = json.refreshToken;
            break;
        case 'denied':
            this.log.error('Authentication denied');
            this.authDenied();
            break;
        default:
            this.authDenied();
            break;
        }
    }

    wsOpen() {
        this.log.info('WebSocket established');
    }

    wsClose(code) {
        switch (code) {
        case 1000:
            this.log.info('Web Socket Closed: Normal closure: ', code);
            break;
        case 1001:
            this.log.info('Web Socket Closed: An endpoint is "going away", such as a server going down or a browser having navigated away from a page:', code);
            break;
        case 1002:
            this.log.info('Web Socket Closed: terminating the connection due to a protocol error: ', code);
            break;
        case 1003:
            this.log.info('Web Socket Closed: terminating the connection because it has received a type of data it cannot accept: ', code);
            break;
        case 1004:
            this.log.info('Web Socket Closed: Reserved. The specific meaning might be defined in the future: ', code);
            break;
        case 1005:
            this.log.info('Web Socket Closed: No status code was actually present: ', code);
            break;
        case 1006:
            this.log.info('Web Socket Closed: The connection was closed abnormally: ', code);
            break;
        case 1007:
            this.log.info('Web Socket Closed: terminating the connection because it has received data within a message that was not consistent with the type of the message: ', code);
            break;
        case 1008:
            this.log.info('Web Socket Closed: terminating the connection because it has received a message that "violates its policy": ', code);
            break;
        case 1009:
            this.log.info('Web Socket Closed: terminating the connection because it has received a message that is too big for it to process: ', code);
            break;
        case 1010:
            this.log.info('Web Socket Closed: client is terminating the connection because it has expected the server to negotiate one or more extension, but the server didn\'t return them in the response message of the WebSocket handshake: ', code);
            break;
        case 1011:
            this.log.info('Web Socket Closed: server is terminating the connection because it encountered an unexpected condition that prevented it from fulfilling the request: ', code);
            break;
        case 1015:
            this.log.info('Web Socket Closed: closed due to a failure to perform a TLS handshake (e.g., the server certificate can\'t be verified): ', code);
            break;
        default:
            this.log.info('Web Socket Closed: unknown error code: ', code);
            break;
        }

        this.ws = null;
        this.token = null;
        this.refreshToken = null;
    }

    login(id, password) {
        this.sendMessage({
            route: this.route,
            type: 'auth',
            component: 'authenticate',
            authentication: {
                id: id,
                password: password
            }
        });
    }

    getPagelet(page) {
        this.sendMessage({
            route: this.route,
            type: 'getPagelet',
            component: page
        });
    }

    showSnackbarMessage(message) {
        let div = document.getElementById('SnackbarDiv');
        div.innerHTML = `${message}<br><span class="close" id="SnackClose">&times;</span>`;
        div.classList.add('show-bar');

        document.getElementById('SnackClose').onclick = () => {
            div.className = div.className.replace('show-bar', '');
        };

        setTimeout(() => {
            div.className = div.className.replace('show-bar', '');
        }, 10000);
    }

    displayHelper(elements, display) {
        elements.forEach((elementID) => {
            let element = document.getElementById(elementID);
            if (typeof element !== 'undefined' && element !== null) {
                element.style.display = display;
            }
        });
    }

    authDenied() {
        this.ws.close();
        this.log.error('Auth denied');
    }

    keepAlive() {
        if (this.ws !== null) {
            if (this.ws.readyState !== WebSocket.CLOSED && this.ws.readyState !== WebSocket.CONNECTING && this.token !== null) {
                // clear previous timeout
                window.clearTimeout(this.timeout);
                window.clearInterval(this.timeout);
                const json = {
                    route: this.route,
                    type: 'keepalive',
                    component: 'active'
                };
                this.sendMessage(json);
            }
            this.timeout = window.setTimeout(this.keepAlive.bind(this), 60000);
        }
    }

    waitThenSendMessage(json) {
        document.addEventListener('WebSocketEstablished', () => {
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
                this.log.error(err);
                window.setTimeout(this.sendMessage(json).bind(this), 250);
            }
        }
    }
}

export { WebSocketClient };
