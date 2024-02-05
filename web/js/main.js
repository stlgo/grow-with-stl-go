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
import { WebSocketClient } from './websocket.js';

class GrowWithSTLGO {
    constructor() {
        this.log = new Log();
        this.ws = new WebSocketClient(this.log);
        this.type = this.constructor.name.toLowerCase();

        document.addEventListener('WebSocketClosed', () => {
            this.displayHelper([ 'RouterDiv', 'NavbarDiv' ], 'none');
            this.displayHelper([ 'LoginDiv' ], '');
        });
    }

    displayHelper(elements, display) {
        elements.forEach((elementID) => {
            let element = document.getElementById(elementID);
            if (typeof element !== 'undefined' && element !== null) {
                element.style.display = display;
            }
        });
    }

    login(id, password) {
        this.ws.login(id, password);
        document.addEventListener('AuthComplete', () => {
            this.displayHelper([ 'LoginDiv' ], 'none');
            this.load('home', '/_home/index.html');
        });
        document.getElementById('IDInput').value = '';
        document.getElementById('PasswordInput').value = '';
    }

    load(type, uri) {
        let div = document.getElementById('RouterDiv');
        div.innerHTML = `<h2>Loading ${type} please wait...`;
        this.displayHelper([ 'RouterDiv', 'NavbarDiv' ], '');

        let request = new XMLHttpRequest();
        request.open('GET', uri, true);
        request.onerror = (e) => {
            this.log.error(`Unable to connect to the backend ${e.target.status}`);
        };
        request.onreadystatechange = () => {
            if (request.readyState === XMLHttpRequest.DONE && request.status === 200) {
                div.innerHTML = request.responseText;
                window.history.replaceState(null, null, `/${type}`);
            }
        };
        request.send();
    }
}

export { GrowWithSTLGO };
