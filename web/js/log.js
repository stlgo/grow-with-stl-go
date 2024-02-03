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

const logLevel = {
    6: 'TRACE',
    5: 'DEBUG',
    4: 'INFO',
    3: 'WARN',
    2: 'ERROR',
    1: 'FATAL'
};

const stackRegex = /\((.*):(\d+):(\d+)\)$/;

class Log {
    constructor() {
        this.level = 5;
    }

    setLogLevel(level) {
        if (level in logLevel) {
            this.level = logLevel[level];
        }
    }

    getLogLevel() {
        return this.level;
    }

    foo(message) {
        console.log(message);
        this.writeToLog(6, message);
    }

    trace(message) {
        this.writeToLog(6, message);
    }

    debug(message) {
        this.writeToLog(5, message);
    }

    info(message) {
        this.writeToLog(4, message);
    }

    warn(message) {
        this.writeToLog(3, message);
    }

    error(message) {
        this.writeToLog(2, message);
    }

    fatal(message) {
        this.writeToLog(1, message);
    }

    writeToLog(level, message) {
        // get the caller function
        if (level <= this.level) {
            let stack = new Error().stack.split('\n');
            // works for chrome
            if (stack[2].includes('log.js')) {
                let match = stackRegex.exec(stack[3]);
                if (match !== null) {
                    console.log(`[grow-with-stl-go] ${ new Date().toLocaleString() } ${ match[1] }:${ match[2] } [${ logLevel[level] }] ${ message }`);
                } else {
                    console.log(`[grow-with-stl-go] ${ new Date().toLocaleString() } [${ logLevel[level] }] ${ message }`);
                }
            } else {
                // works for firefox
                let caller = stack[2].split('@')[1].replace();
                caller = caller.substring(0, caller.lastIndexOf(':'));
                console.log(`[grow-with-stl-go] ${ new Date().toLocaleString() } ${ caller } [${ logLevel[level] }] ${ message }`);
            }
        }
    }
}

export {
    Log
};
