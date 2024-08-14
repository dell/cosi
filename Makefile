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

.PHONY: all
all: clean codegen lint build

include overrides.mk

.PHONY: help
help:	##show help
	@fgrep --no-filename "##" $(MAKEFILE_LIST) | fgrep --invert-match fgrep | sed --expression='s/\\$$//' | sed --expression='s/##//'

.PHONY: clean
clean:	##clean directory
	rm --force pkg/config/*.gen.go
	rm --force pkg/internal/iamapi/mock/*.gen.go
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
build: ##build project
	GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o ${COSI_BUILD_DIR}/cosi ${COSI_BUILD_PATH}

########################################################################
##                             CONTAINER                              ##
########################################################################

.PHONY: build-base-image
build-base-image: vendor download-csm-common
	$(eval include csm-common.mk)
	sh ./scripts/build-ubi-micro.sh $(DEFAULT_BASEIMAGE)
	$(eval BASEIMAGE=cosi-ubimicro:latest)

.PHONY: podman
podman: build-base-image
	@echo "Base Images is set to: $(BASEIMAGE)"
	@echo "Building: $(IMAGENAME):$(IMAGETAG)"
	podman build -t "$(IMAGENAME):$(IMAGETAG)" --build-arg BASEIMAGE=$(BASEIMAGE) --build-arg GOIMAGE=$(DEFAULT_GOIMAGE) .

.PHONY: push
push: podman	##build and push the podman container to repository
	@echo "Pushing: $(REGISTRY)/$(IMAGENAME):$(IMAGETAG)"
	podman tag "$(IMAGENAME):$(IMAGETAG)" "$(REGISTRY)/$(IMAGENAME):$(IMAGETAG)"
	podman push "$(REGISTRY)/$(IMAGENAME):$(IMAGETAG)"

.PHONY: download-csm-common
download-csm-common:
	curl -O -L https://raw.githubusercontent.com/dell/csm/main/config/csm-common.mk

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
