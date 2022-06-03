package models

type Ticket struct {
	TicketPath       string         `json:"ticketPath"`
	TicketType       string         `json:"ticketType"`
	TicketInfo       string         `json:"ticketInfo"`
	TicketRegion     string         `json:"ticketRegion"`
	AwsAssumeRole    AwsAssumeRole  `json:"assumeRoleDetails,omitempty"`
	AwsPermissions   AwsPermissions `json:"awsPermissions,omitempty"`
	OwnersGroup      []string       `json:"ownersGroup,omitempty"`
	SecretData       string         `json:"secretData,omitempty"`
	SourceAddresses  []string       `json:"sourceAddress,omitempty"`
	SAccountPassword string         `json:"sAccountPassword,omitempty"`
	K8sDetails       K8sDetails     `json:"k8sDetails"`
	CreatedAt        string         `json:"createdAt"`
	CreatedBy        string         `json:"createdBy"`
	UpdatedAt        string         `json:"updatedAt"`
	UpdatedBy        string         `json:"updatedBy"`
}

type AwsPermissions struct {
	Effect   string   `json:"effect"`
	Action   []string `json:"action"`
	Resource string   `json:"resource"`
}
type AwsAssumeRole struct {
	RoleArn string `json:"roleArn"`
	Ttl     int32  `json:"ttl"`
}

type K8sDetails struct {
	Server  string `json:"server"`
	Cluster string `json:"cluster"`
	User    string `json:"user"`
	Name    string `json:"name"`
}
