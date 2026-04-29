module "vpc" {
   #   source = "https://github.com/RashRAJ/GCP-terraform-modules/tree/main/network/vpc"
  source = "/Users/rashee/hello/raj/GCP-terraform-modules/network/vpc"
    project_id   = var.account_config.project_id
    vpc_name     = var.network_config.vpc_name
    subnet_cidr   = var.network_config.subnet_cidr
    pods_cidr   = var.network_config.pods_cidr
    services_cidr = var.network_config.services_cidr
    region       = var.account_config.region
}

module "cluster" {
   #   source = "https://github.com/RashRAJ/GCP-terraform-modules/tree/main/compute/gke"
  source = "/Users/rashee/hello/raj/GCP-terraform-modules/compute/gke"
    project_id      = var.account_config.project_id
    cluster_name    = "kafka-f1-cluster"
    create_arm64_pool = true
    arm_node_count = var.node_config.arm_node_count
    arm_machine_type = var.node_config.arm_machine_type
    region       = var.account_config.zones[0]
    network      = module.vpc.network_name
    subnetwork   = module.vpc.subnet_name
  
}

