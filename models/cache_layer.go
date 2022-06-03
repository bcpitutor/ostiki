package models

type CacheLayer interface {
	GetTickets() ([]Ticket, error)
	GetGroups() ([]TicketGroup, error)
	Usable() bool
	CacheType() string
}
