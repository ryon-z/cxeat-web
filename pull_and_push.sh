TRUE=1
FALSE=0

function is_succeeded_git_command() {
    local message=$1

    if [[ ${message} == *error* ]]; then
        return ${FALSE}
    fi

    return ${TRUE}
}

function sleep_with_msg() {
    local sleep_sec=$1
    echo "+-- sleep ${sleep_sec} sec"
    sleep ${sleep_sec}
}

function pull_and_push() {
    echo "+-- pull --rebase from upstream dev"
    result=1

    # Catch stderr of this command
    result=`git pull --rebase upstream dev 2>&1 > /dev/null`
    sleep_with_msg 0.5

    is_succeeded_git_command "${result}"
    is_succeeded=$?
    if [[ is_succeeded -eq ${FALSE} ]]; then
        echo "${result}"
        echo "+-- failed"
        return 9001 
    fi
    echo "echo result :: ${result}"

    # Catch stderr of this command
    echo "+-- push from local dev to upstream dev"
    result=`git push upstream dev 2>&1 > /dev/null`
    sleep_with_msg 0.5

    is_succeeded_git_command "${result}"
    is_succeeded=$?
    if [[ is_succeeded -eq ${FALSE} ]]; then
        echo "${result}"
        echo "+-- failed"
        return 9001
    fi
    echo "echo result :: ${result}"

    # Catch stderr of this command
    echo "+-- push from local dev to origin dev"
    result=`git push origin dev 2>&1 > /dev/null`
    sleep_with_msg 0.5

    is_succeeded_git_command "${result}"
    is_succeeded=$?
    if [[ is_succeeded -eq ${FALSE} ]]; then
        echo "${result}"
        echo "+-- failed"
        return 9001
    fi
    echo "echo result :: ${result}"

    echo "+-- done"
}

pull_and_push
