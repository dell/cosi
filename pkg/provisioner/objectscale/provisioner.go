// Copyright © 2022-2025 Dell Inc. or its subsidiaries. All Rights Reserved.
//
// This software contains the intellectual property of Dell Inc.
// or is licensed to Dell Inc. from third parties. Use of this software
// and the intellectual property contained therein is expressly limited to the
// terms and conditions of the License Agreement under which it is provided by or
// on behalf of Dell Inc. or its subsidiaries.

package objectscale

import (
	"context"
	"time"

	"github.com/dell/cosi/pkg/provisioner/policy"

	"github.com/aws/aws-sdk-go-v2/service/iam"

	"go.opentelemetry.io/otel"

	cosi "sigs.k8s.io/container-object-storage-interface/proto"
)

const (
	defaultTimeout = time.Second * 20
)

func parsePolicyStatement(
	ctx context.Context,
	inputStatements []policy.StatementEntry,
	awsBucketResourceARNs []string,
	awsPrincipalString string,
) []policy.StatementEntry {
	_, span := otel.Tracer(GrantBucketAccessTraceName).Start(ctx, "ObjectscaleParsePolicyStatement")
	defer span.End()

	// check if our policy already exists
	for _, statement := range inputStatements {
		if statement.Sid == PolicySid && statement.Principal["AWS"] == awsPrincipalString {
			return inputStatements
		}
	}

	newStatement := policy.StatementEntry{}
	newStatement.Sid = "cosi"
	newStatement.Resource = awsBucketResourceARNs
	newStatement.Effect = allowEffect
	newStatement.Principal = map[string]string{"AWS": awsPrincipalString}
	newStatement.Action = []string{"*"}
	inputStatements = append(inputStatements, newStatement)

	return inputStatements
}

func assembleCredentials(
	ctx context.Context,
	accessKey *iam.CreateAccessKeyOutput,
	s3Endpoint,
	userName,
	bucketName string,
) map[string]*cosi.CredentialDetails {
	_, span := otel.Tracer(GrantBucketAccessTraceName).Start(ctx, "ObjectscaleAssembeCredentials")
	defer span.End()

	secretsMap := make(map[string]string)
	secretsMap["accessKeyID"] = *accessKey.AccessKey.AccessKeyId
	secretsMap["accessSecretKey"] = *accessKey.AccessKey.SecretAccessKey
	secretsMap["endpoint"] = s3Endpoint
	secretsMap["bucketName"] = bucketName

	log.Debugf("Created secret access key %s for user %s with endpoint %s was created.", *accessKey.AccessKey.AccessKeyId, userName, s3Endpoint)
	span.AddEvent("secret access key for user with endpoint was created")

	credentialDetails := cosi.CredentialDetails{Secrets: secretsMap}
	credentials := make(map[string]*cosi.CredentialDetails)
	credentials["s3"] = &credentialDetails

	log.Infof("Successfully granted access to the bucket %s for user %s", bucketName, userName)
	span.AddEvent("access to the bucket for user successfully granted")

	return credentials
}
