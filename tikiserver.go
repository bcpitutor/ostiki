package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/bcpitutor/ostiki/apiserver"
	"github.com/bcpitutor/ostiki/appconfig"
	"github.com/bcpitutor/ostiki/cache"
	"github.com/bcpitutor/ostiki/db"
	"github.com/bcpitutor/ostiki/logger"
	"github.com/bcpitutor/ostiki/repositories"
	"github.com/bcpitutor/ostiki/services"
	"github.com/bcpitutor/ostiki/version"
	"go.uber.org/dig"
)

func main() {
	config := appconfig.GetAppConfig()

	if len(os.Args) > 1 && os.Args[1] == "--dumpConfig" {
		fmt.Printf("Config for version: %s\n", version.VersionDetails.Version)
		fmt.Printf("%+v\n", config)
		os.Exit(0)
	}

	container := dig.New()

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
