package main

import (
	"context"
	"fmt"
	"os"

	configSchema "github.com/conductorone/baton-sdk/pkg/config"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
	"github.com/conductorone/baton-sdk/pkg/field"
	"github.com/conductorone/baton-sdk/pkg/types"
	"github.com/conductorone/baton-vgs/pkg/client"
	"github.com/conductorone/baton-vgs/pkg/connector"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	version       = "dev"
	connectorName = "baton-vgs"
)

var (
	ServiceAccountClientId     = field.StringField(client.ServiceAccountClientIdName, field.WithRequired(true), field.WithDescription("The VGS client id."))
	ServiceAccountClientSecret = field.StringField(client.ServiceAccountClientSecretName, field.WithRequired(true), field.WithDescription("The VGS client secret."))
	OrganizationId             = field.StringField(client.OrganizationId, field.WithRequired(true), field.WithDescription("The VGS organization id."))
	Vault                      = field.StringField(client.VaultId, field.WithRequired(true), field.WithDescription("The VGS vault id."))
	configurationFields        = []field.SchemaField{Vault, ServiceAccountClientId, ServiceAccountClientSecret, OrganizationId}
)

func main() {
	ctx := context.Background()
	_, cmd, err := configSchema.DefineConfiguration(ctx,
		connectorName,
		getConnector,
		field.NewConfiguration(configurationFields),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	cmd.Version = version
	err = cmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func getConnector(ctx context.Context, cfg *viper.Viper) (types.ConnectorServer, error) {
	l := ctxzap.Extract(ctx)
	cb, err := connector.New(ctx, cfg)
	if err != nil {
		l.Error("error creating connector", zap.Error(err))
		return nil, err
	}

	c, err := connectorbuilder.NewConnector(ctx, cb)
	if err != nil {
		l.Error("error creating connector", zap.Error(err))
		return nil, err
	}

	return c, nil
}
