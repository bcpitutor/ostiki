package main

import (
	"testing"

	"github.com/bcpitutor/ostiki/appconfig"
	"github.com/bcpitutor/ostiki/db"
	"github.com/bcpitutor/ostiki/logger"
	"github.com/bcpitutor/ostiki/repositories"
)

func TestGetAllGroups(t *testing.T) {
	config := appconfig.GetAppConfig()
	logger := logger.GetTikiLogger(config)

	dynamo, err := db.NewDynamoDBDriver(config, logger)
	if err != nil {
		t.Fatalf("Failed to create DynamoDB driver: %v", err)
	}
	grt := repositories.ProvideGroupRepository(dynamo, nil)
	_, err = grt.GroupRepository.GetAllGroups()
	if err != nil {
		t.Fatalf("Failed to get all groups: %v", err)
	}
	t.Logf("All groups have been retrieved")
}

func TestGetGroup(t *testing.T) {
	config := appconfig.GetAppConfig()
	logger := logger.GetTikiLogger(config)

	dynamo, err := db.NewDynamoDBDriver(config, logger)
	if err != nil {
		t.Fatalf("Failed to create DynamoDB driver: %v", err)
	}
	grt := repositories.ProvideGroupRepository(dynamo, nil)
	_, err = grt.GroupRepository.GetGroup("lms-test1-admins")
	if err != nil {
		t.Fatalf("Failed to get group: %v", err)
	} else {
		t.Logf("Group has been retrieved")
	}

	_, err = grt.GroupRepository.GetGroup("non-existent-group")
	if err != nil {
		t.Logf("Group does not exist, which is fine")
	} else {
		t.Logf("Group exists, which is not fine")
	}
}
