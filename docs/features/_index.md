---
title: "Features"
linkTitle: "Features" 
weight: 4
description: Description of COSI Driver features
---

## ObjectScale

| Area              | Core Features          |  Implementation level  |     Status      | Details                                                                                     |
|:------------------|:-----------------------|:----------------------:|:---------------:|---------------------------------------------------------------------------------------------|
| Provisioning      | _Create Bucket_        | Minimum Viable Product |     ✅ Done      | Bucket is created using default settings.                                                   |
|                   |                        | Advanced provisioning  | 📝 Design draft | Extra (non-default) parameters for bucket provisioning are controlled from the BucketClass. |
|                   | _Delete Bucket_        | Minimum Viable Product |     ✅ Done      | Bucket is deleted.                                                                          |
| Access Management | _Grant Bucket Access_  | Minimum Viable Product |    🚧 Doing     | Full access is granted for given bucket.                                                    |
|                   |                        |  Advanced permissions  | 📝 Design draft | More control over permission is done through BucketAccessClass.                             |
|                   | _Revoke Bucket Access_ | Minimum Viable Product |     ⌛ Todo      | Access is revoked.                                                                          |

## ECS

| Area              | Core Features          |  Implementation level  |     Status      | Details                                                                                     |
|:------------------|:-----------------------|:----------------------:|:---------------:|---------------------------------------------------------------------------------------------|
| Provisioning      | _Create Bucket_        | Minimum Viable Product | 📝 Design draft | Bucket is created using default settings.                                                   |
|                   |                        | Advanced provisioning  | 📝 Design draft | Extra (non-default) parameters for bucket provisioning are controlled from the BucketClass. |
|                   | _Delete Bucket_        | Minimum Viable Product | 📝 Design draft | Bucket is deleted.                                                                          |
| Access Management | _Grant Bucket Access_  | Minimum Viable Product | 📝 Design draft | Full access is granted for given bucket.                                                    |
|                   |                        |  Advanced permissions  | 📝 Design draft | More control over permission is done through BucketAccessClass.                             |
|                   | _Revoke Bucket Access_ | Minimum Viable Product | 📝 Design draft | Access is revoked.                                                                          |
