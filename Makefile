# Copyright © 2026 Dell Inc. or its subsidiaries. All Rights Reserved.
#
# Dell Technologies, Dell and other trademarks are trademarks of Dell Inc.
# or its subsidiaries. Other trademarks may be trademarks of their respective
# owners.

include images.mk

COSI_BUILD_DIR   := build
COSI_BUILD_PATH  := ./cmd/

# When developing docs, change it so it points to the right file.
CONFIGURATION_DOCS ?= ./docs/installation/configuration_file.md

.PHONY: all
all: codegen build test

.PHONY: help
help:	##show help
	@fgrep --no-filename "##" $(MAKEFILE_LIST) | fgrep --invert-match fgrep | sed --expression='s/\\$$//' | sed --expression='s/##//'

.PHONY: clean
clean:	##clean directory
	rm --force pkg/config/*.gen.go
	rm --force pkg/internal/iamapi/mock/*.gen.go
	rm --force golangci.yaml csm-common.mk go-code-tester
	rm --force coverage_results.txt coverage.txt cover.out new_coverage.txt
	rm --force --recursive $(COSI_BUILD_DIR)
	go clean

########################################################################
##                                 GO                                 ##
########################################################################

.PHONY: codegen
codegen: ##regenerate files
	go generate -skip="mockery" ./...

.PHONY: build
build: codegen ##build project
	GOOS=linux CGO_ENABLED=0 go build -mod=vendor -ldflags="-s -w" -o ${COSI_BUILD_DIR}/cosi ${COSI_BUILD_PATH}

########################################################################
##                              TESTING                               ##
########################################################################

.PHONY: test
test: unit-test fuzz	##run unit and fuzzy tests

.PHONY: unit-test
unit-test: go-code-tester
	GITHUB_OUTPUT=/dev/null \
	./go-code-tester 90 "." "" "true" "" "" "./pks/provisioner/virtualdriver/fake|./pkg/provisioner/virtualdriver/fake|./pkg/provisioner/objectscale/mocks|./tests/integration/steps|./pkg/provisioner/virtualdriver"

.PHONY: fuzz
fuzz:	##run all possible fuzzy tests for 10s each
	./scripts/fuzz-all.sh

.PHONY: integration-test
integration-test:	##run integration test (Linux only)
	(cd tests/integration; sh run.sh)

########################################################################
##                              TOOLING                               ##
########################################################################

.PHONY: example-config
example-config:	##generate the example configuration file
	./scripts/example-config.pl \
		--filename $(CONFIGURATION_DOCS) \
		--output ./sample-config.generated.yaml
