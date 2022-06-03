package db

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/tiki-systems/tikiserver/models"
)

// Session Table
func (db DynamoDBDriver) CreateSession(session *models.Session) error {
	item, err := attributevalue.MarshalMap(session)
	if err != nil {
		return err
	}

	tableName := db.TableNames["session_table"]
	dynamoDB := db.Client

	_, err = dynamoDB.PutItem(
		context.TODO(),
		&dynamodb.PutItemInput{
			Item:      item,
			TableName: aws.String(tableName),
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (db DynamoDBDriver) GetSessionByRefreshToken(rtoken string) (models.Session, error) {
	var sess models.Session
	tableName := db.TableNames["session_table"]
	dynamoDB := db.Client

	resp, err := dynamoDB.Scan(
		context.TODO(),
		&dynamodb.ScanInput{
			TableName:        &tableName,
			FilterExpression: aws.String("RefreshToken = :refreshToken and IsRevoked = :isRevoked"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":refreshToken": &types.AttributeValueMemberS{Value: rtoken},
				":isRevoked":    &types.AttributeValueMemberBOOL{Value: false},
			},
		},
	)
	if err != nil {
		return sess, err
	}

	if len(resp.Items) == 0 {
		return sess, fmt.Errorf("the token does not exist.")
	}

	if err := attributevalue.UnmarshalMap(resp.Items[0], &sess); err != nil {
		return sess, err
	}

	return sess, nil
}

func (db DynamoDBDriver) UpdateSession(prevToken string, currentToken string, currentTokenExpires int64, refreshToken string) bool {
	session, err := db.GetSessionByRefreshToken(refreshToken)
	if err != nil {
		return false
	}

	renewTime := time.Now().Unix()
	// will be updated variables
	session.PreviousIdToken = prevToken
	session.Details = fmt.Sprintf("Token is renewed at [%d]", renewTime)
	session.RefreshToken = refreshToken
	newTimes := session.Rtimes + 1
	session.Rtimes = newTimes
	session.IdToken = currentToken
	session.Expire = strconv.FormatInt(currentTokenExpires, 10)

	err = db.CreateSession(&session)
	if err != nil {
		return false
	}

	return true
}

func (db DynamoDBDriver) GetSessions(scanType string) ([]models.Session, error) {
	sessions := []models.Session{}

	//sessTable := "tikiserver_sess-dev"
	tableName := db.TableNames["session_table"]
	dynamoDB := db.Client
	input := &dynamodb.ScanInput{}

	switch scanType {
	case "active":
		timeNow := strconv.FormatInt(time.Now().Unix(), 10)
		input = &dynamodb.ScanInput{
			TableName:        aws.String(tableName),
			FilterExpression: aws.String("IsRevoked = :isRevoked and Expire >= :expire"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":isRevoked": &types.AttributeValueMemberBOOL{Value: false},
				":expire":    &types.AttributeValueMemberS{Value: timeNow},
			},
		}

	case "expired":
		timeNow := strconv.FormatInt(time.Now().Unix(), 10)
		input = &dynamodb.ScanInput{
			TableName:        aws.String(tableName),
			FilterExpression: aws.String("IsRevoked = :isRevoked and Expire <= :expire"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":isRevoked": &types.AttributeValueMemberBOOL{Value: false},
				":expire":    &types.AttributeValueMemberS{Value: timeNow},
			},
		}
	case "revoked":
		input = &dynamodb.ScanInput{
			TableName:        aws.String(tableName),
			FilterExpression: aws.String("IsRevoked = :isRevoked"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":isRevoked": &types.AttributeValueMemberBOOL{Value: true},
			},
		}
	default:
		input = &dynamodb.ScanInput{
			TableName: aws.String(tableName),
		}
	}

	resp, err := dynamoDB.Scan(context.TODO(), input)
	if err != nil {
		return sessions, err
	}

	err = attributevalue.UnmarshalListOfMaps(resp.Items, &sessions)
	if err != nil {
		fmt.Printf("failed to unmarshal Dynamodb Scan Items, %v", err)
		return sessions, err
	}

	return sessions, nil
}
