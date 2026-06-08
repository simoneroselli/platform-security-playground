# Pod Security Standards (Enforcing the Non-Root Constraint)

Kubernetes Pod Security Standards (PSS) define three distinct isolation profiles to control the security posture of your workloads. These are built directly into the control plane and enforced natively via namespace labels.

## Pod Security Profiles

* **Privileged:** *Unrestricted.* Allows for known privilege escalations, host namespaces, and raw root access. Typically reserved only for system-level infrastructure components (like CNIs or storage plugins).
* **Baseline:** *Default Least-Privilege.* Prevents known privilege escalations but permits standard container defaults. It blocks host access but does not aggressively restrict user IDs or capabilities.
* **Restricted:** *Hardened Production.* Enforces strict hardening best practices. It explicitly requires containers to run as non-root, drops all default Linux capabilities except those strictly required, and forbids volume types that mount host paths.

---

## Namespace Configuration

To enforce the **Restricted** standard on a namespace, you apply the `pod-security.kubernetes.io` control labels. 

The configuration below ensures that any pod attempting to run as root, or failing other hardened criteria, will be completely blocked from entering the `vote-apps` cluster workspace.

### Example Manifest (`namespace.yaml`)

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: vote-apps
  labels:
    # 🚫 ENFORCE: Block any pods that do not comply with the 'Restricted' profile
    pod-security.kubernetes.io/enforce: restricted
    pod-security.kubernetes.io/enforce-version: latest

    # ⚠️ WARN: Send a user-facing terminal warning during deployment if things look wrong
    pod-security.kubernetes.io/warn: restricted
    pod-security.kubernetes.io/warn-version: latest
```

### Deployment & Verification

1. Apply the Namespace Policy
Create and apply the namespace manifest to initialize your hardened logical perimeter:

```bash
kubectl apply -f namespace.yaml
```

2. Verify Enforcement Posture
You can verify that the labels are correctly active on your namespace by inspecting its configuration metadata:

```bash
kubectl get ns vote-apps --show-labels
```