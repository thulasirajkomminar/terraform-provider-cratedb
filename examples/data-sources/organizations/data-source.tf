terraform {
  required_providers {
    cratedb = {
      source = "komminarlabs/cratedb"
    }
  }
}

data "cratedb_organizations" "all" {}

output "all_organizations" {
  value = data.cratedb_organizations.all
}