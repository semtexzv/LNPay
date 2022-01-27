resource random_password ui_admin_password {
  length = 32
  special = false
  number = true
  upper = true
}

output "acr-repo" {
  value = azurerm_container_registry.lnpay.login_server
}

output "azure-static-storage" {
  value = azurerm_storage_account.lnpay-static-assets.name
}

output ui_admin_password {
  value = random_password.ui_admin_password.result
}

resource "local_file" "outputs" {
  filename = "../vars.env"
  content = <<EOF
ACR_REPO=${azurerm_container_registry.lnpay.login_server}
ASSETS_STORAGE_ACCOUNT=${azurerm_storage_account.lnpay-static-assets.name}
DNS_PREFIX=${local.DNS_PREFIX}
APP_URL=${local.APP_URL}
API_URL=https://api.${local.DNS_PREFIX}lnpay.semtexzv.com
ADMIN_PASSWORD=${random_password.ui_admin_password.result}
NETWORK=${var.network}
EOF

}