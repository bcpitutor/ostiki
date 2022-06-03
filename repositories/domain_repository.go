package repositories

import (
	"github.com/tiki-systems/tikiserver/models"
	"go.uber.org/dig"
)

type DomainRepository struct {
	DBLayer models.DBLayer
}

type DomainRepositoryResult struct {
	dig.Out
	DomainRepository *DomainRepository
}

func ProvideDomainRepository(db models.DBLayer) DomainRepositoryResult {
	return DomainRepositoryResult{
		DomainRepository: &DomainRepository{
			DBLayer: db,
		},
	}
}

func (d *DomainRepository) DoesTicketDomainExist(domainPath string) bool {
	return d.DBLayer.DoesTicketDomainExist(domainPath)
}

func (d *DomainRepository) ListDomains() ([]models.TicketDomain, error) {
	return d.DBLayer.GetAllDomains()
}

func (d *DomainRepository) GetDomain(domainPath string) (models.TicketDomain, error) {
	return d.DBLayer.GetDomain(domainPath)
}

func (d *DomainRepository) CreateDomain(domain models.TicketDomain) error {
	return d.DBLayer.CreateDomain(domain)
}

func (d *DomainRepository) DeleteDomain(domainPath string) error {
	return d.DBLayer.DeleteDomain(domainPath)
}
