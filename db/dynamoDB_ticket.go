package db

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/bcpitutor/ostiki/models"
)

// Ticket Table
func (db DynamoDBDriver) GetAllTickets() ([]models.Ticket, error) {
	tableName := db.TableNames["ticket_table"]

	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	resp, err := db.Client.Scan(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	tickets := []models.Ticket{}

	err = attributevalue.UnmarshalListOfMaps(resp.Items, &tickets)
	if err != nil {
		return nil, err
	}
	return tickets, nil
}

func (db DynamoDBDriver) CanUserPerformTicketOperation(userEmail string, operationType string) bool {
	if db.IsUserInTikiadmins(userEmail) {
		return true
	}

	switch operationType {
	case models.Operation.Delete, models.Operation.Create, models.Operation.SetSecret, models.Operation.GetSecret:
		groups, err := db.GetGroupsOfUser(userEmail)
		if err != nil {
			return false
		}

		for _, group := range groups {
			accPerms := group.AccessPerms
			if accPerms.Ticket[operationType] {
				return true
			}
		}

		return false
	case models.Operation.Show, models.Operation.Info, models.Operation.List:
		return true
	default:
		return false
	}
}

func (db DynamoDBDriver) DeleteTicket(ticketPath string, ticketType string) error {
	tableName := db.TableNames["ticket_table"]

	key := map[string]types.AttributeValue{
		"TicketPath": &types.AttributeValueMemberS{
			Value: ticketPath,
		},
		"TicketType": &types.AttributeValueMemberS{
			Value: ticketType,
		},
	}
	params := dynamodb.DeleteItemInput{Key: key, TableName: &tableName}
	_, err := db.Client.DeleteItem(context.TODO(), &params)
	if err != nil {
		return err
	}

	return nil
}

func (db DynamoDBDriver) IsUserAllowedByDomainScope(userEmail string, ticketOrDomainPath string, domainScopeOperation string) bool {
	isUserAllowed := false

	groups, err := db.GetGroupsOfUser(userEmail)
	if err != nil {
		return false
	}

	var domainScopeList []string
	for _, v := range groups {
		domainScopeList = append(domainScopeList, v.DomainScope.Root)
	}

	switch domainScopeOperation {
	case "createTicket", "deleteTicket", "assumeTicket":
		if db.IsUserInTikiadmins(userEmail) {
			isUserAllowed = true
			return isUserAllowed
		}

		for _, v := range domainScopeList {
			if strings.HasPrefix(ticketOrDomainPath, v) {
				isUserAllowed = true
				return isUserAllowed
			}
		}

		return isUserAllowed

	case "createDomain", "deleteDomain":
		if db.IsUserInTikiadmins(userEmail) {
			return true
		}

		for _, v := range domainScopeList {
			if strings.HasPrefix(ticketOrDomainPath, v) {
				isUserAllowed = true
				return isUserAllowed
			}
		}
		return isUserAllowed

	default:
		return isUserAllowed
	}

}
func (db DynamoDBDriver) CanUserAccessToTicket(userEmail string, ticketPath string) bool {
	if db.IsUserInTikiadmins(userEmail) {
		return true
	}

	groups, err := db.GetGroupNamesOfUser(userEmail)
	if err != nil {
		return false
	}

	ticket, err := db.QueryTicketByPath(ticketPath)
	if err != nil {
		return false
	}

	for _, gOwner := range ticket.OwnersGroup {
		for _, g := range groups {
			if g == gOwner {
				return true
			}
		}
	}

	return false
}

func (db DynamoDBDriver) QueryTicketByPath(ticketPath string) (models.Ticket, error) {
	tableName := db.TableNames["ticket_table"]

	var resultTicket models.Ticket
	out, err := db.Client.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		KeyConditionExpression: aws.String("TicketPath = :ticketPath"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":ticketPath": &types.AttributeValueMemberS{Value: ticketPath},
		},
	})
	if err != nil {
		return resultTicket, err
	}

	if len(out.Items) == 0 {
		return resultTicket, fmt.Errorf("the ticket [%v] does not exist", ticketPath)
	}

	err = attributevalue.UnmarshalMap(out.Items[0], &resultTicket)
	if err != nil {
		return resultTicket, err
	}

	return resultTicket, nil
}

func (db DynamoDBDriver) DoesTicketExist(ticketPath string) bool {
	doesExist := false

	result, err := db.QueryTicketByPath(ticketPath)
	if err != nil {
		return doesExist
	}

	var ticket models.Ticket

	ticketBytes, err := json.Marshal(result)
	if err != nil {
		return doesExist
	}

	if err := json.Unmarshal(ticketBytes, &ticket); err != nil {
		return doesExist
	}

	if ticket.TicketPath == ticketPath {
		doesExist = true
	}

	return doesExist
}

func (db DynamoDBDriver) CreateTicket(ticket models.Ticket) error {
	tableName := db.TableNames["ticket_table"]
	timeNow := time.Now().Unix()
	dynamoDB := db.Client

	ticket.CreatedAt = strconv.FormatInt(timeNow, 10)
	ticket.UpdatedAt = strconv.FormatInt(timeNow, 10)

	item, err := attributevalue.MarshalMap(ticket)
	if err != nil {
		return err
	}

	_, err = dynamoDB.PutItem(
		context.TODO(),
		&dynamodb.PutItemInput{
			Item:      item,
			TableName: aws.String(tableName),
		},
	)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}

func (db DynamoDBDriver) SetTicketSecret(ticketPath string, secretData string) error {
	ticket, err := db.QueryTicketByPath(ticketPath)
	if err != nil {
		return err
	}

	if ticket.TicketType != "secretTicket" {
		return fmt.Errorf("This ticket is not a secret ticket")
	}

	ticket.SecretData = secretData
	err = db.CreateTicket(ticket)
	if err != nil {
		return err
	}

	return nil
}

func (db DynamoDBDriver) GetTicketSecret(ticketPath string) (string, error) {
	ticket, err := db.QueryTicketByPath(ticketPath)
	if err != nil {
		return "", err
	}

	if ticket.TicketType != "secretTicket" {
		return "", fmt.Errorf("This ticket is not a secret ticket")
	}

	return ticket.SecretData, nil
}
