#!/usr/bin/env bash

set -aex

ginkgo \
    --keep-going \
    --race \
    --trace \
    --tags integration \
    --output-dir=../reports/integration \
    ./...
