// Copyright © 2023 Dell Inc. or its subsidiaries. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//      http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package objectscale

import (
	"github.com/aws/aws-sdk-go/aws"

	"github.com/dell/cosi/pkg/config"
)

var invalidBase64 = `💀`

// TODO: config.Objectscale builder.
var (
	validConfig = &config.Objectscale{
		Id:                 "valid.id",
		ObjectscaleGateway: "gateway.objectscale.test",
		ObjectstoreGateway: "gateway.objectstore.test",
		Namespace:          "validnamespace",
		Credentials: config.Credentials{
			Username: "testuser",
			Password: "testpassword",
		},
		Protocols: config.Protocols{
			S3: &config.S3{
				Endpoint: "s3.objectstore.test",
			},
		},
		Tls: config.Tls{
			Insecure: true,
		},
		Region:        aws.String("us-east-1"),
		ObjectscaleId: "valid.objectscale.id",
		ObjectstoreId: "valid.objectstore.id",
	}

	invalidConfigWithHyphens = &config.Objectscale{
		Id:                 "id-with-hyphens",
		ObjectscaleGateway: "gateway.objectscale.test",
		ObjectstoreGateway: "gateway.objectstore.test",
		Credentials: config.Credentials{
			Username: "testuser",
			Password: "testpassword",
		},
		Namespace: "validnamespace",
		Protocols: config.Protocols{
			S3: &config.S3{
				Endpoint: "s3.objectstore.test",
			},
		},
		Tls: config.Tls{
			Insecure: true,
		},
		Region:        aws.String("us-east-1"),
		ObjectscaleId: "valid.objectscale.id",
		ObjectstoreId: "valid.objectstore.id",
	}

	invalidConfigEmptyID = &config.Objectscale{
		Id:                 "",
		ObjectscaleGateway: "gateway.objectscale.test",
		ObjectstoreGateway: "gateway.objectstore.test",
		Namespace:          "validnamespace",
		Credentials: config.Credentials{
			Username: "testuser",
			Password: "testpassword",
		},
		Protocols: config.Protocols{
			S3: &config.S3{
				Endpoint: "s3.objectstore.test",
			},
		},
		Tls: config.Tls{
			Insecure: true,
		},
		Region:        aws.String("us-east-1"),
		ObjectscaleId: "valid.objectscale.id",
		ObjectstoreId: "valid.objectstore.id",
	}

	invalidConfigTLS = &config.Objectscale{
		Id:                 "valid.id",
		ObjectscaleGateway: "gateway.objectscale.test",
		ObjectstoreGateway: "gateway.objectstore.test",
		Namespace:          "validnamespace",
		Credentials: config.Credentials{
			Username: "testuser",
			Password: "testpassword",
		},
		Protocols: config.Protocols{
			S3: &config.S3{
				Endpoint: "s3.objectstore.test",
			},
		},
		Tls: config.Tls{
			Insecure: false,
			RootCas:  &invalidBase64,
		},
		Region:        aws.String("us-east-1"),
		ObjectscaleId: "valid.objectscale.id",
		ObjectstoreId: "valid.objectstore.id",
	}

	emptyNamespaceConfig = &config.Objectscale{
		Id:                 "valid.id",
		ObjectscaleGateway: "gateway.objectscale.test",
		ObjectstoreGateway: "gateway.objectstore.test",
		Namespace:          "",
		Credentials: config.Credentials{
			Username: "testuser",
			Password: "testpassword",
		},
		Protocols: config.Protocols{
			S3: &config.S3{
				Endpoint: "s3.objectstore.test",
			},
		},
		Tls: config.Tls{
			Insecure: true,
		},
		Region:        aws.String("us-east-1"),
		ObjectscaleId: "valid.objectscale.id",
		ObjectstoreId: "valid.objectstore.id",
	}
	emptyUsernameConfig = &config.Objectscale{
		Id:                 "valid.id",
		ObjectscaleGateway: "gateway.objectscale.test",
		ObjectstoreGateway: "gateway.objectstore.test",
		Namespace:          "validnamespace",
		Credentials: config.Credentials{
			Username: "",
			Password: "testpassword",
		},
		Protocols: config.Protocols{
			S3: &config.S3{
				Endpoint: "s3.objectstore.test",
			},
		},
		Tls: config.Tls{
			Insecure: true,
		},
		Region:        aws.String("us-east-1"),
		ObjectscaleId: "valid.objectscale.id",
		ObjectstoreId: "valid.objectstore.id",
	}
	emptyPasswordConfig = &config.Objectscale{
		Id:                 "valid.id",
		ObjectscaleGateway: "gateway.objectscale.test",
		ObjectstoreGateway: "gateway.objectstore.test",
		Namespace:          "validnamespace",
		Credentials: config.Credentials{
			Username: "testuser",
			Password: "",
		},
		Protocols: config.Protocols{
			S3: &config.S3{
				Endpoint: "s3.objectstore.test",
			},
		},
		Tls: config.Tls{
			Insecure: true,
		},
		Region:        aws.String("us-east-1"),
		ObjectscaleId: "valid.objectscale.id",
		ObjectstoreId: "valid.objectstore.id",
	}
	emptyRegionConfig = &config.Objectscale{
		Id:                 "valid.id",
		ObjectscaleGateway: "gateway.objectscale.test",
		ObjectstoreGateway: "gateway.objectstore.test",
		Namespace:          "validnamespace",
		Credentials: config.Credentials{
			Username: "testuser",
			Password: "testpassword",
		},
		Protocols: config.Protocols{
			S3: &config.S3{
				Endpoint: "s3.objectstore.test",
			},
		},
		Tls: config.Tls{
			Insecure: true,
		},
		Region:        aws.String(""),
		ObjectscaleId: "valid.objectscale.id",
		ObjectstoreId: "valid.objectstore.id",
	}
	regionNotSetConfig = &config.Objectscale{
		Id:                 "valid.id",
		ObjectscaleGateway: "gateway.objectscale.test",
		ObjectstoreGateway: "gateway.objectstore.test",
		Namespace:          "validnamespace",
		Credentials: config.Credentials{
			Username: "testuser",
			Password: "testpassword",
		},
		Protocols: config.Protocols{
			S3: &config.S3{
				Endpoint: "s3.objectstore.test",
			},
		},
		Tls: config.Tls{
			Insecure: true,
		},
		Region:        nil,
		ObjectscaleId: "valid.objectscale.id",
		ObjectstoreId: "valid.objectstore.id",
	}
	emptyObjectscaleGatewayConfig = &config.Objectscale{
		Id:                 "valid.id",
		ObjectscaleGateway: "",
		ObjectstoreGateway: "gateway.objectstore.test",
		Namespace:          "validnamespace",
		Credentials: config.Credentials{
			Username: "testuser",
			Password: "testpassword",
		},
		Protocols: config.Protocols{
			S3: &config.S3{
				Endpoint: "s3.objectstore.test",
			},
		},
		Tls: config.Tls{
			Insecure: true,
		},
		Region:        aws.String("us-east-1"),
		ObjectscaleId: "valid.objectscale.id",
		ObjectstoreId: "valid.objectstore.id",
	}
	emptyObjectstoreGatewayConfig = &config.Objectscale{
		Id:                 "valid.id",
		ObjectscaleGateway: "gateway.objectscale.test",
		ObjectstoreGateway: "",
		Namespace:          "validnamespace",
		Credentials: config.Credentials{
			Username: "testuser",
			Password: "testpassword",
		},
		Protocols: config.Protocols{
			S3: &config.S3{
				Endpoint: "s3.objectstore.test",
			},
		},
		Tls: config.Tls{
			Insecure: true,
		},
		Region:        aws.String("us-east-1"),
		ObjectscaleId: "valid.objectscale.id",
		ObjectstoreId: "valid.objectstore.id",
	}
	emptyS3EndpointConfig = &config.Objectscale{
		Id:                 "valid.id",
		ObjectscaleGateway: "gateway.objectscale.test",
		ObjectstoreGateway: "gateway.objectstore.test",
		Namespace:          "validnamespace",
		Credentials: config.Credentials{
			Username: "testuser",
			Password: "testpassword",
		},
		Protocols: config.Protocols{
			S3: &config.S3{
				Endpoint: "",
			},
		},
		Tls: config.Tls{
			Insecure: true,
		},
		Region:        aws.String("us-east-1"),
		ObjectscaleId: "valid.objectscale.id",
		ObjectstoreId: "valid.objectstore.id",
	}
	emptyObjectscaleIDConfig = &config.Objectscale{
		Id:                 "valid.id",
		ObjectscaleGateway: "gateway.objectscale.test",
		ObjectstoreGateway: "gateway.objectstore.test",
		Namespace:          "validnamespace",
		Credentials: config.Credentials{
			Username: "testuser",
			Password: "testpassword",
		},
		Protocols: config.Protocols{
			S3: &config.S3{
				Endpoint: "s3.objectstore.test",
			},
		},
		Tls: config.Tls{
			Insecure: true,
		},
		Region:        aws.String("us-east-1"),
		ObjectscaleId: "",
		ObjectstoreId: "valid.objectstore.id",
	}
	emptyObjectstoreIDConfig = &config.Objectscale{
		Id:                 "valid.id",
		ObjectscaleGateway: "gateway.objectscale.test",
		ObjectstoreGateway: "gateway.objectstore.test",
		Namespace:          "validnamespace",
		Credentials: config.Credentials{
			Username: "testuser",
			Password: "testpassword",
		},
		Protocols: config.Protocols{
			S3: &config.S3{
				Endpoint: "s3.objectstore.test",
			},
		},
		Tls: config.Tls{
			Insecure: true,
		},
		Region:        aws.String("us-east-1"),
		ObjectscaleId: "valid.objectscale.id",
		ObjectstoreId: "",
	}
)
