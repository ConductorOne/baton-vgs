# `baton-vgs` [![Go Reference](https://pkg.go.dev/badge/github.com/conductorone/baton-vgs.svg)](https://pkg.go.dev/github.com/conductorone/baton-vgs) ![main ci](https://github.com/conductorone/baton-vgs/actions/workflows/main.yaml/badge.svg)

`baton-vgs` is a connector for Very Good Security built using the [Baton SDK](https://github.com/conductorone/baton-sdk). It communicates with the Very Good Security API to sync data about users and access groups in your Very Good Security organization.
Check out [Baton](https://github.com/conductorone/baton) to learn more about the project in general.

# Getting Started

## Prerequisites

- Access to the Very Good Security dashboard.
- API key. To get the API key log in to the Very Good Security dashboard and go to User Profile -> API Tokens -> View button of Global API Key
- Email - email used to login to Very Good Security dashboard.
- Account ID

## brew

```
brew install conductorone/baton/baton conductorone/baton/baton-vgs

BATON_ACCOUNT_ID=cloudflareAccountId BATON_API_KEY=cloudflareApiKey BATON_EMAIL=yourEmail baton-vgs
baton resources
```

## docker

```
docker run --rm -v $(pwd):/out -e BATON_ACCOUNT_ID=cloudflareAccountId BATON_API_KEY=cloudflareApiKey BATON_EMAIL=yourEmail ghcr.io/conductorone/baton-vgs:latest -f "/out/sync.c1z"
docker run --rm -v $(pwd):/out ghcr.io/conductorone/baton:latest -f "/out/sync.c1z" resources
```

## source

```
go install github.com/conductorone/baton/cmd/baton@main
go install github.com/conductorone/baton-vgs/cmd/baton-vgs@main

BATON_ACCOUNT_ID=cloudflareAccountId BATON_API_KEY=cloudflareApiKey BATON_EMAIL=yourEmail baton-vgs
baton resources
```

# Data Model

`baton-vgs` will pull down information about the following Very Good Security resources:

- Users
- Access Groups

# Contributing, Support and Issues

We started Baton because we were tired of taking screenshots and manually building spreadsheets. We welcome contributions, and ideas, no matter how small -- our goal is to make identity and permissions sprawl less painful for everyone. If you have questions, problems, or ideas: Please open a Github Issue!

See [CONTRIBUTING.md](https://github.com/ConductorOne/baton/blob/main/CONTRIBUTING.md) for more details.

# `baton-vgs` Command Line Usage

```
baton-vgs

Usage:
  baton-cloudflare-zero-trust [flags]
  baton-cloudflare-zero-trust [command]

Available Commands:
  capabilities       Get connector capabilities
  completion         Generate the autocompletion script for the specified shell
  help               Help about any command

Flags:
      --account-id string      Cloudflare account ID ($BATON_ACCOUNT_ID)
      --api-key string         Cloudflare API key ($BATON_API_KEY)
      --api-token string       Cloudflare API token ($BATON_API_TOKEN)
      --client-id string       The client ID used to authenticate with ConductorOne ($BATON_CLIENT_ID)
      --client-secret string   The client secret used to authenticate with ConductorOne ($BATON_CLIENT_SECRET)
      --email string           Cloudflare account email ($BATON_EMAIL)
  -f, --file string            The path to the c1z file to sync with ($BATON_FILE) (default "sync.c1z")
  -h, --help                   help for baton-cloudflare-zero-trust
      --log-format string      The output format for logs: json, console ($BATON_LOG_FORMAT) (default "json")
      --log-level string       The log level: debug, info, warn, error ($BATON_LOG_LEVEL) (default "info")
  -p, --provisioning           This must be set in order for provisioning actions to be enabled. ($BATON_PROVISIONING)
  -v, --version                version for baton-cloudflare-zero-trust

Use "baton-cloudflare-zero-trust [command] --help" for more information about a command.
```
