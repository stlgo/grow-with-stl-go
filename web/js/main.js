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

import { Admin } from './admin.js';
import { Seeds } from './seeds.js';
import { Log } from './log.js';
import { WebSocketClient } from './websocket.js';

class GrowWithSTLGO {
    constructor() {
        this.log = new Log();
        this.ws = null;
        this.type = this.constructor.name.toLowerCase();

        document.addEventListener('WebSocketClosed', () => {
            this.ws.displayHelper([ 'RouterDiv', 'NavbarDiv', 'AdminNavLink' ], 'none');
            this.ws.displayHelper([ 'LoginDiv' ], '');
        });

        document.querySelectorAll('.loginInput').forEach((input) => {
            input.addEventListener('keypress', (event) => {
                // If the user presses the "Enter" key on the keyboard
                if (event.key === 'Enter') {
                    // Cancel the default action, if needed
                    event.preventDefault();
                    // Trigger the button element with a click
                    document.getElementById('LoginButton').click();
                }
            });
        });
    }

    load(pageType) {
        this.ws.getPagelet(pageType);
    }

    login() {
        let id = document.getElementById('IDInput').value;
        let password = document.getElementById('PasswordInput').value;
        if (id.length > 0 && password.length > 0) {
            if (this.ws === null) {
                this.ws = new WebSocketClient(this.log);
                const admin = new Admin(this.ws, this.log);
                const seeds = new Seeds(this.ws, this.log);
            }

            this.ws.login(id, password);
            document.addEventListener('AuthComplete', () => {
                this.ws.displayHelper([ 'LoginDiv' ], 'none');
                let sessionData = JSON.parse(window.sessionStorage.getItem('grow-with-stlgo'));
                if (sessionData && Object.prototype.hasOwnProperty.call(sessionData, 'pageType')) {
                    this.ws.getPagelet(sessionData.pageType);
                } else {
                    this.ws.getPagelet('home');
                }
            }, {
                once: true
            });
            document.getElementById('IDInput').value = '';
            document.getElementById('PasswordInput').value = '';
        }
    }

    // load(type) {
    //     this.ws.getPagelet()
    //     let div = document.getElementById('RouterDiv');
    //     div.innerHTML = `<h2>Loading ${type} please wait...`;
    //     this.displayHelper([ 'RouterDiv', 'NavbarDiv' ], '');

    //     let request = new XMLHttpRequest();
    //     request.open('GET', uri, true);
    //     request.onerror = (e) => {
    //         this.log.error(`Unable to connect to the backend ${e.target.status}`);
    //     };
    //     request.onreadystatechange = () => {
    //         if (request.readyState === XMLHttpRequest.DONE && request.status === 200) {
    //             div.innerHTML = request.responseText;
    //             window.history.replaceState(null, null, `/${type}`);
    //         }
    //     };
    //     request.send();
    // }
}

export { GrowWithSTLGO };
