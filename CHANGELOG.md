# Changelog

All notable changes to this project will automatically be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## v1.0.0 - 2026-07-10

### Added

* **New Data Source:** `cratedb_projects`
* **New Data Source:** `cratedb_regions`
* The provider `url` is now optional and defaults to `https://console.cratedb.cloud` (trailing slashes are trimmed). It can still be overridden via the `url` attribute or the `CRATEDB_URL` environment variable.
* With `TF_LOG=DEBUG`, every HTTP request/response (including each retry attempt) is logged. Credentials never appear in the logs.
* Acceptance tests for every resource and data source (`make testacc`, see the README for the required environment variables).

### Fixed

* Resources deleted outside of Terraform are now removed from state on refresh (plans a re-create) instead of failing with an API error.
* Deleting a resource that is already gone (HTTP 404) no longer fails the destroy.
* Unexpected API responses (e.g. a proxy returning HTML with status 200) now produce a diagnostic containing the HTTP status and raw response body instead of a provider crash or an unhelpful formatting error.
* `dc.created` / `dc.modified` timestamps now use the RFC3339 type, so semantically equal values (`+00:00` vs `Z`) no longer show as drift, and cleared values no longer linger in state.
* The `cratedb_cluster` data source no longer fails on every read: its schema declared attributes (`health.last_seen`, `health.running_operation`, `last_async_operation`, and a missing `organization_id` mapping) that the provider never populated.
* Credentials are no longer attached to the provider's log context.
* The `cratedb_regions` data source no longer fails with a time parsing error: the API returns `last_seen` timestamps without a timezone offset, which are now decoded tolerantly (interpreted as UTC).

### Changed

* `terraform import` for `cratedb_cluster`, `cratedb_organization`, and `cratedb_project` now takes the resource **id** (previously the name was stored, which never allowed a successful refresh after import).

## v0.4.0 - 2026-02-27

### What's Changed

* chore: modified the cratedb package import url.
* chore(deps): dependency updates.

## v0.3.0 - 2025-10-29

### What's Changed

* chore: moved the repo from org to personal account.
* chore: removed auto assign action.
* chore(deps): dependency updates.

> [!Important]
>
> The older versions of the provider are still available under the `komminarlabs/cratedb` namespace on the Terraform Registry. But the new versions (`v0.3.0` and above) will be available under the `thulasirajkomminar/cratedb` namespace. Please update your provider source accordingly in your Terraform configurations.

## v0.2.0 - 2025-03-11

### What's Changed

* chore: go module upgrades.
* chore(deps): dependency updates.

## v0.1.0 - 2024-10-08

### Added

* **New Data Source:** `cratedb_cluster`
* **New Data Source:** `cratedb_organization`
* **New Data Source:** `cratedb_organizations`
* **New Data Source:** `cratedb_project`
* **New Resource:** `cratedb_cluster`
* **New Resource:** `cratedb_organization`
* **New Resource:** `cratedb_project`
