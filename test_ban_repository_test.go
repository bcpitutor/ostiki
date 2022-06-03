package main

import (
	"strconv"
	"testing"
	"time"

	"github.com/tiki-systems/tikiserver/appconfig"
	"github.com/tiki-systems/tikiserver/db"
	"github.com/tiki-systems/tikiserver/logger"
	"github.com/tiki-systems/tikiserver/models"
	"github.com/tiki-systems/tikiserver/repositories"
)

func TestAddBannedUser(t *testing.T) {
	config := appconfig.GetAppConfig()
	logger := logger.GetTikiLogger(config)

	dynamo, err := db.NewDynamoDBDriver(config, logger)
	if err != nil {
		t.Fatalf("Failed to create DynamoDB driver: %v", err)
	}

	brt := repositories.ProvideBanRepository(dynamo)

	err = brt.BanRepository.DBLayer.AddBannedUser(models.BannedUser{
		UserEmail: "selami@sahin.org",
		Details:   "Created for go tests",
		CreatedAt: strconv.FormatInt(time.Now().Unix(), 10),
		UpdatedAt: strconv.FormatInt(time.Now().Unix(), 10),
		CreatedBy: "go-tests",
		UpdatedBy: "go-tests",
	})
	if err != nil {
		t.Fatalf("Failed to add banned user: %v", err)
	}
	t.Logf("Banned user has been added")

	err = brt.BanRepository.DBLayer.AddBannedUser(models.BannedUser{
		UserEmail: "selami@sahin.org",
		Details:   "Created for go tests",
		CreatedAt: strconv.FormatInt(time.Now().Unix(), 10),
		UpdatedAt: strconv.FormatInt(time.Now().Unix(), 10),
		CreatedBy: "go-tests",
		UpdatedBy: "go-tests",
	})
	if err != nil {
		t.Logf("User is already banned")
	} else {
		t.Fatalf("Banned user has been added even though it was already banned")
	}
}
func TestGetBannedUsers(t *testing.T) {
	config := appconfig.GetAppConfig()
	logger := logger.GetTikiLogger(config)

	dynamo, err := db.NewDynamoDBDriver(config, logger)
	if err != nil {
		t.Fatalf("Failed to create DynamoDB driver: %v", err)
	}

	brt := repositories.ProvideBanRepository(dynamo)
	_, err = brt.BanRepository.GetBannedUsers()
	if err != nil {
		t.Fatalf("Failed to get banned users: %v", err)
	}
	t.Logf("Banned users have been retrieved")
}
func TestGetBannedUserByEmail(t *testing.T) {
	config := appconfig.GetAppConfig()
	logger := logger.GetTikiLogger(config)

	dynamo, err := db.NewDynamoDBDriver(config, logger)
	if err != nil {
		t.Fatalf("Failed to create DynamoDB driver: %v", err)
	}

	brt := repositories.ProvideBanRepository(dynamo)
	_, err = brt.BanRepository.GetBannedUserByEmail("michael@jackson.com")
	if err != nil {
		t.Fatalf("Failed to get banned user by email: %v", err)
	}
	t.Logf("Banned user by email has been retrieved")
}

func TestUnbanUser(t *testing.T) {
	config := appconfig.GetAppConfig()
	logger := logger.GetTikiLogger(config)

	dynamo, err := db.NewDynamoDBDriver(config, logger)
	if err != nil {
		t.Fatalf("Failed to create DynamoDB driver: %v", err)
	}

	brt := repositories.ProvideBanRepository(dynamo)
	err = brt.BanRepository.UnbanUser("selami@sahin.org")
	if err != nil {
		t.Fatalf("Failed to unban user: %v", err)
	}
	t.Logf("User has been unbanned")
}
