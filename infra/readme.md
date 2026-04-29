# F1 Kafka Infra

Terraform config for a GKE cluster on GCP to run the F1 Kafka pipeline.

## Architecture

| Resource | Value |
|---|---|
| Project | `pca-prep-441322` (KipaHQ) |
| Region | `us-central1` |
| Zone | `us-central1-a` |
| VPC | `kafka-f1-vpc` |
| Subnet | `10.1.0.0/24` |
| Pods CIDR | `10.2.0.0/16` |
| Services CIDR | `10.3.0.0/16` |
| Cluster | `kafka-f1-cluster` |
| Node type | `c4a-standard-2` (ARM64) |
| Node count | 1 |

> Zone (not region) is used for the cluster due to GCP quota limits on regional clusters.

## Modules

Uses local Terraform modules from [GCP-terraform-modules](https://github.com/RashRAJ/GCP-terraform-modules):
- `network/vpc` — VPC + subnet with secondary ranges for pods/services
- `compute/gke` — GKE cluster with ARM64 node pool

## Usage

```bash
# Authenticate
gcloud auth application-default login

# Init (first time)
terraform init

# Plan
terraform plan

# Apply
terraform apply --auto-approve
```

## Get cluster credentials

```bash
gcloud container clusters get-credentials kafka-f1-cluster --zone us-central1-a --project pca-prep-441322
```

## Kafka credentials

SASL credentials for Confluent Cloud are **not** stored here. Set them as environment variables:

```bash
export KAFKA_SASL_USERNAME=...
export KAFKA_SASL_PASSWORD=...
```

See [.env.example](../.env.example) at the repo root.
