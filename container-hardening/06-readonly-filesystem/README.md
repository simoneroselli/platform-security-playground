# Read Only root Filesystem

Setting up a read-only root filesystem is one of the most effective container hardening practices in Docker security. It implements the principle of immutability at runtime: even if an attacker manages to exploit an application vulnerability (like Remote Code Execution) and gain shell access, they cannot download malicious binaries, modify application code, or alter system configurations because the underlying filesystem rejects write operations.

## Read-Only + tmpfs

When you lock down the root filesystem, you must selectively provide temporary, memory-backed storage (tmpfs) for paths where the application expects to write data dynamically. Unlike regular volume mounts, tmpfs mounts write directly to the host's system memory (RAM) rather than the host's persistent disk, ensuring that any temporary files disappear completely the moment the container stops.

### Native Docker CLI

```bash
docker run -d \
  --read-only \
  --tmpfs /tmp \
  --tmpfs /var/cache/nginx \
  -p 8080:80 \
  nginx:alpine
```

### Docker Compose

```Dockerfile
version: '3.8'

services:
  secure-api:
    image: my-app-image:latest
    read_only: true
    tmpfs:
      - /tmp:size=64M,uid=1000  # Restrict size and map ownership if necessary
      - /app/logs
    ports:
      - "8080:8080"
    security_opt:
      - no-new-privileges:true
```

### Kubernetes (CKS Blueprint)

The Docker configuration translates directly into the pod's securityContext.
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: hardened-app
spec:
  replicas: 1
  template:
    spec:
      containers:
      - name: web-app
        image: alpine
        command: ["sh", "-c", "echo 'Running secure' && sleep 3600"]
        securityContext:
          readOnlyRootFilesystem: true
          runAsNonRoot: true
          runAsUser: 1000
        volumeMounts:
        - mountPath: /tmp
          name: cache-volume
      volumes:
      - name: cache-volume
        emptyDir:
          medium: Memory
          sizeLimit: "64Mi"
```

## So far so good! ..but, what prevent me to download malicious files directly into tmpfs and execute them?

## Mounting tmpfs with noexec

The absolute best defense against this specific vector is restricting how files inside the tmpfs volume can behave. When you mount a filesystem in Linux, you can pass mount flags. The most critical flag here is noexec.

### Docker Compose

```bash
...
services:
  secure-app:
    image: my-app:latest
    read_only: true
    tmpfs:
      - /tmp:noexec,nosuid,nodev,size=64M
```

### Kubernetes: Pod SecurityContext (I did not test it!)
```yaml
...
        securityContext:
          readOnlyRootFilesystem: true
        volumeMounts:
        - mountPath: /tmp
          name: tmp-noexec
      volumes:
      - name: tmp-noexec
        ephemeral:
          volumeClaimTemplate:
            spec:
              accessModes: [ "ReadWriteOnce" ]
              # Pass your strict Linux mount security parameters here
              mountOptions:
                - noexec
                - nosuid
                - nodev
              resources:
                requests:
                  storage: 64Mi
```