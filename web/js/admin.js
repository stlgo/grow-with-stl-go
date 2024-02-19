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
                navigator.clipboard.writeText(data.password);
                alert('Password has been copied to your clipboard'); // eslint-disable-line no-alert
                break;
            case 'existingUser':
                if (Object.prototype.hasOwnProperty.call(data, 'target')) {
                    document.getElementById(data.target).value = data.password;
                    navigator.clipboard.writeText(data.password);
                    alert('Password has been copied to your clipboard'); // eslint-disable-line no-alert
                    break;
                }
                break;
            default:
                this.log.error(`cannot apply password for user type '${userType}'`);
            }
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
            tr.insertCell(-1);
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
            checkbox.checked = true;
            tr.insertCell(-1).appendChild(checkbox);
            tr.insertCell(-1);
        });
        let div = document.getElementById('CurrentUsersDiv');
        div.innerHTML = '';
        div.appendChild(table);

        let userTable = $('#UsersTable').DataTable({
            deferRender: true,
            orderClasses: false,
            columnDefs: [ {
                targets: 0,
                className: 'details-control',
                orderable: false,
                data: null,
                defaultContent: '<i class="fa fa-plus-circle fa-lg" style="color:green"></i>'
            }, {
                targets: 4,
                className: 'remove-user',
                orderable: false,
                data: null,
                defaultContent: '<i class="fa fa-trash fa-lg" style="color:maroon"></i>'
            } ]
        });

        this.bindActiveToggle();
        this.bindSlideOut(userTable);
        this.bindRemoveUser(userTable);
    }

    bindSlideOut(table) {
        // TODO (aschiefe): replace jquery lookup with pure JS one
        $('#UsersTable tbody').on('click', 'td.details-control', function() {
            const tr = $(this).closest('tr'); // eslint-disable-line no-invalid-this
            const row = table.row(tr);

            if (row.child.isShown()) {
                row.child.hide();
                tr.removeClass('shown');
                this.innerHTML = '<i class="fa fa-plus-circle fa-lg" style="color:green"></i>'; // eslint-disable-line no-invalid-this
            } else {
                const userID = row.data()[1];
                const div = document.createElement('div');
                div.id = `${userID}-detailsDiv`;
                div.innerHTML = `Loading data for user "${userID}" please wait...`;
                row.child(div).show();
                tr.addClass('shown');
                this.innerHTML = '<i class="fa fa-times-circle fa-lg" style="color:maroon"></i>'; // eslint-disable-line no-invalid-this
            }
        });
    }

    bindActiveToggle() {
        document.querySelectorAll('.user-active-chk').forEach((chkBox) => {
            chkBox.onclick = () => {
                const user = chkBox.name;
                let active = false;
                active = chkBox.checked;
                this.log.info(active);
            };
        });
    }

    bindRemoveUser(table) {
        // TODO (aschiefe): replace jquery lookup with pure JS one
        $('#UsersTable tbody').on('click', 'td.remove-user', function() {
            const tr = $(this).closest('tr'); // eslint-disable-line no-invalid-this
            const row = table.row(tr);
            const userID = row.data()[1];
            if (window.confirm(`Confirm removal of user "${userID}"`)) { // eslint-disable-line no-alert
                console.log(userID);
            }
        });
    }

    handleMessage(json) {
        if (Object.prototype.hasOwnProperty.call(json, 'error')) {
            this.log.error(json.error);
        } else {
            switch(json.component) {
            case 'addUser':
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
