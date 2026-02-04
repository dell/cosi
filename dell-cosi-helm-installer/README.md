# Helm Installer for Dell COSI Storage Providers

## Description

This directory provides scripts to install, upgrade, uninstall the COSI driver, and to verify the Kubernetes environment. This includes the driver for:
* [COSI](https://github.com/dell/cosi)

## Dependencies

Installing the Dell COSI Driver requires a few utilities to be installed on the system running the installation.

| Dependency    | Usage  |
| ------------- | ----- |
| `kubectl`     | Kubectl is used to validate that the Kubernetes system meets the requirements of the driver. |
| `helm`        | Helm v3 is used as the deployment tool for Charts. See, [Install Helm 3](https://helm.sh/docs/intro/install/) for instructions to install Helm 3. |

In order to use these tools, a valid `KUBECONFIG` is required. Ensure that either a valid configuration is in the default location or that the `KUBECONFIG` environment variable points to a valid configuration before using these tools.

## Capabilities

This project provides the following capabilities, each one is discussed in detail later in this document.

* Install a driver. When installing a driver, options are provided to specify the target namespace as well as options to control the types of verifications to be performed on the target system.
* Upgrade a driver. Upgrading a driver is an effective way to either deploy a new version of the driver or to modify the parameters used in an initial deployment.
* Uninstall a driver. This removes the driver and any installed storage classes.
* Verify a Kubernetes system for suitability with a driver. These verification steps include verifiying version compatibility, namespace availability, and existence of required secrets. 


Most of these usages require the creation/specification of a values file. These files specify configuration settings that are passed into the driver and configure it for use. To create one of these files, the following steps should be followed:
1. Download a template file for the driver to a new location, naming this new file is at the users discretion. The template files are always found at `https://raw.githubusercontent.com/dell/helm-charts/refs/heads/main/charts/cosi/values.yaml`
2. Edit the file such that it contains the proper configuration settings for the specific environment. These files are yaml formatted so maintaining the file structure is important.

For example, to create a values file for the COSI driver the following steps can be executed
```
# cd to  the installation script directory
cd dell-cosi-helm-installer

# download the template file
 wget -O my-cosi-settings.yaml https://raw.githubusercontent.com/dell/helm-charts/refs/heads/main/charts/cosi/values.yaml

# edit the newly created values file
vi my-cosi-settings.yaml
```

These values files can then be archived for later reference or for usage when upgrading the driver.


### Install A Driver

Installing a driver is performed via the `cosi-install.sh` script. This script requires a few arguments: the target namespace and the user created values file. By default, this will verify the Kubernetes environment and present a list of warnings and/or errors. Errors must be addressed before installing, warning should be examined for their applicability. For example, in order to install the COSI driver into a namespace called "cosi", the following command should be run:
```
./cosi-install.sh --namespace cosi --values ./my-cosi-settings.yaml
```

For usage information:
```
[dell-cosi-helm-installer]# ./cosi-install.sh -h
Help for ./cosi-install.sh

Usage: ./cosi-install.sh options...
Options:
  Required
  --namespace[=]<namespace>                Kubernetes namespace containing the COSI driver
  --values[=]<values.yaml>                 Values file, which defines configuration values
  Optional
  --release[=]<helm release>               Name to register with helm, default value will match the driver name
  --upgrade                                Perform an upgrade of the specified driver, default is false
  --skip-verify                            Skip the kubernetes configuration verification to use the COSI driver, default will run verification
  -h                                       Help
```

### Upgrade A Driver

Upgrading a driver is very similar to installation. The `cosi-install.sh` script is run, with the same required arguments, along with a `--upgrade` argument. For example, to upgrade the previously installed COSI driver, the following command can be supplied:

```
./cosi-install.sh --namespace cosi --values ./my-cosi-settings.yaml --upgrade
```

For usage information:
```
[dell-cosi-helm-installer]# ./cosi-install.sh -h
Help for ./cosi-install.sh

Usage: ./cosi-install.sh options...
Options:
  Required
  --namespace[=]<namespace>                Kubernetes namespace containing the COSI driver
  --values[=]<values.yaml>                 Values file, which defines configuration values
  Optional
  --release[=]<helm release>               Name to register with helm, default value will match the driver name
  --upgrade                                Perform an upgrade of the specified driver, default is false
  --skip-verify                            Skip the kubernetes configuration verification to use the COSI driver, default will run verification
  -h                                       Help
```

### Uninstall A Driver

To uninstall a driver, the `cosi-uninstall.sh` script provides a handy wrapper around the `helm` utility. The only required argument for uninstallation is the namespace name. To uninstall the COSI driver:

```
./cosi-uninstall.sh --namespace cosi
```

For usage information:
```
[dell-cosi-helm-installer]# ./cosi-uninstall.sh -h
Help for ./cosi-uninstall.sh

Usage: ./cosi-uninstall.sh options...
Options:
  Required
  --namespace[=]<namespace>  Kubernetes namespace to uninstall the COSI driver from
  Optional
  --release[=]<helm release> Name to register with helm, default value will match the driver name
  -h                         Help
```

### Verify A Kubernetes Environment

The `verify.sh` script is run, automatically, as part of the installation and upgrade procedures and can also be run by itself. This provides a handy means to validate a Kubernetes system without meaning to actually perform the installation. To verify an environment, run `verify.sh` with the namespace name and values file options.

```
./verify.sh --namespace cosi --values ./my-cosi-settings.yaml
```

For usage information:
```
[dell-cosi-helm-installer]# ./verify.sh -h
Help for ./verify.sh

Usage: ./verify.sh options...
Options:
  Required
  --namespace[=]<namespace>       Kubernetes namespace to install the COSI driver
  --values[=]<values.yaml>        Values file, which defines configuration values
  Optional
  --release[=]<helm release>      Name to register with helm, default value will match the driver name
  -h                              Help                           Help
```
