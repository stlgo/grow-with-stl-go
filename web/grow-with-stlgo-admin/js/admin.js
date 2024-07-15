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

        this.route = this.constructor.name.toLowerCase();
        this.ws.registerHandlers(this.route, this);

        document.addEventListener('admin', () => {
            this.ws.sendMessage({
                route: this.route,
                type: 'pageLoad',
                component: 'pageLoad',
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
                route: this.route,
                type: 'addUser',
                component: userIDInput.value,
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
                route: this.route,
                type: 'updateUser',
                component: userID,
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
            route: this.route,
            type: 'generatePassword',
            component: uesrType,
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
            'Admin',
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


            let activeChk = document.createElement('input');
            activeChk.type = 'checkbox';
            activeChk.id = `${userID}_active_chk`;
            activeChk.name = userID;
            activeChk.classList.add('user-active-chk');
            activeChk.checked = data[userID].active;
            tr.insertCell(-1).appendChild(activeChk);

            let adminChk = document.createElement('input');
            adminChk.type = 'checkbox';
            adminChk.id = `${userID}_admin_chk`;
            adminChk.name = userID;
            adminChk.classList.add('user-admin-chk');
            adminChk.checked = data[userID].admin;
            tr.insertCell(-1).appendChild(adminChk);

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
        this.bindAdminToggle();
        this.bindSlideOut(userTable);
        this.bindRemoveUser(userTable);
    }

    bindSlideOut(table) {
        table.on('click', '.details-control', (e) => {
            let tr = e.target.closest('tr');
            let row = table.row(tr);
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
                        default:
                            this.log.info(`cannot decision ${tagElement.id}`);
                        }

                        const name = tagElement.name;
                        if (name !== undefined && name.length > 0) {
                            tagElement.name = `${userID}-${name}`;
                        }
                    });
                });

                row.child(clone).show();
                details.innerHTML = '<i class="fa fa-times-circle fa-lg details-control" style="color:maroon"></i>';
            }
        });
    }

    bindActiveToggle() {
        document.querySelectorAll('.user-active-chk').forEach((chkBox) => {
            chkBox.onclick = (event) => {
                event.preventDefault();
                const userID = chkBox.name;
                let active = false;
                if (chkBox.checked) {
                    active = true;
                }
                let confirm = window.confirm(`Confirm toggle of user "${userID}" active to ${active}`); // eslint-disable-line no-alert
                if (confirm === true) {
                    this.ws.sendMessage({
                        route: this.route,
                        type: 'updateActive',
                        component: userID,
                        subComponent: `${active}`
                    });
                }
            };
        });
    }

    bindAdminToggle() {
        document.querySelectorAll('.user-admin-chk').forEach((chkBox) => {
            chkBox.onclick = (event) => {
                event.preventDefault();
                const userID = chkBox.name;
                let admin = false;
                if (chkBox.checked) {
                    admin = true;
                }
                let confirm = window.confirm(`Confirm toggle of user "${userID}" admin to ${admin}`); // eslint-disable-line no-alert
                if (confirm === true) {
                    this.ws.sendMessage({
                        route: this.route,
                        type: 'updateAdmin',
                        component: userID,
                        subComponent: `${admin}`
                    });
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
                        route: this.route,
                        type: 'removeUser',
                        component: userID,
                    });
                }
            };
        });
    }

    showVhosts(vhosts) {
        const select = document.createElement('select');
        vhosts.forEach((vhost) => {
            select.options.add(new Option(vhost, vhost));
        });
        document.getElementById('VhostDiv').appendChild(select);
    }

    handleMessage(json) {
        if (Object.prototype.hasOwnProperty.call(json, 'error')) {
            this.log.error(json.error);
            alert(json.error); // eslint-disable-line no-alert
        } else {
            switch(json.type) {
            case 'addUser':
            case 'updateUser':
            case 'updateActive':
            case 'updateAdmin':
            case 'removeUser':
            case 'pageLoad':
                this.showUsers(json.data.users);
                this.showVhosts(json.data.vhosts);
                this.bindButtons();
                break;
            case 'generatePassword':
                this.populatePassword(json.component, json.data);
                break;
            default:
                this.log.error(`Cannot handle type '${json.type}' for ${this.route}`);
                break;
            }
        }
    }
}

export { Admin };
