package connector

import (
	"context"
	"io"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
	"github.com/conductorone/baton-vgs/pkg/client"
	"github.com/spf13/viper"
)

type (
	Connector struct {
		client *client.VGSClient
	}
)

// ResourceSyncers returns a ResourceSyncer for each resource type that should be synced from the upstream service.
func (d *Connector) ResourceSyncers(ctx context.Context) []connectorbuilder.ResourceSyncer {
	return []connectorbuilder.ResourceSyncer{
		userBuilder(d.client),
		orgBuilder(d.client),
		vaultBuilder(d.client),
	}
}

// Asset takes an input AssetRef and attempts to fetch it using the connector's authenticated http client
// It streams a response, always starting with a metadata object, following by chunked payloads for the asset.
func (d *Connector) Asset(ctx context.Context, asset *v2.AssetRef) (string, io.ReadCloser, error) {
	return "", nil, nil
}

// Metadata returns metadata about the connector.
func (d *Connector) Metadata(ctx context.Context) (*v2.ConnectorMetadata, error) {
	return &v2.ConnectorMetadata{
		DisplayName: "VGS Connector",
		Description: "Connector syncing users, organizations and vaults from VGS.",
	}, nil
}

// Validate is called to ensure that the connector is properly configured. It should exercise any API credentials
// to be sure that they are valid.
func (d *Connector) Validate(ctx context.Context) (annotations.Annotations, error) {
	return nil, nil
}

// New returns a new instance of the connector.
func New(ctx context.Context, cfg *viper.Viper) (*Connector, error) {
	var (
		vc             *client.VGSClient
		config         = client.Config{}
		clientId       = cfg.GetString(client.ServiceAccountClientIdName)
		clientSecret   = cfg.GetString(client.ServiceAccountClientSecretName)
		organizationId = cfg.GetString(client.OrganizationId)
		vaultId        = cfg.GetString(client.VaultId)
		err            error
	)

	config.WithServiceAccountClientId(clientId).WithServiceAccountClientSecret(clientSecret)
	config.WithOrganizationId(organizationId).WithVaultId(vaultId)
	if clientId != "" && clientSecret != "" {
		vc, err = client.New(ctx, config)
		if err != nil {
			return nil, err
		}
	}

	return &Connector{
		client: vc,
	}, nil
}
