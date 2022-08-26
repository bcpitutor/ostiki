package repositories

import (
	"context"
	"fmt"

	"github.com/bcpitutor/ostiki/logger"
	"github.com/bcpitutor/ostiki/models"
	"github.com/hazelcast/hazelcast-go-client"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

type SessionRepository struct {
	DBLayer       models.DBLayer
	imoRepository IMORepository
	isCacheReady  *bool
	sugar         *zap.SugaredLogger
}

type SessionRepositoryResult struct {
	dig.Out
	SessionRepository *SessionRepository
}

func ProvideSessionRepository(isCacheReady *bool, logger *logger.TikiLogger, db models.DBLayer, imoRepository *IMORepository) SessionRepositoryResult {
	// logger := params.logger
	// logger.Logger.Sugar().Infof("ProvideSessionRepository with isCacheReady: %v", *params.isCacheReady)

	return SessionRepositoryResult{
		SessionRepository: &SessionRepository{
			DBLayer:       db,
			isCacheReady:  isCacheReady,
			sugar:         logger.Logger.Sugar(),
			imoRepository: *imoRepository,
		},
	}
}

func (sr *SessionRepository) GetSessions(scanType string) ([]models.Session, error) {
	var sessions []models.Session

	var err error
	if *sr.isCacheReady {
		//sr.sugar.Infof("requesting %s type of sessions from cache", scanType)
		sessions, err = sr.imoRepository.GetSessionsFromCache(scanType)
		if err != nil {
			sr.sugar.Errorf("error getting %s type of sessions from cache: %v", scanType, err)
			return nil, err // return from db
		}
		if len(sessions) == 0 {
			sr.sugar.Infof("GetSessions from cache, no sessions found, so reading from DB")
			sessions, err = sr.DBLayer.GetSessions(scanType)
			sr.imoRepository.FillSessionsIntoCache(scanType, sessions)
		}
	} else {
		sr.sugar.Infof("GetSessions from DB")
		sessions, err = sr.DBLayer.GetSessions(scanType)
		if err != nil {
			return nil, err
		}
		sr.imoRepository.FillSessionsIntoCache(scanType, sessions)
	}

	return sessions, nil
}

func (sr *SessionRepository) GetSessionByToken(token string) (models.Session, error) {
	var session models.Session

	//sr.sugar.Infof("GetSessionByToken: %s", token)
	if sr.isCacheReady != nil && *sr.isCacheReady {
		resultFromCache, err := getSessionFromCacheByToken(sr.sugar, sr.imoRepository.activeSessions, token)
		if err != nil {
			sr.sugar.Errorf("error getting session from cache: %v", err)
			resultFromDB, err := sr.DBLayer.GetSessionByToken(token)
			if err != nil {
				sr.sugar.Errorf("error getting session from DB: %v", err)
				return session, err
			}
			sr.sugar.Infof("Got session from DB")
			return resultFromDB, nil
		}
		session = resultFromCache
	} else {
		resultFromDB, err := sr.DBLayer.GetSessionByToken(token)
		if err != nil {
			return session, err
		}
		session = resultFromDB
	}

	return session, nil
}

func (sr *SessionRepository) UpdateSession(prevToken string, currentToken string, currentTokenExpires int64, refreshToken string) bool {

	// if sr.isCacheReady != nil && *sr.isCacheReady {
	// 	removeSessionFromCache(sr.sugar, sr.imoRepository.activeSessions, session.SessID)
	// 	removeSessionFromCache(sr.sugar, sr.imoRepository.allSessions, session.SessID)
	// 	sr.DBLayer.UpdateSession(prevToken, currentToken, currentTokenExpires, refreshToken)
	// 	addSessionToCache(sr.sugar, sr.imoRepository.activeSessions, session)
	// 	addSessionToCache(sr.sugar, sr.imoRepository.allSessions, session)
	// } else {
	// 	sr.DBLayer.UpdateSession(prevToken, currentToken, currentTokenExpires, refreshToken)
	// 	addSessionToCache(sr.sugar, sr.imoRepository.activeSessions, session)
	// 	addSessionToCache(sr.sugar, sr.imoRepository.allSessions, session)
	// }

	done := sr.DBLayer.UpdateSession(prevToken, currentToken, currentTokenExpires, refreshToken)
	if done {
		sr.sugar.Infof("UpdateSession: updated session in DB")
	} else {
		sr.sugar.Errorf("UpdateSession: error updating session in DB")
		return false
	}

	// TODO: implement a better solution for this
	// We are currently removing the session from the cache when we update it
	// and tell imo to retrieve it from the DB again
	if *sr.isCacheReady {
		sr.sugar.Infof("UpdateSession: removing all sessions from cache")
		sr.imoRepository.allSessions.Clear(context.TODO())
		sr.imoRepository.activeSessions.Clear(context.TODO())
		sr.imoRepository.expiredSessions.Clear(context.TODO())
		sr.imoRepository.revokedSessions.Clear(context.TODO())
		sr.sugar.Infof("UpdateSession: all sessions removed, now reading from DB again to fully update cache")
		sr.GetSessions("all")
		sr.GetSessions("active")
		sr.GetSessions("expired")
		sr.GetSessions("revoked")
	}

	return true
}

// func addSessionToCache(sugar *zap.SugaredLogger, sessionMap *hazelcast.Map, session models.Session) {
// 	sugar.Infof("addSessionToCache: %s", session.SessID)
// 	if yes, err := sessionMap.ContainsKey(context.TODO(), session.SessID); err == nil && !yes {
// 		sessionMap.Set(context.TODO(), session.SessID, session)
// 	}
// }

// func removeSessionFromCache(sugar *zap.SugaredLogger, sessionMap *hazelcast.Map, sessionID string) {
// 	sugar.Infof("removing session from cache: %s", sessionID)
// 	_, err := sessionMap.TryRemoveWithTimeout(context.TODO(), sessionID, (15 * time.Second))
// 	if err != nil {
// 		sugar.Errorf("error removing session from cache: %v", err)
// 		sugar.Infof("removing all session from cache due to previous error")
// 		sessionMap.Clear(context.TODO())
// 	}

// 	sugar.Infof("removed session from cache: %s", sessionID)
// }

func (sr *SessionRepository) CreateSession(session *models.Session) error {
	return sr.DBLayer.CreateSession(session)
}

func (sr *SessionRepository) DeleteSession(sessionID string, epoch int64) error {
	return sr.DBLayer.DeleteSession(sessionID, epoch)
}

func (sr *SessionRepository) GetSessionByRefreshToken(rtoken string) (models.Session, error) {
	return sr.DBLayer.GetSessionByRefreshToken(rtoken)
}

func (sr *SessionRepository) GetSessionsByEmail(email string) ([]models.Session, error) {
	return sr.DBLayer.GetSessionsByEmail(email)
}

func getSessionFromCacheByToken(sugar *zap.SugaredLogger, hzMap *hazelcast.Map, token string) (models.Session, error) {
	var session models.Session

	entries, err := hzMap.GetEntrySet(context.TODO())
	if err != nil {
		return session, err
	}
	sugar.Infof("Got %d entries from cache", len(entries))

	for _, entry := range entries {
		s := (entry.Value).(models.Session)
		if s.IdToken == token {
			sugar.Infof("Found session in cache")
			return s, nil
		}
	}

	return session, fmt.Errorf("session not found")
}
