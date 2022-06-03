package db

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/tiki-systems/tikiserver/models"
)

func (db DynamoDBDriver) GetAllDomains() ([]models.TicketDomain, error) {
	tableName := db.TableNames["domain_table"]

	input := &dynamodb.ScanInput{
		TableName: &tableName,
	}
	resp, err := db.Client.Scan(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	domains := []models.TicketDomain{}

	err = attributevalue.UnmarshalListOfMaps(resp.Items, &domains)
	if err != nil {
		return nil, err
	}

	return domains, nil
}

func (db DynamoDBDriver) DoesTicketDomainExist(ticketDomainPath string) bool {
	doesExist := false

	result, err := db.GetDomain(ticketDomainPath)
	if err != nil {
		return doesExist
	}

	var domain models.TicketDomain
	domainBytes, err := json.Marshal(result)
	if err != nil {
		return doesExist
	}

	if err := json.Unmarshal(domainBytes, &domain); err != nil {
		return doesExist
	}

	if domain.DomainPath == ticketDomainPath {
		doesExist = true
	}

	return doesExist
}

func (db DynamoDBDriver) CanUserPerformDomainOperation(userEmail string, operationType string) bool {
	if db.IsUserInTikiadmins(userEmail) {
		return true
	}

	switch operationType {
	case "create", "add", "delete":
		groups, err := db.GetGroupsOfUser(userEmail)
		if err != nil {
			return false
		}
		for _, group := range groups {
			accPerms := group.AccessPerms
			if accPerms.Group[operationType] {
				return true
			}
		}
		return false
	case "show", "info", "list":
		return true
	default:
		return false
	}
}

func (db DynamoDBDriver) GetDomain(domainPath string) (models.TicketDomain, error) {
	tableName := db.TableNames["domain_table"]

	var domain models.TicketDomain
	key := map[string]types.AttributeValue{
		"DomainPath": &types.AttributeValueMemberS{Value: domainPath},
	}

	params := dynamodb.GetItemInput{
		Key:       key,
		TableName: &tableName,
	}

	resp, err := db.Client.GetItem(context.TODO(), &params)
	if err != nil {
		return domain, err
	}

	if len(resp.Item) == 0 {
		return domain, fmt.Errorf("the domain [%v] does not exist", domainPath)
	}

	err = attributevalue.UnmarshalMap(resp.Item, &domain)
	if err != nil {
		return domain, fmt.Errorf("%v", err)
	}

	return domain, nil
}

func (db DynamoDBDriver) CreateDomain(domain models.TicketDomain) error {
	item, err := attributevalue.MarshalMap(domain)
	if err != nil {
		return fmt.Errorf("Error while marshalling domain object: %v", err)
	}

	tableName := db.TableNames["domain_table"]
	timeNow := time.Now().Unix()
	domain.CreatedAt = timeNow
	domain.UpdatedAt = timeNow

	_, err = db.Client.PutItem(
		context.TODO(),
		&dynamodb.PutItemInput{
			Item:      item,
			TableName: aws.String(tableName),
		},
	)

	if err != nil {
		return fmt.Errorf("Error while creating domain: %v", err)
	}

	return nil
}

func (db DynamoDBDriver) DeleteDomain(domainPath string) error {
	tableName := db.TableNames["domain_table"]
	key := map[string]types.AttributeValue{
		"DomainPath": &types.AttributeValueMemberS{Value: domainPath},
	}

	_, err := db.Client.DeleteItem(
		context.TODO(),
		&dynamodb.DeleteItemInput{
			Key:       key,
			TableName: &tableName,
		},
	)

	if err != nil {
		return fmt.Errorf("Error while deleting domain: %v", err)
	}

	return nil
}
