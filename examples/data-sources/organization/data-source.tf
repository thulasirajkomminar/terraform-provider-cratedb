terraform {
  required_providers {
    cratedb = {
      source = "thulasirajkomminar/cratedb"
    }
  }
}

data "cratedb_organization" "default" {
  id = "667796de-3c06-4503-bc3c-a9adc2a849cc"
}

output "default_organization" {
  value = data.cratedb_organization.default
}
