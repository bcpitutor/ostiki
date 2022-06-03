package cache

import (
	"context"
	"fmt"
	"os"

	"github.com/hazelcast/hazelcast-go-client"
	"github.com/hazelcast/hazelcast-go-client/cluster"
	"github.com/tiki-systems/tikiserver/appconfig"
	"github.com/tiki-systems/tikiserver/logger"
	"github.com/tiki-systems/tikiserver/models"
)

type HazelcastDriver struct {
	Type   string
	Client any
}

func NewHazelcastDriver(appconfig *appconfig.AppConfig, gLogger *logger.TikiLogger) (models.CacheLayer, error) {
	fmt.Printf("In NewHazelcastDriver constructor\n")
	logger := gLogger.Logger
	if appconfig.Deployment == "" {
		logger.Error("Can't detect deployment to init Hazelcast")
		return nil, fmt.Errorf("Can't detect deployment to init Hazelcast")
	}

	hostname, err := os.Hostname()
	if err != nil {
		logger.Error("Can't detect hostname to init Hazelcast")
		return nil, fmt.Errorf("Can't detect hostname to init Hazelcast")
	}

	clientname := fmt.Sprintf("tiki-pod-%s", hostname)
	clusterConfig := cluster.Config{}

	if appconfig.TikiInMemoryStoreConfig.HazelcastAddr != "" {
		clusterConfig.Network.Addresses = []string{appconfig.TikiInMemoryStoreConfig.HazelcastAddr}
	} else {
		clusterConfig.Network.Addresses = []string{
			fmt.Sprintf("hz-%s.tiki.svc.cluster.local", appconfig.Deployment),
		}
	}

	hzConfig := hazelcast.Config{
		ClientName: clientname,
		Cluster:    clusterConfig,
	}

	client, err := hazelcast.StartNewClientWithConfig(context.Background(), hzConfig)
	if err != nil {
		logger.Sugar().Errorf("Can't start hazelcast client")
		return nil, fmt.Errorf("Can't start hazelcast client")
	}

	return &HazelcastDriver{
		Type:   "hazelcast",
		Client: client,
	}, nil
}

func (h *HazelcastDriver) GetTickets() ([]models.Ticket, error) {
	return nil, nil
}

func (h *HazelcastDriver) GetGroups() ([]models.TicketGroup, error) {
	return nil, nil
}

func (h *HazelcastDriver) Usable() bool {
	return true
}

func (h *HazelcastDriver) CacheType() string {
	return "Hazelcast"
}

// func InitHazelcast(appconfig *appconfig.AppConfig, gLogger *logger.TikiLogger) (*hazelcast.Client, error) {
// 	logger := gLogger.Logger
// 	if appconfig.Deployment == "" {
// 		logger.Error("Can't detect deployment to init Hazelcast")
// 		return nil, fmt.Errorf("Can't detect deployment to init Hazelcast")
// 	}

// 	hostname, err := os.Hostname()
// 	if err != nil {
// 		logger.Error("Can't detect hostname to init Hazelcast")
// 		return nil, fmt.Errorf("Can't detect hostname to init Hazelcast")
// 	}

// 	clientname := fmt.Sprintf("tiki-pod-%s", hostname)
// 	clusterConfig := cluster.Config{}

// 	if appconfig.TikiInMemoryStoreConfig.HazelcastAddr != "" {
// 		clusterConfig.Network.Addresses = []string{appconfig.TikiInMemoryStoreConfig.HazelcastAddr}
// 	} else {
// 		clusterConfig.Network.Addresses = []string{
// 			fmt.Sprintf("hz-%s.tiki.svc.cluster.local", appconfig.Deployment),
// 		}
// 	}

// 	hzConfig := hazelcast.Config{
// 		ClientName: clientname,
// 		Cluster:    clusterConfig,
// 	}

// 	client, err := hazelcast.StartNewClientWithConfig(context.Background(), hzConfig)
// 	if err != nil {
// 		logger.Sugar().Errorf("Can't start hazelcast client")
// 		return nil, fmt.Errorf("Can't start hazelcast client")
// 	}

// 	return client, nil
// }
