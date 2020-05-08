package flow

import "time"

type Id uint

type Module struct {
}

type Country struct {
	Id          Id     `json:"id"`
	Name        string `json:"name"`
	IsoAlpha2   string `json:"iso_alpha_2"`
	IsoAlpha3   string `json:"iso_alpha_3"`
	CallingCode string `json:"calling_code"`
}

type Organization struct {
	Id                    Id     `json:"id"`
	Name                  string `json:"name"`
	Address               string `json:"address"`
	Zip                   string `json:"zip"`
	City                  string `json:"city"`
	PhoneNumber           string `json:"phone_number"`
	InvoiceDeploymentFees bool   `json:"invoice_deployment_fees"`

	Status struct {
		Id            Id         `json:"id"`
		Name          string     `json:"name"`
		RetentionTime *time.Time `json:"retention_time"`
	} `json:"status"`

	RegisteredModules []Module `json:"registered_modules"`

	Contacts struct {
		Primary   *User  `json:"primary"`
		Billing   *User  `json:"billing"`
		Technical []User `json:"technical"`
	} `json:"contacts"`
}

type User struct {
	Id                    uint           `json:"id"`
	Username              string         `json:"username"`
	FirstName             string         `json:"firstname"`
	LastName              string         `json:"lastname"`
	PhoneNumber           string         `json:"phone_number"`
	AssignedOrganizations []Organization `json:"assigned_organizations"`
	DefaultOrganization   Organization   `json:"default_organization"`
}
