#!/bin/bash

set -e

echo "Building..."
go build -o subMatcher ../../

doTest() {
    local TEST_DIR=$1
    local EXPECTED_RESULT=$2
    local FLAGS=$3
    local BACKUP="tmp/backup"
    mkdir -p tmp

    cp -R $TEST_DIR $BACKUP

    ./subMatcher -y $3 $TEST_DIR > tmp/output.txt
    ls --recursive $TEST_DIR > tmp/result.test

    set +e
    cmp tmp/result.test $EXPECTED_RESULT
    CHECK=$?
    set -e

    rm -rf $TEST_DIR
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
doTest "./cases/case1" "expectations/1.txt" "-r"
doTest "./cases/case2" "expectations/2.txt" "-r"
