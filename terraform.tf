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

variable "name_prefix" {
  type    = string
  default = "hashi"
}

variable "name_count" {
  type    = number
  default = 4
}

resource "random_pet" "name" {
  prefix = var.name_prefix
  length = var.name_count
}

output "name" {
  value = random_pet.name.id
}