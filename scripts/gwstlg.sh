#!/bin/bash
#
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
#

GWSTLG_HOME=$(dirname $(dirname $(readlink -fm $0)))
GWSTLG_LOG=$(dirname ${GWSTLG_HOME})/logs
GWSTLG_ETC=$(dirname ${GWSTLG_HOME})/etc

LOG_LEVEL=6

# sometimes the display can be set and it messes with the output of the script
unset DISPLAY

# Alter umask so that the group owner of the log files has read / write permissions
umask 002

# text display colors
RED=`tput setaf 1`
GREEN=`tput setaf 2`
NORMAL=`tput sgr0`

startServer() {
    # Externalize the log directory so it's not overwritten at time of revision
    if [ ! -d "${GWSTLG_LOG}" ]; then
        mkdir -p ${GWSTLG_LOG}
    fi

    if [ ! -d "${GWSTLG_HOME}/logs" ]; then
        ln -sf ${GWSTLG_LOG} ${GWSTLG_HOME}/logs
    fi

    # Externalize the etc directory so it's not overwritten at time of revision
    if [ ! -d "${GWSTLG_ETC}" ]; then
        mkdir -p ${GWSTLG_ETC}
    fi

    if [ ! -d "${GWSTLG_HOME}/etc" ]; then
        ln -sf ${GWSTLG_ETC} ${GWSTLG_HOME}/etc
    fi

    # Check to see if it's already running
    if [ -r ${GWSTLG_LOG}/gwstlg.pid ]; then
        PID=`cat ${GWSTLG_LOG}/gwstlg.pid`
        if [ `ps -eaf | grep ${PID} | grep -v grep | grep gwstlg | wc -l` != 0 ]; then
            echo -e "${0} ERROR: Grow With STL Go process id ${PID} is still running, unable to issue start: \t\t\t[ ${RED}FAILED ${NORMAL}]"
            exit 1
        fi
    fi

    # Log roll
    if [ -r ${GWSTLG_LOG}/SystemOut.log ]; then
        DATE=$(date +%a-%b-%d-%H.%M.%S-%Y)
        mv ${GWSTLG_LOG}/SystemOut.log ${GWSTLG_LOG}/SystemOut.log.${DATE}
        gzip -9 ${GWSTLG_LOG}/SystemOut.log.${DATE}
    fi

    # Get rid of old stuff, per whatever retention period you define
    find ${GWSTLG_LOG} -mindepth 1 -mtime +30 -name 'SystemOut.log.*' -exec rm {} \;

    nohup ${GWSTLG_HOME}/bin/grow-with-stl-go --loglevel ${LOG_LEVEL} -c ${GWSTLG_ETC}/grow-with-stlg-go.json >> ${GWSTLG_LOG}/SystemOut.log & echo $! > ${GWSTLG_LOG}/gwstlg.pid 2>&1 &
    checkStart
}

checkStart() {
    MAX_WAITS="10"
    STATUS_INTERVAL="3"
    SUCCESS_MESSAGE="Attempting to start webservice on"

    while [ ${MAX_WAITS} != -1 ]; do
        if [ -r ${GWSTLG_LOG}/SystemOut.log ]; then
            if [ `tail ${GWSTLG_LOG}/SystemOut.log | grep "${SUCCESS_MESSAGE}" | wc -l` != 0 ]; then
                echo -e "${0} STATUS: Grow With STL Go has started: \t\t\t\t\t[ ${GREEN}OK ${NORMAL}]"
                break
            else
                if [[ ${MAX_WAITS} -le 1 ]]; then
                    echo -e "${0} ERROR: Grow With STL Go start: \t\t\t\t\t[ ${RED}FAILED ${NORMAL}]"
                    exit 1
                else
                    MAX_WAITS=$((MAX_WAITS - 1))
                    echo -e "${0} STATUS: Waiting for Grow With STL Go to start, ${MAX_WAITS} more attempts before giving up"
                    sleep ${STATUS_INTERVAL}
                fi
            fi
        fi
    done
}

stopServer() {
    STATUS_INTERVAL="10"
    if [ -r ${GWSTLG_LOG}/gwstlg.pid ]; then
        PID=`cat ${GWSTLG_LOG}/gwstlg.pid`
        if [ `ps -eaf | grep ${PID} | grep -v grep | grep gwstlg | wc -l` != 0 ]; then
            echo -e "${0} STATUS: Issuing kill on Grow With STL Go process id ${PID}"
            kill ${PID}
            sleep ${STATUS_INTERVAL}
            if [ `ps -eaf | grep ${PID} | grep -v grep | grep gwstlg | wc -l` != 0 ]; then
                echo -e "${0} STATUS: Grow With STL Go process did not exit normaill, issuing kill -9 on process id ${PID}"
                kill ${PID}
                sleep ${STATUS_INTERVAL}
                if [ `ps -eaf | grep ${PID} | grep -v grep | grep gwstlg | wc -l` != 0 ]; then
                    echo -e "${0} ERROR: Grow With STL Go process id ${PID} failed to stop: \t\t\t[ ${RED}FAILED ${NORMAL}]"
                    exit 1
                fi
            fi
            echo -e "${0} ERROR: Grow With STL Go process id ${PID} has stopped: \t\t\t[ ${GREEN}OK ${NORMAL}]"
            rm -rf ${GWSTLG_LOG}/gwstlg.pid
        fi
    else
        echo -e "${0} STATUS: There are no Grow With STL Go process id files, you may need to manually kill any processes: \t\t\t[ ${RED}FAILED ${NORMAL}]"
    fi
}

case ${1} in
'start')
    startServer
    ;;
'stop')
    stopServer
    ;;
'restart')
    stopServer
    sleep 10
    startServer
    ;;
*)
    echo "usage: ${0} { start | stop | restart }"
    exit 1
    ;;
esac
