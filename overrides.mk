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
# GOIMAGE, REGISTRY, IMAGENAME, IMAGETAG.

# DEFAULT values
# GOIMAGE is the version of Go.
DEFAULT_GOIMAGE:=$(shell sed -En 's/^go (.*)$$/\1/p' go.mod)
# REGISTRY in which COSI image resides.
DEFAULT_REGISTRY=docker.io
# IMAGENAME is COSI image name.
DEFAULT_IMAGENAME=dell/cosi
# IMAGETAG by default is latest commit SHA.
DEFAULT_IMAGETAG:=$(shell git rev-parse HEAD)

# set the GOIMAGE if needed
ifeq ($(GOIMAGE),)
export GOIMAGE="$(DEFAULT_GOIMAGE)"
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
	@echo "GOIMAGE     - The version of Go to build with, default is: $(DEFAULT_GOIMAGE)"
	@echo "              Current setting is: $(GOIMAGE)"
	@echo "REGISTRY    - The registry to push images to, default is: $(DEFAULT_REGISTRY)"
	@echo "              Current setting is: $(REGISTRY)"
	@echo "IMAGENAME   - The image name to be built, defaut is: $(DEFAULT_IMAGENAME)"
	@echo "              Current setting is: $(IMAGENAME)"
	@echo "IMAGETAG    - The image tag to be built, default is an empty string which will determine the tag by examining annotated tags in the repo."
	@echo "              Current setting is: $(IMAGETAG)"
