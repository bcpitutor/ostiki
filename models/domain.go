package models

type TicketDomain struct {
	DomainPath    string `json:"domainPath"`
	OwnerGroup    string `json:"ownerGroup"`
	Parent        string `json:"parent"`
	DomainComment string `json:"domainComment"`
	CreatedAt     int64  `json:"createdAt"`
	CreatedBy     string `json:"createdBy"`
	UpdatedAt     int64  `json:"updatedAt"`
	UpdatedBy     string `json:"updatedBy"`
}
