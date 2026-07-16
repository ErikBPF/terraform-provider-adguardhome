terraform {
  required_providers {
    adguardhome = {
      source  = "ErikBPF/adguardhome"
      version = "0.1.0"
    }
  }
}

# configuration for the provider
provider "adguardhome" {
  host     = "localhost:8080"
  username = "admin"
  password = "SecretP@ssw0rd"
  scheme   = "http" # defaults to https
  timeout  = 5      # in seconds, defaults to 10
}
