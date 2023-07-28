#!/bin/bash

set -e

fuzzTime=${1:-10}

files=$(grep -r --include='**_test.go' --files-with-matches 'func Fuzz' .)
FAILED=0


for file in ${files}
do
    funcs=$(grep -oP 'func \K(Fuzz\w*)' "${file}")
    for func in ${funcs}
    do
        echo "Fuzzing ${func} in ${file}"
        parentDir=$(dirname "${file}")
        if ! go test "${parentDir}" -run="${func}" -fuzz="${func}" -fuzztime="${fuzzTime}"s; then
            FAILED=$((FAILED+1))
        fi
    done
done

if [[ ${FAILED} -ne 0 ]]; then
    echo "${FAILED} fuzzy tests failed!"
    exit 1
fi
