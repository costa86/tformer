
terraform {

  cloud {
    organization = "costa-org"

    workspaces {
      name = "ws-1"
    }
  }

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.31.0"
    }
  }

  required_version = "~> 1.2"
}

