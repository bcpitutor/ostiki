package main

import (
	"testing"

	"github.com/tiki-systems/tikiserver/appconfig"
	"github.com/tiki-systems/tikiserver/db"
	"github.com/tiki-systems/tikiserver/logger"
	"github.com/tiki-systems/tikiserver/models"
	"github.com/tiki-systems/tikiserver/repositories"
)

func TestIsUserInTikiadmins(t *testing.T) {
	config := appconfig.GetAppConfig()
	logger := logger.GetTikiLogger(config)

	dynamo, err := db.NewDynamoDBDriver(config, logger)
	if err != nil {
		t.Fatalf("Failed to create DynamoDB driver: %v", err)
	}
	grt := repositories.ProvideGroupRepository(dynamo, nil)
	trt := repositories.ProvideTicketRepository(dynamo, nil)
	prt := repositories.ProvidePermissionRepository(dynamo, grt.GroupRepository, trt.TicketRepository)

	if prt.PermissionRepository.IsUserInTikiadmins("ozgur.demir@itutor.com") {
		t.Logf("User ozgur.demir@itutor.com is in tikiadmins")
	} else {
		t.Errorf("User is not in tikiadmins, but should be")
	}

	if prt.PermissionRepository.IsUserInTikiadmins("john@itutor.com") {
		t.Errorf("User john@itutor.com is in tikiadmins, but should not be")
	} else {
		t.Logf("User is not in tikiadmins")
	}
}

func TestCanUserPerformTicketOperation(t *testing.T) {
	config := appconfig.GetAppConfig()
	logger := logger.GetTikiLogger(config)

	dynamo, err := db.NewDynamoDBDriver(config, logger)
	if err != nil {
		t.Fatalf("Failed to create DynamoDB driver: %v", err)
	}
	grt := repositories.ProvideGroupRepository(dynamo, nil)
	trt := repositories.ProvideTicketRepository(dynamo, nil)
	prt := repositories.ProvidePermissionRepository(dynamo, grt.GroupRepository, trt.TicketRepository)

	if prt.PermissionRepository.CanUserPerformTicketOperation("ozgur.demir@itutor.com", models.Operation.Delete) {
		t.Logf("User ozgur.demir@itutor.com can perform delete operation")
	} else {
		t.Errorf("User ozgur.demir@itutor.com cannot perform delete operation")
	}

	if !prt.PermissionRepository.CanUserPerformTicketOperation("john@itutor.com", models.Operation.Delete) {
		t.Logf("User john@itutor.com cannot perform delete operation")
	} else {
		t.Errorf("User john@itutor.com cannot perform delete operation")
	}
}

func TestCanUserAccessToTicket(t *testing.T) {
	config := appconfig.GetAppConfig()
	logger := logger.GetTikiLogger(config)

	dynamo, err := db.NewDynamoDBDriver(config, logger)
	if err != nil {
		t.Fatalf("Failed to create DynamoDB driver: %v", err)
	}
	grt := repositories.ProvideGroupRepository(dynamo, nil)
	trt := repositories.ProvideTicketRepository(dynamo, nil)
	prt := repositories.ProvidePermissionRepository(dynamo, grt.GroupRepository, trt.TicketRepository)

	if prt.PermissionRepository.CanUserAccessToTicket("ozgur.demir@itutor.com", "tickets/itutor/infra/prod/admin/aws") {
		t.Logf("User ozgur.demir@itutor.com can access to ticket")
	} else {
		t.Errorf("User ozgur.demir@itutor.com can access to ticket")
	}

	if !prt.PermissionRepository.CanUserAccessToTicket("john@itutor.com", "tickets/itutor/infra/prod/admin/aws") {
		t.Logf("User john@itutor.com cannot access to ticket")
	} else {
		t.Errorf("User john@itutor.com can access to ticket")
	}
}

func TestIsUserAllowedByDomainScope(t *testing.T) {
	config := appconfig.GetAppConfig()
	logger := logger.GetTikiLogger(config)

	dynamo, err := db.NewDynamoDBDriver(config, logger)
	if err != nil {
		t.Fatalf("Failed to create DynamoDB driver: %v", err)
	}
	grt := repositories.ProvideGroupRepository(dynamo, nil)
	trt := repositories.ProvideTicketRepository(dynamo, nil)
	prt := repositories.ProvidePermissionRepository(dynamo, grt.GroupRepository, trt.TicketRepository)

	if prt.PermissionRepository.IsUserAllowedByDomainScope(
		"ozgur.demir@itutor.com",
		"tickets/itutor/infra/prod/admin/aws",
		models.DomainScopeOperation.CreateDomain) {
		t.Logf("User ozgur.demir@itutor.com is allowed by domain scope")
	} else {
		t.Errorf("User ozgur.demir@itutor.com is not allowed by domain scope")
	}

	if !prt.PermissionRepository.IsUserAllowedByDomainScope(
		"john@itutor.com",
		"tickets/itutor/infra/prod/admin/aws",
		models.DomainScopeOperation.CreateDomain) {
		t.Logf("User john@itutor.com is allowed by domain scope")
	} else {
		t.Errorf("User john@itutor.com is not allowed by domain scope")
	}

}
