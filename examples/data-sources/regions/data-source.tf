terraform {
  required_providers {
    cratedb = {
      source = "thulasirajkomminar/cratedb"
    }
  }
}

data "cratedb_regions" "all" {}

# Optionally include an organization's edge regions.
data "cratedb_regions" "org" {
  organization_id = "667796de-3c06-4503-bc3c-a9adc2a849cc"
}

output "region_names" {
  value = [for region in data.cratedb_regions.all.regions : region.name]
}
