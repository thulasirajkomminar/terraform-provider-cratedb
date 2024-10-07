terraform {
  required_providers {
    cratedb = {
      source = "komminarlabs/cratedb"
    }
  }
}

provider "cratedb" {}

resource "cratedb_cluster" "default" {
  organization_id = "667796de-3c06-4503-bc3c-a9adc2a849cc"
  crate_version   = "5.8.2"
  name            = "default-cluster"
  product_name    = "cr4"
  product_tier    = "default"
  project_id      = "a99eb2a8-bcf5-418c-866f-67e65a8ada40"
  subscription_id = "7c156ae9-9c07-4106-8f42-df93855876c1"
  username        = "admin"
  password        = "zyTChd9mfcGBFLb72nJkNeVj6"
}

output "default_cluster" {
  value     = cratedb_cluster.default.health
  sensitive = true
}
