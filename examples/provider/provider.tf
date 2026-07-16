# configuration for the provider
provider "adguardhome" {
  host     = "localhost:8080"
  username = "admin"
  password = "SecretP@ssw0rd"
  scheme   = "http" # defaults to https
  timeout  = 5      # in seconds, defaults to 10
  insecure = false  # when `true` will skip TLS validation
}
