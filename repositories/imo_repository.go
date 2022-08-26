package repositories

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/bcpitutor/ostiki/appconfig"
	"github.com/bcpitutor/ostiki/logger"
	"github.com/bcpitutor/ostiki/models"
	"github.com/hazelcast/hazelcast-go-client"
	"github.com/hazelcast/hazelcast-go-client/cluster"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

type IMORepository struct {
	sugar    *zap.SugaredLogger
	config   *appconfig.AppConfig
	hzClient *hazelcast.Client

	activeSessions  *hazelcast.Map
	expiredSessions *hazelcast.Map
	revokedSessions *hazelcast.Map
	allSessions     *hazelcast.Map
}

type IMORepositoryResult struct {
	dig.Out
	IMORepository *IMORepository
}

// type IMODiscoveryResult struct {
// 	done bool
// }

func ProvideIMORepository(isCacheReady *bool, appconfig *appconfig.AppConfig, logger *logger.TikiLogger) IMORepositoryResult {
	imo := IMORepository{}
	imo.sugar = logger.Logger.Sugar()
	imo.config = appconfig

	// imo.activeSessions = &hazelcast.Map{}
	// imo.allSessions = &hazelcast.Map{}
	// imo.expiredSessions = &hazelcast.Map{}
	// imo.revokedSessions = &hazelcast.Map{}

	//imo.sessionRepository = sr

	// init the hazelcast connection
	hz, err := imo.InitHazelcast()
	if err != nil {
		imo.sugar.Errorf("failed to init hazelcast connection: %v", err)
	} else {
		imo.sugar.Infof("hazelcast connection initialized")
		(*isCacheReady) = true
	}
	imo.hzClient = hz

	return IMORepositoryResult{
		IMORepository: &imo,
	}
}

func (imo *IMORepository) InitHazelcast() (*hazelcast.Client, error) {
	// init a hazelcast client

	// clusterConfig := cluster.Config{
	// 	Name: "ostiki-hz-cluster",
	// }

	clusterConfig := cluster.Config{}
	hzConfig := hazelcast.Config{
		ClientName: "localhost",
		Cluster:    clusterConfig,
	}

	client, err := hazelcast.StartNewClientWithConfig(context.TODO(), hzConfig)
	if err != nil {
		imo.sugar.Errorf("failed to start hazelcast client: %v", err)
		return nil, err
	}
	imo.hzClient = client

	client.AddLifecycleListener(func(event hazelcast.LifecycleStateChanged) {
		imo.sugar.Infof("hazelcast lifecycle event: %+v", event)
		imo.sugar.Infof("State: %v", event.State.String())
		imo.sugar.Infof("State: %v", event.EventName())
	})

	// client.AddDistributedObjectListener(
	// 	context.TODO(),
	// 	func(e hazelcast.DistributedObjectNotified) {
	// 		imo.sugar.Infof("EventType: %s", e.EventType)
	// 		imo.sugar.Infof("ObjectName: %s", e.ObjectName)
	// 	},
	// )

	imo.sugar.Infof("hazelcast client started")
	if imo.hzClient == nil {
		imo.sugar.Infof("hazelcast client is nil, what can I do sometimes?")
	}

	if imo.hzClient != nil {
		allObj, err := imo.hzClient.GetMap(context.TODO(), "sessions-all")
		if err != nil {
			imo.sugar.Errorf("failed to get map sessions-all: %v", err)
			//return nil, err
		}

		expiredObj, err := imo.hzClient.GetMap(context.TODO(), "sessions-expired")
		if err != nil {
			imo.sugar.Errorf("failed to get map sessions-expired: %v", err)
		}

		activeObj, err := imo.hzClient.GetMap(context.TODO(), "sessions-active")
		if err != nil {
			imo.sugar.Errorf("failed to get map sessions-active: %v", err)
		}

		revokedObj, err := imo.hzClient.GetMap(context.TODO(), "sessions-revoked")
		if err != nil {
			imo.sugar.Errorf("failed to get map sessions-revoked: %v", err)
		}

		imo.activeSessions = activeObj
		imo.expiredSessions = expiredObj
		imo.revokedSessions = revokedObj
		imo.allSessions = allObj
	}

	return client, nil
}

func (imo *IMORepository) GetHZClient() *hazelcast.Client {
	return imo.hzClient
}

func (imo *IMORepository) FillSessionsIntoCache(scanType string, sessions []models.Session) error {
	imo.sugar.Infof("Starting FillSessionsIntoCache with scanType %s, Got %d sessions", scanType, len(sessions))

	if imo.allSessions == nil {
		return errors.New("imo.allSessions is nil")
	}

	timeNow := strconv.FormatInt(time.Now().Unix(), 10)
	for _, session := range sessions {
		yes, _ := imo.allSessions.ContainsKey(context.TODO(), session.SessID)
		if !yes {
			imo.allSessions.Put(context.TODO(), session.SessID, session)
		}

		if session.Expire > timeNow && session.IsRevoked == false {
			imo.activeSessions.Put(context.TODO(), session.SessID, session)
			continue
		}
		if session.Expire <= timeNow && session.IsRevoked == false {
			imo.expiredSessions.Put(context.TODO(), session.SessID, session)
			continue
		}
		if session.IsRevoked == true {
			imo.revokedSessions.Put(context.TODO(), session.SessID, session)
		}
	}

	imo.sugar.Infof("finished adding sessions to cache")

	return nil
}

func (imo *IMORepository) GetSessionsFromCache(scanType string) ([]models.Session, error) {
	result := []models.Session{}

	switch scanType {
	case "active":
		activeSessions, _ := imo.activeSessions.GetEntrySet(context.TODO())
		for _, session := range activeSessions {
			result = append(result, session.Value.(models.Session))
		}
	case "expired":
		expiredSessions, _ := imo.expiredSessions.GetEntrySet(context.TODO())
		for _, session := range expiredSessions {
			result = append(result, session.Value.(models.Session))
		}
	case "revoked":
		revokedSessions, _ := imo.revokedSessions.GetEntrySet(context.TODO())
		for _, session := range revokedSessions {
			result = append(result, session.Value.(models.Session))
		}
	case "all":
		if imo.allSessions == nil {
			return nil, errors.New("all sessions cache not initialized")
		}

		allSessions, _ := imo.allSessions.GetEntrySet(context.TODO())
		for _, session := range allSessions {
			result = append(result, session.Value.(models.Session))
		}
	default:
		return nil, fmt.Errorf("invalid scan type %s", scanType)
	}

	return result, nil
}

func (imo *IMORepository) GetCacheObject(objname string) (*hazelcast.Map, error) {
	obj, err := imo.hzClient.GetMap(context.TODO(), objname)
	if err != nil {
		imo.sugar.Errorf("failed to get map %s: %v", objname, err)
		return nil, err
	}
	return obj, nil
}
