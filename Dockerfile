# Copyright © 2023-2026 Dell Inc. or its subsidiaries. All Rights Reserved.
# Dell Technologies, Dell, and other trademarks are trademarks of Dell Inc. or its subsidiaries.
# Other trademarks may be trademarks of their respective owners.

ARG BASEIMAGE
ARG GOIMAGE
ARG VERSION="1.0.0"

FROM $GOIMAGE as builder
ARG VERSION

WORKDIR /workspace
COPY go.mod go.mod
COPY go.sum go.sum
COPY vendor/ vendor/
COPY overrides.mk overrides.mk
COPY images.mk images.mk
COPY helper.mk helper.mk
COPY Makefile Makefile
COPY cmd/main.go cmd/main.go
COPY pkg/ pkg/
RUN make build

FROM ${BASEIMAGE} AS final
ARG VERSION

WORKDIR /dell
COPY --from=builder /workspace/build/cosi /dell/cosi

# Create a non-root user and set permissions on the binary.
RUN echo "cosi:*:1001:cosi-user" >> /etc/group && \
    echo "cosi-user:*:1001:1001::/cosi:/bin/false" >> /etc/passwd && \
    chown 1001:1001 /dell/cosi && \
    chmod 0550 /dell/cosi && \
    mkdir -p /var/lib/cosi /cosi && \
    chown -R 1001:1001 /var/lib/cosi /cosi

USER cosi-user

# Set volume mount point for app socket and config file.
VOLUME [ "/var/lib/cosi", "/cosi" ]

LABEL vendor="Dell Technologies" \
    maintainer="Dell Technologies" \
    name="cosi" \
    summary="COSI Driver for Dell Storage Systems" \
    description="COSI Driver for provisioning object storage from Dell Storage Systems" \
    release="1.16.0" \
    version=$VERSION \
    license="Dell CSM Operator Apache License"

COPY licenses /licenses

ENTRYPOINT ["/dell/cosi"]
CMD []
