package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/bcpitutor/ostiki/apiserver"
	"github.com/bcpitutor/ostiki/appconfig"
	"github.com/bcpitutor/ostiki/db"
	"github.com/bcpitutor/ostiki/imo"
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

	//imoSender := imo.NewIMOSender()
	// imoListener := imo.NewIMOListener(8671, imoSender)

	container := dig.New()

	// Handle other db types, if implemented.
	if strings.ToLower(config.TikiDBConfig.DbType) == "dynamodb" {
		container.Provide(db.NewDynamoDBDriver)
	}

	container.Provide(appconfig.GetAppConfig)
	container.Provide(imo.NewIMOSender)
	container.Provide(imo.NewIMOListener)

	// container.Provide(imoListener)
	container.Provide(logger.GetTikiLogger)
	container.Provide(services.GetAWS)

	container.Provide(repositories.ProvideSessionRepository)
	container.Provide(repositories.ProvideBanRepository)
	container.Provide(repositories.ProvideDomainRepository)
	container.Provide(repositories.ProvideGroupRepository)
	container.Provide(repositories.ProvideTicketRepository)
	container.Provide(repositories.ProvidePermissionRepository)
	container.Provide(repositories.ProvideIMORepository)
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
