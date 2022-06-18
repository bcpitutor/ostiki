package models

type OPERATION struct {
	Create    string
	Delete    string
	AddMember string
	DelMember string
	Show      string
	Info      string
	List      string
	SetSecret string
	GetSecret string
}

var Operation = OPERATION{
	Create:    "create",
	Delete:    "delete",
	AddMember: "addMember",
	DelMember: "delMember",
	List:      "list",
	Show:      "show",
	Info:      "info",
	SetSecret: "setSecret",
	GetSecret: "getSecret",
}

type DOMAINSCOPEOPERATION struct {
	CreateTicket string
	DeleteTicket string
	CreateDomain string
	DeleteDomain string
}

var DomainScopeOperation = DOMAINSCOPEOPERATION{
	CreateTicket: "createTicket",
	DeleteTicket: "deleteTicket",
	CreateDomain: "createDomain",
	DeleteDomain: "deleteDomain",
}

type KMSOPERATION struct {
	Encrypt string
	Decrypt string
}

var KmsOperation = KMSOPERATION{
	Encrypt: "encrypt",
	Decrypt: "decrypt",
}
