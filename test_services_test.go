package main

import (
	"testing"

	"github.com/bcpitutor/ostiki/appconfig"
	"github.com/bcpitutor/ostiki/services"
)

func TestAwsService(t *testing.T) {
	config := appconfig.GetAppConfig()

	aws, err := services.GetAWS(config)

	if err != nil {
		t.Errorf("Failed to get AWS service: %v", err)
	}

	t.Logf("AWS service has been created: %+v", aws)

	buckets, err := aws.GetS3BucketList()
	if err != nil {
		t.Errorf("Failed to get S3 buckets: %v", err)
	}

	t.Logf("Got %d S3 buckets", len(buckets))
}
