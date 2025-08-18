/* eslint-disable no-undef */
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
        this.typeaheadData = null;

        this.locationTypeAhead = null;

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

        const vhosts = [];
        for (const row of document.getElementById('VhostTable').rows) {
            const cell = row.cells[1].children[0];
            if (cell.checked) {
                vhosts.push(cell.value);
            }
        }

        if (userIDInput.validity.valid && passwordInput.value.length >= 10) {
            this.ws.sendMessage({
                route: this.route,
                type: 'addUser',
                component: userIDInput.value,
                data: {
                    id: userIDInput.value,
                    password: passwordInput.value,
                    vhosts: vhosts
                }
            });
            userIDInput.value = '';
            passwordInput.value = '';
        }
    }

    updateUser(userID) {
        let passwordInput = document.getElementById(`${userID}-GeneratedPasswd`);
        let password = null;
        if (passwordInput.value.length > 0) {
            password = passwordInput.value;
        }
        let location = null;
        if (document.getElementById(`${userID}-LocationInput`).value.length > 0) {
            location = document.getElementById(`${userID}-LocationInput`).value;
        }

        const vhosts = [];
        for (const row of document.getElementById(`${userID}-VhostTable`).rows) {
            const cell = row.cells[1].children[0];
            if (cell.checked) {
                vhosts.push(cell.value);
            }
        }

        this.ws.sendMessage({
            route: this.route,
            type: 'updateUser',
            component: userID,
            data: {
                id: userID,
                location: location,
                password: password,
                vhosts: vhosts
            }
        });
        passwordInput.value = '';
    }

    generatePassword(userType) {
        this.ws.sendMessage({
            route: this.route,
            type: 'generatePassword',
            component: userType,
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
            this.ws.showSnackbarMessage('Password has been copied to your clipboard');
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
                targets: [ 0, 4 ],
                orderable: false,
            } ]
        });

        this.bindActiveToggle(userTable);
        this.bindAdminToggle(userTable);
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
                this.ws.sendMessage({
                    route: this.route,
                    type: 'getUserDetails',
                    component: userID,
                });

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
                        case `${userID}-UserLegend`:
                            tagElement.innerHTML = `Update ${userID}`;
                            break;
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
                            this.log.trace(`cannot decision ${tagElement.id}`);
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

    bindActiveToggle(userTable) {
        userTable.on('click', 'td:nth-child(4)', (event) => {
            let row = userTable.row(event.target.closest('tr'));
            const userID = row.data()[1];
            let chkBox = document.getElementById(`${userID}_active_chk`);
            let dialog = document.createElement('dialog');
            dialog.appendChild(document.createTextNode(`Confirm toggle of user "${userID}" active to ${chkBox.checked}`));
            let dialogTable = document.createElement('table');
            dialogTable.createTBody();
            let tr = dialogTable.getElementsByTagName('tbody')[0].insertRow(-1);
            let button = document.createElement('button');
            button.classList.add('btn', 'btn-outline-success');
            button.innerText = 'Confirm';
            button.onclick = () => {
                this.ws.sendMessage({
                    route: this.route,
                    type: 'updateActive',
                    component: userID,
                    subComponent: chkBox.checked ? 'true' : 'false'
                });
                dialog.close();
            };
            tr.insertCell(-1).appendChild(button);
            button = document.createElement('button');
            button.classList.add('btn', 'btn-outline-secondary');
            button.innerText = 'Cancel ';
            button.onclick = () => {
                dialog.close();
            };
            tr.insertCell(-1).appendChild(button);
            dialog.appendChild(dialogTable);
            document.body.appendChild(dialog);
            dialog.showModal();
        });
    }

    bindAdminToggle(userTable) {
        userTable.on('click', 'td:nth-child(5)', (event) => {
            let row = userTable.row(event.target.closest('tr'));
            const userID = row.data()[1];
            let chkBox = document.getElementById(`${userID}_admin_chk`);

            let dialog = document.createElement('dialog');
            dialog.appendChild(document.createTextNode(`Confirm toggle of user "${userID}" admin to ${chkBox.checked}`));
            let dialogTable = document.createElement('table');
            dialogTable.createTBody();
            let tr = dialogTable.getElementsByTagName('tbody')[0].insertRow(-1);
            let button = document.createElement('button');
            button.classList.add('btn', 'btn-outline-success');
            button.innerText = 'Confirm';
            button.onclick = () => {
                this.ws.sendMessage({
                    route: this.route,
                    type: 'updateAdmin',
                    component: userID,
                    subComponent: chkBox.checked ? 'true' : 'false'
                });
                dialog.close();
            };
            tr.insertCell(-1).appendChild(button);
            button = document.createElement('button');
            button.classList.add('btn', 'btn-outline-secondary');
            button.innerText = 'Cancel ';
            button.onclick = () => {
                dialog.close();
            };
            tr.insertCell(-1).appendChild(button);
            dialog.appendChild(dialogTable);
            document.body.appendChild(dialog);
            dialog.showModal();
        });
    }

    bindRemoveUser(userTable) {
        userTable.on('click', 'td:nth-child(6)', (event) => {
            let row = userTable.row(event.target.closest('tr'));
            const userID = row.data()[1];
            let dialog = document.createElement('dialog');
            dialog.appendChild(document.createTextNode(`Confirm removal of user "${userID}"`));
            let dialogTable = document.createElement('table');
            dialogTable.createTBody();
            let tr = dialogTable.getElementsByTagName('tbody')[0].insertRow(-1);
            let button = document.createElement('button');
            button.classList.add('btn', 'btn-outline-success');
            button.innerText = 'Confirm';
            button.onclick = () => {
                this.ws.sendMessage({
                    route: this.route,
                    type: 'removeUser',
                    component: userID,
                });
                dialog.close();
            };
            tr.insertCell(-1).appendChild(button);
            button = document.createElement('button');
            button.classList.add('btn', 'btn-outline-secondary');
            button.innerText = 'Cancel ';
            button.onclick = () => {
                dialog.close();
            };
            tr.insertCell(-1).appendChild(button);
            dialog.appendChild(dialogTable);
            document.body.appendChild(dialog);
            dialog.showModal();
        });
    }

    showVhosts(vhosts) {
        const table = document.createElement('table');
        table.setAttribute('width', '99%');
        table.id = 'VhostTable';
        table.classList.add('display', 'responsive');
        let th = table.createTHead();
        let tr = th.insertRow(-1);
        tr.insertCell(-1).innerHTML = 'Vhost';
        tr.insertCell(-1).innerHTML = 'Enabled';
        let tb = table.createTBody();
        vhosts.forEach((vhost) => {
            tr = tb.insertRow(-1);
            tr.insertCell(-1).innerHTML = vhost;
            let checkbox = document.createElement('input');
            checkbox.type = 'checkbox';
            checkbox.name = vhost;
            checkbox.value = vhost;
            checkbox.id = vhost;
            tr.insertCell(-1).appendChild(checkbox);
        });

        let div = document.getElementById('VhostDiv');
        div.innerHTML = '';
        div.appendChild(table);

        $('#VhostTable').DataTable({
            deferRender: true,
            orderClasses: false,
            columnDefs: [ {
                targets: 1,
                orderable: false,
            } ]
        });
    }

    showUserInfo(user, data) {
        data.vhosts.forEach((vhost) => {
            let element = document.getElementById(`${user}-${vhost}`);
            if (typeof element !== 'undefined' && element !== null) {
                element.checked = true;
            }
        });
        if (Object.prototype.hasOwnProperty.call(data, 'location')) {
            document.getElementById(`${user}-LocationInput`).value = data.location;
        }

        $(`#${user}-LocationInput`).typeahead({
            hint: true,
            highlight: true,
            minLength: 3
        },
        {
            name: 'Locations',
            source: this.typeaheadData
        });
    }

    bindTypeAhead(data) {
        this.typeaheadData = new Bloodhound({
            datumTokenizer: Bloodhound.tokenizers.whitespace,
            queryTokenizer: Bloodhound.tokenizers.whitespace,
            // `states` is an array of state names defined in "The Basics"
            local: data
        });

        $('#LocationInput').typeahead({
            hint: true,
            highlight: true,
            minLength: 3
        },
        {
            name: 'Locations',
            source: this.typeaheadData
        });
    }

    handleMessage(json) {
        if (Object.prototype.hasOwnProperty.call(json, 'error')) {
            this.log.error(json.error);
            this.ws.showSnackbarMessage(json.error);
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
                if (Object.prototype.hasOwnProperty.call(json.data, 'version')) {
                    document.getElementById('VersionDiv').innerHTML = `Current Version: ${json.data.version}`;
                }
                if (Object.prototype.hasOwnProperty.call(json.data, 'zipCodes')) {
                    this.bindTypeAhead(json.data.zipCodes);
                }
                break;
            case 'generatePassword':
                this.populatePassword(json.component, json.data);
                break;
            case 'getUserDetails':
                this.showUserInfo(json.component, json.data);
                break;
            default:
                this.log.error(`Cannot handle type '${json.type}' for ${this.route}`);
                break;
            }
        }
    }
}

export { Admin };
