resource "heroku_app" "onlyone" {
  name = "${var.heroku_app_name}"
  region = "${var.heroku_app_region}"

  config_vars {
    AUTH_PASS = "${var.env_AUTH_PASS}"
    AUTH_USER = "${var.env_AUTH_USER}"
    DATABASE_MAX_CONNS = "${var.env_DATABASE_MAX_CONNS}"
  }

  provisioner "local-exec" {
    command = "git remote show heroku > /dev/null || git remote add heroku ${self.git_url}"
  }
}

resource "heroku_addon" "db" {
  app = "${heroku_app.onlyone.name}"
  plan = "heroku-postgresql:hobby-dev"
}
