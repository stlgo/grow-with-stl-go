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


class Seeds {
    constructor(ws, log) {
        this.ws = ws;
        this.log = log;

        this.type = this.constructor.name.toLowerCase();
        this.ws.registerHandlers(this.type, this);

        document.addEventListener('seeds', () => {
            this.ws.sendMessage({
                type: this.type,
                component: 'getInventory',
                subComponent: 'getInventory',
            });
        });
    }

    showSeeds(data) {
        const div = document.getElementById('SeedsDiv');
        Object.keys(data).forEach((category) => {
            let heading = document.createElement('h2');
            heading.innerHTML = category;
            div.appendChild(heading);
            if (Object.prototype.hasOwnProperty.call(data[category], 'items')) {
                let table = document.createElement('table');
                table.setAttribute('width', '99%');
                table.classList.add('display', 'responsive', 'seeds', 'pictures');
                let tb = table.createTBody();
                let tr = tb.insertRow(-1);
                data[category].items.forEach((seed) => {
                    let seedDiv = document.createElement('div');
                    let h3 = document.createElement('h3');
                    h3.innerHTML = seed.commonName;

                    let img = document.createElement('img');
                    img.classList.add('rounded');
                    img.src = seed.image;

                    seedDiv.appendChild(h3);
                    seedDiv.appendChild(img);
                    tr.insertCell(-1).appendChild(seedDiv);
                });
                div.appendChild(table);
            } else {
                let p = document.createElement('paragraph');
                p.innerHTML = `No ${category} found in inventory`;
                div.appendChild(heading);
            }
        });
        // const table = document.createElement('table');
        // table.setAttribute('width', '99%');
        // table.id = 'UsersTable';
        // table.classList.add('display', 'responsive', 'seeds');
        // let th = table.createTHead();
        // let tr = th.insertRow(-1);
        // headers.forEach((header) => {
        //     tr.insertCell(-1).innerHTML = header;
        // });
        // let tb = table.createTBody();
        // Object.keys(data).forEach((userID) => {
        //     tr = tb.insertRow(-1);
        //     tr.insertCell(-1).innerHTML = `<div id="${userID}-details"><i class="fa fa-plus-circle fa-lg details-control" style="color:green"></i></div>`;
        //     tr.insertCell(-1).innerHTML = userID;
        //     if (data[userID].lastLogin === null) {
        //         tr.insertCell(-1).innerHTML = 'Never Logged In';
        //     } else {
        //         tr.insertCell(-1).innerHTML = new Date(data[userID].lastLogin).toLocaleString('en-us', { timeZoneName: 'short' });
        //     }


        //     let activeChk = document.createElement('input');
        //     activeChk.type = 'checkbox';
        //     activeChk.id = `${userID}_active_chk`;
        //     activeChk.name = userID;
        //     activeChk.classList.add('user-active-chk');
        //     activeChk.checked = data[userID].active;
        //     tr.insertCell(-1).appendChild(activeChk);

        //     let adminChk = document.createElement('input');
        //     adminChk.type = 'checkbox';
        //     adminChk.id = `${userID}_admin_chk`;
        //     adminChk.name = userID;
        //     adminChk.classList.add('user-admin-chk');
        //     adminChk.checked = data[userID].admin;
        //     tr.insertCell(-1).appendChild(adminChk);

        //     tr.insertCell(-1).innerHTML = '<i class="fa fa-trash fa-lg remove-user" style="color:maroon"></i>';
        // });
        // let div = document.getElementById('CurrentUsersDiv');
        // div.innerHTML = '';
        // div.appendChild(table);

        // let userTable = $('#UsersTable').DataTable({
        //     deferRender: true,
        //     orderClasses: false,
        //     columnDefs: [ {
        //         targets: 0,
        //         orderable: false,
        //     }, {
        //         targets: 4,
        //         orderable: false,
        //     } ]
        // });

        // this.bindActiveToggle();
        // this.bindAdminToggle();
        // this.bindSlideOut(userTable);
        // this.bindRemoveUser(userTable);
    }

    handleMessage(json) {
        if (Object.prototype.hasOwnProperty.call(json, 'error')) {
            this.log.error(json.error);
            alert(json.error); // eslint-disable-line no-alert
        } else {
            switch(json.component) {
            case 'getInventory':
                this.showSeeds(json.data);
                break;
            default:
                this.log.error(`Cannot handle component '${json.component}' for ${this.type}`);
                break;
            }
        }
    }
}

export { Seeds };
