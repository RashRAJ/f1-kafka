variable "account_config" {
  description = "account details"
  type = object({
    region       = string
    project_id   = string
    project_name = string
    zones        = list(string)
  })
}

variable "network_config" {
  description = "Network and vpc details"
  type = object({
    subnet_cidr   = string
    vpc_name      = string
    pods_cidr     = string
    services_cidr = string
  })
}

variable "node_config" {
  description = "Node configuration"
  type = object({
    arm_machine_type = string
    arm_node_count   = number
  })
  
}
