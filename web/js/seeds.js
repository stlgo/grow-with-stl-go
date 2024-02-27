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
            this.ws.displayHelper([ 'SeedsDiv' ], 'none');
            this.ws.displayHelper([ 'LoadingSeeds' ], '');
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
                Object.keys(data[category].items).forEach((id) => {
                    let seed = data[category].items[id];
                    let seedDiv = document.createElement('div');
                    let h3 = document.createElement('h3');
                    h3.innerHTML = seed.cultivar === undefined ? seed.commonName : `${seed.cultivar} ${seed.commonName}`;

                    let img = document.createElement('img');
                    img.classList.add('rounded');
                    img.src = seed.image;

                    let packets = parseInt(seed.packets);
                    let input = document.createElement('input');
                    input.classList.add('seed-input');
                    input.min = 1;
                    input.type = 'number';
                    input.max = packets;
                    input.value = 1;
                    input.onkeyup = () => {
                        let oldValue = input.min;
                        Array.from(input.value).forEach((c) => {
                            if (!isNaN) {
                                oldValue = oldValue + c;
                            }
                        });
                        if (input.value !== '') {
                            if (parseInt(input.value) < parseInt(input.min)) {
                                input.value = input.min;
                            }
                            if (parseInt(input.value) > parseInt(input.max)) {
                                input.value = input.max;
                            }
                        } else {
                            input.value = oldValue;
                        }
                    };

                    let itemTable = document.createElement('table');
                    let itemBody = itemTable.createTBody();
                    let itemTR = itemBody.insertRow(-1);
                    let itemCell = itemTR.insertCell(-1);
                    let infoButton = document.createElement('button');
                    infoButton.innerHTML = 'Info';
                    infoButton.classList.add('btn', 'btn-info', 'seed-button');
                    infoButton.onclick = () => {
                        this.ws.sendMessage({
                            type: this.type,
                            component: 'getDetail',
                            subComponent: seed.category,
                            data: {
                                id: id
                            }
                        });
                    };
                    itemCell.appendChild(infoButton);
                    itemTR.insertCell(-1).appendChild(input);

                    let requestButton = document.createElement('button');
                    requestButton.innerHTML = 'Request';
                    requestButton.classList.add('btn', 'btn-danger', 'seed-button');
                    requestButton.onclick = () => {
                        this.ws.sendMessage({
                            type: this.type,
                            component: 'requestSeeds',
                            subComponent: seed.id,
                            data: {
                                id: seed.id,
                                quantity: input.value
                            }
                        });
                    };
                    itemTR.insertCell(-1).appendChild(requestButton);

                    seedDiv.appendChild(h3);
                    seedDiv.appendChild(img);
                    seedDiv.appendChild(itemTable);

                    tr.insertCell(-1).appendChild(seedDiv);
                });
                div.appendChild(table);
            } else {
                let p = document.createElement('paragraph');
                p.innerHTML = `No ${category} found in inventory`;
                div.appendChild(heading);
            }
        });
        this.ws.displayHelper([ 'SeedsDiv' ], '');
        this.ws.displayHelper([ 'LoadingSeeds' ], 'none');
    }

    showDetail(data) {
        let table = document.createElement('table');
        let tb = table.createTBody();
        let tr = tb.insertRow(-1);

        let cell = tr.insertCell(-1);
        let img = document.createElement('img');
        img.classList.add('rounded-lg', 'pictures-img');
        img.src = data.image;
        cell.appendChild(img);

        let detailTable = document.createElement('table');
    }

    handleMessage(json) {
        if (Object.prototype.hasOwnProperty.call(json, 'error')) {
            this.log.error(json.error);
            alert(json.error); // eslint-disable-line no-alert
        } else {
            switch(json.component) {
            case 'getDetail':
                console.log(JSON.stringify(json.data));
                break;
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
