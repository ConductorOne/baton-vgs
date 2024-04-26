package client

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/conductorone/baton-sdk/pkg/uhttp"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
)

type VGSClient struct {
	httpClient      *uhttp.BaseHttpClient
	token           *JWT
	serviceEndpoint string
	organizationId  string
	vaultId         string
}

func WithBody(body string) uhttp.RequestOption {
	return func() (io.ReadWriter, map[string]string, error) {
		var buffer bytes.Buffer
		_, err := buffer.WriteString(body)
		if err != nil {
			return nil, nil, err
		}

		_, headers, err := WithContentTypeFormHeader()()
		if err != nil {
			return nil, nil, err
		}

		return &buffer, headers, nil
	}
}

func WithJSONBodyV2(body interface{}) uhttp.RequestOption {
	return func() (io.ReadWriter, map[string]string, error) {
		buffer := new(bytes.Buffer)
		err := json.NewEncoder(buffer).Encode(body)
		if err != nil {
			return nil, nil, err
		}

		_, headers, err := WithContentTypeVndHeader()()
		if err != nil {
			return nil, nil, err
		}

		return buffer, headers, nil
	}
}

func WithContentTypeFormHeader() uhttp.RequestOption {
	return func() (io.ReadWriter, map[string]string, error) {
		return nil, map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		}, nil
	}
}

func WithContentTypeVndHeader() uhttp.RequestOption {
	return func() (io.ReadWriter, map[string]string, error) {
		return nil, map[string]string{
			"Content-Type": "application/vnd.api+json",
		}, nil
	}
}

func WithAcceptVndJSONHeader() uhttp.RequestOption {
	return func() (io.ReadWriter, map[string]string, error) {
		return nil, map[string]string{
			"Accept": "application/vnd.api+json",
		}, nil
	}
}

func WithAuthorizationBearerHeader(token string) uhttp.RequestOption {
	return uhttp.WithHeader("Authorization", "Bearer "+token)
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func WithSetBasicAuthHeader(username, password string) uhttp.RequestOption {
	return uhttp.WithHeader("Authorization", "Basic "+basicAuth(username, password))
}

func New(ctx context.Context, clientId, clientSecret, orgId, vaultId string) (*VGSClient, error) {
	var jwt = &JWT{}
	uri, err := url.Parse("https://auth.verygoodsecurity.com/auth/realms/vgs/protocol/openid-connect/token")
	if err != nil {
		return nil, err
	}

	httpClient, err := uhttp.NewClient(ctx, uhttp.WithLogger(true, ctxzap.Extract(ctx)))
	if err != nil {
		return nil, err
	}

	cli := uhttp.NewBaseHttpClient(httpClient)
	req, err := cli.NewRequest(ctx,
		http.MethodPost,
		uri,
		uhttp.WithAcceptJSONHeader(),
		WithBody(`grant_type=client_credentials`),
		WithSetBasicAuthHeader(clientId, clientSecret),
	)
	if err != nil {
		return nil, err
	}

	resp, err := cli.Do(req, uhttp.WithJSONResponse(&jwt))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("token is not valid")
	}

	vc := VGSClient{
		httpClient: cli,
		token: &JWT{
			AccessToken:      jwt.AccessToken,
			ExpiresIn:        jwt.ExpiresIn,
			RefreshExpiresIn: jwt.RefreshExpiresIn,
			Scope:            jwt.Scope,
			TokenType:        jwt.TokenType,
			NotBeforePolicy:  jwt.NotBeforePolicy,
		},
		serviceEndpoint: "https://accounts.apps.verygoodsecurity.com",
		organizationId:  orgId,
		vaultId:         vaultId,
	}

	return &vc, nil
}

func (v *VGSClient) GetToken() string {
	return v.token.AccessToken
}

func (v *VGSClient) GetOrganizationId() string {
	return v.organizationId
}

func (v *VGSClient) GetVaultId() string {
	return v.vaultId
}

func (v *VGSClient) ListOrganizations(ctx context.Context) ([]Organization, error) {
	var (
		organizations        []Organization
		organizationsAPIData organizationsAPIData
	)
	strUrl, err := url.JoinPath(v.serviceEndpoint, "/organizations")
	if err != nil {
		return nil, err
	}

	uri, err := url.Parse(strUrl)
	if err != nil {
		return nil, err
	}

	req, err := v.httpClient.NewRequest(ctx,
		http.MethodGet,
		uri,
		WithAcceptVndJSONHeader(),
		WithAuthorizationBearerHeader(v.GetToken()),
	)
	if err != nil {
		return nil, err
	}

	resp, err := v.httpClient.Do(req, uhttp.WithJSONResponse(&organizationsAPIData))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	for _, org := range organizationsAPIData.Data {
		organizations = append(organizations, Organization{
			Id:        org.Id,
			Name:      org.Attributes.Name,
			State:     org.Attributes.State,
			CreatedAt: org.Attributes.CreatedAt,
			UpdatedAt: org.Attributes.UpdatedAt,
		})
	}

	return organizations, nil
}

