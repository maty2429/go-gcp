# 1. Configurar los proveedores que usaremos (GCP en este caso)
terraform {
  required_version = ">= 1.0.0"
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
  }
}

# 2. Conectar Terraform con tu proyecto real de Google Cloud
provider "google" {
  project = "project-4bf3e3b5-7b1bdsd-4eb3-8a3" # Tu ID de proyecto real
  region  = "southamerica-west1"            # Santiago de Chile
}

# 3. Tu recurso automatizado corregido con seguridad corporativa
resource "google_storage_bucket" "mi_balde_automatizado" {
  name          = "tudu-bucket-tf-project-4bf3e3bds5"
  location      = "SOUTHAMERICA-WEST1"
  force_destroy = true

  public_access_prevention = "enforced" # Seguridad 1: Evita que sea público hacia internet

  # Seguridad 2: ¡La perilla que soluciona el error! Obliga a usar solo IAM corporativo
  uniform_bucket_level_access = true
}