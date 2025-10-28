terraform {
  required_providers {
    cratedb = {
      source = "thulasirajkomminar/cratedb"
    }
  }
}

data "cratedb_cluster" "default" {
  id = "156e7f96-0f6e-4fcc-8940-6e2a52efcee3"
}

output "default_cluster" {
  value     = data.cratedb_cluster.default
  sensitive = true
}
