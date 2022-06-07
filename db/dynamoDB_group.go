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
	"github.com/bcpitutor/ostiki/models"
)

// Group Table
func (db DynamoDBDriver) DoesGroupExist(groupName string) bool {
	doesExist := false

	result, err := db.GetGroup(groupName)
	if err != nil {
		return doesExist
	}

	var group models.TicketGroup
	groupBytes, err := json.Marshal(result)
	if err != nil {
		return doesExist
	}

	if err := json.Unmarshal(groupBytes, &group); err != nil {
		return doesExist
	}

	if group.GroupName == groupName {
		doesExist = true
	}

	return doesExist
}

func (db DynamoDBDriver) GetAllGroups() ([]models.TicketGroup, error) {
	tableName := db.TableNames["group_table"]
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	resp, err := db.Client.Scan(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	groups := []models.TicketGroup{}

	err = attributevalue.UnmarshalListOfMaps(resp.Items, &groups)
	if err != nil {
		return nil, err
	}

	return groups, nil
}

func (db DynamoDBDriver) GetGroup(groupName string) (models.TicketGroup, error) {
	tableName := db.TableNames["group_table"]

	var group models.TicketGroup
	key := map[string]types.AttributeValue{
		"GroupName": &types.AttributeValueMemberS{Value: groupName},
	}
	params := dynamodb.GetItemInput{
		Key:       key,
		TableName: &tableName,
	}

	resp, err := db.Client.GetItem(context.TODO(), &params)
	if err != nil {
		return group, err
	}

	if len(resp.Item) == 0 {
		return group, fmt.Errorf("the group [%s] does not exist", groupName)
	}

	err = attributevalue.UnmarshalMap(resp.Item, &group)
	if err != nil {
		return group, err
	}

	return group, nil
}

func (db DynamoDBDriver) DeleteGroup(groupName string) error {
	tableName := db.TableNames["group_table"]

	key := map[string]types.AttributeValue{
		"GroupName": &types.AttributeValueMemberS{Value: groupName},
	}
	params := dynamodb.DeleteItemInput{
		Key:       key,
		TableName: &tableName,
	}

	_, err := db.Client.DeleteItem(context.TODO(), &params)
	if err != nil {
		return fmt.Errorf("failed to delete group item on DynamoDB: %v", err)
	}

	return nil
}

func (db DynamoDBDriver) CreateGroup(newGroup models.TicketGroup) error {
	tableName := db.TableNames["group_table"]

	item, err := attributevalue.MarshalMap(newGroup)
	if err != nil {
		return err
	}

	timeNow := time.Now().Unix()
	newGroup.CreatedAt = timeNow
	newGroup.UpdatedAt = timeNow

	_, err = db.Client.PutItem(
		context.TODO(),
		&dynamodb.PutItemInput{
			Item:      item,
			TableName: aws.String(tableName),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create group item on DynamoDB: %v", err)
	}

	return nil
}

func (db DynamoDBDriver) CanUserPerformGroupOperation(userEmail string, operationType string) bool {
	if db.IsUserInTikiadmins(userEmail) {
		return true
	}

	switch operationType {
	case "delete", "create", "show", "addMember", "delMember", "info":
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
	case "list":
		return true
	default:
		return false
	}
}

func (db DynamoDBDriver) IsUserInTikiadmins(userEmail string) bool {
	groupList, err := db.GetGroupNamesOfUser(userEmail)
	if err != nil {
		return false
	}

	for _, group := range groupList {
		if group == "tikiadmins" {
			return true
		}
	}

	return false
}

func (db DynamoDBDriver) GetGroupNamesOfUser(userEmail string) ([]string, error) {
	var userGroupNamesList []string
	allGroups, err := db.GetAllGroups()
	if err != nil {
		return nil, err
	}

	var gMembers []string
	for _, group := range allGroups {
		gMembers = group.GroupMembers

		for _, member := range gMembers {
			if member == userEmail {
				userGroupNamesList = append(userGroupNamesList, group.GroupName)
			}
		}
	}
	return userGroupNamesList, nil
}

func (db DynamoDBDriver) GetGroupsOfUser(userEmail string) ([]models.TicketGroup, error) {
	var groups []models.TicketGroup
	allGroups, err := db.GetAllGroups()
	if err != nil {
		return nil, err
	}

	var gMembers []string
	for _, group := range allGroups {
		gMembers = group.GroupMembers

		for _, member := range gMembers {
			if member == userEmail {
				groups = append(groups, group)
			}
		}
	}
	return groups, nil
}

func (db DynamoDBDriver) GetGroupMembers(groupName string) ([]string, error) {
	group, err := db.GetGroup(groupName)
	if err != nil {
		return nil, err
	}

	return group.GroupMembers, nil
}

func (db DynamoDBDriver) IsUserMemberOfGroup(member string, groupName string) bool {
	isMember := false

	members, err := db.GetGroupMembers(groupName)
	if err != nil {
		return false
	}

	for _, v := range members {
		if v == member {
			isMember = true
			break
		}
	}

	return isMember
}

func (db DynamoDBDriver) AddMemberToGroup(newMember string, groupName string, changedBy string) error {
	tableName := db.TableNames["group_table"]

	group, err := db.GetGroup(groupName)
	if err != nil {
		return fmt.Errorf("Error while retrieving group: %v", err)
	}

	members := group.GroupMembers

	alreadyAMember := false
	for _, v := range members {
		if v == newMember {
			alreadyAMember = true
			break
		}
	}

	if alreadyAMember {
		return fmt.Errorf("User %s is already a member of group %s", newMember, groupName)
	}

	members = append(members, newMember)
	newMemberValues := []types.AttributeValue{}
	for _, v := range members {
		newMemberValues = append(newMemberValues, &types.AttributeValueMemberS{Value: v})
	}

	key := map[string]types.AttributeValue{
		"GroupName": &types.AttributeValueMemberS{
			Value: groupName,
		},
	}
	timeNowString := fmt.Sprintf("%v", time.Now().Unix())

	updateSection := map[string]types.AttributeValueUpdate{
		"GroupMembers":    {Action: types.AttributeActionPut, Value: &types.AttributeValueMemberL{Value: newMemberValues}},
		"UpdatedBy":       {Action: types.AttributeActionPut, Value: &types.AttributeValueMemberS{Value: changedBy}},
		"UpdatedAt":       {Action: types.AttributeActionPut, Value: &types.AttributeValueMemberN{Value: timeNowString}},
		"LastAddedMember": {Action: types.AttributeActionPut, Value: &types.AttributeValueMemberS{Value: newMember}},
	}

	input := dynamodb.UpdateItemInput{
		Key:              key,
		TableName:        aws.String(tableName),
		AttributeUpdates: updateSection,
	}

	_, err = db.Client.UpdateItem(context.TODO(), &input)
	if err != nil {
		return fmt.Errorf("Error while updating group: %v", err)
	}

	return nil
}

func (db DynamoDBDriver) DelMemberFromGroup(memberToDelete string, groupName string, changedBy string) error {
	tableName := db.TableNames["group_table"]

	group, err := db.GetGroup(groupName)
	if err != nil {
		return fmt.Errorf("Error while retrieving group: %v", err)
	}

	members := group.GroupMembers
	alreadyAMember := false
	var idx int

	for k, v := range members {
		if v == memberToDelete {
			alreadyAMember = true
			idx = k
			break
		}
	}

	if !alreadyAMember {
		return fmt.Errorf("User %s is not a member of group %s", memberToDelete, groupName)
	}

	members = append(members[:idx], members[idx+1:]...)
	key := map[string]types.AttributeValue{
		"GroupName": &types.AttributeValueMemberS{
			Value: groupName,
		},
	}

	timeNowString := fmt.Sprintf("%v", time.Now().Unix())

	updateSection := map[string]types.AttributeValueUpdate{
		"GroupMembers":    {Action: types.AttributeActionPut, Value: &types.AttributeValueMemberSS{Value: members}},
		"UpdatedBy":       {Action: types.AttributeActionPut, Value: &types.AttributeValueMemberS{Value: changedBy}},
		"UpdatedAt":       {Action: types.AttributeActionPut, Value: &types.AttributeValueMemberN{Value: timeNowString}},
		"LastAddedMember": {Action: types.AttributeActionPut, Value: &types.AttributeValueMemberS{Value: memberToDelete}},
	}

	params := dynamodb.UpdateItemInput{
		Key:       key,
		TableName: &tableName, AttributeUpdates: updateSection,
	}

	_, err = db.Client.UpdateItem(context.TODO(), &params)
	if err != nil {
		return fmt.Errorf("Error while updating group: %v", err)
	}

	return nil
}
