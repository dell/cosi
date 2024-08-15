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
# GOIMAGE is a Go version used for bulding driver.
ARG GOIMAGE

# First stage: building binary of the driver.
FROM $GOIMAGE as builder

WORKDIR /workspace

# Copy the Go Modules manifests.
COPY go.mod go.mod
COPY go.sum go.sum

# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer.
# FIXME: this should be added after we remove dependency on private goobjectscale
# RUN go mod download
COPY vendor/ vendor/

# Copy the go source.
COPY overrides.mk overrides.mk
COPY Makefile Makefile
COPY cmd/main.go cmd/main.go
COPY pkg/ pkg/

# Build.
RUN make build

# Second stage: building final environment for running the driver.
FROM ${BASEIMAGE} AS final

WORKDIR /dell

COPY --from=builder /workspace/build/cosi /dell/cosi

# Create a non-root user and set permissions on the binary.
RUN echo "cosi:*:1001:cosi-user" >> /etc/group && \
    echo "cosi-user:*:1001:1001::/cosi:/bin/false" >> /etc/passwd && \
    chown 1001:1001 /dell/cosi && \
    chmod 0550 /dell/cosi && \
    mkdir -p /var/lib/cosi /cosi && \
    chown -R 1001:1001 /var/lib/cosi /cosi

# Run as non-root
USER cosi-user

# Set volume mount point for app socket and config file.
VOLUME [ "/var/lib/cosi", "/cosi" ]

# Disable healthcheck.
HEALTHCHECK NONE

# Set the entrypoint.
ENTRYPOINT ["/dell/cosi"]
CMD []
