package repositories

import (
	"github.com/bcpitutor/ostiki/models"
	"go.uber.org/dig"
)

type SessionRepository struct {
	DBLayer models.DBLayer
}

type SessionRepositoryResult struct {
	dig.Out
	SessionRepository *SessionRepository
}

func ProvideSessionRepository(db models.DBLayer) SessionRepositoryResult {
	return SessionRepositoryResult{
		SessionRepository: &SessionRepository{
			DBLayer: db,
		},
	}
}

func (sr *SessionRepository) CreateSession(session *models.Session) error {
	return sr.DBLayer.CreateSession(session)
}

func (sr *SessionRepository) UpdateSession(prevToken string, currentToken string, currentTokenExpires int64, refreshToken string) bool {
	return sr.DBLayer.UpdateSession(prevToken, currentToken, currentTokenExpires, refreshToken)
}

func (sr *SessionRepository) GetSessionByRefreshToken(rtoken string) (models.Session, error) {
	return sr.DBLayer.GetSessionByRefreshToken(rtoken)
}

func (sr *SessionRepository) GetSessions(scanType string) ([]models.Session, error) {
	return sr.DBLayer.GetSessions(scanType)
}
