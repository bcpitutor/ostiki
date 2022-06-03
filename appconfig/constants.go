package appconfig

type TableNames struct {
	GroupTableName    string
	DomainTableName   string
	SessTableName     string
	TicketTableName   string
	RegisterTableName string
	BannedTableName   string
}

var TNames TableNames

func (tnames *TableNames) SetTableNames(names map[string]string) {
	tnames.GroupTableName = names["group"]
	tnames.DomainTableName = names["domain"]
	tnames.SessTableName = names["sess"]
	tnames.TicketTableName = names["ticket"]
	tnames.RegisterTableName = names["register"]
	tnames.BannedTableName = names["banned"]
}

func (tnames *TableNames) GetTableNames() *TableNames {
	return tnames
}
