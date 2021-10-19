#!/usr/bin/env bash
SERVICE_URL="http://cueat.kr/"

function health_check() {
    # 요청
    status_code=`curl -s -o /dev/null -w "%{http_code}" ${SERVICE_URL} `
    echo "status_code: ${status_code}"


    # 정상인지 검사
    share=$(( status_code / 100 ))
    mod=$(( status_code % 100 ))
    if [[ $(( share - mod )) != 2 ]];then
        curl -X POST https://hooks.slack.com/services/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 \
            -H 'Content-type: application/json' \
            --data '{"text":"health_check 실패 -  url: cueat.kr, status_code: '${status_code}'"}'
    fi
}

function alarmDiskFull() {
   limit=80
   remain_percent=`df -h | grep /dev/root | awk '{ print $5 }'`
   remain_percent=${remain_percent%%%}

   # 서버 사용 용량이 제한 용량보다 많은 경우
   if [[ ${remain_percent} -gt ${limit} ]]; then
       message="서버 용량이 거의 다찼습니다. 현재 사용 용량 퍼센트: ${remain_percent}"
       curl -X POST https://hooks.slack.com/services/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 \
            -H 'Content-type: application/json' \
            --data "{\"text\":\"${message}\"}"
   fi
}

echo "Started health_check"
while :
do
    health_check
    alarmDiskFull
    sleep 3m
done

# 실행 명령어.
# nohup bash health_checker.sh > /dev/null 2>&1 &
