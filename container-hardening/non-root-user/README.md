# Non-Root User Execution

This example demonstrates how to configure and enforce non-root user execution within a container. By default, Docker containers run as the `root` user (UID 0). This lab shows how to transition to a restricted, unprivileged user to drastically minimize the impact of a container breakout.

## Security Benefits
- **Host Protection:** Prevents an attacker from automatically gaining root privileges on the underlying host machine if they manage to escape the container.
- **Lateral Movement Mitigation:** Restricts what an attacker can do within the container filesystem (e.g., they cannot modify system binaries, install malicious packages, or bind to privileged ports below 1024).
- **Compliance Ready:** Satisfies strict enterprise compliance standards and CIS benchmarks that explicitly forbid running containers as root in production.

## Why This Matters

If a vulnerability (like a Remote Code Execution) is exploited in your application, the attacker inherits the exact permissions of the process running it. If that process is `root`, they control the container's mini-OS entirely.

### The Distroless & Kubernetes Caveat (The Text Name vs. UID Trap)

When using Google's secure Distroless images, a common mistake is adding `USER nonroot:nonroot` to the Dockerfile. Because Distroless images are stripped bare to reduce the attack surface, they lack traditional Linux user management files (like `/etc/passwd`). 

Without these files:
1. **Docker cannot resolve the text string** `"nonroot"` to a number during build time, causing the image's static metadata config to be left entirely blank (`""`).
2. **Kubernetes admission controllers** (like Pod Security Admissions or Kyverno) will reject the Pod. When a cluster enforces `runAsNonRoot: true`, it scans the image metadata before launching it. If it sees a blank user field, it assumes the image defaults to root and blocks it with a `CreateContainerConfigError`.

By explicitly hardcoding the numeric ID **`USER 65532:65532`**, we completely bypass the local lookup, stamp the security settings directly into the image blueprint, and ensure 100% compatibility with hardened Kubernetes clusters.

---

## How to Run This Lab

### 1. Build the Secure Image
Build the image using the provided Dockerfile which explicitly enforces the numeric UID:
```bash
docker build -t security-nonroot-demo .

docker run --rm -p 9095:8080 security-nonroot-demo

curl http://localhost:9095

```

### Result


Welcome to the Security Playground!<BR>
Current Process Running UID: 65532<BR>
✅ SUCCESS:Running as NON-ROOT user.
