terraform {
  required_providers {
    cratedb = {
      source = "thulasirajkomminar/cratedb"
    }
  }
}

data "cratedb_projects" "all" {}

output "project_names" {
  value = [for project in data.cratedb_projects.all.projects : project.name]
}
