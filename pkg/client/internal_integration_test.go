package client

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	ctx             = context.Background()
	clientId, _     = os.LookupEnv("BATON_SERVICE_ACCOUNT_CLIENT_ID")
	clientSecret, _ = os.LookupEnv("BATON_SERVICE_ACCOUNT_CLIENT_SECRET")
	vaultId, _      = os.LookupEnv("BATON_VAULT")
	orgId, _        = os.LookupEnv("BATON_ORGANIZATION_ID")
	cfg             = Config{
		serviceAccountClientId:     clientId,
		serviceAccountClientSecret: clientSecret,
		organizationId:             orgId,
		vaultId:                    vaultId,
	}
)

const (
	authType = "Bearer "
	baseUrl  = "https://accounts.apps.verygoodsecurity.com"
)

func TestOrganizationResources(t *testing.T) {
	if clientId == "" && clientSecret == "" && orgId == "" && vaultId == "" {
		t.Skip()
	}

	tests := []struct {
		name     string
		resource string
	}{
		{
			name:     "Checking Members",
			resource: "members",
		},
		{
			name:     "Checking Users",
			resource: "users",
		},
		{
			name:     "Checking Invites",
			resource: "invites",
		},
	}

	cfg := Config{
		serviceAccountClientId:     clientId,
		serviceAccountClientSecret: clientSecret,
		organizationId:             orgId,
		vaultId:                    vaultId,
	}
	cli, err := getClientForTesting(ctx, cfg)
	assert.Nil(t, err)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			endpointUrl, err := url.JoinPath(baseUrl, "organizations", orgId, test.resource)
			assert.Nil(t, err)

			uri, err := url.Parse(endpointUrl)
			assert.Nil(t, err)

			req, err := getRequestForTesting(cli, uri)
			assert.Nil(t, err)

			resp, err := cli.httpClient.Do(req)
			assert.Nil(t, err)

			defer resp.Body.Close()
			res, err := io.ReadAll(resp.Body)
			assert.Nil(t, err)
			assert.NotNil(t, res)

			var data any
			err = json.Unmarshal(res, &data)
			assert.Nil(t, err)
		})
	}
}

func TestVaultMembers(t *testing.T) {
	if clientId == "" && clientSecret == "" && orgId == "" && vaultId == "" {
		t.Skip()
	}

	cli, err := getClientForTesting(ctx, cfg)
	assert.Nil(t, err)

	endpointUrl, err := url.JoinPath(baseUrl, "vaults", vaultId, "members")
	assert.Nil(t, err)

	uri, err := url.Parse(endpointUrl)
	assert.Nil(t, err)

	req, err := getRequestForTesting(cli, uri)
	assert.Nil(t, err)

	resp, err := cli.httpClient.Do(req)
	assert.Nil(t, err)

	defer resp.Body.Close()
	res, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.NotNil(t, res)

	var data any
	err = json.Unmarshal(res, &data)
	assert.Nil(t, err)

	// -- force cache response --
	resp1, err := cli.httpClient.Do(req)
	assert.Nil(t, err)

	defer resp1.Body.Close()
	res1, err := io.ReadAll(resp1.Body)
	assert.Nil(t, err)
	assert.NotNil(t, res1)

	var data1 any
	err = json.Unmarshal(res1, &data1)
	assert.Nil(t, err)
}

func getClientForTesting(ctx context.Context, cfg Config) (*VGSClient, error) {
	cli, err := New(ctx, cfg)
	return cli, err
}

func TestVaults(t *testing.T) {
	if clientId == "" && clientSecret == "" && orgId == "" && vaultId == "" {
		t.Skip()
	}

	cli, err := getClientForTesting(ctx, cfg)
	assert.Nil(t, err)

	endpointUrl, err := url.JoinPath(baseUrl, "vaults")
	assert.Nil(t, err)

	uri, err := url.Parse(endpointUrl)
	assert.Nil(t, err)

	req, err := getRequestForTesting(cli, uri)
	assert.Nil(t, err)

	resp, err := cli.httpClient.Do(req)
	assert.Nil(t, err)

	defer resp.Body.Close()
	res, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.NotNil(t, res)

	var data any
	err = json.Unmarshal(res, &data)
	assert.Nil(t, err)
}

func getRequestForTesting(cli *VGSClient, uri *url.URL) (*http.Request, error) {
	req, err := cli.httpClient.NewRequest(ctx,
		http.MethodGet,
		uri,
		WithAcceptVndJSONHeader(),
		WithAuthorizationBearerHeader(cli.GetToken()),
	)
	return req, err
}
