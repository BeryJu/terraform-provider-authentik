# To get the complete users list

data "authentik_service_connection_kubernetes" "local" {
  name = "Local Kubernetes Cluster"
}

# Then use `data.authentik_service_connection_kubernetes.local.id`
