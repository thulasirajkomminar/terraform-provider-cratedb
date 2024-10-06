terraform {
  required_providers {
    cratedb = {
      source = "komminarlabs/cratedb"
    }
  }
}

provider "cratedb" {}

resource "cratedb_organization" "default" {
  name = "default"
}

output "default_organization" {
  value = cratedb_organization.default
}
