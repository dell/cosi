---
title: "COSI Driver"
linkTitle: "COSI Driver"
description: About Dell Technologies (Dell) COSI Driver 
weight: 3
---

The COSI Driver by Dell implements an interface between [COSI](https://github.com/container-object-storage-interface/container-object-storage-interface.github.io/tree/master/docs). It is a plug-in that is installed into Kubernetes to provide object storage using Dell storage systems.

Dell COSI Driver is a multi-backend driver, meaning that it can connect to multiple Object Storage Platform (OSP) Instances and provide access to them using the same COSI interface.

## Features and capabilities

### COSI Driver Capabilities

| Features               | ObjectStore | ECS | PowerScale |
|------------------------|:-----------:|:---:|:----------:|
| Bucket Creation        |     yes     | no  |     no     |
| Bucket Deletion        |     yes     | no  |     no     |
| Bucket Access Granting |     yes     | no  |     no     |
| Bucket Access Revoking |     yes     | no  |     no     |

## Bucket Lifecycle Workflow

1. Create Bucket &rarr; Delete Bucket
1. Create Bucket &rarr; Grant Access &rarr; Revoke Access &rarr; Delete Bucket
