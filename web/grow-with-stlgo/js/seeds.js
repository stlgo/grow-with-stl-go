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

        this.route = this.constructor.name.toLowerCase();
        this.ws.registerHandlers(this.route, this);

        document.addEventListener('seeds', () => {
            this.ws.displayHelper([ 'SeedsDiv' ], 'none');
            this.ws.displayHelper([ 'LoadingSeeds' ], '');
            this.ws.sendMessage({
                route: this.route,
                type: 'getInventory'
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
                    input.id = `${seed.id}-quantity`;
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
                    itemTable.classList.add('seed-info-table');
                    let itemBody = itemTable.createTBody();

                    let columns = {
                        '<b>Seeds Per Packet</b>': seed.perPacketCount,
                        '<b>Price Per Packet</b>': `$${seed.price}`
                    };

                    Object.keys(columns).forEach((key) => {
                        let itemTR = itemBody.insertRow(-1);
                        let itemCell = itemTR.insertCell(-1);
                        itemCell.colSpan = '2';
                        itemCell.innerHTML = key;
                        itemTR.insertCell(-1).innerHTML = columns[key];
                    });

                    let itemTR = itemBody.insertRow(-1);
                    let itemCell = itemTR.insertCell(-1);
                    itemCell.colSpan = '2';
                    itemCell.innerHTML = '<b>Packets Available</b>';
                    let availableDiv = document.createElement('div');
                    availableDiv.id = `${seed.id}-available`;
                    availableDiv.innerHTML = seed.packets;
                    itemTR.insertCell(-1).appendChild(availableDiv);

                    itemTR = itemBody.insertRow(-1);
                    itemCell = itemTR.insertCell(-1);
                    let infoButton = document.createElement('button');
                    infoButton.innerHTML = 'Info';
                    infoButton.classList.add('btn', 'btn-info', 'seed-button');
                    infoButton.onclick = () => {
                        this.ws.sendMessage({
                            route: this.route,
                            type: 'getDetail',
                            component: seed.category,
                            subComponent: id,
                        });
                    };
                    itemCell.appendChild(infoButton);
                    itemTR.insertCell(-1).appendChild(input);

                    let purchaseButton = document.createElement('button');
                    purchaseButton.id = `${seed.id}-purchase`;
                    purchaseButton.innerHTML = 'Purchase';
                    purchaseButton.classList.add('btn', 'btn-danger', 'seed-button');
                    purchaseButton.onclick = () => {
                        this.ws.sendMessage({
                            route: this.route,
                            type: 'purchase',
                            component: seed.category,
                            subComponent: seed.id,
                            data: {
                                id: seed.id,
                                quantity: input.value
                            }
                        });
                    };
                    itemTR.insertCell(-1).appendChild(purchaseButton);

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
        detailTable.classList.add('seed-info-table');
        let dtb = detailTable.createTBody();

        let columns = {
            '<b>Category</b>': data.category,
            '<b>Genus</b>': data.genus,
            '<b>Species</b>': data.species,
            '<b>Cultivar</b>': data.cultivar === undefined ? 'N/A' : data.cultivar,
            '<b>Common Name</b>': data.commonName,
            '<b>Description</b>': data.description,
            '<b>Seeds Per Packet</b>': data.perPacketCount,
        };

        Object.keys(columns).forEach((key) => {
            let dtr = dtb.insertRow(-1);
            dtr.insertCell(-1).innerHTML = key;
            dtr.insertCell(-1).innerHTML = columns[key];
        });

        let dtr = dtb.insertRow(-1);
        dtr.insertCell(-1).innerHTML = '<b>Packets Available</b>';
        let availableDiv = document.createElement('div');
        availableDiv.id = `${data.id}-available`;
        availableDiv.innerHTML = data.packets;
        dtr.insertCell(-1).appendChild(availableDiv);
        dtr = dtb.insertRow(-1);
        dtr.insertCell(-1).innerHTML = '<b>Price</b>';
        dtr.insertCell(-1).innerHTML = `$${data.price}`;
        dtr = dtb.insertRow(-1);
        let element = document.getElementById(`${data.id}-quantity`);
        let parent = element.parentElement;
        parent.removeChild(element);
        dtr.insertCell(-1).appendChild(element);
        element = document.getElementById(`${data.id}-purchase`);
        parent = element.parentElement;
        parent.removeChild(element);
        dtr.insertCell(-1).appendChild(element);

        tr.insertCell(-1).appendChild(detailTable);

        let backHref = document.createElement('a');
        backHref.textContent = 'Return to seeds';
        backHref.href = '/seeds';
        backHref.onclick = (event) => {
            event.preventDefault();
            this.ws.getPagelet('seeds');
        };

        const div = document.getElementById('SeedsDiv');
        div.innerHTML = '';

        div.appendChild(backHref);
        div.appendChild(table);
    }

    updateSeed(data) {
        document.getElementById(`${data.id}-available`).innerHTML = data.packets;
        document.getElementById(`${data.id}-quantity`).max = data.packets;
    }

    handleMessage(json) {
        if (Object.prototype.hasOwnProperty.call(json, 'error')) {
            this.log.error(json.error);
            alert(json.error); // eslint-disable-line no-alert
        } else {
            switch(json.type) {
            case 'getDetail':
                this.showDetail(json.data);
                break;
            case 'getInventory':
                this.showSeeds(json.data);
                break;
            case 'purchase':
                this.updateSeed(json.data);
                this.log.info(`\n${JSON.stringify(json, null, 4)}`);
                break;
            default:
                this.log.error(`Cannot handle component '${json.component}' for ${this.type}`);
                break;
            }
        }
    }
}

export { Seeds };
