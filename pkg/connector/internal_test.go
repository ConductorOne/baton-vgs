package connector

import (
	"context"
	"os"
	"testing"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	"github.com/conductorone/baton-vgs/pkg/client"
	"github.com/stretchr/testify/assert"
)

var (
	clientId, _     = os.LookupEnv("BATON_SERVICE_ACCOUNT_CLIENT_ID")
	clientSecret, _ = os.LookupEnv("BATON_SERVICE_ACCOUNT_CLIENT_SECRET")
	vaultId, _      = os.LookupEnv("BATON_VAULT")
	orgId, _        = os.LookupEnv("BATON_ORGANIZATION_ID")
	ctx             = context.Background()
)

func TestUserResourceTypeList(t *testing.T) {
	if clientId == "" && clientSecret == "" && orgId == "" && vaultId == "" {
		t.Skip()
	}

	user := &userResourceType{
		resourceType: &v2.ResourceType{},
		client:       getClientForTesting(ctx, clientId, clientSecret, orgId, vaultId),
	}
	rs, _, _, err := user.List(ctx, &v2.ResourceId{}, &pagination.Token{})
	assert.Nil(t, err)
	assert.NotNil(t, rs)
}

func TestOrgResourceTypeList(t *testing.T) {
	if clientId == "" && clientSecret == "" && orgId == "" && vaultId == "" {
		t.Skip()
	}

	org := &orgResourceType{
		resourceType: &v2.ResourceType{},
		client:       getClientForTesting(ctx, clientId, clientSecret, orgId, vaultId),
	}
	rs, _, _, err := org.List(ctx, &v2.ResourceId{}, &pagination.Token{})
	assert.Nil(t, err)
	assert.NotNil(t, rs)
}

func TestVaultResourceTypeList(t *testing.T) {
	if clientId == "" && clientSecret == "" && orgId == "" && vaultId == "" {
		t.Skip()
	}

	vault := &vaultResourceType{
		resourceType: &v2.ResourceType{},
		client:       getClientForTesting(ctx, clientId, clientSecret, orgId, vaultId),
	}
	rs, _, _, err := vault.List(ctx, &v2.ResourceId{}, &pagination.Token{})
	assert.Nil(t, err)
	assert.NotNil(t, rs)
}

func getClientForTesting(ctx context.Context, clientId, clientSecret, orgId, vaultId string) *client.VGSClient {
	cli, _ := client.New(ctx, clientId, clientSecret, orgId, vaultId)
	return cli
}

func TestClient(t *testing.T) {
	if clientId == "" && clientSecret == "" && orgId == "" && vaultId == "" {
		t.Skip()
	}

	cli, err := client.New(ctx, clientId, clientSecret, orgId, vaultId)
	assert.Nil(t, err)
	assert.NotNil(t, cli)
}

func TestListVaults(t *testing.T) {
	if clientId == "" && clientSecret == "" && orgId == "" && vaultId == "" {
		t.Skip()
	}

	cliTest := getClientForTesting(ctx, clientId, clientSecret, orgId, vaultId)
	lv, err := cliTest.ListVaults(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, lv)
}

func TestListVaultUsers(t *testing.T) {
	if clientId == "" && clientSecret == "" && orgId == "" && vaultId == "" {
		t.Skip()
	}

	cliTest := getClientForTesting(ctx, clientId, clientSecret, orgId, vaultId)
	lvu, err := cliTest.ListVaultUsers(ctx, vaultId)
	assert.Nil(t, err)
	assert.NotNil(t, lvu)
}

func TestListUsers(t *testing.T) {
	if clientId == "" && clientSecret == "" && orgId == "" && vaultId == "" {
		t.Skip()
	}

	cliTest := getClientForTesting(ctx, clientId, clientSecret, orgId, vaultId)
	lu, err := cliTest.ListUsers(ctx, orgId, vaultId)
	assert.Nil(t, err)
	assert.NotNil(t, lu)
}

func TestListUserInvites(t *testing.T) {
	if clientId == "" && clientSecret == "" && orgId == "" && vaultId == "" {
		t.Skip()
	}

	cliTest := getClientForTesting(ctx, clientId, clientSecret, orgId, vaultId)
	lui, err := cliTest.ListUserInvites(ctx, orgId)
	assert.Nil(t, err)
	assert.NotNil(t, lui)
}

func TestListOrganizations(t *testing.T) {
	if clientId == "" && clientSecret == "" && orgId == "" && vaultId == "" {
		t.Skip()
	}

	cliTest := getClientForTesting(ctx, clientId, clientSecret, orgId, vaultId)
	lo, err := cliTest.ListOrganizations(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, lo)
}

func TestUpdateVault(t *testing.T) {
	if clientId == "" && clientSecret == "" && orgId == "" && vaultId == "" {
		t.Skip()
	}

	cliTest := getClientForTesting(ctx, clientId, clientSecret, orgId, vaultId)
	err := cliTest.UpdateUserAccessVault(ctx, vaultId, "ID9hRKLhcc6RWBvaHQ7L1Uan", "write")
	assert.Nil(t, err)
}

func TestRevokeVault(t *testing.T) {
	if clientId == "" && clientSecret == "" && orgId == "" && vaultId == "" {
		t.Skip()
	}

	cliTest := getClientForTesting(ctx, clientId, clientSecret, orgId, vaultId)
	err := cliTest.RevokeUserAccessVault(ctx, vaultId, "IDjSP9BVbJ3RnPr2FonGxXp5")
	assert.Nil(t, err)
}
