terraform {
  backend "azurerm" {

  }
}


provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy = true
    }
  }
}

provider "azuread" {}
provider "random" {}
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "lnpay" {
  name = "lnpay"
  location = "West US"
}

resource "azuread_application" "lnpay" {
  name = "LnPay"
}

resource "azuread_service_principal" "lnpay" {
  application_id = azuread_application.lnpay.application_id
}


resource "random_password" "principal-pass" {
  length = 32
}

resource "azuread_service_principal_password" "lnpay" {
  service_principal_id = azuread_service_principal.lnpay.id
  value = random_password.principal-pass.result
  end_date = "2030-01-01T01:01:01Z"
}

resource "azurerm_role_assignment" "lnpay-dns" {
  principal_id = azuread_service_principal.lnpay.id
  role_definition_name = "DNS Zone Contributor"
  scope = azurerm_dns_zone.semtexzv.id
}

resource "azurerm_storage_account" "lnpay-static-assets" {
  location = azurerm_resource_group.lnpay.location
  resource_group_name = azurerm_resource_group.lnpay.name
  account_replication_type = "LRS"
  account_kind = "StorageV2"
  account_tier = "Standard"
  name = "lnpaystaticassets"

  static_website {
    index_document = "index.html"
  }
  blob_properties {
    cors_rule {
      allowed_headers = [
        "*"]
      allowed_methods = [
        "GET",
        "POST",
        "PUT",
        "DELETE"]
      allowed_origins = [
        "*"]
      exposed_headers = [
        "*"]
      max_age_in_seconds = 3600
    }
  }
}

resource "azurerm_cdn_profile" "lnpay" {
  location = azurerm_resource_group.lnpay.location
  resource_group_name = azurerm_resource_group.lnpay.name

  name = "lnpaycdn"
  sku = "Standard_Microsoft"
}

resource "azurerm_cdn_endpoint" "lnpay" {
  location = azurerm_resource_group.lnpay.location
  resource_group_name = azurerm_resource_group.lnpay.name

  name = "lnpay"
  profile_name = azurerm_cdn_profile.lnpay.name

  origin_host_header = azurerm_storage_account.lnpay-static-assets.primary_web_host

  origin {
    host_name = azurerm_storage_account.lnpay-static-assets.primary_web_host
    name = "lnpay"
  }


  delivery_rule {
    order = 2
    name = "CacheExpiration"
    url_path_condition {
      operator = "Any"
    }
    cache_expiration_action {
      behavior = "Override"
      duration = "00:00:30"
    }
  }

  delivery_rule {
    name = "EnforceHTTPS"
    order = 1
    request_scheme_condition {
      operator = "Equal"
      match_values = [
        "HTTP"]
    }
    url_redirect_action {
      redirect_type = "Found"
      protocol = "Https"
    }
  }
}

/*
resource "azurerm_key_vault" "lnpay" {
  location = azurerm_resource_group.lnpay.location
  name = "lnpay-vault"
  resource_group_name = azurerm_resource_group.lnpay.name
  sku_name = "standard"
  tenant_id = data.azurerm_client_config.current.tenant_id
  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id
    key_permissions = [
      "create",
      "get"]
    secret_permissions = [
      "set",
      "get"]
  }
}
*/