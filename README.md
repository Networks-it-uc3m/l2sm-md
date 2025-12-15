# L2S-CES Layer 2 Secure Slices

**L2S-CES** is a powerful gRPC and CLI client designed to centrally manage multi-domain L2 networks and inter-cluster overlays. It simplifies complex multi-cluster network configurations, enabling efficient management of inter-cluster connectivity and L2 network resources from a single control point.

---

### Key Features:
- **Single Point Management:** Create and manage inter-cluster overlays and L2Networks effortlessly.
- **Integration with Kubernetes:** Leverages Kubernetes Custom Resources (CRs) for resource management.
- **Inter Domain Connectivity:** Supports IDCO (Inter Domain Connectivity Orchestrator) for unified network orchestration across clusters.

For more in-depth concepts about L2S-M and inter-domain functionalities, visit the [main repository](https://www.github.com/Networks-it-uc3m/L2S-M).

---

## 🚀 Quick Deployment

Deploy directly from GitHub:
```bash
kubectl apply -f https://github.com/Networks-it-uc3m/l2sc-es/raw/refs/heads/development/deployments/l2sces-deployment.yaml
```

Or locally:
```bash
kubectl apply -f ./deployments/l2sces-deployment.yaml
```

### Prerequisites
- **Available Ports:** By default, L2S-M MD uses NodePort from Kubernetes to expose its microservices. However, for improved scalability, security, and manageability—particularly in production environments—we recommend using a **LoadBalancer service** or implementing an **Ingress Controller** instead of NodePort.
    - **30663 (TCP)**: OpenFlow communication with NEDs in other clusters.
    - **30808 (TCP)**: Communication with IDCO Provider for network status and network creation requests.
    - **30053 (TCP/UDP)**: DNS resolution requests from managed clusters (if DNS is used).
    - **30818 (TCP)**: DNS gRPC service for managing DNS entries (if DNS is used).
    - **30051 (TCP)**: gRPC server interface for external network management and provisioning (optional).

- **Minimum two clusters** with the base L2S-M installation ([installation guide](https://www.github.com/Networks-it-uc3m/L2S-M/deployment)).

---

## ⚙️ Configuration

Default deployment includes DNS and IDCO services in namespace `l2sm-system`. Customize your deployment settings, ports, namespace, and microservices using [Kustomize](./config/default).  By default DNS and IDCO Provider microservices are deployed.

---

## 🎯 Usage

Once deployed, interact via the `l2sm-grpc-server`. Reference the gRPC API spec in [`./api/v1/l2sces.proto`](./api/v1/l2sces.proto) and explore client examples in [`./test/client.go`](./test/client.go).

### Configuring Managed Clusters
Prepare the following:
- **Cluster Name**: Unique identifier for each cluster (e.g., `sample-cluster`).
- **API Key**: Kubernetes API endpoint (e.g., `https://api.sample-cluster.local:6443`).
- **Bearer Token**: Token created from Kubernetes (`kubectl create token -n l2sm-system l2sm-controller-manager --context sample-cluster --duration 1000h`).
- **Public Key**: CA Certificate from the cluster, stored as a secret using provided CLI:

```bash
make build
kubectl config view -o jsonpath='{.cluster.certificate-authority-data}' --raw | base64 -d > sample-cluster.key
./bin/apply-cert --namespace l2sm-system --kubeconfig control-plane-kc --clustername sample-cluster sample-cluster.key
```

---

## 📌 Examples

### Creating an Inter-Cluster Slice
Define your slice with IDCO provider details and participating clusters:

```yaml
provider:
  name: test-slice
  domain: "<control plane domain>"
  dns_port: 30053
  sdn_port: 30808
  of_port: 30663
  dns_grpc_port: 30818

clusters:
  - name: "kind-worker-cluster-1"
    apiKey: "https://172.20.0.3:6443"
    bearerToken: "<your-bearer-token>"
    nodes: ["worker-cluster-1-control-plane"]
    gatewayNode:
      name: "worker-cluster-1-control-plane"
      ipAddress: "172.20.0.3"

  - name: "kind-worker-cluster-2"
    apiKey: "https://172.20.0.4:6443"
    bearerToken: "<your-bearer-token>"
    nodes: ["worker-cluster-2-control-plane"]
    gatewayNode:
      name: "worker-cluster-2-control-plane"
      ipAddress: "172.20.0.4"
```

### Creating an Inter-Domain L2Network
Define your L2Network clearly for effective management:

```yaml
name: "l2network-sample"
pod_cidr: "10.1.0.0/16"
type: "vnet"
provider:
  name: test-slice
  domain: "<control plane domain>"
clusters:
  - name: "kind-worker-cluster-1"
    apiKey: "https://172.20.0.3:6443"
    bearerToken: "<your-bearer-token>"
  - name: "kind-worker-cluster-2"
    apiKey: "https://172.20.0.4:6443"
    bearerToken: "<your-bearer-token>"
```

---

For additional support, detailed architectural understanding, and further customization options, refer to our comprehensive [Architecture Guide](https://www.github.com/Networks-it-uc3m/L2S-M).
