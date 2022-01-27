
resource "azurerm_kubernetes_cluster" "lnpay" {
  name = "lnpay"
  dns_prefix = "lnpay"
  location = azurerm_resource_group.lnpay.location
  resource_group_name = azurerm_resource_group.lnpay.name
  identity {
    type = "SystemAssigned"
  }

  addon_profile {
    aci_connector_linux {
      enabled = false
    }
    azure_policy {
      enabled = false
    }
    http_application_routing {
      enabled = false
    }
    oms_agent {
      enabled = false

    }
    kube_dashboard {
      enabled = true
    }
  }

  default_node_pool {
    name = "default"
    vm_size = "standard_b2s"
    os_disk_size_gb = 30
    node_count = 1
  }
}

resource "azurerm_container_registry" "lnpay" {
  location = azurerm_resource_group.lnpay.location
  name = "lnpay"
  resource_group_name = azurerm_resource_group.lnpay.name
  sku = "Basic"
  admin_enabled = true
}

resource "null_resource" "login-acr" {
  depends_on = [
    azurerm_container_registry.lnpay
  ]

  provisioner "local-exec" {
    command = "docker login ${azurerm_container_registry.lnpay.login_server} -u ${azurerm_container_registry.lnpay.admin_username} -p ${azurerm_container_registry.lnpay.admin_password}"
  }
}

# add the role to the identity the kubernetes cluster was assigned
resource "azurerm_role_assignment" "kubweb_to_acr" {
  scope = azurerm_container_registry.lnpay.id
  role_definition_name = "AcrPull"
  principal_id = azurerm_kubernetes_cluster.lnpay.kubelet_identity[0].object_id
}
