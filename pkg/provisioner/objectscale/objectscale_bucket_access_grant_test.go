package objectscale

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/dell/cosi-driver/pkg/iamfaketoo"
	"github.com/dell/goobjectscale/pkg/client/fake"
	"github.com/dell/goobjectscale/pkg/client/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	cosi "sigs.k8s.io/container-object-storage-interface-spec"
)

var _ iamiface.IAMAPI = (*iamfaketoo.IAMAPI)(nil)

const (
	namespace     = "namespace"
	testID        = "test.id"
	objectScaleID = "objectscale"
	objectStoreID = "objectstore"
)

func TestValidAccessGranting(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100) // Magic value,... abra kadabra
	defer cancel()

	// That's how we can mock the objectscale IAM api client
	IAMClient := iamfaketoo.NewIAMAPI(t)
	IAMClient.On("CreateUserWithContext", mock.Anything, mock.Anything).Return(
		&iam.CreateUserOutput{
			User: &iam.User{
				UserName: aws.String("namespace-user-valid"), // This mocked response is based on `namesapce` from server and bucketId from request
			},
		}, nil).Once()
	IAMClient.On("GetUser", mock.Anything).Return(nil, nil).Once()
	IAMClient.On("CreateAccessKey", mock.Anything).Return(&iam.CreateAccessKeyOutput{AccessKey: &iam.AccessKey{AccessKeyId: aws.String("acc"), SecretAccessKey: aws.String("sec")}}, nil).Once()

	server := Server{
		mgmtClient: fake.NewClientSet(&model.Bucket{ // That's how we can mock the objectscale bucket api client
			Name:      "valid", // This is based on "bucket-valid" BucketId from request
			Namespace: namespace,
		}),
		iamClient:     IAMClient, // Inject mocked IAM client
		namespace:     namespace,
		backendID:     testID,
		objectScaleID: objectScaleID,
		objectStoreID: objectStoreID,
	}

	req := &cosi.DriverGrantBucketAccessRequest{
		BucketId:           "bucket-valid",
		Name:               "bucket-access-valid",
		AuthenticationType: cosi.AuthenticationType_Key,
		Parameters: map[string]string{
			"X-TEST/Buckets/UpdatePolicy/force-success": "true", // This is mocking response from objectscale bucket api client
		},
	}

	response, err := server.DriverGrantBucketAccess(ctx, req)
	assert.ErrorIs(t, err, nil, err)
	assert.NotNil(t, response)
}

func TestInvalidAccessKeyCreation(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100) // Magic value,... abra kadabra
	defer cancel()

	// That's how we can mock the objectscale IAM api client
	IAMClient := iamfaketoo.NewIAMAPI(t)
	IAMClient.On("CreateUserWithContext", mock.Anything, mock.Anything).Return(
		&iam.CreateUserOutput{
			User: &iam.User{
				UserName: aws.String("namespace-user-valid"), // This mocked response is based on `namesapce` from server and bucketId from request
			},
		}, nil).Once()
	IAMClient.On("GetUser", mock.Anything).Return(nil, nil).Once()
	IAMClient.On("CreateAccessKey", mock.Anything).Return(nil, errors.New("failed to create access key")).Once()

	server := Server{
		mgmtClient: fake.NewClientSet(&model.Bucket{ // That's how we can mock the objectscale bucket api client
			Name:      "valid", // This is based on "bucket-valid" BucketId from request
			Namespace: namespace,
		}),
		iamClient:     IAMClient, // Inject mocked IAM client
		namespace:     namespace,
		backendID:     testID,
		objectScaleID: objectScaleID,
		objectStoreID: objectStoreID,
	}

	req := &cosi.DriverGrantBucketAccessRequest{
		BucketId:           "bucket-valid",
		Name:               "bucket-access-valid",
		AuthenticationType: cosi.AuthenticationType_Key,
		Parameters: map[string]string{
			"X-TEST/Buckets/UpdatePolicy/force-success": "true", // This is mocking response from objectscale bucket api client
		},
	}

	response, err := server.DriverGrantBucketAccess(ctx, req)
	assert.ErrorIs(t, err, status.Error(codes.Internal, "failed to create access key"), err)
	assert.Nil(t, response)
}

