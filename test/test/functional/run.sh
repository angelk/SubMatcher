#!/bin/bash

set -e

doTest() {
    local TEST_DIR=$1
    local EXPECTED_RESULT=$2
    local BACKUP="tmp/backup"
    mkdir -p tmp

    cp -R $TEST_DIR $BACKUP

    ls --recursive $TEST_DIR > tmp/result.test

    set +e
    cmp tmp/result.test $EXPECTED_RESULT
    CHECK=$?
    set -e

    rm -rf $TEST_DIR
    echo "mv $BACKUP $TEST_DIR"
    mv $BACKUP $TEST_DIR


    if [ $CHECK -ne "0" ]; then
        echo "$TEST_DIR test failed"
        echo "--- Expected "
        cat $EXPECTED_RESULT
        echo "+++ Got"
        cat tmp/result.test
        exit 1
    fi

    echo "$TEST_DIR ok"
}

doTest "./cases/case1" "expectations/1.txt"
