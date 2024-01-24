
variable "name_two" {
  type    = string
  default = "two"
}
variable "name_new" {
  type    = string
  default = "new"
}


variable "name_count_two" {
  type    = number
  default = 22
}

resource "random_pet" "name_two" {
  prefix = var.name_two
  length = var.name_count_two
}

output "name_two" {
  value = random_pet.name_two
}