func TestInvalidUserCreation(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100) // Magic value,... abra kadabra
	defer cancel()

	// That's how we can mock the objectscale IAM api client
	IAMClient := iamfaketoo.NewIAMAPI(t)
	IAMClient.On("GetUser", mock.Anything).Return(nil, nil).Once()
	IAMClient.On("CreateUserWithContext", mock.Anything, mock.Anything).Return(nil, errors.New("failed to create user")).Once()

	server := Server{
		mgmtClient: fake.NewClientSet(&model.Bucket{ // That's how we can mock the objectscale bucket api client
			Name:      "valid", // This is based on "bucket-valid" BucketId from request
			Namespace: namespace,
		}),
		iamClient:     IAMClient, // Inject mocked IAM client
		namespace:     namespace,
		backendID:     testID,
		objectScaleID: objectScaleID,
		objectStoreID: objectStoreID,
	}

	req := &cosi.DriverGrantBucketAccessRequest{
		BucketId:           "bucket-valid",
		Name:               "bucket-access-valid",
		AuthenticationType: cosi.AuthenticationType_Key,
		Parameters: map[string]string{
			"X-TEST/Buckets/UpdatePolicy/force-success": "true", // This is mocking response from objectscale bucket api client
		},
	}

	response, err := server.DriverGrantBucketAccess(ctx, req)
	assert.ErrorIs(t, err, status.Error(codes.Internal, "cannot create user namespace-user-valid"), err)
	assert.Nil(t, response)
}

func TestInvalidUserRetrieval(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100) // Magic value,... abra kadabra
	defer cancel()

	// That's how we can mock the objectscale IAM api client
	IAMClient := iamfaketoo.NewIAMAPI(t)
	IAMClient.On("GetUser", mock.Anything).Return(nil, errors.New("failed to retrieve user")).Once()

	server := Server{
		mgmtClient: fake.NewClientSet(&model.Bucket{ // That's how we can mock the objectscale bucket api client
			Name:      "valid", // This is based on "bucket-valid" BucketId from request
			Namespace: namespace,
		}),
		iamClient:     IAMClient, // Inject mocked IAM client
		namespace:     namespace,
		backendID:     testID,
		objectScaleID: objectScaleID,
		objectStoreID: objectStoreID,
	}

	req := &cosi.DriverGrantBucketAccessRequest{
		BucketId:           "bucket-valid",
		Name:               "bucket-access-valid",
		AuthenticationType: cosi.AuthenticationType_Key,
		Parameters: map[string]string{
			"X-TEST/Buckets/UpdatePolicy/force-success": "true", // This is mocking response from objectscale bucket api client
		},
	}

	response, err := server.DriverGrantBucketAccess(ctx, req)
	assert.ErrorIs(t, err, status.Error(codes.Internal, "failed to check for user existence"), err)
	assert.Nil(t, response)
}

func TestInvalidBucketPolicyUpdate(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100) // Magic value,... abra kadabra
	defer cancel()

	// That's how we can mock the objectscale IAM api client
	IAMClient := iamfaketoo.NewIAMAPI(t)
	IAMClient.On("GetUser", mock.Anything).Return(nil, nil).Once()
	IAMClient.On("CreateUserWithContext", mock.Anything, mock.Anything).Return(
		&iam.CreateUserOutput{
			User: &iam.User{
				UserName: aws.String("namespace-user-valid"), // This mocked response is based on `namesapce` from server and bucketId from request
			},
		}, nil).Once()

	server := Server{
		mgmtClient: fake.NewClientSet(&model.Bucket{ // That's how we can mock the objectscale bucket api client
			Name:      "valid", // This is based on "bucket-valid" BucketId from request
			Namespace: namespace,
		}),
		iamClient:     IAMClient, // Inject mocked IAM client
		namespace:     namespace,
		backendID:     testID,
		objectScaleID: objectScaleID,
		objectStoreID: objectStoreID,
	}

	req := &cosi.DriverGrantBucketAccessRequest{
		BucketId:           "bucket-valid",
		Name:               "bucket-access-valid",
		AuthenticationType: cosi.AuthenticationType_Key,
		Parameters: map[string]string{
			"X-TEST/Buckets/UpdatePolicy/force-fail": "true", // This is mocking response from objectscale bucket api client
		},
	}

	response, err := server.DriverGrantBucketAccess(ctx, req)
	assert.ErrorIs(t, err, status.Error(codes.Internal, "failed to update bucket policy"), err)
	assert.Nil(t, response)
}

