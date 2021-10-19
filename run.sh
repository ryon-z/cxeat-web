#!/usr/bin/env bash
FILE_PATH=`realpath ${0}`
WORKING_DIR_PATH=`dirname ${FILE_PATH}`

function run() {
    yelloment_pid=`ps aux | grep "yelloment-api" | awk '{if($11=="'${WORKING_DIR_PATH}'/yelloment-api") {print $2}}'`
    if [[ ${yelloment_pid} != "" ]];then
        echo "Could not found yelloment_api's pid"

        echo "yelloment_api's pid: ${yelloment_pid}"
        kill -9 ${yelloment_pid}
        echo "yelloment_api is killed"
    fi
    
    go build
    nohup "${WORKING_DIR_PATH}/yelloment-api" &
    echo "yelloment_api is restarted"
}

run
