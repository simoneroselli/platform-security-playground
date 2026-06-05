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

## Restricting Pod Access using Kubernetes RBAC (Least Privilege Example)

To configure the `vote-frontend-sa` service account so that it can **ONLY** view/list a specific target pod (such as a `redis` pod) and is denied access to everything else, you need to implement a strict `Role` and `RoleBinding` that limits permissions to that specific resource name within the namespace.

### 1. Define the Manifests

Create a file named `rbac-vote-apps.yaml` with the following configuration. Notice the use of `resourceNames` to restrict permissions to the exact pod instance.

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  namespace: vote-apps
  name: vote-frontend-redis-reader
rules:
- apiGroups: [""]
  resources: ["pods"]
  verbs: ["get", "list"]
  resourceNames: ["redis"] # Strict limitation to only the redis pod
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  namespace: vote-apps
  name: vote-frontend-redis-read-binding
subjects:
- kind: ServiceAccount
  name: vote-frontend-sa
  namespace: vote-apps
roleRef:
  kind: Role
  name: vote-frontend-redis-reader
  apiGroup: rbac.authorization.k8s.io
```

## 2. Apply the Configuration

Deploy the RBAC resources to your cluster via kubectl:
```bash
kubectl apply -f rbac-vote-apps.yaml
```

## 3. Verify and Audit Permissions

You can confirm that your logical boundaries and the principle of least privilege are functioning as expected by running cross-boundary checks via kubectl auth can-i:

### Scenario A: vote-frontend CAN see the redis pod

```bash
kubectl auth can-i get pod/redis --as=system:serviceaccount:vote-apps:vote-frontend-sa -n vote-apps
```
```plaintext
yes
```

### Scenario B: redis CANNOT see the vote-frontend pod
(Ensures that if the Redis container is compromised, the attacker cannot scrape or discover your frontend infrastructure)

```bash
kubectl auth can-i get pod/vote-frontend --as=system:serviceaccount:vote-apps:redis-sa -n vote-apps
```
```plaintext
no
```

### Scenario C: vote-frontend CANNOT see cluster-wide or blanket pods
(Verifies that wildcard namespace-wide queries targeting any other unapproved pods are dropped by the API server)

```bash
kubectl auth can-i get pods --as=system:serviceaccount:vote-apps:vote-frontend-sa -n vote-apps
```
```plaintext
no
```