func TestEmptyBucketID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100) // Magic value,... abra kadabra
	defer cancel()

	server := Server{
		mgmtClient: fake.NewClientSet(&model.Bucket{ // That's how we can mock the objectscale bucket api client
			Name:      "valid", // This is based on "bucket-valid" BucketId from request
			Namespace: namespace,
		}),
		iamClient:     iamfaketoo.NewIAMAPI(t), // Inject mocked IAM client
		namespace:     namespace,
		backendID:     testID,
		objectScaleID: objectScaleID,
		objectStoreID: objectStoreID,
	}

	req := &cosi.DriverGrantBucketAccessRequest{
		BucketId:           "",
		Name:               "bucket-access-valid",
		AuthenticationType: cosi.AuthenticationType_Key,
		Parameters: map[string]string{
			"X-TEST/Buckets/UpdatePolicy/force-fail": "true", // This is mocking response from objectscale bucket api client
		},
	}

	response, err := server.DriverGrantBucketAccess(ctx, req)
	assert.ErrorIs(t, err, status.Error(codes.InvalidArgument, "empty bucketID"), err)
	assert.Nil(t, response)
}

func TestEmptyName(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100) // Magic value,... abra kadabra
	defer cancel()

	server := Server{
		mgmtClient: fake.NewClientSet(&model.Bucket{ // That's how we can mock the objectscale bucket api client
			Name:      "valid", // This is based on "bucket-valid" BucketId from request
			Namespace: namespace,
		}),
		iamClient:     iamfaketoo.NewIAMAPI(t), // Inject mocked IAM client
		namespace:     namespace,
		backendID:     testID,
		objectScaleID: objectScaleID,
		objectStoreID: objectStoreID,
	}

	req := &cosi.DriverGrantBucketAccessRequest{
		BucketId:           "bucket-valid",
		Name:               "",
		AuthenticationType: cosi.AuthenticationType_Key,
		Parameters: map[string]string{
			"X-TEST/Buckets/UpdatePolicy/force-fail": "true", // This is mocking response from objectscale bucket api client
		},
	}

	response, err := server.DriverGrantBucketAccess(ctx, req)
	assert.ErrorIs(t, err, status.Error(codes.InvalidArgument, "empty bucket access name"), err)
	assert.Nil(t, response)
}

func TestInvalidAuthenticationType(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100) // Magic value,... abra kadabra
	defer cancel()

	server := Server{
		mgmtClient: fake.NewClientSet(&model.Bucket{ // That's how we can mock the objectscale bucket api client
			Name:      "valid", // This is based on "bucket-valid" BucketId from request
			Namespace: namespace,
		}),
		iamClient:     iamfaketoo.NewIAMAPI(t), // Inject mocked IAM client
		namespace:     namespace,
		backendID:     testID,
		objectScaleID: objectScaleID,
		objectStoreID: objectStoreID,
	}

	req := &cosi.DriverGrantBucketAccessRequest{
		BucketId:           "bucket-valid",
		Name:               "bucket-access-valid",
		AuthenticationType: cosi.AuthenticationType_UnknownAuthenticationType,
		Parameters: map[string]string{
			"X-TEST/Buckets/UpdatePolicy/force-fail": "true", // This is mocking response from objectscale bucket api client
		},
	}

	response, err := server.DriverGrantBucketAccess(ctx, req)
	assert.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid authentication type"), err)
	assert.Nil(t, response)
}

func TestIAMNotImplemented(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100) // Magic value,... abra kadabra
	defer cancel()

	server := Server{
		mgmtClient: fake.NewClientSet(&model.Bucket{ // That's how we can mock the objectscale bucket api client
			Name:      "valid", // This is based on "bucket-valid" BucketId from request
			Namespace: namespace,
		}),
		iamClient:     iamfaketoo.NewIAMAPI(t), // Inject mocked IAM client
		namespace:     namespace,
		backendID:     testID,
		objectScaleID: objectScaleID,
		objectStoreID: objectStoreID,
	}

	req := &cosi.DriverGrantBucketAccessRequest{
		BucketId:           "bucket-valid",
		Name:               "bucket-access-valid",
		AuthenticationType: cosi.AuthenticationType_IAM,
		Parameters: map[string]string{
			"X-TEST/Buckets/UpdatePolicy/force-fail": "true", // This is mocking response from objectscale bucket api client
		},
	}

	response, err := server.DriverGrantBucketAccess(ctx, req)
	assert.ErrorIs(t, err, status.Error(codes.Unimplemented, "authentication type IAM not implemented"), err)
	assert.Nil(t, response)
}

