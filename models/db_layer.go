package models

type DBLayer interface {
	// GetTickets() ([]Ticket, error)
	// GetGroups() ([]TicketGroup, error)
	// Usable() bool
	DBType() string

	// Ban Table
	AddBannedUser(bannedUser BannedUser) error
	GetBannedUsers() ([]BannedUser, error)
	GetBannedUserByEmail(userEmail string) (BannedUser, error)
	UnbanUser(userEmail string) error

	// Session Table
	CreateSession(session *Session) error
	UpdateSession(prevToken string,
		currentToken string,
		currentTokenExpires int64,
		refreshToken string,
	) bool
	DeleteSession(sessionID string, epoch int64) error
	GetSessionByRefreshToken(refreshToken string) (Session, error)
	GetSessions(scanType string) ([]Session, error)

	// Group Table
	GetAllGroups() ([]TicketGroup, error)
	GetGroup(groupName string) (TicketGroup, error)
	CreateGroup(group TicketGroup) error
	DeleteGroup(groupName string) error
	DoesGroupExist(groupName string) bool
	IsUserInTikiadmins(userEmail string) bool
	GetGroupMembers(groupName string) ([]string, error)
	IsUserMemberOfGroup(member string, groupName string) bool
	CanUserPerformGroupOperation(userEmail string, operation string) bool
	GetGroupNamesOfUser(userEmail string) ([]string, error)
	GetGroupsOfUser(userEmail string) ([]TicketGroup, error)
	CanUserAccessToTicket(userEmail string, ticketPath string) bool
	AddMemberToGroup(newMember string, groupName string, changedBy string) error
	DelMemberFromGroup(memberToDelete string, groupName string, changedBy string) error

	// Ticket Table
	GetAllTickets() ([]Ticket, error)
	QueryTicketByPath(ticketPath string) (Ticket, error)
	DoesTicketExist(ticketPath string) bool
	CreateTicket(ticket Ticket) error
	DeleteTicket(ticketPath string, ticketType string) error
	IsUserAllowedByDomainScope(userEmail string, ticketOrDomainPath string, domainScopeOperation string) bool
	SetTicketSecret(ticketPath string, secretData string) error
	GetTicketSecret(ticketPath string) (string, error)
	CanUserPerformTicketOperation(userEmail string, operationType string) bool

	// Domain Table
	CanUserPerformDomainOperation(userEmail string, operation string) bool
	DoesTicketDomainExist(ticketDomainPath string) bool
	GetAllDomains() ([]TicketDomain, error)
	GetDomain(domainPath string) (TicketDomain, error)
	CreateDomain(domain TicketDomain) error
	DeleteDomain(domainPath string) error
	IsUserBanned(userEmail string) bool
}
