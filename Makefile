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
all: clean format build

COSI_BUILD_DIR   := build
COSI_BUILD_PATH  := ./cmd/

include overrides.mk

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# Tag parameters
MAJOR=1
MINOR=0
PATCH=0
NOTES=-beta
TAGMSG="COSI Spec 1.0"

help:	##show help
	@fgrep --no-filename "##" $(MAKEFILE_LIST) | fgrep --invert-match fgrep | sed --expression='s/\\$$//' | sed --expression='s/##//'

format:	##run gofmt
	gofmt -w -s .

clean:	##clean directory
	rm --force core/core.gen.go
	go clean

.PHONY: build
build: ##build project
	GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o ${COSI_BUILD_DIR}/cosi-driver ${COSI_BUILD_PATH}

# Tags the release with the Tag parameters set above
tag:	##tag the release
	go generate ./...
	go run core/semver/semver.go -f mk >semver.mk
	-git tag --delete v$(MAJOR).$(MINOR).$(PATCH)$(NOTES)
	git tag --annotate --message=$(TAGMSG) v$(MAJOR).$(MINOR).$(PATCH)$(NOTES)

# Generates the docker container (but does not push)
docker: ##generate the docker container
	make --file=docker.mk docker

# Pushes container to the repository
push: docker	##generate and push the docker container to repository
	make --file=docker.mk push

# Windows or Linux; requires no hardware
unit-test:	##run unit tests
	( go clean -cache; CGO_ENABLED=0 go test -v -coverprofile=c.out ./...)

# Linux only; populate env.sh with the hardware parameters
integration-test:	##run integration test (Linux only)
	( cd tests/integration; sh run.sh )

release:	##generate and push release
	BUILD_TYPE="R" $(MAKE) clean build docker push
