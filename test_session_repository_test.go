package main

import (
	"testing"

	"github.com/bcpitutor/ostiki/appconfig"
	"github.com/bcpitutor/ostiki/db"
	"github.com/bcpitutor/ostiki/logger"
	"github.com/bcpitutor/ostiki/repositories"
)

func TestGetSessions(t *testing.T) {
	config := appconfig.GetAppConfig()
	logger := logger.GetTikiLogger(config)

	dynamo, err := db.NewDynamoDBDriver(config, logger)
	if err != nil {
		t.Fatalf("Failed to create DynamoDB driver: %v", err)
	}
	srt := repositories.ProvideSessionRepository(dynamo)

	sTypes := []string{"all", "active", "expired", "revoked"}

	for _, sType := range sTypes {
		_, err = srt.SessionRepository.GetSessions(sType)
		if err != nil {
			t.Fatalf("Failed to get %s sessions: %v", sType, err)
		}
		t.Logf("%s sessions have been retrieved", sType)
	}
}
