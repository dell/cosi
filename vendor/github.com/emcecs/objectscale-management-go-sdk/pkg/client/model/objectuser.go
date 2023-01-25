package model

type BlobUser struct {
	UserID    string `json:"userid"`
	Namespace string `json:"namespace"`
}

// ObjectUserList contains an array of object users.
type ObjectUserList struct {
	BlobUser []BlobUser `json:"blobuser"`
}

// ObjectUserInfo contains information about an object user.
type ObjectUserInfo struct {
	Namespace string   `json:"namespace"`
	Name      string   `json:"name"`
	Locked    bool     `json:"locked"`
	Created   string   `json:"created"`
	Tags      []string `json:"tags"`
}

// ObjectUserSecret contains information about object user's secrets.
type ObjectUserSecret struct {
	SecretKey1          string `json:"secret_key_1"`
	KeyTimestamp1       string `json:"key_timestamp_1"`
	KeyExpiryTimestamp1 string `json:"key_expiry_timestamp_1"`
	SecretKey2          string `json:"secret_key_2"`
	KeyTimestamp2       string `json:"key_timestamp_2"`
	KeyExpiryTimestamp2 string `json:"key_expiry_timestamp_2"`
	Link                Link   `json:"link"`
}

// ObjectUserSecretKeyCreateReq to marshal ObjectUserSecretKey create req
type ObjectUserSecretKeyCreateReq struct {
	SecretKey          string `json:"secretkey"`
	Namespace          string `json:"namespace"`
	ExistingKeyExpTime string `json:"existing_key_expiry_time_mins,omitempty"`
}

// ObjectUserSecretKeyDeleteReq to marshal ObjectUserSecretKey delete req
type ObjectUserSecretKeyDeleteReq struct {
	SecretKey string `json:"secret_key"`
	Namespace string `json:"namespace"`
}

// ObjectUserSecretKeyRes to unmarshal ObjectUserSecretKey create resp
type ObjectUserSecretKeyCreateRes struct {
	SecretKey          string `json:"secret_key"`
	KeyTimeStamp       string `json:"key_timestamp"`
	KeyExpiryTimestamp string `json:"key_expiry_timestamp"`
	Link               Link   `json:"link"`
}
