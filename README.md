# Onlyone

Go service for unique bookmarks for a host

## Setup

App is deployed to https://fly.io with database hosted on https://supabase.com

To deploy, install [`flyctl`](https://fly.io/docs/hands-on/install-flyctl/), auth with `flyctl auth login` and run `flyctl deploy`.

## Old heroku Setup

* Install [`Terraform`](https://terraform.io/downloads.html) to `bin/`.
* Setup vars in `terraform/terraform.tfvars`, see vars in `terraform/variables.tf`.
* Terraform it with `make tfplan`, `make tfapply`.
* Deploy with `git push heroku master`.
* Import DB structure `heroku pg:psql --app [APP_NAME] < structure.sql`.
* Enjoy!
