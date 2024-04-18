package main

import (
	"context"
	"errors"

	"github.com/conductorone/baton-sdk/pkg/cli"
	"github.com/spf13/cobra"
)

// config defines the external configuration required for the connector to run.
type config struct {
	cli.BaseConfig               `mapstructure:",squash"` // Puts the base config options in the same place as the connector options
	Vault                        string                   `mapstructure:"vault"`
	ServiceAccountClientId       string                   `mapstructure:"service-account-client-id"`
	ServiceAccountClientSecret   string                   `mapstructure:"service-account-client-secret"`
	ServiceAccountOrganizationId string                   `mapstructure:"serviceorganizationid"`
}

// validateConfig is run after the configuration is loaded, and should return an error if it isn't valid.
func validateConfig(ctx context.Context, cfg *config) error {
	if cfg.Vault == "" {
		return errors.New("vault is required")
	}

	if cfg.ServiceAccountClientId == "" {
		return errors.New("service-account-client-id is required")
	}

	if cfg.ServiceAccountClientSecret == "" {
		return errors.New("service-account-client-secret is required")
	}

	if cfg.ServiceAccountOrganizationId == "" {
		return errors.New("service-account-organization-id is required")
	}

	return nil
}

// cmdFlags sets the cmdFlags required for the connector.
func cmdFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().String("vault", "", "The VGS vault id. ($BATON_VAULT)")
	cmd.PersistentFlags().String("service-account-client-id", "", "The VGS client id. ($BATON_SERVICE_ACCOUNT_CLIENT_ID)")
	cmd.PersistentFlags().String("service-account-client-secret", "", "The VGS client secret. ($BATON_SERVICE_ACCOUNT_CLIENT_SECRET)")
	cmd.PersistentFlags().String("service-account-organization-id", "", "The VGS organization id. ($BATON_SERVICE_ACCOUNT_ORGANIZATION_ID)")
}
