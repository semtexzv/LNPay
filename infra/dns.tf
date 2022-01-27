variable network {
  type = string
}

locals {
  DNS_PREFIX = "${var.network}${var.network == "" ? "" : "."}"
  APP_URL = "${azurerm_dns_cname_record.lnpay.name}.${azurerm_dns_zone.semtexzv.name}"
}

resource "azurerm_dns_zone" "semtexzv" {
  name = "semtexzv.com"
  resource_group_name = azurerm_resource_group.lnpay.name
}

resource "azurerm_dns_cname_record" "lnpay" {
  zone_name = azurerm_dns_zone.semtexzv.name
  resource_group_name = azurerm_resource_group.lnpay.name
  name = "${local.DNS_PREFIX}lnpay"
  target_resource_id = azurerm_cdn_endpoint.lnpay.id
  ttl = 300
}

resource null_resource cnd_domain_apply {
  provisioner "local-exec" {
    command = <<EOF
az cdn custom-domain create \
--name lnpay-cdn \
--resource-group ${azurerm_resource_group.lnpay.name} \
--endpoint-name ${azurerm_cdn_endpoint.lnpay.name} \
--profile-name ${azurerm_cdn_profile.lnpay.name} \
--hostname ${local.APP_URL}
EOF
  }

  provisioner "local-exec" {
    command = <<EOF
az cdn custom-domain enable-https \
--name lnpay-cdn \
--resource-group ${azurerm_resource_group.lnpay.name} \
--endpoint-name ${azurerm_cdn_endpoint.lnpay.name} \
--profile-name ${azurerm_cdn_profile.lnpay.name} \
--min-tls-version 1.2
EOF
  }
}

/*
resource "azurerm_dns_cname_record" "lnpay-app" {
  zone_name = azurerm_dns_zone.semtexzv.name
  resource_group_name = azurerm_resource_group.lnpay.name
  name = "app.lnpay"
  target_resource_id = azurerm_cdn_endpoint.lnpay.id
  ttl = 300
}*/