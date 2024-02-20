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

    bindButtons() {
        let element = document.getElementById('UserIDInput');
        element.onblur = () => {
            if (element.validity.valid && document.getElementById('PasswordInput').value === '') {
                this.generatePassword('newUser');
            }
        };
        document.getElementById('RefreshPassword').onclick = () => {
            this.generatePassword('newUser');
        };
        document.getElementById('AddUserButton').onclick = () => {
            this.addUser();
        };
    }

    addUser() {
        let userIDInput = document.getElementById('UserIDInput');
        let passwordInput = document.getElementById('GeneratedPasswd');
        if (userIDInput.validity.valid && passwordInput.value.length >= 10) {
            this.ws.sendMessage({
                type: this.type,
                component: 'addUser',
                subComponent: userIDInput.value,
                data: {
                    id: userIDInput.value,
                    password: passwordInput.value
                }
            });
            userIDInput.value = '';
            passwordInput.value = '';
        }
    }

    updateUser(userID) {
        let passwordInput = document.getElementById(`${userID}-GeneratedPasswd`);
        if (passwordInput.value.length >= 10) {
            this.ws.sendMessage({
                type: this.type,
                component: 'updateUser',
                subComponent: userID,
                data: {
                    id: userID,
                    password: passwordInput.value
                }
            });
            passwordInput.value = '';
        }
    }

    generatePassword(uesrType) {
        this.ws.sendMessage({
            type: this.type,
            component: 'generatePassword',
            subComponent: uesrType,
        });
    }

    populatePassword(userType, data) {
        if (Object.prototype.hasOwnProperty.call(data, 'password')) {
            switch(userType) {
            case 'newUser':
                document.getElementById('GeneratedPasswd').value = data.password;
                break;
            default:
                document.getElementById(`${userType}-GeneratedPasswd`).value = data.password;
                break;
            }
            navigator.clipboard.writeText(data.password);
            alert('Password has been copied to your clipboard'); // eslint-disable-line no-alert
        }
    }

    showUsers(data) {
        const headers = [
            '',
            'ID',
            'Last Login',
            'Enabled',
            ''
        ];
        const table = document.createElement('table');
        table.setAttribute('width', '99%');
        table.id = 'UsersTable';
        table.classList.add('display', 'responsive');
        let th = table.createTHead();
        let tr = th.insertRow(-1);
        headers.forEach((header) => {
            tr.insertCell(-1).innerHTML = header;
        });
        let tb = table.createTBody();
        Object.keys(data).forEach((userID) => {
            tr = tb.insertRow(-1);
            tr.insertCell(-1).innerHTML = `<div id="${userID}-details"><i class="fa fa-plus-circle fa-lg details-control" style="color:green"></i></div>`;
            tr.insertCell(-1).innerHTML = userID;
            if (data[userID].lastLogin === null) {
                tr.insertCell(-1).innerHTML = 'Never Logged In';
            } else {
                tr.insertCell(-1).innerHTML = new Date(data[userID].lastLogin).toLocaleString('en-us', { timeZoneName: 'short' });
            }


            let checkbox = document.createElement('input');
            checkbox.type = 'checkbox';
            checkbox.id = `${userID}_active_chk`;
            checkbox.name = userID;
            checkbox.classList.add('user-active-chk');
            checkbox.checked = data[userID].active;
            tr.insertCell(-1).appendChild(checkbox);

            tr.insertCell(-1).innerHTML = '<i class="fa fa-trash fa-lg remove-user" style="color:maroon"></i>';
        });
        let div = document.getElementById('CurrentUsersDiv');
        div.innerHTML = '';
        div.appendChild(table);

        let userTable = $('#UsersTable').DataTable({
            deferRender: true,
            orderClasses: false,
            columnDefs: [ {
                targets: 0,
                orderable: false,
            }, {
                targets: 4,
                orderable: false,
            } ]
        });

        this.bindActiveToggle();
        this.bindSlideOut(userTable);
        this.bindRemoveUser(userTable);
    }

    bindSlideOut(table) {
        document.querySelectorAll('.details-control').forEach((slideButton) => {
            slideButton.onclick = () => {
                const row = table.row(slideButton.parentNode.parentNode);
                const userID = row.data()[1];
                const details = document.getElementById(`${userID}-details`);

                if (details.innerHTML.includes('color:maroon')) {
                    row.child.hide();
                    details.innerHTML = '<i class="fa fa-plus-circle fa-lg details-control" style="color:green"></i>';
                } else {
                    const clone = document.getElementById('NewUserFieldset').cloneNode(true);
                    clone.id = `${userID}-details`;

                    const elementTypes = [ 'button', 'input', 'fieldset', 'legend', 'table' ];
                    elementTypes.forEach((element) => {
                        Array.from(clone.getElementsByTagName(element)).forEach((tagElement) => {
                            const id = tagElement.id;
                            if (id.length > 0) {
                                tagElement.id = `${userID}-${id}`;
                            }

                            switch (tagElement.id) {
                            case `${userID}-UserIDInput`:
                                tagElement.value = userID;
                                tagElement.disabled = true;
                                break;
                            case `${userID}-AddUserButton`:
                                tagElement.innerHTML = 'Update User';
                                tagElement.onclick = () => {
                                    this.updateUser(userID);
                                };
                                break;
                            case `${userID}-RefreshPassword`:
                                tagElement.value = '';
                                tagElement.onclick = () => {
                                    this.generatePassword(userID);
                                };
                                break;
                            }

                            const name = tagElement.name;
                            if (name !== undefined && name.length > 0) {
                                tagElement.name = `${userID}-${name}`;
                            }
                        });
                    });

                    row.child(clone).show();
                    details.innerHTML = '<i class="fa fa-times-circle fa-lg details-control" style="color:maroon"></i>';

                    // rebind click
                    this.bindSlideOut(table);
                }
            };
        });
    }

    bindActiveToggle() {
        document.querySelectorAll('.user-active-chk').forEach((chkBox) => {
            chkBox.onclick = () => {
                const userID = chkBox.name;
                let active = false;
                if (chkBox.checked) {
                    active = true;
                }
                if (window.confirm(`Confirm toggle of user "${userID}"`)) { // eslint-disable-line no-alert
                    this.ws.sendMessage({
                        type: this.type,
                        component: 'updateActive',
                        subComponent: userID,
                        data: {
                            enabled: active
                        }
                    });
                    this.log.info(active);
                }
            };
        });
    }

    bindRemoveUser(table) {
        document.querySelectorAll('.remove-user').forEach((removeButton) => {
            removeButton.onclick = () => {
                const row = table.row(removeButton.parentNode.parentNode);
                const userID = row.data()[1];
                if (window.confirm(`Confirm removal of user "${userID}"`)) { // eslint-disable-line no-alert
                    this.ws.sendMessage({
                        type: this.type,
                        component: 'removeUser',
                        subComponent: userID,
                    });
                }
            };
        });
    }

    handleMessage(json) {
        if (Object.prototype.hasOwnProperty.call(json, 'error')) {
            this.log.error(json.error);
            alert(json.error); // eslint-disable-line no-alert
        } else {
            switch(json.component) {
            case 'addUser':
            case 'updateUser':
            case 'updateActive':
            case 'removeUser':
            case 'pageLoad':
                this.showUsers(json.data);
                this.bindButtons();
                break;
            case 'generatePassword':
                this.populatePassword(json.subComponent, json.data);
                break;
            default:
                this.log.error(`Cannot handle component '${json.component}' for ${this.type}`);
                break;
            }
        }
    }
}

export { Admin };
