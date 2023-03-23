package iamfake

import (
	"errors"

	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
)


type fakeIAMClient struct {
	iamiface.IAMAPI
	GetUserOutputs map[string]*iam.GetUserOutput
	CreateAccessKeyOutputs map[string]*iam.CreateAccessKeyOutput
	DeleteAccesKeyOutputs map[string]*iam.DeleteAccessKeyOutput
	ListAccessKeysOutput map[string]*iam.ListAccessKeysOutput
}

func (fakeIAM *fakeIAMClient) GetUser(input *iam.GetUserInput) (*iam.GetUserOutput, error) {
	switch *input.UserName {
	case "success":
		return fakeIAM.GetUserOutputs["success"], nil
	case "fail":
		return nil, errors.New("user does not exist")
	default:
		return nil, errors.New("failed to check user existence")
	}
}

func (fakeIAM *fakeIAMClient) CreateAccessKey(*iam.CreateAccessKeyInput) (*iam.CreateAccessKeyOutput, error)) {
	switch *input.UserName {
	case "success":
		return fakeIAM.CreateAccessKeyOutputs["success"], nil
	case "fail":
		return nil, errors.New("failed to create access key")
	default:
		return nil, errors.New("unexpected error")
	}
}

func (fakeIAM *fakeIAMClient) DeleteAccessKey(*iam.DeleteAccessKeyInput) (*iam.DeleteAccessKeyOutput, error) {
	switch *input.UserName {
	case "success":
		return fakeIAM.DeleteAccesKeyOutputs["success"], nil
	case "fail":
		return nil, errors.New("key was not deleted")
	default:
		return nil, errors.New("unexpected error")
	}
}

func (fakeIAM *fakeIAMClient) ListAccessKeys(*iam.ListAccessKeysInput) (*iam.ListAccessKeysOutput, error) {
	switch *input.UserName {
	case "success":
		return fakeIAM.ListAccessKeysOutput["success"], nil
	case "fail":
		return nil, errors.New("access key does not exist")
	default:
		return nil, errors.New("failed to check user access key existence")
	}
}