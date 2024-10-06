terraform {
  required_providers {
    cratedb = {
      source = "komminarlabs/cratedb"
    }
  }
}

provider "cratedb" {}
