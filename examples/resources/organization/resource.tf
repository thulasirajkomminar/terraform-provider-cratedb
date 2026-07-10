terraform {
  required_providers {
    cratedb = {
      source = "thulasirajkomminar/cratedb"
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
