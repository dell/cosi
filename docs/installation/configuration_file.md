---
title: Configuration File
linktitle: Configuration File
weight: 1
Description: Description of configuration file for ObjectScale
---

## Dell COSI Driver Configuration Schema

This configuration file is used to specify the settings for the Dell COSI Driver, which is responsible for managing connections to the Dell ObjectScale platform. The configuration file is written in YAML format and based on the JSON schema and adheres to it's specification.

YAML files can have comments, which are lines in the file that begin with the `#` character. Comments can be used to provide context and explanations for the data in the file, and they are ignored by parsers when reading the YAML data.

### Properties

- **cosi-endpoint** (_string_): Path to the COSI socket. Default value is `unix:///var/lib/cosi/cosi.sock`.
- **log-level** (_string_): Defines how verbose logs should be. Valid values are: `"fatal"`, `"error"`, `"warning"`, `"info"`, `"debug"`, and `"trace"`. Default value is `"info"`.
- **connections** (_array_): List of connections to object storage platforms that is used for object storage provisioning. (**required field**)
  - **objectscale** (_object_): Configuration specific to the Dell ObjectScale platform.
    - **id** (_string_): Default, unique identifier for the single connection. All hyphens '`-`' are replaced with underscores '`_`', and may cause issues, thus it is adviced to not use them. (**required field**)
    - **credentials** (_object_): Credentials used for authentication to object storage provider. (**required field**)
      - **username** (_string_): Base64 encoded username. (**required field**)
      - **password** (_string_): Base64 encoded password. (**required field**)
    - **objectscale-gateway** (_string_): Endpoint of the ObjectScale Gateway Internal service. (**required field**)
    - **objectstore-gateway** (_string_): Endpoint of the ObjectScale ObjectStore Management Gateway service. (**required field**)
    - **region** (_string_): Identity and Access Management (IAM) API specific field, points to the region in which object storage provider is installed. 
    - **protocols** (_object_): Protocols supported by the connection. (**required field**)
      - **s3** (_object_): S3 configuration. (**required field**)
        - **endpoint** (_string_): Endpoint of the S3 service. (**required field**)
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
      insecure: false
      client-cas: |-
        LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURJVENDQWdtZ0F3SUJBZ0lJQ3BVOEZFQXNS
        MDh3RFFZSktvWklodmNOQVFFTEJRQXdGVEVUTUJFR0ExVUUKQXhNS2EzVmlaWEp1WlhSbGN6QWVG
        dzB5TWpFd01ETXhNek0xTlRaYUZ3MHlNekV3TURNeE16TTFOVGhhTURReApGekFWQmdOVkJBb1RE
        bk41YzNSbGJUcHRZWE4wWlhKek1Sa3dGd1lEVlFRREV4QnJkV0psY201bGRHVnpMV0ZrCmJXbHVN
        SUlCSWpBTkJna3Foa2lHOXcwQkFRRUZBQU9DQVE4QU1JSUJDZ0tDQVFFQStkQmVqdWFJczZpRm5n
        OFQKRUFldkkwMzdEVDBTZ2E3WEdISmUvaGdSYTloWm1HSlA4SHhKY0NvZWJMYi9iTUNBa09zbjJN
        Q0VFamV4U0NaagpUMjN2Q2toVjVIL2gxdTluWVRidTRjV25sbUl0RG53eU1tMS96QW5hZkd4clVC
        VDVMaENRS2IvZkNzK280M0gvCkZRZHdta2wySXFBNmRHL2ljcGszZGl2c0M3WlRxL2p5VmhUdHN4
        Z2l4YjFzZTc4Y2xlQmJoYzVocG9qN2I1MHcKQ1RjZWh1NDhWVjE3VjhkRFAraGc3T3NPOWhXcGdl
        RmZrb1JnQm1DSjNvRlc0N1dDQzRaWUQzd3hsTlBJOUZ4LwptSG9nbHFlU3J1Sjc4UUM4eWprOUpB
        eFJ1ZlRWTExKSnIwUWJUc2I5dC96NlVZcU91VGVsTVd2UUVtQ1RaZzg0CmtmUWc5d0lEQVFBQm8x
        WXdWREFPQmdOVkhROEJBZjhFQkFNQ0JhQXdFd1lEVlIwbEJBd3dDZ1lJS3dZQkJRVUgKQXdJd0RB
        WURWUjBUQVFIL0JBSXdBREFmQmdOVkhTTUVHREFXZ0JTcjRPQW5jZDJpQU0wVVMwa0NVUDNrNWZV
        ZgoxakFOQmdrcWhraUc5dzBCQVFzRkFBT0NBUUVBYXVTemRCckVTVkkzWHhIZ3FhT1A3Y1UyT1Vs
        MytXV21SVk9DCkVETWtJVitoc3d4NmxUZmt1SXF3ZnZDaWVIZURoMXh5T3ZkYmNGdDl4WENIalor
        bStucDN1cnFXdm9yYjNRL2kKSGRXUEtqTnVzNmcrbTExczI3cy9CdFMrVjdxWkNINENiek5hS0hD
        dzZobmlhYk9jTnFjSlpXYXgxdFUwbzkvZApMWm8vYmlCelRtajN3dG41UVhxMG04NHdrOHl6Mzh2
        WW5NamRzZjRpNlIzcGFwVVVEamJMK3N2QS8zdXh3ZG5vCktYbzR0TjU3dGRHaEYwQ1ZnZjJldEIr
        TGtjTVpOTHlSQzcrTVdXYTJRSnNmV1JoYzJyeTBBd1lKejhjNFdXWkgKblUwb3VTc1lUcHV4NmNL
        b1U4cHZNREV2Y0pBOVhyc214SnRUeEJTQjQ1NWRpNUwvNEE9PQotLS0tLUVORCBDRVJUSUZJQ0FU
        RS0tLS0tCg==
```
