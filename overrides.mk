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

# overrides.mk file
# This file, included from the Makefile, will overlay default values with environment variables:
# BASEIMAGE, DIGEST, GOVERSION, REGISTRY, IMAGENAME, IMAGETAG.

# DEFAULT values
# Image from https://catalog.redhat.com/software/containers/ubi9/ubi-micro/615bdf943f6014fa45ae1b58
DEFAULT_BASEIMAGE=registry.access.redhat.com/ubi9/ubi-micro
# DIGEST is the digest for version 9.2-13 of ubi-micro.
DEFAULT_DIGEST=sha256:630cf7bdef807f048cadfe7180d6c27eb3aaa99323ffc3628811da230ed3322a
# GOVERSION is a build version for driver.
DEFAULT_GOVERSION:=$(shell sed -En 's/^go (.*)$$/\1/p' go.mod)
# REGISTRY in which COSI image resides.
DEFAULT_REGISTRY=sample_registry
# IMAGENAME is COSI image name.
DEFAULT_IMAGENAME=cosi
# IMAGETAG by default is latest commit SHA.
DEFAULT_IMAGETAG:=$(shell git rev-parse HEAD)

# set the BASEIMAGE if needed
ifeq ($(BASEIMAGE),)
export BASEIMAGE="$(DEFAULT_BASEIMAGE)"
endif

# set the IMAGEDIGEST if needed
ifeq ($(DIGEST),)
export DIGEST="$(DEFAULT_DIGEST)"
endif

# set the GOVERSION if needed
ifeq ($(GOVERSION),)
export GOVERSION="$(DEFAULT_GOVERSION)"
endif

# set the REGISTRY if needed
ifeq ($(REGISTRY),)
export REGISTRY="$(DEFAULT_REGISTRY)"
endif

# set the IMAGENAME if needed
ifeq ($(IMAGENAME),)
export IMAGENAME="$(DEFAULT_IMAGENAME)"
endif

#set the IMAGETAG if needed
ifeq ($(IMAGETAG),) 
export IMAGETAG="$(DEFAULT_IMAGETAG)"
endif

# target to print some help regarding these overrides and how to use them
overrides-help:
	@echo
	@echo "The following environment variables can be set to control the build"
	@echo
	@echo "GOVERSION   - The version of Go to build with, default is: $(DEFAULT_GOVERSION)"
	@echo "              Current setting is: $(GOVERSION)"
	@echo "BASEIMAGE   - The base container image to build from, default is: $(DEFAULT_BASEIMAGE)"
	@echo "              Current setting is: $(BASEIMAGE)"
	@echo "IMAGEDIGEST - The digest of baseimage, default is: $(DEFAULT_DIGEST)"
	@echo "              Current setting is: $(DIGEST)"
	@echo "REGISTRY    - The registry to push images to, default is: $(DEFAULT_REGISTRY)"
	@echo "              Current setting is: $(REGISTRY)"
	@echo "IMAGENAME   - The image name to be built, defaut is: $(DEFAULT_IMAGENAME)"
	@echo "              Current setting is: $(IMAGENAME)"
	@echo "IMAGETAG    - The image tag to be built, default is an empty string which will determine the tag by examining annotated tags in the repo."
	@echo "              Current setting is: $(IMAGETAG)"
