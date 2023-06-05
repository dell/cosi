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
#
all: clean generate format build

COSI_BUILD_DIR   := build
COSI_BUILD_PATH  := ./cmd/

include overrides.mk

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

help:	##show help
	@fgrep --no-filename "##" $(MAKEFILE_LIST) | fgrep --invert-match fgrep | sed --expression='s/\\$$//' | sed --expression='s/##//'

format:	##run gofmt
	gofmt -w -s .

clean:	##clean directory
	rm --force core/core.gen.go
	rm --force semver.mk
	go clean

.PHONY: generate
generate: ##regenerate files
	go generate ./...

.PHONY: build
build:	##build project
	GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o ${COSI_BUILD_DIR}/cosi-driver ${COSI_BUILD_PATH}

# Builds dockerfile
docker:	##generate the docker container
	@echo "Base Images is set to: $(BASEIMAGE)"
	@echo "Building: $(IMAGENAME):$(IMAGETAG)"
	docker build -t "$(IMAGENAME):$(IMAGETAG)" --build-arg BASEIMAGE=$(BASEIMAGE) --build-arg GOVERSION=$(GOVERSION) --build-arg DIGEST=$(DIGEST) .

# Pushes container to the repository
push: docker	##generate and push the docker container to repository
	@echo "Pushing: $(REGISTRY)/$(IMAGENAME):$(IMAGETAG)"
	docker tag "$(IMAGENAME):$(IMAGETAG)" "$(REGISTRY)/$(IMAGENAME):$(IMAGETAG)"
	docker push "$(REGISTRY)/$(IMAGENAME):$(IMAGETAG)"

# Windows or Linux; requires no hardware
unit-test:	##run unit tests
	( go clean -cache; CGO_ENABLED=0 go test -v -coverprofile=c.out ./...)

# Windows or Linux; requires no hardware, requires CGO
unit-test-race:	##run unit tests with race condition reporting
	( go clean -cache; CGO_ENABLED=1 go test -race -v -coverprofile=c.out ./...)

# Linux only; populate env.sh with the hardware parameters
integration-test:	##run integration test (Linux only)
	( cd tests/integration; sh run.sh )

.PHONY: lint
lint: 
	golangci-lint run

.PHONY: gofumpt
gofumpt:
	gofumpt -w cmd pkg tests   
