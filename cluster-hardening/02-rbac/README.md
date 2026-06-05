# Namespaces & RBAC Boundaries

This directory implements logical workload separation and identity isolation for our core application layers. We configure a dedicated namespace and explicit, minimum-privilege identities to mitigate the blast radius of container-level compromises.

---

## 🔒 Hardening Strategy Overview

We are locking down the logical perimeters for two components:
1. **`vote-frontend`** (Public-facing web routing)
2. **`redis`** (Internal data store)

By ensuring these applications run inside distinct namespace boundaries and use isolated credentials, we prevent horizontal privilege escalation.

---

## 🏗️ Architecture & Configuration Components

### 1️⃣ Dedicated Workspace Isolation (`01-namespace.yaml`)
* **The Rule:** Never run production workloads in the `default` namespace. 
* **The Reason:** The `default` namespace contains standard configurations, making it a soft target. Shared namespaces lack administrative boundaries, allowing compromised applications to describe or manipulate neighboring resources. Move workloads to a distinct logical tenant.

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: vote-app
```

### 2️⃣ Explicit Application Identities (`02-serviceaccounts.yaml`)
* **The Rule:** Give every microservice its own identity, and explicitly declare `automountServiceAccountToken: false`.
* **The Reason:** By default, Kubernetes mounts a secret API token volume inside every pod. If an attacker compromises the web container, they can steal this token to talk to the cluster control plane. Since these services only need network connectivity—not API access—disabling token mounting strips the container of any ambient cluster credentials.

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: vote-frontend-sa
  namespace: vote-app
automountServiceAccountToken: false # 💡 Staged for production isolation
```
## Workload Identification Binding (03-deployments.yaml)

* **The Rule:** Explicitly link your custom identity within the pod specification using serviceAccountName.

* **The Reason:** Forcing an exact identity mapping ensures the pod runs under the designated vote-frontend-sa context instead of implicitly falling back to the shared namespace default account.

```bash
apiVersion: apps/v1
kind: Deployment
metadata:
  name: vote-frontend
  namespace: vote-app
spec:
  template:
    spec:
      serviceAccountName: vote-frontend-sa # 👈 Binding the hardened identity
      containers:
      - name: frontend
        image: dockersamples/visualizer:stable
```