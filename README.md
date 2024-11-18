# Kubernetes Learning Module

This repository contains hands-on exercises and solutions for learning Kubernetes core concepts and components.

## Prerequisites

- [kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation) - For creating local Kubernetes clusters
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) - Kubernetes command-line tool
- Docker installed and running

## Cluster Setup

Create a multi-node cluster using kind:
```bash
kind create cluster --config kind-cluster-config/kind-multi-node-cluster.yaml
```

To delete the cluster:
```bash
kind delete cluster
```

**Tip**: Use `kubectl get pods -w` in a separate terminal to monitor pod state changes in real-time.

## Module Contents

### 1. Kubernetes Components
Learn about core Kubernetes control plane components.

**Tasks:**
1. Examine control plane components:
   ```bash
   # Inspect API Server
   kubectl describe pod kube-apiserver-kind-control-plane -n kube-system
   
   # Inspect Controller Manager
   kubectl describe pod kube-controller-manager-kind-control-plane -n kube-system
   
   # Inspect Scheduler
   kubectl describe pod kube-scheduler-kind-control-plane -n kube-system
   ```

2. Explore controller manager configuration options in the [official documentation](https://kubernetes.io/docs/reference/command-line-tools-reference/kube-controller-manager/)

### 2. Workloads
Practice working with different types of Kubernetes workloads.

**Tasks:**
1. Basic Pod Operations:
   ```bash
   # Create a pod that pings Google DNS
   kubectl run ping-google-pod --image=alpine ping 8.8.8.8
   
   # Clean up
   kubectl delete pod ping-google-pod
   ```

2. Init Containers:
   ```bash
   kubectl apply -f workloads/init-container.yaml
   kubectl logs static-file-app-pod
   
   # Experiment with container lifecycle
   kubectl exec static-file-app-pod -c static-file-app-container -- rm /tmp/intro.txt
   ```
   
   **Note**: The container will attempt to restart after deleting `/tmp/intro.txt` but fail because init containers only run on pod start/restart.

3. DaemonSets:
   ```bash
   kubectl apply -f workloads/nginx-daemonset.yaml
   ```

4. DaemonSet Rolling Updates:
   ```bash
   # Apply DaemonSet
   kubectl apply -f workloads/nginx-daemonset.yaml
   
   # Trigger failed update (for learning purposes)
   # Change image tag to 'latest-not-available' in yaml file
   kubectl apply -f workloads/nginx-daemonset.yaml
   
   # Monitor rollout
   kubectl rollout status daemonset.apps/nginx-app
   
   # Rollback
   kubectl rollout undo daemonset.apps/nginx-app
   ```

5. StatefulSets:
   ```bash
   kubectl apply -f workloads/postgresql-statefulset.yaml
   ```

6. Scheduler Modification:
   ```bash
   # Copy modified scheduler config
   docker cp workloads/kube-scheduler-modified.yaml kind-control-plane:/etc/kubernetes/manifests/
   
   # Apply changes
   docker exec -it kind-control-plane /bin/bash
   cd /etc/kubernetes/manifests/
   mv kube-scheduler-modified.yaml kube-scheduler.yaml
   ```

### 3. Kubernetes Objects
Learn about Kubernetes object management.

**Task:**
```bash
# Deploy with replicas
kubectl apply -f k8s-objects/nginx-deployment-replicas.yaml

# Experiment with label changes
# Modify pod template labels in nginx-deployment-replicas.yaml
```

### 4. Networking
Explore Kubernetes networking concepts.

**Tasks:**
1. Service Discovery within namespace:
   ```bash
   kubectl apply -f networking/nginx-deployment.yaml
   kubectl run test-pod-ns-one -n ns-one --image=busybox:latest ping 8.8.8.8
   
   # Get service endpoint
   kubectl describe service -n ns-one nginx-service
   
   # Test connectivity
   kubectl exec -it -n ns-one test-pod-ns-one -- /bin/sh
   wget http://<endpoint>
   ```

2. Cross-namespace communication
3. Network Policies:
   ```bash
   # Create test deployment
   kubectl create deployment nginx --image=nginx
   kubectl expose deployment nginx --port=80
   
   # Test connectivity before policy
   kubectl run busybox --rm -ti --image=busybox -- /bin/sh
   wget --spider --timeout=1 nginx
   
   # Apply network policy
   kubectl apply -f networking/nginx-policy.yaml
   
   # Test with and without required labels
   kubectl run busybox --rm -ti --labels="access=true" --image=busybox -- /bin/sh
   ```

### 5. Configuration
Work with ConfigMaps and Secrets.

**Tasks:**
1. ConfigMaps:
   ```bash
   kubectl apply -f configuration/configmap-pod.yaml
   ```

2. Secrets:
   ```bash
   kubectl apply -f configuration/secret-pod.yaml
   ```

3. Secret Creation Methods:
   ```bash
   # From file
   echo -n '1f2d1e2e67df' > configuration/password.txt
   kubectl create secret generic pass-example --from-file=password=configuration/password.txt
   
   # From literal
   kubectl create secret generic pass-example-literal --from-literal=password='1f2d1e2e67df'
   ```

### 6. Storage
Learn about Kubernetes storage options.

**Tasks:**
1. Study [Dynamic Provisioning](https://kubernetes.io/docs/concepts/storage/dynamic-provisioning/)
2. Shared Volumes between containers:
   ```bash
   kubectl apply -f storage/pod-with-two-containers.yaml
   ```
3. HostPath volumes:
   ```bash
   kubectl apply -f storage/hostPath-example.yaml
   ```

### 7. Security
Understand Kubernetes security concepts.

**Tasks:**
1. RBAC:
   ```bash
   kubectl apply -f security/kanister-pod.yaml
   kubectl create rolebinding pod-reader-pod --role=pod-reader --serviceaccount=default:default
   ```
2. Authorization Modes
3. Service Accounts

### 8. Final Project: Library Application

Deploy a complete library application with the following components:

**Prerequisites:**
- Docker Hub account (if building custom image)
- Pre-built image available: `shlokc/library_application:1.0.8`

**Deployment Steps:**
```bash
# Optional: Build and push custom image
cd ./library-app
docker login
docker build -f Dockerfile . -t <image-name:tag>
docker push <image-name:tag>

# Deploy application
kubectl apply -f k8s-manifests/mysql-database.yaml

# Deploy API server
kubectl run lib-api-server -n library-app \
  --env="MYSQL_SERVER_IP=mysql-service" \
  --env="MYSQL_SERVER_PORT=3306" \
  --env="MYSQL_SERVER_USER=root" \
  --env="MYSQL_SERVER_PASSWORD=password" \
  --port=8080 \
  --labels="app=library" \
  --image=<image-name:tag> \
  --command -- go run main.go

# Deploy service
kubectl apply -f k8s-manifests/library-service.yaml
```

