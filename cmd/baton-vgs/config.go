package main

import (
	"context"
	"errors"

	"github.com/conductorone/baton-sdk/pkg/cli"
	"github.com/spf13/cobra"
)

// config defines the external configuration required for the connector to run.
type config struct {
	cli.BaseConfig `mapstructure:",squash"` // Puts the base config options in the same place as the connector options
	ClientId       string                   `mapstructure:"clientid"`
	ClientSecret   string                   `mapstructure:"clientsecret"`
}

// validateConfig is run after the configuration is loaded, and should return an error if it isn't valid.
func validateConfig(ctx context.Context, cfg *config) error {
	if cfg.ClientId == "" {
		return errors.New("clientid is required")
	}

	if cfg.ClientSecret == "" {
		return errors.New("clientsecret is required")
	}

	return nil
}

// cmdFlags sets the cmdFlags required for the connector.
func cmdFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().String("clientid", "", "The VGS client id. ($BATON_VGS_CLIENT_ID)")
	cmd.PersistentFlags().String("clientsecret", "", "The VGS client secret. ($BATON_VGS_CLIENT_SECRET)")
}
