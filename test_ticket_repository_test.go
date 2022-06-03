package main

import (
	"testing"

	"github.com/tiki-systems/tikiserver/appconfig"
	"github.com/tiki-systems/tikiserver/db"
	"github.com/tiki-systems/tikiserver/logger"
	"github.com/tiki-systems/tikiserver/models"
	"github.com/tiki-systems/tikiserver/repositories"
	"github.com/tiki-systems/tikiserver/services"
)

func TestGetAllTickets(t *testing.T) {
	config := appconfig.GetAppConfig()
	logger := logger.GetTikiLogger(config)

	dynamo, err := db.NewDynamoDBDriver(config, logger)
	if err != nil {
		t.Fatalf("Failed to create DynamoDB driver: %v", err)
	}

	trr := repositories.ProvideTicketRepository(dynamo, nil)
	_, err = trr.TicketRepository.GetAllTickets()
	if err != nil {
		t.Fatalf("Failed to get all tickets: %v", err)
	}
	t.Logf("All tickets have been retrieved")
}
func TestQueryTicketByPath(t *testing.T) {
	config := appconfig.GetAppConfig()
	logger := logger.GetTikiLogger(config)

	dynamo, err := db.NewDynamoDBDriver(config, logger)
	if err != nil {
		t.Fatalf("Failed to create DynamoDB driver: %v", err)
	}

	trr := repositories.ProvideTicketRepository(dynamo, nil)
	_, err = trr.TicketRepository.QueryTicketByPath("tickets/itutor/infra/development/admin/aws")
	if err != nil {
		t.Fatalf("Failed to query ticket by path: %v", err)
	}
	t.Logf("Ticket found")
}
func TestDoesTicketExist(t *testing.T) {
	config := appconfig.GetAppConfig()
	logger := logger.GetTikiLogger(config)

	dynamo, err := db.NewDynamoDBDriver(config, logger)
	if err != nil {
		t.Fatalf("Failed to create DynamoDB driver: %v", err)
	}

	trr := repositories.ProvideTicketRepository(dynamo, nil)
	exists := trr.TicketRepository.DoesTicketExist("tickets/itutor/infra/development/admin/aws")
	if !exists {
		t.Fatalf("Ticket does not exist")
	}
	t.Logf("Ticket exists")
}
func TestCreateTicket(t *testing.T) {
	config := appconfig.GetAppConfig()
	logger := logger.GetTikiLogger(config)

	dynamo, err := db.NewDynamoDBDriver(config, logger)
	if err != nil {
		t.Fatalf("Failed to create DynamoDB driver: %v", err)
	}

	trr := repositories.ProvideTicketRepository(dynamo, nil)
	ticket := models.Ticket{
		TicketPath: "tickets/itutor/test-tickets",
		TicketType: "awsTicket",
		AwsAssumeRole: models.AwsAssumeRole{
			RoleArn: "arn:aws:iam::123456789012:role/tiki-test-role",
			Ttl:     3600,
		},
		AwsPermissions: models.AwsPermissions{
			Effect:   "Allow",
			Action:   []string{"ec2:DescribeInstances", "ec2:DescribeImages"},
			Resource: "*",
		},
		OwnersGroup:  []string{"tiki-test-group"},
		TicketInfo:   "Test ticket",
		TicketRegion: "us-west-1",
	}

	err = trr.TicketRepository.CreateTicket(ticket)
	if err != nil {
		t.Fatalf("Failed to create ticket: %v", err)
	}
	t.Logf("Ticket created")
}
func TestDeleteTicket(t *testing.T) {
	config := appconfig.GetAppConfig()
	logger := logger.GetTikiLogger(config)

	dynamo, err := db.NewDynamoDBDriver(config, logger)
	if err != nil {
		t.Fatalf("Failed to create DynamoDB driver: %v", err)
	}

	trr := repositories.ProvideTicketRepository(dynamo, nil)
	err = trr.TicketRepository.DeleteTicket("tickets/itutor/test-tickets", "awsTicket")
	if err != nil {
		t.Fatalf("Failed to delete ticket: %v", err)
	}
	t.Logf("Ticket deleted")
}
func TestSetTicketSecret(t *testing.T) {
	config := appconfig.GetAppConfig()
	logger := logger.GetTikiLogger(config)
	aws, err := services.GetAWS()
	if err != nil {
		t.Fatalf("Failed to create AWS service: %v", err)
	}

	dynamo, err := db.NewDynamoDBDriver(config, logger)
	if err != nil {
		t.Fatalf("Failed to create DynamoDB driver: %v", err)
	}

	trr := repositories.ProvideTicketRepository(dynamo, nil)
	ticket := models.Ticket{
		TicketPath:   "tickets/itutor/test-secret-ticket",
		TicketType:   "secretTicket",
		OwnersGroup:  []string{"tiki-test-group"},
		TicketInfo:   "Test ticket",
		TicketRegion: "us-west-1",
	}

	err = trr.TicketRepository.CreateTicket(ticket)
	if err != nil {
		t.Fatalf("Failed to create secret ticket: %v", err)
	}

	encryptedData, err := aws.GetEncryptedSecret("veryVeryS3cret")
	if err != nil {
		t.Fatalf("Failed to encrypt secret: %v", err)
	}

	err = trr.TicketRepository.SetTicketSecret("tickets/itutor/test-secret-ticket", encryptedData)
	if err != nil {
		t.Fatalf("Failed to set secret ticket: %v", err)
	}
	t.Logf("Ticket secret set")
}

func TestGetTicketSecret(t *testing.T) {
	config := appconfig.GetAppConfig()
	logger := logger.GetTikiLogger(config)
	aws, err := services.GetAWS()
	if err != nil {
		t.Fatalf("Failed to create AWS service: %v", err)
	}

	dynamo, err := db.NewDynamoDBDriver(config, logger)
	if err != nil {
		t.Fatalf("Failed to create DynamoDB driver: %v", err)
	}

	trr := repositories.ProvideTicketRepository(dynamo, nil)

	encrypted, err := trr.TicketRepository.GetTicketSecret("tickets/itutor/test-secret-ticket")
	if err != nil {
		t.Fatalf("Failed to get secret ticket: %v", err)
	}

	decrypted, err := aws.GetDecryptedText(encrypted)
	if err != nil {
		t.Fatalf("Failed to decrypt secret: %v", err)
	}

	if decrypted == "veryVeryS3cret" {
		t.Logf("Ticket secret retrieved")
	} else {
		t.Fatalf("Ticket secret mismatch")
	}
}
