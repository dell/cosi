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
# limitations under the License

# BASEIMAGE is a base image for final COSI-Driver container.
ARG BASEIMAGE
# DIGEST is a hash-version of a used BASEIMAGE.
ARG DIGEST
# GOVERSION is a Go version used for bulding driver.
ARG GOVERSION

# First stage: building binary of the driver.
FROM golang:${GOVERSION} as builder
WORKDIR /cosi-driver
COPY . /cosi-driver/
RUN make build

# Second stage: building final environment for running the driver.
FROM ${BASEIMAGE}@${DIGEST} AS final
RUN echo 'cosi:x:625:625:cosi:/cosi:/bin/nologin' >> /etc/passwd
USER cosi
WORKDIR /cosi
CMD [ "whoami" ]
WORKDIR /cosi-driver
COPY --from=builder /cosi-driver/build/cosi-driver .
ENTRYPOINT ["./cosi-driver"]
