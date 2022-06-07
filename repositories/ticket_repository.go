package repositories

import (
	"github.com/bcpitutor/ostiki/models"
	"go.uber.org/dig"
)

type TicketRepository struct {
	DBLayer    models.DBLayer
	CacheLayer models.CacheLayer
}

type TicketRepositoryResult struct {
	dig.Out
	TicketRepository *TicketRepository
}

func ProvideTicketRepository(db models.DBLayer, cl models.CacheLayer) TicketRepositoryResult {
	return TicketRepositoryResult{
		TicketRepository: &TicketRepository{
			DBLayer:    db,
			CacheLayer: cl,
		},
	}
}

func (t *TicketRepository) GetAllTickets() ([]models.Ticket, error) {
	return t.DBLayer.GetAllTickets()
}

func (t *TicketRepository) DoesTicketExist(ticketPath string) bool {
	return t.DBLayer.DoesTicketExist(ticketPath)
}

func (t *TicketRepository) QueryTicketByPath(ticketPath string) (models.Ticket, error) {
	return t.DBLayer.QueryTicketByPath(ticketPath)
}

func (t *TicketRepository) CreateTicket(ticket models.Ticket) error {
	return t.DBLayer.CreateTicket(ticket)
}

func (t *TicketRepository) DeleteTicket(ticketPath string, ticketType string) error {
	return t.DBLayer.DeleteTicket(ticketPath, ticketType)
}

func (t *TicketRepository) SetTicketSecret(ticketPath string, secretData string) error {
	return t.DBLayer.SetTicketSecret(ticketPath, secretData)
}

func (t *TicketRepository) GetTicketSecret(ticketPath string) (string, error) {
	return t.DBLayer.GetTicketSecret(ticketPath)
}

// func (t *TicketRepository) ObtainTickets(ticketPaths []string) ([]models.Ticket, error) {
// 	return t.DBLayer.ObtainTickets(ticketPaths)
// }

// func (t *TicketRepository) CanUserAccessToTicket(userEmail string, ticketPath string) bool {
// 	return t.DBLayer.CanUserAccessToTicket(userEmail, ticketPath)
// }

// func (t *TicketRepository) CanUserPerformTicketOperation(userEmail string, operation string) {
// 	t.DBLayer.CanUserPerformTicketOperation(userEmail, operation)
// }
