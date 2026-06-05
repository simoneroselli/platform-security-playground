# Network Policies: Zero-Trust Isolation

This directory contains the production-grade Network Policy architecture used to transition our application runtime environment from an open multi-tenant routing scheme to a strict, least-privilege Zero-Trust layout.

## 📊 Traffic Communication Flow

When a network perimeter is locked down, all communications must be explicit and decoupled. The direct connection required for this playground segment is limited to the public-facing frontend component and the transient data cache, backed by the cluster's internal resolution infrastructure.

```text
                     [ CoreDNS ] 
                          ^
                          │ (UDP/TCP:53)
                          │
  ┌───────────────────────┴────────────────────────┐
  │                                                │
  │  ┌─────────────────┐             ┌──────────┐  │
  └──┤  vote-frontend  ├────────────>│  redis   ├──┘
     └─────────────────┘ (TCP:6379)  └──────────┘

```
## Traffic Partitioning Strategy

To keep configurations clean and avoid duplicated logic (bloat) across every application deployment manifest, traffic boundaries are systematically layered in a strict, sequential order:

## 1️⃣ Core Resolution Anchor (`01-global-dns.yaml`)
* **The Logic:** Network plumbing is treated as namespace infrastructure, not individual application logic. Before an application name can be targeted, pods must be permitted to figure out where that target lives.
* **The Blueprint:** A broad egress allowance targets all pods in the local namespace (`podSelector: {}`). It permits traffic to escape exclusively toward the `kube-system` namespace to query the `kube-dns` pods on port 53 (supporting both UDP and TCP). This universally handles local service discovery.

---

## 2️⃣ Application Microsegmentation (`02-vote-frontend.yaml` / `03-redis.yaml`)
* **The Logic:** Security contexts are established matching the absolute minimum operational requirements of the software architecture. Unused cross-pod communication pathways are entirely omitted.
* **The Blueprint:**
  * **Frontend Egress:** An explicit outgoing rule attaches to pods carrying the `app: vote-frontend` label, permitting them to initiate outbound connections targeting `app: redis` on TCP port 6379.
  * **Redis Ingress:** A corresponding inbound firewall rule isolates pods carrying the `app: redis` label, forcing them to drop any incoming packet unless it explicitly originates from a container matching the `app: vote-frontend` label.

---

## 3️⃣ Perimeter Lockdown: Global Fallback (`00-default-deny.yaml`)
* **The Logic:** Production clusters must operate under a default-closed model. Any protocol, pathway, or port not explicitly opened by the preceding rules must be denied by default to fully eliminate potential lateral movement vectors.
* **The Blueprint:** A catch-all wildcard policy targets every pod in the active namespace. It turns off default open-routing and enforces an absolute **Deny-All** stance on both incoming (`Ingress`) and outgoing (`Egress`) traffic.