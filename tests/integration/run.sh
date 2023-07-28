#!/usr/bin/env bash

set -aex

if [ ! -z "${CI}" ]; then
    NO_COLOR='--no-color'
fi

ginkgo \
    ${NO_COLOR} \
    -vv \
    --keep-going \
    --race \
    --trace \
    --tags integration \
    --label-filter "objectscale" \
    --output-dir=../reports/integration \
    ./...
