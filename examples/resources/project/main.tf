terraform {
  required_providers {
    cratedb = {
      source = "komminarlabs/cratedb"
    }
  }
}

provider "cratedb" {}

resource "cratedb_project" "default" {
  name            = "default-project"
  organization_id = "667796de-3c06-4503-bc3c-a9adc2a849cc"
  region          = "eks1.eu-west-1.aws"
}

output "default_project" {
  value = cratedb_project.default
}
