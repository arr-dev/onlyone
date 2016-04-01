# Onlyone

Go service for unique bookmarks for a host

## Setup

* Install [`Terraform`](https://terraform.io/downloads.html) to `bin/`.
* Setup vars in `terraform/terraform.tfvars`, see vars in `terraform/variables.tf`.
* Terraform it with `make tfplan`, `make tfapply`.
* Deploy with `git push heroku master`.
* Import DB structure `heroku pg:psql --app [APP_NAME] < structure.sql`.
* Enjoy!
