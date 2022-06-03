package cache

import "github.com/tiki-systems/tikiserver/models"

type NoCacheDriver struct {
	Type string
}

func NewNoCacheDriver() (models.CacheLayer, error) {
	return &HazelcastDriver{
		Type: "no-cache",
	}, nil
}

func (h *NoCacheDriver) GetTickets() ([]models.Ticket, error) {
	return nil, nil
}

func (h *NoCacheDriver) GetGroups() ([]models.TicketGroup, error) {
	return nil, nil
}

func (h *NoCacheDriver) Usable() bool {
	return true
}

func (h *NoCacheDriver) CacheType() string {
	return "no-cache"
}
