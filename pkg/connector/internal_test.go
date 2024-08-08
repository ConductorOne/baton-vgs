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

	cli, err := getClientForTesting(ctx)
	assert.Nil(t, err)

	user := &userResourceType{
		resourceType: &v2.ResourceType{},
		client:       cli,
	}
	rs, _, _, err := user.List(ctx, &v2.ResourceId{}, &pagination.Token{})
	assert.Nil(t, err)
	assert.NotNil(t, rs)
}

func TestOrgResourceTypeList(t *testing.T) {
	if clientId == "" && clientSecret == "" && orgId == "" && vaultId == "" {
		t.Skip()
	}

	cli, err := getClientForTesting(ctx)
	assert.Nil(t, err)

	org := &orgResourceType{
		resourceType: &v2.ResourceType{},
		client:       cli,
	}
	rs, _, _, err := org.List(ctx, &v2.ResourceId{}, &pagination.Token{})
	assert.Nil(t, err)
	assert.NotNil(t, rs)
}

func TestVaultResourceTypeList(t *testing.T) {
	if clientId == "" && clientSecret == "" && orgId == "" && vaultId == "" {
		t.Skip()
	}

	cli, err := getClientForTesting(ctx)
	assert.Nil(t, err)

	vault := &vaultResourceType{
		resourceType: &v2.ResourceType{},
		client:       cli,
	}
	rs, _, _, err := vault.List(ctx, &v2.ResourceId{}, &pagination.Token{})
	assert.Nil(t, err)
	assert.NotNil(t, rs)
}

func getClientForTesting(ctx context.Context) (*client.VGSClient, error) {
	cfg := client.Config{}
	cfg.WithVaultId(vaultId).
		WithOrganizationId(orgId).
		WithServiceAccountClientId(clientId).
		WithServiceAccountClientSecret(clientSecret)
	cli, err := client.New(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return cli, nil
}

func TestClient(t *testing.T) {
	if clientId == "" && clientSecret == "" && orgId == "" && vaultId == "" {
		t.Skip()
	}

	cli, err := getClientForTesting(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, cli)
}

func TestListVaults(t *testing.T) {
	if clientId == "" && clientSecret == "" && orgId == "" && vaultId == "" {
		t.Skip()
	}

	cliTest, err := getClientForTesting(ctx)
	assert.Nil(t, err)

	lv, err := cliTest.ListVaults(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, lv)
}

func TestListVaultUsers(t *testing.T) {
	if clientId == "" && clientSecret == "" && orgId == "" && vaultId == "" {
		t.Skip()
	}

	cliTest, err := getClientForTesting(ctx)
	assert.Nil(t, err)

	lvu, err := cliTest.ListVaultUsers(ctx, vaultId)
	assert.Nil(t, err)
	assert.NotNil(t, lvu)
}

func TestListUsers(t *testing.T) {
	if clientId == "" && clientSecret == "" && orgId == "" && vaultId == "" {
		t.Skip()
	}

	cliTest, err := getClientForTesting(ctx)
	assert.Nil(t, err)

	lu, err := cliTest.ListUsers(ctx, orgId, vaultId)
	assert.Nil(t, err)
	assert.NotNil(t, lu)
}

func TestListUserInvites(t *testing.T) {
	if clientId == "" && clientSecret == "" && orgId == "" && vaultId == "" {
		t.Skip()
	}

	cliTest, err := getClientForTesting(ctx)
	assert.Nil(t, err)

	lui, err := cliTest.ListUserInvites(ctx, orgId)
	assert.Nil(t, err)
	assert.NotNil(t, lui)
}

func TestListOrganizations(t *testing.T) {
	if clientId == "" && clientSecret == "" && orgId == "" && vaultId == "" {
		t.Skip()
	}

	cliTest, err := getClientForTesting(ctx)
	assert.Nil(t, err)

	lo, err := cliTest.ListOrganizations(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, lo)
}

func TestUpdateVault(t *testing.T) {
	if clientId == "" && clientSecret == "" && orgId == "" && vaultId == "" {
		t.Skip()
	}

	cliTest, err := getClientForTesting(ctx)
	assert.Nil(t, err)

	err = cliTest.UpdateUserAccessVault(ctx, vaultId, "ID9hRKLhcc6RWBvaHQ7L1Uan", "write")
	assert.Nil(t, err)
}

func TestRevokeVault(t *testing.T) {
	if clientId == "" && clientSecret == "" && orgId == "" && vaultId == "" {
		t.Skip()
	}

	cliTest, err := getClientForTesting(ctx)
	assert.Nil(t, err)

	err = cliTest.RevokeUserAccessVault(ctx, vaultId, "IDjSP9BVbJ3RnPr2FonGxXp5")
	assert.Nil(t, err)
}
