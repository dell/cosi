# COSI Driver

**Repository for COSI Driver for Dell Container Storage Modules**

## Description
COSI Driver is part of the [CSM (Container Storage Modules)](https://github.com/dell/csm) open-source suite of Kubernetes storage enablers for Dell products. COSI Driver is a Container Object Storage Interface (COSI) driver that provides support for provisioning persistent storage using Dell storage array.

<!-- It supports COSI specification version v0.2.1. -->

This project may be compiled as a stand-alone binary using Golang that, when run, provides a valid COSI endpoint. It also can be used as a precompiled container image.

## Building
This project is a Go module (see golang.org Module information for explanation).
The dependencies for this project are in the go.mod file.

To build the source, execute `make vendor build`.

To build an image, execute `make push`.

To run unit tests, execute `make vendor unit-test`.

Default parameters for building an image are defined in overrides.mk. Run `make -f overrides.mk overrides-help` to display current values.

## Documentation
For more detailed information on the driver, please refer to [Container Storage Modules documentation](https://dell.github.io/csm-docs/).