func (v *VGSClient) ListUsers(ctx context.Context, orgId, vaultId string) ([]OrganizationUser, error) {
	var (
		users                    []OrganizationUser
		organizationUsersAPIData organizationUsersAPIData
	)
	if !strings.Contains(v.token.Scope, "organization-users:read") {
		return nil, fmt.Errorf("organization-users:read scope not found")
	}

	strUrl, err := url.JoinPath(v.serviceEndpoint, "/organizations/", orgId, "/members")
	if err != nil {
		return nil, err
	}

	uri, err := url.Parse(strUrl)
	if err != nil {
		return nil, err
	}

	req, err := v.httpClient.NewRequest(ctx,
		http.MethodGet,
		uri,
		WithAcceptVndJSONHeader(),
		WithAuthorizationBearerHeader(v.GetToken()),
	)
	if err != nil {
		return nil, err
	}

	resp, err := v.httpClient.Do(req, uhttp.WithJSONResponse(&organizationUsersAPIData))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	for _, userAPI := range organizationUsersAPIData.Data {
		users = append(users, OrganizationUser{
			Id:        userAPI.Id,
			Name:      userAPI.Attributes.Name,
			Email:     userAPI.Attributes.EmailAddress,
			CreatedAt: userAPI.Attributes.CreatedAt,
			UpdatedAt: userAPI.Attributes.UpdatedAt,
		})
	}

	return users, nil
}

func (v *VGSClient) ListUserInvites(ctx context.Context, orgId string) ([]OrganizationUser, error) {
	var (
		userInvites                []OrganizationUser
		organizationInvitesAPIData organizationInvitesAPIData
	)
	if !strings.Contains(v.token.Scope, "organization-users:read") {
		return nil, fmt.Errorf("organization-users:read scope not found")
	}

	strUrl, err := url.JoinPath(v.serviceEndpoint, "/organizations/", orgId, "/invites")
	if err != nil {
		return nil, err
	}

	uri, err := url.Parse(strUrl)
	if err != nil {
		return nil, err
	}

	req, err := v.httpClient.NewRequest(ctx,
		http.MethodGet,
		uri,
		WithAcceptVndJSONHeader(),
		WithAuthorizationBearerHeader(v.GetToken()),
	)
	if err != nil {
		return nil, err
	}

	resp, err := v.httpClient.Do(req, uhttp.WithJSONResponse(&organizationInvitesAPIData))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	for _, inviteAPI := range organizationInvitesAPIData.Data {
		if inviteAPI.Attributes.InviteStatus != "EXPIRED" {
			userInvites = append(userInvites, OrganizationUser{
				Id:        inviteAPI.Attributes.InviteId,
				Name:      inviteAPI.Attributes.InvitedBy,
				Email:     inviteAPI.Attributes.UserEmail,
				CreatedAt: inviteAPI.Attributes.CreatedAt,
			})
		}
	}

	return userInvites, nil
}

func (v *VGSClient) ListVaultUsers(ctx context.Context, vaultId string) ([]vaultUserAPI, error) {
	var vaultUsersAPIData vaultUsersAPIData
	if !strings.Contains(v.token.Scope, "organization-users:read") {
		return nil, fmt.Errorf("organization-users:read scope not found")
	}

	strUrl, err := url.JoinPath(v.serviceEndpoint, "/vaults/", vaultId, "/members")
	if err != nil {
		return nil, err
	}

	uri, err := url.Parse(strUrl)
	if err != nil {
		return nil, err
	}

	req, err := v.httpClient.NewRequest(ctx,
		http.MethodGet,
		uri,
		WithAcceptVndJSONHeader(),
		WithAuthorizationBearerHeader(v.GetToken()),
	)
	if err != nil {
		return nil, err
	}

	resp, err := v.httpClient.Do(req, uhttp.WithJSONResponse(&vaultUsersAPIData))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return vaultUsersAPIData.Data, nil
}

func (v *VGSClient) ListVaults(ctx context.Context) ([]Vault, error) {
	var (
		organizationVaults        []Vault
		organizationVaultsAPIData organizationVaultsAPIData
	)
	strUrl, err := url.JoinPath(v.serviceEndpoint, "/vaults")
	if err != nil {
		return nil, err
	}

	uri, err := url.Parse(strUrl)
	if err != nil {
		return nil, err
	}

	req, err := v.httpClient.NewRequest(ctx,
		http.MethodGet,
		uri,
		WithAcceptVndJSONHeader(),
		WithAuthorizationBearerHeader(v.GetToken()),
	)
	if err != nil {
		return nil, err
	}

	resp, err := v.httpClient.Do(req, uhttp.WithJSONResponse(&organizationVaultsAPIData))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	for _, vault := range organizationVaultsAPIData.Data {
		organizationVaults = append(organizationVaults, Vault{
			Id:          vault.Attributes.Identifier,
			Name:        vault.Attributes.Name,
			Environment: vault.Attributes.Environment,
			CreatedAt:   vault.Attributes.CreatedAt,
			UpdatedAt:   vault.Attributes.UpdatedAt,
		})
	}

	return organizationVaults, nil
}

func (v *VGSClient) UpdateVault(ctx context.Context, vaultIdentifier, userId, role string) error {
	var (
		body    Body
		payload = []byte(fmt.Sprintf(`{"data":{"attributes":{"role":"%s"}}}`, role))
	)
	if !strings.Contains(v.token.Scope, "organization-users:write") {
		return fmt.Errorf("organization-users:write scope not found")
	}

	strUrl, err := url.JoinPath(v.serviceEndpoint, "/vaults/", vaultIdentifier, "/members/", userId)
	if err != nil {
		return err
	}

	uri, err := url.Parse(strUrl)
	if err != nil {
		return err
	}

	err = json.Unmarshal(payload, &body)
	if err != nil {
		return err
	}

	req, err := v.httpClient.NewRequest(ctx,
		http.MethodPut,
		uri,
		WithAcceptVndJSONHeader(),
		WithAuthorizationBearerHeader(v.GetToken()),
		WithJSONBodyV2(body),
	)

	if err != nil {
		return err
	}

	resp, err := v.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		return errors.New("user details not updated")
	}

	return nil
}
