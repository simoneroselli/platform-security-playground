# Platform Security Playground 🧪

Welcome to this personal playground and research repository. This space serves as a hands-on notebook where I collect deep-dive analyses, practical experiments, and notes regarding various security scenarios encountered in the cloud-native ecosystem.

## Repository Structure

The repo is organized into focused, self-contained lab directories. Each directory contains the exact code, configuration files, and a dedicated write-up detailing my findings, gotchas, and implementation notes.

### Explored Scenarios

| Category | Deep-Dive Scenario | Key Findings & Quirks Covered |
| :--- | :--- | :--- |
| **Container Hardening** | [Multistage Builds](./container-hardening/01-multistage-builds/README.md) | Eliminating build-tool attack surface and minimizing production image footprints. |
| **Container Hardening** | [Non-Root User Execution](./container-hardening/02-non-root-user/README.md) | Enforcing least privilege, handling macOS UID mapping, and bypassing Distroless metadata quirks for K8s Admission Controllers.

---

## Local Setup & Prerequisites

These are the primary tools I use locally to spin up and test these environments:

- **Docker Desktop** (Container runtime & local VM environment)
- **Go** (For compiling minimal test binaries)
- **k3d** (For lightweight, multi-node Kubernetes cluster simulation)
- **Terraform** (For automating infrastructure provisioning and IaC security states)

## How to Navigate

Each subfolder contains its own local `README.md` which acts as my lab log—detailing the risk, the specific configuration changes required to fix it, and any unexpected platform behaviors (like Docker Desktop or Kubernetes engine constraints) discovered along the way.