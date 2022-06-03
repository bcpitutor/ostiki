package repositories

import (
	"github.com/tiki-systems/tikiserver/models"
	"go.uber.org/dig"
)

type PermissionRepository struct {
	DBLayer          models.DBLayer
	GroupRepository  GroupRepository
	TicketRepository TicketRepository
}

type PermissionRepositoryResult struct {
	dig.Out
	PermissionRepository *PermissionRepository
}

func ProvidePermissionRepository(db models.DBLayer, gr *GroupRepository, tr *TicketRepository) PermissionRepositoryResult {
	return PermissionRepositoryResult{
		PermissionRepository: &PermissionRepository{
			DBLayer:          db,
			GroupRepository:  *gr,
			TicketRepository: *tr,
		},
	}
}

func (pr *PermissionRepository) IsUserInTikiadmins(userEmail string) bool {
	return pr.GroupRepository.IsUserInTikiadmins(userEmail)
}

func (pr *PermissionRepository) CanUserPerformTicketOperation(userEmail string, operation string) bool {
	return pr.TicketRepository.DBLayer.CanUserPerformTicketOperation(userEmail, operation)
}

func (pr *PermissionRepository) CanUserPerformGroupOperation(userEmail string, operation string) bool {
	return pr.GroupRepository.DBLayer.CanUserPerformGroupOperation(userEmail, operation)
}

func (pr *PermissionRepository) CanUserPerformDomainOperation(userEmail string, operation string) bool {
	return pr.GroupRepository.DBLayer.CanUserPerformDomainOperation(userEmail, operation)
}

func (pr *PermissionRepository) CanUserAccessToTicket(userEmail string, ticketPath string) bool {
	return pr.TicketRepository.DBLayer.CanUserAccessToTicket(userEmail, ticketPath)
}

func (pr *PermissionRepository) IsUserAllowedByDomainScope(userEmail string, ticketOrDomainPath string, domainScopeOperation string) bool {
	return pr.TicketRepository.DBLayer.IsUserAllowedByDomainScope(userEmail, ticketOrDomainPath, domainScopeOperation)
}
