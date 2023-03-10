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
#s
# docker makefile, included from Makefile, will build/push images with docker
#

# To set the base image, use the BASEIMAGE environment variable
# To set the image name, use the IMAGENAME environment variable
# To set the image tag, use the IMAGETAG environment variable

# Build with docker
docker:
	@echo "Base Images is set to: $(BASEIMAGE)"
	@echo "Building: $(REGISTRY)/$(IMAGENAME):$(IMAGETAG)"
	docker build -t "$(REGISTRY)/$(IMAGENAME):$(IMAGETAG)" --build-arg BASEIMAGE=$(BASEIMAGE) --build-arg GOVERSION=$(GOVERSION) --build-arg DIGEST=$(DIGEST) .

push:   
	@echo "Pushing: $(REGISTRY)/$(IMAGENAME):$(IMAGETAG)"
	docker push "$(REGISTRY)/$(IMAGENAME):$(IMAGETAG)"
