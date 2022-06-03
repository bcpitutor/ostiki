package db

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/tiki-systems/tikiserver/appconfig"
	"github.com/tiki-systems/tikiserver/logger"
	"github.com/tiki-systems/tikiserver/models"
)

type DynamoDBDriver struct {
	Client     *dynamodb.Client
	TableNames map[string]string
}

func NewDynamoDBDriver(appconfig *appconfig.AppConfig, tikilogger *logger.TikiLogger) (models.DBLayer, error) {
	staticProvider := credentials.NewStaticCredentialsProvider(
		appconfig.TikiDBConfig.DbProfileId,
		appconfig.TikiDBConfig.DbProfileSecret,
		"",
	)

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(appconfig.TikiDBConfig.DbRegion),
		config.WithCredentialsProvider(staticProvider),
	)

	if err != nil {
		fmt.Println("unable to load SDK config, " + err.Error())
		return nil, err
	}

	suffix := appconfig.Deployment
	if suffix == "local" {
		suffix = "dev"
	}

	return DynamoDBDriver{
		Client: dynamodb.NewFromConfig(cfg),
		TableNames: map[string]string{
			"session_table": fmt.Sprintf("tikiserver_sess-%s", suffix),
			"ban_table":     fmt.Sprintf("tikiserver_banned-%s", suffix),
			"group_table":   fmt.Sprintf("tikiserver_group-%s", suffix),
			"ticket_table":  fmt.Sprintf("tikiserver_ticket-%s", suffix),
			"domain_table":  fmt.Sprintf("tikiserver_domain-%s", suffix),
		},
	}, nil
}

func (db DynamoDBDriver) DBType() string {
	return "DynamoDB"
}
