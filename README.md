# Go-Kafka DevOps Case Study

A complete DevOps pipeline implementation for a Go-based microservice that processes Kafka messages and stores data in MongoDB. The application demonstrates modern containerization, Kubernetes deployment, GitOps workflow, and monitoring practices.

## Architecture

┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   GitHub Repo   │───▶│  GitHub Actions │───▶│  Docker Registry│
│                 │    │   (CI/CD)       │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                                        │
┌─────────────────┐    ┌─────────────────┐             │
│     ArgoCD      │───▶│   Kubernetes    │◄────────────┘
│   (GitOps)      │    │    Cluster      │
└─────────────────┘    └─────────────────┘
                              │
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Prometheus    │───▶│   Go App Pods   │───▶│    MongoDB      │
│   + Grafana     │    │     + Kafka     │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘


### Application Stack

- **Go Web Service**: REST API with health endpoints
- **Kafka**: Message queue for async processing (KRaft mode)
- **MongoDB**: Database for data persistence
- **Docker**: Containerization

### Infrastructure Stack

- **Kubernetes**: Container orchestration
- **ArgoCD**: GitOps continuous deployment
- **GitHub Actions**: CI/CD pipeline
- **Helm**: Package manager for Kubernetes
- **Prometheus + Grafana**: Monitoring and observability

## API Endpoints

| Endpoint | Method | Description |
| --- | --- | --- |
| `/ping` | GET | Health check endpoint |
| `/message` | POST | Send message to Kafka |

## Quick Start

### Prerequisites

- Docker and Docker Compose
- Kubernetes cluster
- kubectl configured
- Helm 3.x

### Local Development

1. **Clone Repository**
    
    ```bash
    git clone https://github.com/yenisehirli/devops-case-netmera.git
    cd devops-case-netmera
    
    ```
    
2. **Run with Docker Compose**
    
    ```bash
    docker-compose up -d
    
    ```
    
3. **Test the Application**
    
    ```bash
    # Health check
    curl http://localhost:8080/ping
    
    # Send a message
    curl -X POST http://localhost:8080/message \
      -H "Content-Type: application/json" \
      -d '{"content":"Hello DevOps!","user_id":"test-user"}'
    
    ```
    

### Production Deployment (Kubernetes)

1. **Deploy KRaft Kafka**
    
    ```bash
    helm repo add bitnami https://charts.bitnami.com/bitnami
    helm install kafka bitnami/kafka \
      --set kraft.enabled=true \
      --set zookeeper.enabled=false \
      --set listeners.client.protocol=PLAINTEXT \
      --set replicaCount=3 \
      --namespace default
    
    ```
    
2. **Create Main Kafka Service** (if not exists)
    
    ```bash
    kubectl expose service kafka-controller-headless --name=kafka --port=9092 --target-port=9092
    
    ```
    
3. **Deploy with ArgoCD**
    
    ```bash
    # Install ArgoCD
    kubectl create namespace argocd
    kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
    
    # Create Application
    kubectl apply -f argocd/application.yaml
    
    ```
    
4. **Create Kafka Topic**
    
    ```bash
    kubectl exec -it deployment/kafka -- kafka-topics.sh \
      --bootstrap-server kafka:9092 \
      --create --topic messages \
      --partitions 1 --replication-factor 1
    
    ```
    

## Application Configuration

### Environment Variables

| Variable | Description | Default |
| --- | --- | --- |
| `KAFKA_BROKERS` | Kafka broker endpoints | `kafka:9092` |
| `KAFKA_TOPIC` | Topic name for messages | `messages` |
| `MONGODB_URI` | MongoDB connection string | `mongodb://root:password@mongodb:27017` |
| `MONGODB_DATABASE` | Database name | `devops_case` |

### Message Format

```json
{
  "content": "Your message content",
  "user_id": "unique-user-identifier"
}

```

## Development Workflow

1. **Code Changes**: Make changes to Go application
2. **CI/CD Pipeline**: GitHub Actions builds Docker image and pushes to Docker Hub
3. **GitOps Deployment**: ArgoCD detects Git changes and syncs to Kubernetes
4. **Zero Downtime**: Rolling updates with health checks

## Monitoring & Observability

### Grafana Dashboards

Access Grafana at `http://localhost:3000` (admin/prom-operator)

**Key Metrics:**

- Pod CPU Usage
- Pod Memory Usage
- Pod Restart Count
- Service Status
- HTTP Request Rate
- Kafka Message Rate

### Prometheus Queries

```
# HTTP Request Rate
rate(http_requests_total[5m])

# Error Rate
rate(http_requests_total{status=~"5.."}[5m]) / rate(http_requests_total[5m])

# Response Time
histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))

```

## Message Processing Flow

1. **HTTP Request**: Client sends POST to `/message`
2. **Kafka Producer**: Go app publishes message to Kafka topic
3. **Kafka Consumer**: Same app consumes message after 10s delay
4. **MongoDB Storage**: Processed message stored in database
5. **Response**: HTTP 200 with message ID and status

## Troubleshooting

### Common Issues

**Pod CrashLoopBackOff**

```bash
kubectl logs -l app=go-kafka-app --tail=50
kubectl describe pod <pod-name>

```

**Kafka Connection Issues**

```bash
# Test Kafka connectivity
kubectl exec -it deployment/kafka -- kafka-topics.sh --bootstrap-server kafka:9092 --list

# Check service endpoints
kubectl get svc kafka

```

**MongoDB Connection Failed**

```bash
kubectl logs deployment/mongodb
kubectl exec -it deployment/mongodb -- mongosh -u root -p password

```

## Architecture Decisions

### Why KRaft Mode?

- **Modern Kafka**: ZooKeeper is deprecated, KRaft is the future
- **Simplified Deployment**: No external coordination service needed
- **Better Performance**: Reduced latency and improved scalability

### Why ArgoCD?

- **GitOps Best Practice**: Declarative and auditable deployments
- **Drift Detection**: Automatic sync and self-healing
- **Multi-Environment**: Easy promotion across environments

### Why Helm?

- **Package Management**: Reusable and configurable charts
- **Template Engine**: Environment-specific configurations
- **Release Management**: Rollback and upgrade capabilities

## Security Considerations

- **Non-root Containers**: All containers run as non-root users
- **Resource Limits**: CPU and memory limits enforced
- **Network Policies**: Pod-to-pod communication restricted
- **Secret Management**: Credentials stored in Kubernetes secrets
- **RBAC**: Role-based access control for service accounts

## Performance Tuning

### Resource Allocation

```yaml
resources:
  requests:
    cpu: 250m
    memory: 256Mi
  limits:
    cpu: 500m
    memory: 512Mi

```

### Kafka Optimization

- Batch processing for better throughput
- Appropriate partition count for parallelism
- Replication factor for availability

### MongoDB Optimization

- Connection pooling enabled
- Read/write concern optimized
- Persistent storage with proper I/O

## Testing

### Unit Tests

```bash
go test ./...

```

### Integration Tests

```bash
# Test HTTP endpoints
curl http://localhost:8080/ping
curl -X POST http://localhost:8080/message -d '{"content":"test","user_id":"test"}'

```

### Load Testing

```bash
# Using hey tool
hey -n 1000 -c 10 http://localhost:8080/ping

```

## Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

## License

This project is part of a DevOps case study for educational purposes.

## Contact

**Repository**: https://github.com/yenisehirli/devops-case-netmera

**Docker Hub**: https://hub.docker.com/r/yenisehirli/go-kafka-app

Blog Site: https://blog.mustafayenisehirli.me (you will see this project soon here)
