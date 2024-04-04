package awskms

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"sync"
)

// AWSKMSKeyProvider implements the KeyProvider interface for retrieving keys from AWS KMS.
type AWSKMSKeyProvider struct {
	sync.RWMutex

	// AWS KMS specific fields
	region string

	keyId string

	svc *kms.KMS

	isEncKeyExists bool
	encKey         *kms.GenerateDataKeyOutput
}

func NewAWSKMSKeyProvider(region string, keyId string) (*AWSKMSKeyProvider, error) {
	// Initialize DynamoDB client
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return nil, err
	}
	svc := kms.New(sess)

	return &AWSKMSKeyProvider{
		region: region,
		keyId:  keyId,
		svc:    svc,
	}, nil
}

// GenerateKey generates a new encryption key using AWS KMS.
func (kp *AWSKMSKeyProvider) GenerateKey() ([]byte, []byte, []byte, error) {
	var resp *kms.GenerateDataKeyOutput
	var err error

	if kp.isEncKeyExists {
		resp = kp.encKey
	} else {
		// Call AWS KMS API to generate a new data key
		resp, err = kp.svc.GenerateDataKey(&kms.GenerateDataKeyInput{
			KeyId:         aws.String(kp.keyId),
			NumberOfBytes: aws.Int64(64),
		})
		if err != nil {
			return nil, nil, nil, err
		}

		kp.Lock()
		kp.encKey = resp
		kp.isEncKeyExists = true
		kp.Unlock()
	}

	return resp.Plaintext[:32], resp.Plaintext[32:], resp.CiphertextBlob, nil
}

// RetrieveKey retrieves the encryption key from AWS KMS.
func (kp *AWSKMSKeyProvider) RetrieveKey(ctBlob []byte) ([]byte, []byte, error) {
	// Implement logic to retrieve the key from AWS KMS

	// Call AWS KMS API to decrypt the encrypted key
	resp, err := kp.svc.Decrypt(&kms.DecryptInput{
		CiphertextBlob: ctBlob,
	})
	if err != nil {
		return nil, nil, err
	}

	return resp.Plaintext[:32], resp.Plaintext[32:], nil
}
