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

export const LogLevel = {
    TRACE: 6,
    DEBUG: 5,
    INFO: 4,
    WARN: 3,
    ERROR: 2,
    FATAL: 1
};

export class Log {
    constructor() {
        this.level = 5;
    }

    setLogLevel(level) {
        if (level in LogLevel) {
            this.level = LogLevel[level];
        }
    }

    static trace(message) {
        this.writeToLog(LogLevel.Trace, message);
    }

    static debug(message) {
        this.writeToLog(LogLevel.Debug, message);
    }

    static info(message) {
        this.writeToLog(LogLevel.Info, message);
    }

    static warn(message) {
        this.writeToLog(LogLevel.Warn, message);
    }

    static error(message) {
        this.writeToLog(LogLevel.Error, message);
    }

    static fatal(message) {
        this.writeToLog(LogLevel.Fatal, message);
    }

    static writeToLog(level, message) {
        if (level <= this.Level) {
            console.log(
                `[grow-with-stl-go][${ LogLevel[level] }] ${ new Date().toLocaleString() } - ${
                    message.className } - ${ message.message }: `, message.logMessage);
        }
    }
}
