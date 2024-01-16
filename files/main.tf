terraform {
  cloud {
    organization = "costa-org"

    workspaces {
      name = "sample"
    }
  }
}