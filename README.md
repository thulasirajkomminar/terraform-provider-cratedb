# terraform-provider-cratedb
Terraform provider to manage CrateDB

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.20

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider

Add the below code to your configuration.

```terraform
terraform {
  required_providers {
    cratedb = {
      source = "thulasirajkomminar/cratedb"
    }
  }
}
```

Initialize the provider

```terraform
provider "cratedb" {
  api_key    = "*******"
  api_secret = "*******"
  url        = "https://console.cratedb.cloud/"
}
```

## Available functionalities

### Data Sources

* `cratedb_cluster`
* `cratedb_organization`
* `cratedb_organizations`
* `cratedb_project`

### Resources

* `cratedb_cluster`
* `cratedb_organization`
* `cratedb_project`

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `make docs`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```
