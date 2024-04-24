package client

type JWT struct {
	AccessToken      string `json:"access_token,omitempty"`
	IDToken          string `json:"id_token,omitempty"`
	ExpiresIn        int    `json:"expires_in,omitempty"`
	RefreshExpiresIn int    `json:"refresh_expires_in,omitempty"`
	RefreshToken     string `json:"refresh_token,omitempty"`
	TokenType        string `json:"token_type,omitempty"`
	NotBeforePolicy  int    `json:"not-before-policy,omitempty"`
	SessionState     string `json:"session_state,omitempty"`
	Scope            string `json:"scope,omitempty"`
}

type Organization struct {
	Id           string        `json:"id"`
	Name         string        `json:"name"`
	State        string        `json:"state"`
	CreatedAt    string        `json:"created_at"`
	UpdatedAt    string        `json:"updated_at"`
	Users        []User        `json:"users"`
	Environments []Environment `json:"environments"`
}

type User struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type Environment struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Identifier string `json:"identifier"`
	Region     string `json:"region"`
}

type OrganizationUser struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type Vault struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Environment string `json:"env_identifier"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type organizationsAPIData struct {
	Data []organizationAPI `json:"data,omitempty"`
}

type organizationUsersAPIData struct {
	Data []organizationUserAPI `json:"data,omitempty"`
}

type organizationInvitesAPIData struct {
	Data []organizationInviteAPI `json:"data,omitempty"`
}

type organizationVaultsAPIData struct {
	Data []organizationVaultAPI `json:"data,omitempty"`
}

type vaultUsersAPIData struct {
	Data []vaultUserAPI `json:"data,omitempty"`
}

type organizationVaultAPI struct {
	Id            string                         `json:"id,omitempty"`
	Type          string                         `json:"type,omitempty"`
	Attributes    organizationVaultAPIAttributes `json:"attributes,omitempty"`
	Relationships vaultAPIRelationships          `json:"relationships,omitempty"`
	Links         vaultAPILinks                  `json:"links,omitempty"`
}

type vaultAPILinks struct {
	Self               string `json:"self,omitempty"`
	ReverseProxy       string `json:"reverse_proxy,omitempty"`
	ForwardProxy       string `json:"forward_proxy,omitempty"`
	VaultManagementApi string `json:"vault_management_api,omitempty"`
	VaultApi           string `json:"vault_api,omitempty"`
}

type vaultAPIRelationships struct {
	Organization organizationRelationshipsAPIData `json:"organization,omitempty"`
}

type organizationRelationshipsAPIData struct {
	Links Links `json:"links,omitempty"`
	Data  Data  `json:"data,omitempty"`
}

type Data struct {
	Id   string `json:"id,omitempty"`
	Type string `json:"type,omitempty"`
}

type Links struct {
	Self    string `json:"self,omitempty"`
	Related string `json:"related,omitempty"`
}

type organizationInviteAPI struct {
	Id         string                          `json:"id,omitempty"`
	Type       string                          `json:"type,omitempty"`
	Attributes organizationInviteAPIAttributes `json:"attributes,omitempty"`
}

type organizationAPI struct {
	Id         string                    `json:"id,omitempty"`
	Type       string                    `json:"type,omitempty"`
	Attributes organizationAPIAttributes `json:"attributes,omitempty"`
}

type organizationUserAPI struct {
	Id         string                        `json:"id,omitempty"`
	Type       string                        `json:"type,omitempty"`
	Attributes organizationUserAPIAttributes `json:"attributes,omitempty"`
}

type vaultUserAPI struct {
	Id         string                 `json:"id,omitempty"`
	Type       string                 `json:"type,omitempty"`
	Attributes vaultUserAPIAttributes `json:"attributes,omitempty"`
}

type organizationAPIAttributes struct {
	InternalId  string   `json:"internal_id,omitempty"`
	Identifier  string   `json:"identifier,omitempty"`
	Name        string   `json:"name,omitempty"`
	Active      bool     `json:"active,omitempty"`
	State       string   `json:"state,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
	CreatedAt   string   `json:"created_at,omitempty"`
	UpdatedAt   string   `json:"updated_at,omitempty"`
}

type vaultUserAPIAttributes struct {
	Id    string `json:"id,omitempty"`
	Email string `json:"email,omitempty"`
	Role  string `json:"role"`
}

type organizationUserAPIAttributes struct {
	CreatedAt    string               `json:"created_at,omitempty"`
	UpdatedAt    string               `json:"updated_at,omitempty"`
	Id           string               `json:"id,omitempty"`
	Name         string               `json:"name,omitempty"`
	EmailAddress string               `json:"email_address,omitempty"`
	Permissions  []string             `json:"permissions"`
	Vaults       []vaultAPIAttributes `json:"vaults"`
	Role         string               `json:"role"`
	LastLogin    any                  `json:"last_login"`
	LastIP       any                  `json:"last_ip"`
}

type organizationVaultAPIAttributes struct {
	Identifier  string `json:"identifier,omitempty"`
	Environment string `json:"environment,omitempty"`
	Role        string `json:"role"`
	UpdatedAt   string `json:"updated_at,omitempty"`
	Name        string `json:"name,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
}

type organizationInviteAPIAttributes struct {
	InviteId     string               `json:"invite_id,omitempty"`
	InviteStatus string               `json:"invite_status,omitempty"`
	UserEmail    string               `json:"user_email,omitempty"`
	InvitedBy    string               `json:"invited_by"`
	CreatedAt    string               `json:"created_at,omitempty"`
	Role         string               `json:"role"`
	Vaults       []vaultAPIAttributes `json:"vaults"`
}

type vaultAPIAttributes struct {
	Id          string   `json:"id,omitempty"`
	Name        string   `json:"name,omitempty"`
	Role        string   `json:"role,omitempty"`
	Environment string   `json:"env_identifier,omitempty"`
	Identifier  string   `json:"identifier,omitempty"`
	State       string   `json:"state,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
	CreatedAt   string   `json:"created_at,omitempty"`
	UpdatedAt   string   `json:"updated_at,omitempty"`
}
