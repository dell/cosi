#!/usr/bin/env bash

set -aex

if [ -n "${CI}" ]; then
    NO_COLOR='--no-color'
fi

# shellcheck disable=SC2086
ginkgo \
    ${NO_COLOR} \
    --keep-going \
    --race \
    --trace \
    --tags integration \
    --label-filter "objectscale" \
    --output-dir=../reports/integration \
    ./...
