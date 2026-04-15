# Food Delivery Platform

Production-grade CDC data pipeline built on Kubernetes with GitOps.

## Architecture

```
Order UI → Order Service (Go) → PostgreSQL → WAL → Debezium → Kafka → Consumer Pods → Elasticsearch → Search Service (Go) → Search UI
```

## Live Pipeline Demo

The UI features a live animated pipeline that lights up step by step:
- **Place Order**: UI → Order Service → PostgreSQL → Debezium → Kafka → Consumer → Elasticsearch (orange glow, turns green when complete)
- **Search**: Search Service → Elasticsearch (blue glow, shows results)

## Tech Stack

| Component | Technology |
|-----------|-----------|
| Microservices | Go (Gin framework) |
| Database | PostgreSQL (WAL logical replication) |
| CDC | Debezium (Kafka Connect) |
| Messaging | Kafka (Strimzi operator) |
| Search | Elasticsearch |
| Dashboards | Kibana + Redpanda Console (Kowl) |
| CI/CD | GitHub Actions → Docker Hub |
| GitOps | FluxCD (Kustomization + HelmRelease) |
| Orchestration | Kubernetes (Minikube) |
| Ingress | NGINX Ingress Controller |
| Containerization | Docker |

## Project Structure

```
food-delivery-platform/
├── .github/workflows/       # GitHub Actions CI/CD
│   └── build.yaml           # Builds and pushes Docker images
├── charts/                  # Helm charts (FluxCD deploys these via HelmRelease)
│   ├── order-service/
│   │   ├── Chart.yaml
│   │   ├── values.yaml
│   │   └── templates/
│   │       ├── deployment.yaml
│   │       └── service.yaml
│   ├── order-consumer/
│   │   ├── Chart.yaml
│   │   ├── values.yaml
│   │   └── templates/
│   │       └── deployment.yaml
│   └── search-service/
│       ├── Chart.yaml
│       ├── values.yaml
│       └── templates/
│           ├── deployment.yaml
│           └── service.yaml
├── clusters/                # FluxCD config (what to watch and deploy)
│   └── minikube/
│       ├── apps-infra.yaml        # KS: watches infra/ folder
│       ├── apps-services.yaml     # HR: deploys Helm charts
│       └── flux-system/           # FluxCD system files
├── infra/                   # Infrastructure manifests (FluxCD applies via KS)
│   ├── kustomization.yaml
│   ├── debezium.yaml
│   ├── elasticsearch.yaml
│   ├── kafka-cluster.yaml
│   ├── kibana.yaml
│   ├── kowl.yaml
│   └── ingress.yaml
├── order-service/           # Go source code + Dockerfile
│   ├── main.go
│   ├── go.mod
│   └── Dockerfile
├── order-consumer/          # Go source code + Dockerfile
│   ├── main.go
│   ├── go.mod
│   └── Dockerfile
├── search-service/          # Go source code + Dockerfile
│   ├── main.go
│   ├── go.mod
│   └── Dockerfile
└── ui/                      # Frontend
    └── index.html           # Modern dark UI with live pipeline animation
```

## How It Works

1. User places order via UI
2. Order Service (Go) writes to PostgreSQL
3. PostgreSQL writes to WAL (Write-Ahead Log) with `wal_level=logical`
4. Debezium reads WAL via logical replication slot
5. Debezium publishes change event to Kafka topic `fooddelivery.public.orders`
6. Consumer pods (2 replicas, same consumer group) read from Kafka
7. Consumer pods index data to Elasticsearch
8. User searches via UI → Search Service queries Elasticsearch
9. Results returned in real-time

## GitOps Flow

```
Developer pushes code → GitHub Actions builds Docker images → pushes to Docker Hub
                       → FluxCD detects changes → auto-deploys to Kubernetes
```

- **KS (Kustomization)**: Watches `infra/` folder, applies plain YAML manifests
- **HR (HelmRelease)**: Watches `charts/` folder, deploys Helm charts with values

## Access Points

| Service | URL |
|---------|-----|
| Order API | http://food.local/api/orders |
| Search API | http://food.local/api/search |
| Kibana | http://kibana.local |
| Kowl (Kafka UI) | http://kowl.local |
| UI | file://ui/index.html |

Requires `minikube tunnel` running and `/etc/hosts` entries:
```
127.0.0.1 food.local kibana.local kowl.local
```

## Quick Start

```bash
# Start Minikube
minikube start --memory=8192 --cpus=4 --driver=docker

# Create namespace
kubectl create namespace food-delivery

# Deploy PostgreSQL
helm install postgresql bitnami/postgresql -n food-delivery \
  --set auth.postgresPassword=postgres123 \
  --set primary.extendedConfiguration="wal_level = logical"

# Install Strimzi (Kafka operator)
helm install strimzi strimzi/strimzi-kafka-operator -n food-delivery

# Apply Kafka cluster
kubectl apply -f infra/kafka-cluster.yaml

# Apply Debezium
kubectl apply -f infra/debezium.yaml

# Apply Elasticsearch
kubectl apply -f infra/elasticsearch.yaml

# Enable Ingress
minikube addons enable ingress

# Bootstrap FluxCD (auto-deploys everything else)
flux bootstrap github --owner=hemant6939 --repository=food-delivery-platform \
  --branch=main --path=clusters/minikube --personal

# Start tunnel
minikube tunnel
```

## Pods Running

```
order-service          - Go microservice (writes to PG)
order-consumer (x2)    - Kafka consumer (indexes to ES)
search-service         - Go microservice (queries ES)
postgresql             - Database with WAL logical replication
kafka-controller       - Kafka broker (Strimzi)
kafka-connect          - Debezium CDC connector
elasticsearch          - Search engine
kibana                 - ES dashboard
kowl                   - Kafka dashboard
strimzi-operator       - Manages Kafka lifecycle
flux-system pods       - GitOps controllers
```
