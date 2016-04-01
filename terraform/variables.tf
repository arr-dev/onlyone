variable "heroku_app_name" {
  description = "Name of the Heroku app"
  default = "onlyone"
}
variable "heroku_app_region" {
  description = "Region of the Heroku app"
  default = "eu"
}

variable "env_AUTH_PASS" {
  description = "sets AUTH_PASS ENV var"
}
variable "env_AUTH_USER" {
  description = "sets AUTH_USER ENV var"
}
variable "env_DATABASE_MAX_CONNS" {
  description = "sets DATABASE_MAX_CONNS var"
  default = 15
}

variable "heroku_email" {
  description = "Heroku email"
}
variable "heroku_api_key" {
  description = "Heroku api key"
}

provider "heroku" {
  email = "${var.heroku_email}"
  api_key = "${var.heroku_api_key}"
}

