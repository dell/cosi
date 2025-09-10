# COSI Driver

**Repository for COSI Driver for Dell Container Storage Modules**

## Description
COSI Driver is part of the [CSM (Container Storage Modules)](https://github.com/dell/csm) open-source suite of Kubernetes storage enablers for Dell products. COSI Driver is a Container Object Storage Interface (COSI) driver that provides support for provisioning persistent storage using Dell storage array. 

<!-- It supports CSI specification version 1.5. -->

# :lock: **Important Notice**
Starting with the release of **Container Storage Modules v1.16.0**, this repository will no longer be maintained as an open source project. Future development will continue under a closed source model. This change reflects our commitment to delivering even greater value to our customers by enabling faster innovation and more deeply integrated features with the Dell storage portfolio.<br>
For existing customers using Dell’s Container Storage Modules, you will continue to receive:
* **Ongoing Support & Community Engagement**<br>
       You will continue to receive high-quality support through Dell Support and our community channels. Your experience of engaging with the Dell community remains unchanged.
* **Streamlined Deployment & Updates**<br>
        Deployment and update processes will remain consistent, ensuring a smooth and familiar experience.
* **Access to Documentation & Resources**<br>
       All documentation and related materials will remain publicly accessible, providing transparency and technical guidance.
* **Continued Access to Current Open Source Version**<br>
       The current open-source version will remain available under its existing license for those who rely on it.

Moving to a closed source model allows Dell’s development team to accelerate feature delivery and enhance integration across our Enterprise Kubernetes Storage solutions ultimately providing a more seamless and robust experience.<br>
We deeply appreciate the contributions of the open source community and remain committed to supporting our customers through this transition.<br>

For questions or access requests, please contact the maintainers via [Dell Support](https://www.dell.com/support/kbdoc/en-in/000188046/container-storage-interface-csi-drivers-and-container-storage-modules-csm-how-to-get-support).

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

