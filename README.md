# 🚀 Go Test App on Kubernetes

This project contains a sample Go web application packaged into a Docker container and deployed on a local Kubernetes cluster using Minikube.

---

## 🧱 Project Structure

```bash
.
├── Dockerfile          # Builds the Go application container
├── deployment.yaml     # Kubernetes Deployment configuration
├── service.yaml        # Kubernetes Service configuration
└── README.md           # You're here!
```

## 📊 Load Testing and Auto-scaling Results

### Load Test Configuration
The application was stress tested using `wrk`, a modern HTTP benchmarking tool, with the following parameters:
- Threads: 8
- Connections: 200
- Duration: 60 seconds
- Target endpoint: http://192.168.49.2/tasks

### Performance Results
```bash
Thread Stats:
- Average Latency: 37.45ms
- Latency Std Dev: 104.69ms
- Max Latency: 1.23s
- Latency Distribution: 97.45% within standard deviation

Throughput:
- Average Requests/sec: 11,686.53
- Total Requests: 726,794 in 1.04 minutes
- Data Transfer: 108.28MB total (1.74MB/sec)
```

### Horizontal Pod Autoscaling (HPA) Behavior
The application is configured with HPA (Horizontal Pod Autoscaler) with the following specifications:
- Target CPU Utilization: 70%
- Min Pods: 2
- Max Pods: 5

During the load test, we observed the following scaling behavior:

1. Initial State:
   - 2 replicas running
   - CPU utilization: 0%

2. Under Load:
   - CPU usage spiked to 115%
   - HPA triggered scaling from 2 to 4 replicas

3. Load Distribution:
   - After scaling, CPU usage decreased to 76%
   - Eventually stabilized at 0% as load normalized

This demonstrates that our HPA configuration successfully handled the increased load by automatically scaling the application from 2 to 4 pods when CPU utilization exceeded the 70% threshold.
