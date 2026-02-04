#!/usr/bin/env bash
# Copyright © 2023-2026 Dell Inc. or its subsidiaries. All Rights Reserved.
#
# This software contains the intellectual property of Dell Inc.
# or is licensed to Dell Inc. from third parties. Use of this software
# and the intellectual property contained therein is expressly limited to the
# terms and conditions of the License Agreement under which it is provided by or
# on behalf of Dell Inc. or its subsidiaries.

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
