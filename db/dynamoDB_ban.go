package db

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/bcpitutor/ostiki/models"
)

// Ban Table
func (db DynamoDBDriver) AddBannedUser(bannedUser models.BannedUser) error {
	tableName := db.TableNames["ban_table"]
	dynamoDB := db.Client

	if db.IsUserBanned(bannedUser.UserEmail) {
		return fmt.Errorf("User %s is already banned", bannedUser.UserEmail)
	}

	item, err := attributevalue.MarshalMap(bannedUser)
	if err != nil {
		return fmt.Errorf("marshall error: %v", err)
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

func (db DynamoDBDriver) GetBannedUserByEmail(userEmail string) (models.BannedUser, error) {
	var bannedUser models.BannedUser
	tableName := db.TableNames["ban_table"]

	input := &dynamodb.ScanInput{
		TableName:        aws.String(tableName),
		FilterExpression: aws.String("UserEmail = :userEmail"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":userEmail": &types.AttributeValueMemberS{Value: userEmail},
		},
	}

	resp, err := db.Client.Scan(context.TODO(), input)
	if err != nil {
		return bannedUser, err
	}

	if resp.Count <= 0 {
		return bannedUser, fmt.Errorf("404")
	}

	err = attributevalue.UnmarshalMap(resp.Items[0], &bannedUser)
	if err != nil {
		fmt.Printf("failed to unmarshal dynamodb scan items, %v", err)
		return bannedUser, err
	}

	return bannedUser, nil
}

func (db DynamoDBDriver) GetBannedUsers() ([]models.BannedUser, error) {
	tableName := db.TableNames["ban_table"]
	var bannedUsers []models.BannedUser

	dynamoDB := db.Client

	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	resp, err := dynamoDB.Scan(context.TODO(), input)
	if err != nil {
		return bannedUsers, err
	}

	err = attributevalue.UnmarshalListOfMaps(resp.Items, &bannedUsers)
	if err != nil {
		fmt.Printf("failed to unmarshal Dynamodb Scan Items, %v", err)
		return bannedUsers, err
	}

	return bannedUsers, nil
}

func (db DynamoDBDriver) UnbanUser(userEmail string) error {
	tableName := db.TableNames["ban_table"]

	if !db.IsUserBanned(userEmail) {
		return fmt.Errorf("User %s is not banned", userEmail)
	}

	key := map[string]types.AttributeValue{
		"UserEmail": &types.AttributeValueMemberS{Value: userEmail},
	}
	params := dynamodb.DeleteItemInput{
		Key:       key,
		TableName: &tableName,
	}

	_, err := db.Client.DeleteItem(context.TODO(), &params)
	if err != nil {
		return err
	}

	return nil
}

func (db DynamoDBDriver) IsUserBanned(userEmail string) bool {
	bannedUser, err := db.GetBannedUserByEmail(userEmail)
	if err != nil {
		return false
	}

	if bannedUser.UserEmail == userEmail {
		return true
	} else {
		return false
	}
}

func (db DynamoDBDriver) HasUserAccessToBanInfo(userEmail string) bool {
	return db.IsUserInTikiadmins(userEmail)
}
