# `baton-vgs` [![Go Reference](https://pkg.go.dev/badge/github.com/conductorone/baton-vgs.svg)](https://pkg.go.dev/github.com/conductorone/baton-vgs) ![main ci](https://github.com/conductorone/baton-vgs/actions/workflows/main.yaml/badge.svg)

`baton-vgs` is a connector for Very Good Security built using the [Baton SDK](https://github.com/conductorone/baton-sdk). It communicates with the Very Good Security API to sync data about users and access groups in your Very Good Security organization.
Check out [Baton](https://github.com/conductorone/baton) to learn more about the project in general.

# Getting Started

## Prerequisites

- Access to the Very Good Security dashboard.
- clientId
- clientSecret 
- organizationId
- vault

You'll need a Service account for your organization with the scope organization-users:read. To do this, you will need to use the VGS-CLI. You can find more info about it here(https://www.verygoodsecurity.com/docs/vault/developer-tools/cli-and-client-libraries), and specifically how to add a service account here(https://www.verygoodsecurity.com/docs/vgs-cli/service-account/#create). After doing that you will get the clientId and clientSecret for that account.

For simplicity, just run the following script. 
```
vgs apply service-account -O <ORG_ID> -f ./pkg/config/service_account.yaml
```

It will create a service account for you. Then you can use the clientId and clientSecret provided by running the script.  

## brew

```
brew install conductorone/baton/baton conductorone/baton/baton-vgs

BATON_SERVICE_ACCOUNT_CLIENT_ID=vgs_clientid BATON_SERVICE_ACCOUNT_CLIENT_SECRET=vgs_clientsecret BATON_ORGANIZATION_ID=vgs_organizationid BATON_VAULT=vgs_vault baton-vgs
baton resources
```

## docker

```
docker run --rm -v $(pwd):/out -e BATON_SERVICE_ACCOUNT_CLIENT_ID=vgs_clientid BATON_SERVICE_ACCOUNT_CLIENT_SECRET=vgs_clientsecret BATON_ORGANIZATION_ID=vgs_organizationid BATON_VAULT=vgs_vault ghcr.io/conductorone/baton-vgs:latest -f "/out/sync.c1z"
docker run --rm -v $(pwd):/out ghcr.io/conductorone/baton:latest -f "/out/sync.c1z" resources
```

## source

```
go install github.com/conductorone/baton/cmd/baton@main
go install github.com/conductorone/baton-vgs/cmd/baton-vgs@main

BATON_SERVICE_ACCOUNT_CLIENT_ID=vgs_clientid BATON_SERVICE_ACCOUNT_CLIENT_SECRET=vgs_clientsecret BATON_ORGANIZATION_ID=vgs_organizationid BATON_VAULT=vgs_vault baton-vgs
baton resources
```

# Data Model

`baton-vgs` will pull down information about the following Very Good Security resources:

- Users
- Organizations

# Contributing, Support and Issues

We started Baton because we were tired of taking screenshots and manually building spreadsheets. We welcome contributions, and ideas, no matter how small -- our goal is to make identity and permissions sprawl less painful for everyone. If you have questions, problems, or ideas: Please open a Github Issue!

See [CONTRIBUTING.md](https://github.com/ConductorOne/baton/blob/main/CONTRIBUTING.md) for more details.

# `baton-vgs` Command Line Usage

```
baton-vgs

Usage:
  baton-vgs [flags]
  baton-vgs [command]

Available Commands:
  capabilities       Get connector capabilities
  completion         Generate the autocompletion script for the specified shell
  help               Help about any command

Flags:
      --client-id string                       The client ID used to authenticate with ConductorOne ($BATON_CLIENT_ID)
      --client-secret string                   The client secret used to authenticate with ConductorOne ($BATON_CLIENT_SECRET)
  -f, --file string                            The path to the c1z file to sync with ($BATON_FILE) (default "sync.c1z")
  -h, --help                                   help for baton-vgs
      --log-format string                      The output format for logs: json, console ($BATON_LOG_FORMAT) (default "json")
      --log-level string                       The log level: debug, info, warn, error ($BATON_LOG_LEVEL) (default "info")
      --organization-id string                 The VGS organization id. ($BATON_ORGANIZATION_ID)
  -p, --provisioning                           This must be set in order for provisioning actions to be enabled. ($BATON_PROVISIONING)
      --service-account-client-id string       The VGS client id. ($BATON_SERVICE_ACCOUNT_CLIENT_ID)
      --service-account-client-secret string   The VGS client secret. ($BATON_SERVICE_ACCOUNT_CLIENT_SECRET)
      --vault string                           The VGS vault id. ($BATON_VAULT)
  -v, --version                                version for baton-vgs

Use "baton-vgs [command] --help" for more information about a command.
```
