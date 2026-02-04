# Offline Installation of Dell COSI Storage Providers

## Description

The `cosi-offline-bundle.sh` script can be used to create a package for the offline installation of Dell COSI storage providers for deployment via Helm.  

This includes the following driver:
* [COSI](https://github.com/dell/cosi)

## Dependencies

Multiple linux based systems may be required to create and process an offline bundle for use.
* One linux based system, with internet access, will be used to create the bundle. This involved the user cloning a git repository hosted on github.com and then invoking a script that utilizes `docker` or `podman` to pull and save container images to file.
* One linux based system, with access to an image registry, to invoke a script that uses `docker` or `podman` to restore container images from file and push them to a registry

If one linux system has both internet access and access to an internal registry, that system can be used for both steps.

Preparing an offline bundle requires the following utilities:

| Dependency            | Usage |
| --------------------- | ----- |
| `docker` or `podman`  | `docker` or `podman` will be used to pull images from public image registries, tag them, and push them to a private registry.  |
|                       | One of these will be required on both the system building the offline bundle as well as the system preparing for installation. |
|                       | Tested version(s) are `docker` 19.03+ and `podman` 1.6.4+
| `git`                 | `git` will be used to manually clone one of the above repos in order to create and offline bundle.
|                       | This is only needed on the system preparing the offline bundle.
|                       | Tested version(s) are `git` 1.8+ but any version should work.

## Workflow

To perform an offline installation of the COSI driver with helm, the following steps should be performed:
1. Build an offline bundle
2. Unpacking an offline bundle and preparing for installation
3. Perform a Helm installation

### Building an offline bundle

This needs to be performed on a linux system with access to the internet as a git repo will need to be cloned, and container images pulled from public registries.

To build an offline bundle, the following steps are needed:
1. Perform a `git clone` of the desired repository. For a Helm based install, the specific driver repo should be cloned. For an Operator based deployment, the Dell CSM Operator repo should be cloned
2. Run the offline bundle script with an argument of `-c` in order to create an offline bundle
  - For Helm installs, the `cosi-offline-bundle.sh` script will be found in the `dell-cosi-helm-installer` directory

The script will perform the following steps:
  - Determine required images by parsing either the driver Helm charts (if run from a cloned COSI Driver git repository) or the Dell CSM Operator configuration files (if run from a clone of the Dell CSM Operator repository)
  - Perform an image `pull` of each image required
  - Save all required images to a file by running `docker save` or `podman save`
  - Build a `tar.gz` file containing the images as well as files required to installer the driver and/or Operator

The resulting offline bundle file can be copied to another machine, if necessary, to gain access to the desired image registry.

### Unpacking an offline bundle and preparing for installation

This needs to be performed on a linux system with access to an image registry that will host container images. If the registry requires `login`, that should be done before proceeding.

To prepare for driver or Operator installation, the following steps need to be performed:
1. Copy the offline bundle file to a system with access to an image registry available to your Kubernetes/OpenShift cluster
2. Expand the bundle file by running `tar xvfz <filename>`
3. Run the `cosi-offline-bundle.sh` script and supply the `-p` option as well as the path to the internal registry with the `-r` option

The script will then perform the following steps:
  - Load the required container images into the local system
  - Tag the images according to the user supplied registry information
  - Push the newly tagged images to the registry
  - Modify the Helm charts or Operator configuration to refer to the newly tagged/pushed images


### Perform a Helm installation

Now that the required images have been made available and the Helm Charts configuration updated, installation can proceed by following the instructions that are documented within the driver repo.
