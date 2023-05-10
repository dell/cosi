#!/usr/bin/env bash

set -aex

if [ ! -z "${CI}" ]; then
    NO_COLOR='--no-color'
fi

ginkgo \
    "${NO_COLOR}" \
    --keep-going \
    --race \
    --trace \
    --tags integration \
    --output-dir=../reports/integration \
    ./...