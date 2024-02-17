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


class Admin {
    constructor(ws, log) {
        this.ws = ws;
        this.log = log;

        this.type = this.constructor.name.toLowerCase();
        this.ws.registerHandlers(this.type, this);

        document.addEventListener('admin', () => {
            this.ws.sendMessage({
                type: this.type,
                component: 'pageLoad',
                subComponent: 'pageLoad',
            });
        });
    }

    handleMessage(json) {
        if (Object.prototype.hasOwnProperty.call(json, 'error')) {
            this.log.error(json.error);
        } else {
            switch(json.component) {
            case 'addUser':
            case 'pageLoad':
                this.log.info(json.data);
                break;
            default:
                this.log.error(`Cannot handle component '${json.component}' for ${this.type}`);
                break;
            }
        }
    }
}

export { Admin };
