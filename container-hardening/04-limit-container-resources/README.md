# Limit Container Resources (Prevent DoS Attacks)

This example demonstrates how to restrict the memory and CPU consumption of a container. Without these boundaries, a compromised application or a targeted Denial of Service (DoS) attack can trigger a resource exhaustion loop, starving adjacent containers and crashing the entire underlying host node.

---

## Problem description

By default, containers share the host machine's hardware pool with zero restrictions. If an attacker manages to exploit your application (e.g., via a ReDoS regular expression bypass, a memory leak exploit, or a cryptomining script), the container process will aggressively consume CPU threads and RAM capacity.

On a shared infrastructure like a Kubernetes node, this causes "Noisy Neighbor" syndrome. The kernel's **OOM (Out Of Memory) Killer** will step in to save the operating system. Because your container didn't declare bounds, the OS might accidentally kill critical infrastructure processes or completely lock up the worker node, leading to a cluster-wide outage.

### Docker CLI
To enforce instant boundaries on a single running container container, pass the strict hardware restraint flags explicitly:

```bash
docker run -d --name secure-cli-demo --memory="512m" --cpus="0.5" my-app
```

### Docker Compose
```bash
# Hardware boundaries enforced at the container level
    deploy:
      resources:
        limits:
          cpus: '0.50'
          memory: 512M
        reservations:
          cpus: '0.25'
          memory: 256M
```

### Kubernetes deployment.yaml
```bash
resources:
          # Guaranteed minimums allocated to the container
          requests:
            memory: "256Mi"
            cpu: "250m"
          # Absolute ceilings the container cannot exceed
          limits:
            memory: "512Mi"
            cpu: "500m"
```