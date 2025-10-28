terraform {
  required_providers {
    cratedb = {
      source = "thulasirajkomminar/cratedb"
    }
  }
}

data "cratedb_project" "default" {
  id = "2f310566-f171-4bf6-bf2e-46e045ff3708"
}

output "default_project" {
  value = data.cratedb_project.default
}
