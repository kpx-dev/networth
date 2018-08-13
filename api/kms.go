package main

import (
	"encoding/base64"

	"github.com/aws/aws-sdk-go-v2/service/kms"
)

// KMSClient kms client struct
type KMSClient struct {
	*kms.KMS
}

// NewKMSClient new instance of the kms client
func NewKMSClient() *KMSClient {
	cfg := loadAWSConfig()
	svc := kms.New(cfg)

	return &KMSClient{svc}
}

// Decrypt decrypt kms token
func (k KMSClient) Decrypt(token string) string {
	decoded, _ := base64.StdEncoding.DecodeString(token)
	input := &kms.DecryptInput{
		CiphertextBlob: []byte(decoded),
	}

	req := k.DecryptRequest(input)
	res, err := req.Send()

	if err != nil {
		panic(err)
	}

	return string(res.Plaintext)
}
