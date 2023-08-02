# Copyright © 2023 Dell Inc. or its subsidiaries. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#      http://www.apache.org/licenses/LICENSE-2.0
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

COSI_BUILD_DIR   := build
COSI_BUILD_PATH  := ./cmd/

# When developing docs, chage it so it points to the right file.
CONFIGURATION_DOCS ?= ./docs/installation/configuration_file.md

include overrides.mk

.PHONY: help
help:	##show help
	@fgrep --no-filename "##" $(MAKEFILE_LIST) | fgrep --invert-match fgrep | sed --expression='s/\\$$//' | sed --expression='s/##//'

.PHONY: all
all: clean codegen lint build

.PHONY: clean
clean:	##clean directory
	rm --force pkg/config/*.gen.go
	rm --force core/core.gen.go
	rm --force semver.mk
	go clean

########################################################################
##                                 GO                                 ##
########################################################################

.PHONY: codegen
codegen: clean	##regenerate files
	go generate ./...

# FIXME: remove this target after we remove dependency on private goobjectscale.
.PHONY: vendor
vendor:	##generate the vendor directory
	go mod vendor

.PHONY: build
build:	##build project
	GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o ${COSI_BUILD_DIR}/cosi-driver ${COSI_BUILD_PATH}

########################################################################
##                             CONTAINER                              ##
########################################################################

.PHONY: docker
docker: vendor	##build the docker container
	@echo "Base Images is set to: $(BASEIMAGE)"
	@echo "Building: $(IMAGENAME):$(IMAGETAG)"
	docker build -t "$(IMAGENAME):$(IMAGETAG)" --build-arg BASEIMAGE=$(BASEIMAGE) --build-arg GOVERSION=$(GOVERSION) --build-arg DIGEST=$(DIGEST) .

.PHONY: push
push: docker	##build and push the docker container to repository
	@echo "Pushing: $(REGISTRY)/$(IMAGENAME):$(IMAGETAG)"
	docker tag "$(IMAGENAME):$(IMAGETAG)" "$(REGISTRY)/$(IMAGENAME):$(IMAGETAG)"
	docker push "$(REGISTRY)/$(IMAGENAME):$(IMAGETAG)"

########################################################################
##                              TESTING                               ##
########################################################################

.PHONY: test
test: unit-test fuzz	##run unit and fuzzy tests

.PHONY: unit-test
unit-test:	##run unit tests (Windows or Linux; requires no hardware)
	( go clean -cache; CGO_ENABLED=0 go test -v -coverprofile=c.out ./...)

.PHONY: unit-test-race
unit-test-race:	##run unit tests with race condition reporting (Windows or Linux; requires no hardware, requires CGO)
	( go clean -cache; CGO_ENABLED=1 go test -race -v -coverprofile=c.out ./...)

.PHONY: fuzz
fuzz:	##run all possible fuzzy tests for 10s each
	./scripts/fuzz-all.sh

.PHONY: integration-test
integration-test:	##run integration test (Linux only)
	( cd tests/integration; sh run.sh )

########################################################################
##                              TOOLING                               ##
########################################################################

.PHONY: lint
lint:	##run golangci-lint over the repository
	golangci-lint run

.PHONY: example-config
example-config:	##generate the example configuration file
	./scripts/example-config.pl \
		--filename $(CONFIGURATION_DOCS) \
		--output ./sample-config.generated.yaml
