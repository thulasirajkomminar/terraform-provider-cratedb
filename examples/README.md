# Examples

This directory contains examples that are mostly used for documentation, but can also be run/tested manually via the Terraform CLI.

Example directories are named after the resource or data source **without** the `cratedb_` prefix (e.g. `resources/cluster` for `cratedb_cluster`). The templates in [`templates/`](../templates) resolve these paths when generating documentation with `make docs`:

* **provider/provider.tf** example file for the provider index page
* **data-sources/`name`/data-source.tf** example file for the `cratedb_<name>` data source page
* **resources/`name`/resource.tf** example file for the `cratedb_<name>` resource page
* **resources/`name`/import.sh** import example for the `cratedb_<name>` resource page

All other *.tf files are ignored by the documentation tool. This is useful for creating examples that can run and/or are testable even if some parts are not relevant for the documentation.
