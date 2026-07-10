# terraform-provider-cratedb
Terraform provider to manage CrateDB

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.8
- [Go](https://golang.org/doc/install) >= 1.25

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
}
```

Every provider attribute can also be set with an environment variable: `CRATEDB_API_KEY`, `CRATEDB_API_SECRET`, and `CRATEDB_URL`. The `url` attribute is optional and defaults to `https://console.cratedb.cloud`.

## Available functionalities

### Data Sources

* `cratedb_cluster`
* `cratedb_organization`
* `cratedb_organizations`
* `cratedb_project`
* `cratedb_projects`
* `cratedb_regions`

### Resources

* `cratedb_cluster`
* `cratedb_organization`
* `cratedb_project`

## Debugging

Run Terraform with `TF_LOG=DEBUG` (or `TF_LOG_PROVIDER=DEBUG`) to log every HTTP request and response the provider sends to the CrateDB Cloud API, including individual retry attempts. Credentials never appear in the logs: the `Authorization` header is added below the logging transport, and credential values are masked if they show up in request or response bodies.

```shell
TF_LOG=DEBUG TF_LOG_PATH=terraform.log terraform apply
```

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `make docs`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
export CRATEDB_API_KEY="*******"
export CRATEDB_API_SECRET="*******"
make testacc
```

The acceptance tests run against the real CrateDB Cloud API. `CRATEDB_API_KEY` and `CRATEDB_API_SECRET` are always required; tests that need additional fixtures are skipped unless the corresponding environment variable is set:

| Environment Variable | Used by | Description |
| -------------------- | ------- | ----------- |
| `CRATEDB_ORGANIZATION_ID` | project tests | An existing organization to create test resources in. |
| `CRATEDB_REGION` | project tests | Optional. The region to create test projects in. When unset, the first non-deprecated, non-edge region reported by the API is used. |
| `CRATEDB_PROJECT_ID` | cluster resource test | An existing project to deploy the test cluster into. |
| `CRATEDB_SUBSCRIPTION_ID` | cluster resource test | The subscription to bill the test cluster to. |
| `CRATEDB_CRATE_VERSION` | cluster resource test | The CrateDB version to deploy, e.g. `5.10.11`. |
| `CRATEDB_PRODUCT_NAME` / `CRATEDB_PRODUCT_TIER` | cluster resource test | Optional. Default to the free tier (`crfree` / `default`). |
| `CRATEDB_CLUSTER_ID` | cluster data source test | An existing cluster to read. |

To run a single test, pass `TESTARGS`:

```shell
make testacc TESTARGS='-run=TestAccOrganizationResource'
```
