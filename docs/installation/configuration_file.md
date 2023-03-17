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
      - **root-cas** (_string_): Base64 encoded content of the clients's certificate file.
      - **root-cas** (_string_): Base64 encoded content of the clients's key certificate file.
      - **root-cas** (_string_): Base64 encoded content of the root certificate authority file.

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
      client-cert: |-
        LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUVORENDQWh5Z0F3SUJBZ0lSQU9JSlZ2NnB3
        a0lIK0p1NTNKSEFuam93RFFZSktvWklodmNOQVFFTEJRQXcKRWpFUU1BNEdBMVVFQXhNSGRHVnpk
        QzFqWVRBZUZ3MHlNekF6TVRjeE1qTTJNelphRncweU5EQTVNVGN4TWpRegpNamxhTUJFeER6QU5C
        Z05WQkFNVEJtTnNhV1Z1ZERDQ0FTSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnRVBBRENDCkFRb0Nn
        Z0VCQU5LVFNHeEEyV2RyNmtCR0N3RjY5c1JVZElPV0xqeTUvN3QyRktKWDVVenNyMDlFWW9tS0sr
        bVQKdWF2eWJIMWhsbTYzdG5kb3VFOHFIQnVhYmYvUGIzSlRTQ0twR0NRdHR2NmQzeGc3MHFZVWIx
        cUZKT2o5andlNgpRZW0xb2RIVFpLc0xMc2J1N1Fzei91MGtseUovMHNYcFQ5K2JXK1M0OHMrL3pK
        dHNDR21SdVhlRjE2Y1FqOWErCkFFejNqVzhrdExMYi9nS25GWGRSS2FiY2RWLzNzN2RLNWx0SXpS
        ZlRvUWw0bzBpckpOa3Z4eXIrYUtMMTR4NUQKc3g2Wm9DUHJhRFYrWWlRS0ZSenFjQ1RYcWdRb3BY
        LzFINFRMV3RkeG14M25IdmhZdzB0VlBZSXZsa245NmpJUwpKdVE2K1VMbVAzZDNzNWJadlhQeUZD
        bENKSENxaWZNQ0F3RUFBYU9CaFRDQmdqQU9CZ05WSFE4QkFmOEVCQU1DCkE3Z3dIUVlEVlIwbEJC
        WXdGQVlJS3dZQkJRVUhBd0VHQ0NzR0FRVUZCd01DTUIwR0ExVWREZ1FXQkJTRWVIOTEKVnBhdDlV
        SWlrRUdkc1ljdUI2dWxOakFmQmdOVkhTTUVHREFXZ0JSbDk4cG1valVUQ3RZb3phTDl6L0hSYkhS
        RwpkREFSQmdOVkhSRUVDakFJZ2daamJHbGxiblF3RFFZSktvWklodmNOQVFFTEJRQURnZ0lCQUQv
        TnZVNWRSajlHCmMzYzVVQ3VLcDI4U2lacjAySE40M091WU5QVlk4L1c5QnZUSk5yMXRPRDFscnhE
        eFMzTkpVdzdGaTNidmU5enMKSzA0a09peUxpVjRLd0g2eitpVm8xZU9GUzJLd1BRaGxsaDlobVBB
        dXZ4Zm5Fd2k2ZEdXZm5nNExmQ1FvbXFkTgpmbkFCODJBbTViZTBubGJvaGdLcFJUWnVBZjR4dVY4
        SWxlQ1pjVHdFL1hBbERhNVhHaDNvWlE3REYrQnFLSkNUCk1pYS9MT0JPYXRoRVh5ZGJmbndOUUhy
        UWlQZzk4c2NMc3FTZEFQMFNGYjMrMmdscFJZT1JrQlFvOWRoa1pGZXkKc2tUakVhbk9YaUhqWldq
        aXZRS2Z2WEUvK1l2eGpCcEJqREE2NnYyeUgzSlJqZEM5ZTR2cnE2R0t6VXZML3ltOQpVOGdVWnho
        L2ZmeFp4TVA5UmxXajQ0U1NGUVpZNGxUNFF5U2lteFpGdVBTamwzV29QME12UHVvUzFUUzhQUk5s
        CnVGeXBVell5SEtlbHpLUnRJZmlnWG9XQi9uR2hSV0RMN2FZS0xYZWRIU0ZrdXBmZm9YM1hHQThM
        ZVAwQ01PaEsKUUJaUkxIeXU0VjhvRG1lakFIcFoyVjlpY2E1emtmcnJWVXFvSzF1VjYvdHd3cEZG
        WDErN0w1bk0ybDJDQWxvegpaVHFUZzNCdVdYd2VkYzZQbkpuU2xQSDNadFhqcGFJUWhXdU85TUlG
        WFVtVFBlSkZ2WGxKeWRsdUxtMlQzanVqCldiVENGcEhyMXBrMGk3K1J4ZVRBcFY0RTk2S09DOXEw
        ZGREOG1waTM0cnkyZjFmQ2RZekhQM0s4bW5od3BPWmkKaG1Xd3VWVDV3em5kVWVBRGNWYUY2UlhU
        UENKSElLd24KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=
      client-key: |-
        LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFcFFJQkFBS0NBUUVBMHBOSWJFRFpa
        MnZxUUVZTEFYcjJ4RlIwZzVZdVBMbi91M1lVb2xmbFRPeXZUMFJpCmlZb3I2Wk81cS9Kc2ZXR1di
        cmUyZDJpNFR5b2NHNXB0Lzg5dmNsTklJcWtZSkMyMi9wM2ZHRHZTcGhSdldvVWsKNlAyUEI3cEI2
        YldoMGROa3F3c3V4dTd0Q3pQKzdTU1hJbi9TeGVsUDM1dGI1TGp5ejcvTW0yd0lhWkc1ZDRYWApw
        eENQMXI0QVRQZU5ieVMwc3R2K0FxY1ZkMUVwcHR4MVgvZXp0MHJtVzBqTkY5T2hDWGlqU0tzazJT
        L0hLdjVvCm92WGpIa096SHBtZ0krdG9OWDVpSkFvVkhPcHdKTmVxQkNpbGYvVWZoTXRhMTNHYkhl
        Y2UrRmpEUzFVOWdpK1cKU2YzcU1oSW01RHI1UXVZL2QzZXpsdG05Yy9JVUtVSWtjS3FKOHdJREFR
        QUJBb0lCQUJFSVVzSlcySDd5RHFlVwpRc3VpMjVUejA5elU1L2FIZ1BUenp5VjJnSmloU0dqYitq
        QnYyYTl5QUlHMUFTdC9Ha0RvWVR6MVhuc2d4OWMvCnZZZ0VpbG92L0ZTNVlyZUNieHZYUHpWaG1W
        OVBwZFlua04yN3JMY09UTWlQcFlBb1hpc3JvMlA1N1hpTGd5SkIKWkd3bzlLNkhlYXQza0k1R20z
        Vk1hVXRsQ0tVcE84cUwzcEZ4S1AwMVVwbGh6ZjhMbXJpTUJQMDlxdFFJejBydQpiR1l5eUdVdk9a
        a0RKZFJycmlSWGJWK0RNMFlmbVpqU1Q4aEI0UDlsOEhwMEZRNUp2TWVGREpzRFFaZjVBZnJmClFI
        WE55SlFUeTNTeXJ1bGd5N0p4MGY1T2JpVWRMRWViQVRpN3VLR3Y5UEZRRUJmSzdFdE4vZ1ZibGsx
        MzRzNUIKWEhkNXU1a0NnWUVBNDBVMjhONko4QXIwY2puYnNLUUJtOGhURWlJSjk3TEJPOU5kOTlJ
        M1dJYklZVzIzVE5wVwo0M2R4K1JHelA4eVMzYzZhN00wbzR1dUl6TXFDSkV3cVNJUjAvVGZaWWdx
        cGtwcFZPalp2VFdCUDFtSUlKUFpwCll1SFk0UVRJdkdhcVFNNnFWQXA4MW9YdXoxTmNmQWpTLzNJ
        Z1BWdGVZeDNKd0pmNWVqenZQclVDZ1lFQTdUSEwKR3VCTWpqTWVhaWk1ZU1sU1BndkYxMHJISUs0
        RzZCZUJDTFFXU2ViNmNOT2x2a1RaOTNqdlFiWko1L3JBTGNWNgpaTVdqbWY5Tkl0NWdDdyt2K2dM
        Qm9BZXM3WEk2K2Rpdk1DYXE0dUFmWkhJWjBYbXpIOGx1a0o5ZUtyK2NyR2FzClNhWkdKRnlyQTZz
        WGdOc1ZJUm85RkFsR3V1dGZnd2hSUmo1eFp3Y0NnWUVBZ241MWcyeGtDMTVlNlU5clkwdG8KV1Fo
        M0dreE5LTnFNdFVzeUExL0N3NlB3WG5EZTlOUFJYQjV6WkszVEhHamNVMXVUL1MvM3NBUEpzcno4
        YU5jSwoyRVNsMzljM2pHSE82QXlScnpFZVMzRm5waEwzMWpGZVpaYUVMdi9PT3M5QUpxSURqdW5P
        c0dhS3JxU1F6KzlKCko3OWgzNWtjNHhCeGpaSTFmd2lKM3BrQ2dZRUFwUnBOMkExYy9IWlVxMnho
        ZmRRVXJSK2d2TFZPV2s4SWU3RXcKbmhCTW0zQnR6dTlqcFVkanVVQ3l1YmpiUk9CanVQaUdzM0pt
        NktDdTNxQ1BsZU43aUxrMmNlQWwzTG53bDB6ZQoxTk4xaTZxWjcxOEUzYXlxcEd1ZnpJZENFdHVC
        Z1BlTzRVMGQ4ZDJYSkZ5SlphWVoxUXJnalB2UUFmZ29hWnIyCmg4Q2JTeTBDZ1lFQW1VQ3BqR0JW
        MGNpVnlmUXNmOGdsclNOdWx6NzBiaVJWQzVSeno0dVJEMkhsYVM2eC8wc0IKQzltSUhpdWgwR0Zp
        dEVFRlg4TzdlZ1ppNWJKMGFuQWYyakk1R1RnTjJOYzFpVlZnWldxcHh2aXpuckpKcENSYgpaejB0
        M2thTkkyNjg0WTNxS2JxeG8ramRNK05hMG1qd2ErTEFOcEdCUDNwb2c0RHJ4eTNNSFdZPQotLS0t
        LUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQo=
      root-cas: |-
        LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUU1RENDQXN5Z0F3SUJBZ0lCQVRBTkJna3Fo
        a2lHOXcwQkFRc0ZBREFTTVJBd0RnWURWUVFERXdkMFpYTjAKTFdOaE1CNFhEVEl6TURNeE56RXlN
        ek16TTFvWERUSTBNRGt4TnpFeU5ETXpNRm93RWpFUU1BNEdBMVVFQXhNSApkR1Z6ZEMxallUQ0NB
        aUl3RFFZSktvWklodmNOQVFFQkJRQURnZ0lQQURDQ0Fnb0NnZ0lCQU9oUmc1Um95UXdxCmVtQ1VN
        TDU3cXVLSXJjMWZXdGdlSGRpbVRSamFsVERQMStqYUhGeG56d2M1MTRwOWNLNzcxRWZ2bDRjZW9Q
        VWsKWnRhNSsxckRxdlBkd25BMnE2TXI5cFB2aWQyRkRiZVZPdXNIaHNQSG1kMDVxa1pnNGNXUGdp
        eXlSM3BmNTF0bApVYkxyNU1tL0FIK0JvRHVMbFo1UG5SVUw1b0hFd1hQa3BXc0UyMXJDc2xSdmJv
        WWZJYlplUzlsOHhlYURMVmdDCk53UmFHRjgxTFpoZjVrTDA0SXJUV0dETzdlbVF0S2tpN0dSZ1Ex
        bHIxRHR3SXZpa0puakhBeEJiOTJ3WDN1WnoKcGdMQksxU2RsUlY1bjY2VTZtUklzMGo1MkVyTG1h
        TDdUSHJxRVZHRXNvczFIbEZFQ2NJMlNhQjZZdmltaTdZawpmT1lOS2NPaE5BcXlXcWhlUERHQ0dq
        d3l4RHR3OWN2Z2FJSTlTOFFUa2w5Z1JiL056dFlMREptejlEYXZiRWNjCjRDelZBdUVmdUVtWUNi
        aFRrUVUyWitZczlKdXgwdmc4WXFFTExlRzlNZHc1cmZJQkkwNmRMRDVkU0JUVFc1Y08KYjRNN0h1
        ODhrZUdIWnlNZXU2cVMyR2czUUFTVEM3RkpFcWFYTkRDc095aCs2Uk14UnkyZy9idEZMRm5VdmlG
        QQo0NktKZHk0QWVjOEpXVkc1OFlLYkd2QlJrekkzY1BNWE1oWFpDS3pZb0tnUWoxMnFOMWM0SkVp
        TUFPK2F2ZW9RCjB0dnJmd3MxMlF3d3ZIZm40SCtYVnlDZGpMcDE5dlhlY0FSRFJyaGlkRW1CbEFD
        cVJVdTFLSGhzejZ2TmxzUzIKSlZiWU9BYW5ISzYzNzdYT211OUthL2x1TmxSVDdmckxBZ01CQUFH
        alJUQkRNQTRHQTFVZER3RUIvd1FFQXdJQgpCakFTQmdOVkhSTUJBZjhFQ0RBR0FRSC9BZ0VBTUIw
        R0ExVWREZ1FXQkJSbDk4cG1valVUQ3RZb3phTDl6L0hSCmJIUkdkREFOQmdrcWhraUc5dzBCQVFz
        RkFBT0NBZ0VBNUVxL09ocGs0MUxLa3R2ajlmSWNFWXI5Vi9QUit6Z1UKQThRTUtvOVBuZSsrdW1N
        dEZEK1R1M040b1lxV0srTmg0ajdHTm80RFJXejNnbWpZdUlDdklhMGo5emppT1FkWgo1Q2xVQkdk
        YUlScFoyRG5CblBUM2tJbnhSd0hmU0JMVlVTRXRTcXh4YkV2dk5LWkZWY0lsRUV5ODZodnJ5OUZD
        CjhFOWRXWEw5VDhMd29uVXpqSjBxZ242cGRjNHpjdEtUMDFjaDQvWGw2UjBVQkR5Q1NoSGFyU29C
        eTkvSk1NTXIKajBoeEZSN3Izb052a2N3QWl6T1RsQ3BWdTZaNHF2cng3NndCc0hIanV6elNiODJL
        dUxnelJUNElWbjFjbzRrVQpSaTlBRkNaRlh6QklaQlEwTUZ6NU03bzJkN0ovN3ZMOFhYRlhwWlpy
        K3RibWE1L3BCSmZhcXliK3FPRXViWGdUCjFsSDZGeFNVcWt0TktQNlZoeWdQY2ZSMlR4YWtHZ0cw
        Ny9qVWZWRmhpVXM5aFBlejh6Sjg2RWMrd283VEVQbEsKSlRnMHZmMDM4MTROR3ZuWmlpTnBFWVBM
        S0ZhcHlDMWJONVdFTGFTWFVBaVFPZDJjK01xVHAyN21vV1RZa29TOApzRFczRTMraEN6c1djdmFY
        RW1nMjZJTjQybmVUWFBuNS9QajNpcUVoT0pQYkJsY3l6dDBZL1BYeU1jR3JtbUs1CkhxOUMzTndl
        VUV3M09rY09BOXlCdC9kLzZ5S3c3QmovSlFQZGI0aDlWWjNGN09wemFpeXQ5cFhvSXRQMHNUSHUK
        S2ZKbDBCRUFYV29SR2lWM2EyeUlUcGp0a0pkQVBoS0xpSkkrWWowZEVEU05WZnlENFhJTXdQSmpV
        eFpsd2FROQorQUtkVDFBdlplbz0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=
```
