package main

import (
	"fmt"
	"strings"

	"github.com/tiki-systems/tikiserver/apiserver"
	"github.com/tiki-systems/tikiserver/appconfig"
	"github.com/tiki-systems/tikiserver/cache"
	"github.com/tiki-systems/tikiserver/db"
	"github.com/tiki-systems/tikiserver/logger"
	"github.com/tiki-systems/tikiserver/repositories"
	"github.com/tiki-systems/tikiserver/services"
	"go.uber.org/dig"
)

func main() {
	container := dig.New()

	config := appconfig.GetAppConfig()

	// Handle other db types, if implemented.
	if strings.ToLower(config.TikiDBConfig.DbType) == "dynamodb" {
		container.Provide(db.NewDynamoDBDriver)
	}

	container.Provide(appconfig.GetAppConfig)
	container.Provide(logger.GetTikiLogger)
	container.Provide(services.GetAWS)

	container.Provide(repositories.ProvideSessionRepository)
	container.Provide(repositories.ProvideBanRepository)
	container.Provide(repositories.ProvideDomainRepository)
	container.Provide(repositories.ProvideGroupRepository)
	container.Provide(repositories.ProvideTicketRepository)
	container.Provide(repositories.ProvidePermissionRepository)

	// Handle other cache types, if implemented.
	if strings.ToLower(config.TikiInMemoryStoreConfig.StoreType) == "hazelcast" {
		fmt.Printf("Loading Hazelcast driver into the container\n")
		container.Provide(cache.NewHazelcastDriver)
	} else {
		fmt.Printf("No cache will be used\n")
		container.Provide(cache.NewNoCacheDriver)
	}

	container.Provide(apiserver.ProvideServer)

	// initiate server
	err := container.Invoke(
		func(server *apiserver.Server) {
			server.Run()
		},
	)
	if err != nil {
		fmt.Println("Error starting server: ", err)
	}
}
