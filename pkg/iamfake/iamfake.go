package iamfake

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
)

// FakeIAMClient is a set of expected outputs for fake client methods.
type FakeIAMClient struct {
	iamiface.IAMAPI
	getUserOutput         *iam.GetUserOutput
	createAccessKeyOutput *iam.CreateAccessKeyOutput
	deleteAccessKeyOutput *iam.DeleteAccessKeyOutput
	listAccessKeysOutput  *iam.ListAccessKeysOutput
	createUserOutput      *iam.CreateUserOutput
	deleteUserOutput      *iam.DeleteUserOutput
}

// NewFakeIAMClient returns new FakeIAMClient based on provided parameters.
func NewFakeIAMClient(objs ...interface{}) *FakeIAMClient {
	var (
		getUserOutput         *iam.GetUserOutput
		createAccessKeyOutput *iam.CreateAccessKeyOutput
		deleteAccessKeyOutput *iam.DeleteAccessKeyOutput
		listAccessKeysOutput  *iam.ListAccessKeysOutput
		createUserOutput      *iam.CreateUserOutput
		deleteUserOutput      *iam.DeleteUserOutput
	)
	for _, o := range objs {
		switch object := o.(type) {
		case *iam.GetUserOutput:
			getUserOutput = object
		case *iam.CreateAccessKeyOutput:
			createAccessKeyOutput = object
		case *iam.DeleteAccessKeyOutput:
			deleteAccessKeyOutput = object
		case *iam.ListAccessKeysOutput:
			listAccessKeysOutput = object
		case *iam.CreateUserOutput:
			createUserOutput = object
		case *iam.DeleteUserOutput:
			deleteUserOutput = object
		default:
			panic(fmt.Sprintf("Fake client set doesn't support %T type", o))
		}
	}

	return &FakeIAMClient{
		getUserOutput:         getUserOutput,
		createAccessKeyOutput: createAccessKeyOutput,
		deleteAccessKeyOutput: deleteAccessKeyOutput,
		listAccessKeysOutput:  listAccessKeysOutput,
		createUserOutput:      createUserOutput,
		deleteUserOutput:      deleteUserOutput,
	}

}

// GetUser returns GetUserOutput or error depending on provided GetUserInput.UserName.
func (fakeIAM *FakeIAMClient) GetUser(input *iam.GetUserInput) (*iam.GetUserOutput, error) {
	switch *input.UserName {
	case "success":
		return fakeIAM.getUserOutput, nil
	case "fail":
		return nil, errors.New(iam.ErrCodeNoSuchEntityException)
	default:
		return nil, errors.New(iam.ErrCodeServiceFailureException)
	}
}

// CreateAccessKey returns CreateAccessKeyOutput or error depending on provided CreateAccessKeyInput.UserName.
func (fakeIAM *FakeIAMClient) CreateAccessKey(input *iam.CreateAccessKeyInput) (*iam.CreateAccessKeyOutput, error) {
	switch *input.UserName {
	case "success":
		return fakeIAM.createAccessKeyOutput, nil
	case "fail":
		return nil, errors.New(iam.ErrCodeNoSuchEntityException)
	default:
		return nil, errors.New(iam.ErrCodeServiceFailureException)
	}
}

// DeleteAccessKey returns DeleteAccessKeyOutput or error depending on provided DeleteAccessKeyInput.UserName.
func (fakeIAM *FakeIAMClient) DeleteAccessKey(input *iam.DeleteAccessKeyInput) (*iam.DeleteAccessKeyOutput, error) {
	switch *input.UserName {
	case "success":
		return fakeIAM.deleteAccessKeyOutput, nil
	case "fail":
		return nil, errors.New(iam.ErrCodeNoSuchEntityException)
	default:
		return nil, errors.New(iam.ErrCodeServiceFailureException)
	}
}

// ListAccessKeys returns ListAccessKeysOutput or error depending on provided ListAccessKeysInput.UserName.
func (fakeIAM *FakeIAMClient) ListAccessKeys(input *iam.ListAccessKeysInput) (*iam.ListAccessKeysOutput, error) {
	switch *input.UserName {
	case "success":
		return fakeIAM.listAccessKeysOutput, nil
	case "fail":
		return nil, errors.New(iam.ErrCodeNoSuchEntityException)
	default:
		return nil, errors.New(iam.ErrCodeServiceFailureException)
	}
}

// CreateUserWithContext returns CreateUserOutput or error depending on provided CreateUserInput.UserName.
func (fakeIAM *FakeIAMClient) CreateUserWithContext(_ aws.Context, input *iam.CreateUserInput, _ ...request.Option) (*iam.CreateUserOutput, error) {
	switch *input.UserName {
	case "success":
		return fakeIAM.createUserOutput, nil
	case "fail":
		return nil, errors.New(iam.ErrCodeEntityAlreadyExistsException)
	default:
		return nil, errors.New(iam.ErrCodeServiceFailureException)
	}
}

// DeleteUser returns DeleteUserOutput or error depending on provided DeleteUserInput.UserName.
func (fakeIAM *FakeIAMClient) DeleteUser(input *iam.DeleteUserInput) (*iam.DeleteUserOutput, error) {
	switch *input.UserName {
	case "success":
		return fakeIAM.deleteUserOutput, nil
	case "fail":
		return nil, errors.New(iam.ErrCodeNoSuchEntityException)
	default:
		return nil, errors.New(iam.ErrCodeServiceFailureException)
	}
}
