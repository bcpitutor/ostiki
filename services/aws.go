package services

import (
	"context"
	b64 "encoding/base64"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type AWSService struct {
	STSClient *sts.Client
	S3Client  *s3.Client
	KMSClient *kms.Client
}

type AwsCredentials struct {
	AccessKeyId     string
	SecretAccessKey string
	SessionToken    string
	Region          string
}

func GetAWS() (*AWSService, error) {
	region := os.Getenv("AWS_REGION")
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load config, %v", err)
	}

	awsService := AWSService{}
	awsService.STSClient = sts.NewFromConfig(cfg)
	awsService.S3Client = s3.NewFromConfig(cfg)
	awsService.KMSClient = kms.NewFromConfig(cfg)
	return &awsService, nil
}

func (as *AWSService) ObtainAWSRoleWithToken(token string, sessionName string, roleArn string, ttl int32, region string) (AwsCredentials, error) {
	creds := AwsCredentials{}

	params := &sts.AssumeRoleWithWebIdentityInput{
		RoleArn:          aws.String(roleArn),
		RoleSessionName:  aws.String(sessionName),
		WebIdentityToken: aws.String(token),
		DurationSeconds:  aws.Int32(ttl),
	}

	resp, err := as.STSClient.AssumeRoleWithWebIdentity(
		context.TODO(),
		params,
	)
	if err != nil {
		return creds, err
	}

	creds.AccessKeyId = *resp.Credentials.AccessKeyId
	creds.SecretAccessKey = *resp.Credentials.SecretAccessKey
	creds.SessionToken = *resp.Credentials.SessionToken
	creds.Region = region

	return creds, nil
}

func (as *AWSService) GetS3BucketList() ([]types.Bucket, error) {
	input := &s3.ListBucketsInput{}
	var buckets []types.Bucket

	resp, err := as.S3Client.ListBuckets(
		context.TODO(),
		input,
	)
	if err != nil {
		return buckets, fmt.Errorf("Failed to list buckets: %v", err)
	}

	return resp.Buckets, nil
}

func (as *AWSService) GetEncryptedSecret(ticketData string) (string, error) {
	aliasesInput := kms.ListAliasesInput{}
	result, err := as.KMSClient.ListAliases(context.TODO(), &aliasesInput)
	if err != nil {
		return "", err
	}

	var KMSSecretKeyID string
	for _, v := range result.Aliases {
		if *v.AliasName == "alias/TikiKey" { // TODO: PARAMETERIZE?
			KMSSecretKeyID = *v.TargetKeyId
		}
	}

	input := &kms.EncryptInput{
		KeyId:     &KMSSecretKeyID,
		Plaintext: []byte(ticketData),
	}

	kmsResult, err := as.KMSClient.Encrypt(context.Background(), input)
	if err != nil {
		return "", err
	}
	blobString := b64.StdEncoding.EncodeToString(kmsResult.CiphertextBlob)

	return blobString, nil
}

func (as *AWSService) GetDecryptedText(ticketData string) (string, error) {
	aliasesInput := kms.ListAliasesInput{}
	result, err := as.KMSClient.ListAliases(context.TODO(), &aliasesInput)
	if err != nil {
		return "", err
	}

	var KMSSecretKeyID string
	for _, v := range result.Aliases {
		if *v.AliasName == "alias/TikiKey" { // TODO: PARAMETERIZE?
			KMSSecretKeyID = *v.TargetKeyId
		}
	}

	blob, err := b64.StdEncoding.DecodeString(ticketData)
	if err != nil {
		return "", err
	}

	input := &kms.DecryptInput{
		KeyId:          &KMSSecretKeyID, // Necessary??
		CiphertextBlob: blob,
	}
	decryptedOutput, err := as.KMSClient.Decrypt(context.Background(), input)
	if err != nil {
		return "", err
	}

	return string(decryptedOutput.Plaintext), nil
}
