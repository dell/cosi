---
title: Configuration File
linktitle: Configuration File
weight: 1
Description: Description of configuration file for ObjectScale
---

## Dell COSI Driver Configuration Schema

This configuration file is used to specify the settings for the Dell COSI Driver, which is responsible for managing connections to the Dell ObjectScale platform. The configuration file is written in JSON format and adheres to the JSON Schema draft-07 specification.

### Format

Configuration file is written using YAML format. It is a human-readable data serialization language that is commonly used for configuration files and data exchange between programming languages.

YAML files can have comments, which are lines in the file that begin with the `#` character. Comments can be used to provide context and explanations for the data in the file, and they are ignored by parsers when reading the YAML data.

### Properties

- **cosi-endpoint** (_string_): Path to the COSI socket. Default value is `unix:///var/lib/cosi/cosi.sock`.
- **log-level** (_string_): Defines how verbose should logs be. Valid values are: `"fatal"`, `"error"`, `"warning"`, `"info"`, `"debug"`, and `"trace"`. Default value is `"info"`.
- **connections** (_array_): List of connections based on which the sub-drivers are constructed. (**required field**)
  - **objectscale** (_object_): Configuration specific to the Dell ObjectScale platform.
    - **credentials** (_object_): Credentials used for authentication to OSP. (**required field**)
      - **username** (_string_): Base64 encoded username. (**required field**)
      - **password** (_string_): Base64 encoded password. (**required field**)
    - **id** (_string_): Default, unique identifier for the sub-driver/platform. All hyphens '`-`' are replaced with underscores '`_`', and may cause issues, thus it is adviced to not use them. (**required field**)
    - **objectscale-gateway** (_string_): Endpoint of the ObjectScale Gateway Internal service. (**required field**)
    - **objectstore-gateway** (_string_): Endpoint of the ObjectScale ObjectStore Management Gateway service. (**required field**)
    - **region** (_string_): IAM API specific field, points to the region in which ObjectScale system is installed. 
    - **protocols** (_object_): Protocols used by the driver. (**required field**)
      - **s3** (_object_): S3 configuration. (**required field**)
        - **endpoint** (_string_): Endpoint of the ObjectStore S3 service. (**required field**)
    - **tls** (_object_): TLS configuration details. (**required field**)
      - **insecure** (_boolean_): Controls whether a client verifies the server's certificate chain and host name. (**required field**)
      - **client-cas** (_string_): Base64 encoded certificate file.

## Configuration file example

```yaml
cosi-endpoint: unix:///var/lib/cosi/cosi.sock
log-level: info
connections:
- objectscale:
    credentials:
      username: dGVzdHVzZXIK          # testuser
      password: dGVzdHBhc3N3b3JkCg==  # testpassword
    id: example.id # do not include hyphens
    objectscale-gateway: gateway.objectscale.test
    objectstore-gateway: gateway.objectstore.test
    protocols:
      s3:
        endpoint: s3.objectstore.test
    tls:
      insecure: true
```
