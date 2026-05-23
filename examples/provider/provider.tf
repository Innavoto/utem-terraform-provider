terraform {
  required_providers {
    utem = {
      source  = "Innavoto/utem"
      version = "~> 0.1"
    }
  }
}

provider "utem" {
  # Set via environment variables:
  #   UTEM_BASE_URL  (default: https://utem.innavoto.com)
  #   UTEM_API_KEY   (required)
  #   UTEM_TENANT_ID (default: 1)
}
