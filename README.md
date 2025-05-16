# COSI Driver

**Repository for COSI Driver for Dell Container Storage Modules**

## Description
COSI Driver is part of the [CSM (Container Storage Modules)](https://github.com/dell/csm) open-source suite of Kubernetes storage enablers for Dell products. COSI Driver is a Container Object Storage Interface (COSI) driver that provides support for provisioning persistent storage using Dell storage array. 

<!-- It supports CSI specification version 1.5. -->

This project may be compiled as a stand-alone binary using Golang that, when run, provides a valid COSI endpoint. It also can be used as a precompiled container image.

## Table of Contents

* [Code of Conduct](https://github.com/dell/csm/blob/main/docs/CODE_OF_CONDUCT.md)
* [Maintainer Guide](https://github.com/dell/csm/blob/main/docs/MAINTAINER_GUIDE.md)
* [Committer Guide](https://github.com/dell/csm/blob/main/docs/COMMITTER_GUIDE.md)
* [Contributing Guide](https://github.com/dell/csm/blob/main/docs/CONTRIBUTING.md)
* [List of Adopters](https://github.com/dell/csm/blob/main/docs/ADOPTERS.md)
* [Support](#support)
* [Security](https://github.com/dell/csm/blob/main/docs/SECURITY.md)
* [Building](#building)
* [Runtime Dependecies](#runtime-dependencies)
* [Documentation](#documentation)

## Support
For any issues, questions or feedback, please contact [Dell support](https://www.dell.com/support/incidents-online/en-us/contactus/product/container-storage-modules).

## Building
This project is a Go module (see golang.org Module information for explanation). 
The dependencies for this project are in the go.mod file.

To build the source, execute `make build`.

To build an image, execute `make podman`.

To run unit tests, execute `make unit-test`.

Default parameters for building an image are defined in overrides.mk. Run `make -f overrides.mk overrides-help` to display current values.

<!-- You can run an integration test on a Linux system by populating the file `env.sh` with values for your Dell PowerMax systems and then run "`make integration-test`". -->

<!-- ## Runtime Dependencies -->
<!-- Both the Controller and the Node portions of the driver can only be run on nodes which have network connectivity to a “`Unisphere for PowerMax`” server (which is used by the driver). 

If you are using ISCSI, then the Node portion of the driver can only be run on nodes that have the iscsi-initiator-utils package installed. -->

## Documentation
For more detailed information on the driver, please refer to [Container Storage Modules documentation](https://dell.github.io/csm-docs/).

*NOTICE*: the COSI driver code is linted with the phenomenal `golangci-lint`. For a detailed list 
of the linters used and their configuration, please refer to the `.golangci.yml` in the root of the project.
