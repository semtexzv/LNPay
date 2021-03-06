provider "kubernetes" {
  host = azurerm_kubernetes_cluster.lnpay.kube_config.0.host
  client_certificate = base64decode(azurerm_kubernetes_cluster.lnpay.kube_config.0.client_certificate)
  client_key = base64decode(azurerm_kubernetes_cluster.lnpay.kube_config.0.client_key)
  cluster_ca_certificate = base64decode(azurerm_kubernetes_cluster.lnpay.kube_config.0.cluster_ca_certificate)
  username = azurerm_kubernetes_cluster.lnpay.kube_config.0.username
  password = azurerm_kubernetes_cluster.lnpay.kube_config.0.password
}

provider "helm" {
  kubernetes {
    host = azurerm_kubernetes_cluster.lnpay.kube_config.0.host
    client_certificate = base64decode(azurerm_kubernetes_cluster.lnpay.kube_config.0.client_certificate)
    client_key = base64decode(azurerm_kubernetes_cluster.lnpay.kube_config.0.client_key)
    cluster_ca_certificate = base64decode(azurerm_kubernetes_cluster.lnpay.kube_config.0.cluster_ca_certificate)
    username = azurerm_kubernetes_cluster.lnpay.kube_config.0.username
    password = azurerm_kubernetes_cluster.lnpay.kube_config.0.password
  }
}


resource "helm_release" "csi-secrets-azure" {
  count = 0
  namespace = "kube-system"
  repository = "https://raw.githubusercontent.com/Azure/secrets-store-csi-driver-provider-azure/master/charts"
  chart = "csi-secrets-store-provider-azure"
  name = "csi-secrets-azure"
}

resource "kubernetes_namespace" "ingress" {
  metadata {
    name = "ingress"
  }
}

resource "helm_release" "ingress-nginx" {
  repository = "https://kubernetes.github.io/ingress-nginx"
  chart = "ingress-nginx"
  name = "ingress-nginx"
  namespace = kubernetes_namespace.ingress.id
}

resource "kubernetes_namespace" "cert-manager" {
  metadata {
    name = "cert-manager"
  }
}
resource "helm_release" "cert-manager" {
  repository = "https://charts.jetstack.io"
  chart = "cert-manager"
  name = "cert-manager"
  namespace = kubernetes_namespace.cert-manager.id
  version = "1.2.0"
  replace = true

  set {
    name = "installCRDs"
    value = "true"
  }
}

resource "kubernetes_secret" "azure-config-file" {
  metadata {
    name      = "azure-config-file"
    namespace = kubernetes_namespace.ingress.id
  }
  data = {
    "azure.json" = <<EOF
      {
        "tenantId"        : "${data.azurerm_client_config.current.tenant_id}",
        "subscriptionId"  : "${data.azurerm_client_config.current.subscription_id}",
        "resourceGroup"   : "${azurerm_resource_group.lnpay.name}",
        "aadClientId"     : "${azuread_application.lnpay.application_id}",
        "aadClientSecret" : "${azuread_service_principal_password.lnpay.value}"
      }
    EOF
  }
}
resource "helm_release" "externaldns" {
  repository = "https://charts.bitnami.com/bitnami"
  chart = "external-dns"
  name = "external-dns"
  namespace = kubernetes_namespace.ingress.id
  replace = true
  set {
    name = "provider"
    value = "azure"
  }
  set {
    name  = "azure.secretName"
    value = "azure-config-file"
  }
  set {
    name = "azure.resourceGroup"
    value = azurerm_resource_group.lnpay.name
  }
  set {
    name = "azure.tenantId"
    value = data.azurerm_client_config.current.tenant_id
  }

  set {
    name = "azure.subscriptionId"
    value = data.azurerm_client_config.current.subscription_id
  }

  set {
    name = "interval"
    value = "1m"
  }
  set {
    name = "policy"
    value = "sync"
  }
  set {
    name = "logLevel"
    value = "debug"
  }
  set {
    name = "txtOwnerId"
    value = "default"
  }
}

resource "kubernetes_service" "cdn-svc" {
  metadata {
    name = "frontend-cdn"
  }
  spec {
    type = "ExternalName"
    external_name = azurerm_storage_account.lnpay-static-assets.primary_web_host
  }
}

resource "random_password" "lnd-password" {
  length = 64
}

resource "tls_private_key" "tls-key" {
  algorithm = "ECDSA"
  ecdsa_curve = "P256"
}

resource "tls_self_signed_cert" "tls-cert" {
  key_algorithm = "ECDSA"
  is_ca_certificate = true

  private_key_pem = tls_private_key.tls-key.private_key_pem
  set_subject_key_id = true
  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "server_auth",
    "cert_signing"
  ]
  validity_period_hours = 24*356*21
  subject {
    organization = "lnd autogenerated cert"
    common_name = "lnd"
  }
  dns_names = ["*", "lnd", "lnd.default.svc.cluster.local", "localhost", "unix", "unixpacket"]
}


resource "kubernetes_secret" "lnd-secret" {
  metadata {
    name = "lnd-secret"
  }

  data = {
    password = random_password.lnd-password.result
    tls_key = tls_private_key.tls-key.private_key_pem
    tls_cert = tls_self_signed_cert.tls-cert.cert_pem
  }
}
