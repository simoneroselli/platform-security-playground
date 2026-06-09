# Platform Security Playground 🧪

Welcome to this personal playground and research repository. This space serves as a hands-on notebook where I collect deep-dive analyses, practical experiments, and notes regarding various security scenarios encountered in the cloud-native ecosystem.

## Repository Structure

The repo is organized into focused, self-contained lab directories. Each directory contains the exact code, configuration files, and a dedicated write-up detailing my findings, gotchas, and implementation notes.

### Explored Scenarios

| Category | Deep-Dive Scenario | Key Findings & Quirks Covered |
| :--- | :--- | :--- |
| **Container Hardening** | [Multistage Builds](./container-hardening/01-multistage-builds/README.md) | Eliminating build-tool attack surface and minimizing production image footprints. |
| **Container Hardening** | [Non-Root User Execution](./container-hardening/02-non-root-user/README.md) | Enforcing least privilege, handling macOS UID mapping, and bypassing Distroless metadata quirks for K8s Admission Controllers. |
| **Container Hardening** | [Prevent Container Privilege Escalation](./container-hardening/03-prevent-privilege-escalation/README.md) | Blocking runtime privilege escalation via SecurityContext controls to prevent capability elevation. |
| **Container Hardening** | [Limit Container Resources](./container-hardening/04-limit-container-resources/README.md) | Setting CPU and memory resource requests and limits to prevent resource exhaustion and DoS attacks. |
| **Container Hardening** | [Automate Vulnerability Scanning](./container-hardening/05-automate-vulnerability-scanning/README.md) | Automating container image scanning to detect vulnerabilities earlier in the pipeline. |
| **Cluster Hardening** | [Zero-Trust Network Policies](./cluster-hardening/01-network-policies/README.md) | Enforcing namespace-wide microsegmentation and least-privilege egress limits to prevent lateral movement. |
| **Cluster Hardening** | [Namespaces & RBAC Boundaries](./cluster-hardening/02-namespaces-rbac/README.md) | Establishing isolated workload environments and stripping ambient API tokens to enforce absolute control-plane least privilege. |
| **Cluster Hardening** | [Pod Security Standards](./cluster-hardening/03-pod-security-standard/README.md) | Enforcing native Admission Control via Restricted profiles to block root containers and strip dangerous Linux capabilities. |
| **Cluster Hardening** | [External Secrets Operator & LocalStack](./cluster-hardening/04-secret-management-isolation/README.md) | Bridging cloud APIs into local clusters by setting up a 15s synchronization loop between ESO and a mocked AWS Secrets Manager, fixing SDK token signature validation loops via custom Helm overrides. |

---

## Local Setup & Prerequisites

These are the primary tools I use locally to spin up and test these environments:

- **Docker Desktop** (Container runtime & local VM environment)
- **Go** (For compiling minimal test binaries)
- **k3d** (For lightweight, multi-node Kubernetes cluster simulation)
- **Terraform** (For automating infrastructure provisioning and IaC security states)

## How to Navigate

Each subfolder contains its own local `README.md` which acts as my lab log—detailing the risk, the specific configuration changes required to fix it, and any unexpected platform behaviors (like Docker Desktop or Kubernetes engine constraints) discovered along the way.