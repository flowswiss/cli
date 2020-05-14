package flow

type User struct {
	Id                    uint            `json:"id"`
	Username              string          `json:"username"`
	FirstName             string          `json:"firstname"`
	LastName              string          `json:"lastname"`
	PhoneNumber           string          `json:"phone_number"`
	AssignedOrganizations []*Organization `json:"assigned_organizations"`
	DefaultOrganization   *Organization   `json:"default_organization"`
}
