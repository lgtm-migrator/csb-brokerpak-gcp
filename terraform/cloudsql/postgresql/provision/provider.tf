provider "google" {
  credentials = var.credentials
  project     = var.project
}

provider "postgresql" {
  host      = google_sql_database_instance.instance.first_ip_address
  port      = 5432
  username  = google_sql_user.admin_user.name
  password  = google_sql_user.admin_user.password
  superuser = false
  database  = google_sql_database.database.name
  sslmode   = var.use_tls ? "require" : "disable"
  clientcert {
      cert = "${path.module}/client_ca_cert.pem"
      key  = "${path.module}/client_private_key.pem"
  }
}

provider "local" {}