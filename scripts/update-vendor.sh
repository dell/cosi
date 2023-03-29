#!/usr/bin/env bash

set -a
set -e

# Open the file for reading
file="go.mod"

GOPRIVATE="github.com/dell/objectscale"

# Define flags
in_require=0
skip_line=0

# Read the file line by line
while read line
do
    # Check if the line contains "require ("
    if echo "$line" | grep -q "require ("
    then
        in_require=1
        skip_line=0
        continue
    fi

    # Check if the line contains ")"
    if echo "$line" | grep -q ")"
    then
        in_require=0
        skip_line=0
        continue
    fi

    # Check if the line contains "// indirect"
    if echo "$line" | grep -q "// indirect"
    then
        skip_line=1
        continue
    fi

    # Update the dependency
    if [ "$in_require" -eq 1 ] && [ "$skip_line" -eq 0 ]
    then
        dependency=$(echo "$line" | cut -d ' ' -f 1)

        echo "updating $dependency"
        go get -u "$dependency"
    fi
done < "$file"

# run go mod tidy
echo "updating sum file"
go mod tidy

# run go mod vendor
echo "updating vendor directory"
go mod vendor
