# Food Delivery Platform

Production-grade CDC data pipeline built on Kubernetes.

## Architecture




## Tech Stack

- **Microservices:** Go (Gin)
- **Database:** PostgreSQL (WAL logical replication)
- **CDC:** Debezium (Kafka Connect)
- **Messaging:** Kafka (Strimzi operator)
- **Search:** Elasticsearch + Kibana
- **Monitoring:** Redpanda Console (Kowl)
- **CI/CD:** GitHub Actions → Docker Hub
- **GitOps:** FluxCD (HelmRelease)
- **Orchestration:** Kubernetes (Minikube)
- **Ingress:** NGINX Ingress Controller

## Project Structure


## How It Works

1. User places order via UI
2. Order Service writes to PostgreSQL
3. Debezium captures WAL changes → sends to Kafka
4. Consumer pods read from Kafka → index to Elasticsearch
5. User searches via UI → Search Service queries Elasticsearch
