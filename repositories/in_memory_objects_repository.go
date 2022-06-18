package repositories

import (
	"github.com/bcpitutor/ostiki/models"
	"go.uber.org/dig"
)

type IMORepository struct {
	DBLayer           models.DBLayer
	GroupRepository   GroupRepository
	TicketRepository  TicketRepository
	DomainRepository  DomainRepository
	SessionRepository SessionRepository
}

type IMORepositoryResult struct {
	dig.Out
	IMORepository *IMORepository
}

func ProvideIMORepository(db models.DBLayer, gr *GroupRepository, tr *TicketRepository, dr *DomainRepository, sr *SessionRepository) IMORepositoryResult {
	return IMORepositoryResult{
		IMORepository: &IMORepository{
			DBLayer:           db,
			GroupRepository:   *gr,
			TicketRepository:  *tr,
			DomainRepository:  *dr,
			SessionRepository: *sr,
		},
	}
}
