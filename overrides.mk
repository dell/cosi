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
# Image from https://catalog.redhat.com/software/containers/ubi8/ubi-minimal/5c359a62bed8bd75a2c3fba8
DEFAULT_BASEIMAGE=registry.access.redhat.com/ubi8/ubi-micro
# DIGEST is the digest for version 8.7-1085 of ubi-micro.
DEFAULT_DIGEST=sha256:c530ca8805f8cae3e469e12d985531f8d2ac259e150f935abf80339ea055ccbe
# GOVERSION is a build version for driver.
DEFAULT_GOVERSION=1.19
# REGISTRY in which COSI-Driver image resides.
DEFAULT_REGISTRY=sample_registry
# IMAGENAME is COSI-Driver image name.
DEFAULT_IMAGENAME=cosi-driver
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