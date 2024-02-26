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
                    h3.innerHTML = seed.cultivar === undefined ? seed.commonName : `${seed.cultivar} ${seed.commonName}`;

                    let img = document.createElement('img');
                    img.classList.add('rounded');
                    img.src = seed.image;

                    let packets = parseInt(seed.packets);
                    let input = document.createElement('input');
                    input.classList.add('form-control');
                    input.type = 'number';
                    input.min = 1;
                    input.max = packets;
                    input.onkeyup = () => {
                        let value = input.value;
                        console.log(value);
                        let result = '';
                        console.log(Array.from(value));
                        Array.from(value).forEach((char) => {
                            console.log(char);
                            if (!isNaN(char)) {
                                result = result + char;
                            } else {
                                console.log(char);
                            }
                        });
                        console.log(result);
                        value = parseInt(result);
                        console.log(value);
                        switch (value) {
                        case value < 0:
                            input.value = 1;
                            break;
                        case value > packets:
                            input.value = packets;
                            break;
                        default:
                            input.value = value;
                            break;
                        }
                    };
                    input.value = 1;


                    seedDiv.appendChild(h3);
                    seedDiv.appendChild(img);
                    seedDiv.appendChild(input);

                    tr.insertCell(-1).appendChild(seedDiv);
                });
                div.appendChild(table);
            } else {
                let p = document.createElement('paragraph');
                p.innerHTML = `No ${category} found in inventory`;
                div.appendChild(heading);
            }
        });
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