func TestFailToGetBucket(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100) // Magic value,... abra kadabra
	defer cancel()

	server := Server{
		mgmtClient: fake.NewClientSet(
			&model.Bucket{ // That's how we can mock the objectscale bucket api client
				Name:      "valid", // This is based on "bucket-valid" BucketId from request
				Namespace: namespace,
			},
		),
		iamClient:     iamfaketoo.NewIAMAPI(t), // Inject mocked IAM client
		namespace:     namespace,
		backendID:     testID,
		objectScaleID: objectScaleID,
		objectStoreID: objectStoreID,
	}

	req := &cosi.DriverGrantBucketAccessRequest{
		BucketId:           "bucket-valid",
		Name:               "bucket-access-valid",
		AuthenticationType: cosi.AuthenticationType_Key,
		Parameters: map[string]string{
			"X-TEST/Buckets/Get/force-fail": "true", // This is mocking response from objectscale bucket api client
		},
	}

	response, err := server.DriverGrantBucketAccess(ctx, req)
	assert.ErrorIs(t, err, status.Error(codes.Internal, "an unexpected error occurred"), err)
	assert.Nil(t, response)
}

func TestBucketNotFound(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100) // Magic value,... abra kadabra
	defer cancel()

	server := Server{
		mgmtClient:    fake.NewClientSet(),     // That's how we can mock the objectscale bucket api client
		iamClient:     iamfaketoo.NewIAMAPI(t), // Inject mocked IAM client
		namespace:     namespace,
		backendID:     testID,
		objectScaleID: objectScaleID,
		objectStoreID: objectStoreID,
	}

	req := &cosi.DriverGrantBucketAccessRequest{
		BucketId:           "bucket-invalid",
		Name:               "bucket-access-valid",
		AuthenticationType: cosi.AuthenticationType_Key,
		Parameters:         map[string]string{},
	}

	response, err := server.DriverGrantBucketAccess(ctx, req)
	assert.ErrorIs(t, err, status.Error(codes.NotFound, "bucket not found"), err)
	assert.Nil(t, response)
}

func TestValidButUserAlreadyExists(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100) // Magic value,... abra kadabra
	defer cancel()

	IAMClient := iamfaketoo.NewIAMAPI(t)
	IAMClient.On("GetUser", mock.Anything).Return(&iam.GetUserOutput{
		User: &iam.User{
			UserName: aws.String("namespace-user-valid"),
		},
	}, nil).Once()
	IAMClient.On("CreateAccessKey", mock.Anything).Return(&iam.CreateAccessKeyOutput{AccessKey: &iam.AccessKey{AccessKeyId: aws.String("acc"), SecretAccessKey: aws.String("sec")}}, nil).Once()

	server := Server{
		mgmtClient: fake.NewClientSet(&model.Bucket{ // That's how we can mock the objectscale bucket api client
			Name:      "valid", // This is based on "bucket-valid" BucketId from request
			Namespace: namespace,
		}),
		iamClient:     IAMClient, // Inject mocked IAM client
		namespace:     namespace,
		backendID:     testID,
		objectScaleID: objectScaleID,
		objectStoreID: objectStoreID,
	}

	req := &cosi.DriverGrantBucketAccessRequest{
		BucketId:           "bucket-valid",
		Name:               "bucket-access-valid",
		AuthenticationType: cosi.AuthenticationType_Key,
		Parameters: map[string]string{
			"X-TEST/Buckets/UpdatePolicy/force-success": "true", // This is mocking response from objectscale bucket api client
		},
	}

	response, err := server.DriverGrantBucketAccess(ctx, req)
	assert.ErrorIs(t, err, nil, err)
	assert.NotNil(t, response)
}
