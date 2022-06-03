package models

type TicketGroup struct {
	GroupName    string   `json:"groupName"`
	GroupMembers []string `json:"groupMembers"`
	InUseBy      []string `json:"inUseBy,omitempty"`
	CreatedAt    int64    `json:"createdAt,omitempty"`
	CreatedBy    string   `json:"createdBy,omitempty"`
	UpdatedAt    int64    `json:"updatedAt,omitempty"`
	UpdatedBy    string   `json:"updatedBy,omitempty"`
	AccessPerms  Aperms   `json:"accessPerms,omitempty"`
	DomainScope  DScope   `json:"domainScope"`
}

type Aperms struct {
	Group        map[string]bool
	Domain       map[string]bool
	Ticket       map[string]bool
	SecretTicket map[string]bool
}

type DScope struct {
	Root string `json:"root"`
	Info string `json:"info"`
}
