# MPA-K8s (Multidimensional Pod Autoscaler)

MPA-K8s is a Kubernetes controller that provides **multidimensional autoscaling**
by managing application scaling through a **single control loop**.

It is inspired by Google’s Multidimensional Pod Autoscaler (MPA), but implemented
as a **Kubernetes-native controller** using Custom Resources.

> 🚧 **Status**: MVP / Developer Preview (v0.1)

---

## Why MPA?

Kubernetes currently provides:

* **HPA** – Horizontal scaling only (replicas)
* **VPA** – Vertical scaling only (CPU / memory requests)

MPA introduces a **single autoscaling brain** that can:

* Decide *when to scale vertically*
* Decide *when to scale horizontally*
* Avoid conflicts between multiple autoscalers

---

## Features (v0.1)

* Custom Resource Definition: `MultidimensionalPodAutoscaler`
* Kubernetes controller-based architecture
* Horizontal scaling (replicas) management
* No HPA or VPA required
* Kubernetes-native (no admission webhook)

---

## Architecture Overview

```
User YAML (MPA)
      |
      v
+--------------------------+
|   MPA Controller Pod     |
|--------------------------|
| - Watches MPA resources  |
| - Patches Deployments    |
+--------------------------+
```

MPA is installed **once per cluster**, then applied **per workload**.

---

## Installation (Users)

> This section is for **platform teams / cluster administrators**.

### Prerequisites

* Kubernetes cluster
* `kubectl` configured (`KUBECONFIG` set)

### Install MPA into the cluster

Install MPA **once per cluster**:

```sh
kubectl apply -k github.com/nulldoot2k/mpa-k8s/config/default?ref=v0.1.0
```

This will:

* Install the MPA CRD
* Deploy the MPA controller
* Configure required RBAC

### Verify installation

```sh
kubectl get crd | grep multidimensional
kubectl get pods -n mpa-k8s-system
```

You should see the MPA controller pod in `Running` state.

---

## Usage (Users)

Apply a `MultidimensionalPodAutoscaler` resource to your workload.

### Example: `mpa.yaml`

```yaml
apiVersion: autoscaling.hacker-mpa.io/v1alpha1
kind: MultidimensionalPodAutoscaler
metadata:
  name: example-mpa
  namespace: default
spec:
  targetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: example-app

  horizontal:
    minReplicas: 2
    maxReplicas: 5
    targetCPU: 60
```

Apply it:

```sh
kubectl apply -f mpa.yaml
```

From this point, the MPA controller will manage the target Deployment automatically.

---

## How it works (High-level)

1. MPA controller watches `MultidimensionalPodAutoscaler` resources
2. For each MPA:

   * Reads the target workload (Deployment)
   * Evaluates scaling rules
   * Patches the Deployment if scaling is required
3. Status is updated on the MPA resource

---

## Development (Contributors)

> This section is for developers who want to modify or extend MPA.

### Prerequisites

* Go **1.21.x**
* kubebuilder **v3.14.x**
* Docker

### Clone the repository

```sh
git clone https://github.com/nulldoot2k/mpa-k8s.git
cd mpa-k8s
```

### Development workflow

If you modify CRD types or RBAC annotations:

```sh
make generate
make manifests
```

Run unit tests (no Kubernetes cluster required):

```sh
go test ./api/... ./cmd/... ./test/utils/...
```

> Note: Controller integration tests using envtest are optional and not required
> for MVP development.

## Install Minikube

```
curl -Lo minikube https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
chmod +x ./minikube
sudo mv ./minikube /usr/local/bin/

minikube config set driver docker
minikube start
kubectl config use-context minikube
```

## Install Kubectl

```
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
chmod +x ./kubectl
sudo mv ./kubectl /usr/local/bin/
```

BUILD IMAGE MPA CHO MINIKUBE
```
make docker-build IMG=vumanhdat2k/mpa-k8s:v0.1.0
```

PUSH IMAGE MPA CHO MINIKUBE

```
make docker-push IMG=vumanhdat2k/mpa-k8s:v0.1.0
```

Modify config/manager/kustomization.yaml
```
resources:
- manager.yaml

images:
- name: controller
  newName: vumanhdat2k/mpa-k8s
  newTag: v0.1.0
```

Apply manifest
```
kubectl apply -f config/default
```

---

## Roadmap

* CPU-based autoscaling using metrics-server
* Vertical scaling (CPU / memory requests)
* Cooldown and anti-thrashing logic
* Decision engine (vertical-first vs horizontal-first)
* Helm chart support

---

## License

Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at:

[http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0)
