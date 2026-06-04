# Prevent Container Privilege Escalation

- [Prevent Container Privilege Escalation](#prevent-container-privilege-escalation)
  - [Problem description](#problem-description)
  - [distroless solve this by design](#distroless-solve-this-by-design)
  - [Secure runtime!](#secure-runtime)
    - [Mitigation](#mitigation)
      - [docker](#docker)
      - [docker-compose](#docker-compose)
      - [kubernetes deployment](#kubernetes-deployment)
  - [References](#references)


This example demonstrates how to prevent an unprivileged container process from dynamically elevating its privileges to `root` at runtime. Even when running a container as a non-root user, certain application vulnerabilities or system binaries can allow an attacker to hijack the process execution thread and gain full root control over the container environment.

---

## Problem description

Imagine you successfully locked down your container to run as an unprivileged user (e.g., standard user `10001`). If an attacker exploits a remote code execution vulnerability inside your application, they land inside your container inheriting that unprivileged user account.

However, the attacker isn't stuck. They can search the filesystem for **SUID (Set User ID)** binaries—special executable files that always run with the permissions of the file's owner (typically `root`) rather than the person running them. Alternatively, the attacker can download a pre-compiled exploit or wrapper script from the internet, flip its SUID flags, and execute it. 

If successful, the attacker's process hooks into the operating system and dynamically elevates itself back up to `UID 0` (root). From there, they can modify system binaries, tap host networks, or initiate a container breakout.

---

## distroless solve this by design

Google's Distroless images mitigate a significant portion of this risk out of the box because **they are stripped of all traditional Linux operating system binaries**. 

There is no `sudo`, no `su`, no package managers, and no standard shell-level utilities that have SUID flags enabled. Because the attack surface is completely empty, an attacker lacks the baseline internal operating system tools required to easily map a path toward privilege elevation.

---

## Secure runtime!

If you are **not** using a Distroless image (e.g., you are building on standard Alpine, Ubuntu, or Debian images because your application requires specific OS packages), relying solely on file elimination isn't enough. If an attacker downloads an exploit from the internet, a static file-stripping step during your Docker build phase will fail to stop them.

The ultimate defense is **Runtime Enforcement via the Linux Kernel**. 

By passing explicit security configuration flags, we tell the kernel to set the `PR_SET_NO_NEW_PRIVS` bit on the container's execution thread. Once this kernel bit is enabled:
1. It is permanently inherited by the container application and any child processes it ever spawns.
2. The Linux kernel completely ignores SUID and SGID execution bits on **all** files—including any malicious scripts or binaries an attacker attempts to download and execute dynamically.

The application is permanently locked into its unprivileged execution bounds.

### Mitigation
#### docker
```bash
docker run -ti --security-opt=no-new-privileges:true my_image#
```

#### docker-compose
```bash
...
    security_opt:
      - no-new-privileges:true
```


#### kubernetes deployment
```bash
allowPrivilegeEscalation: false
```

## References
[PR_SET_NO_NEW_PRIVS](https://man7.org/linux/man-pages/man2/PR_SET_NO_NEW_PRIVS.2const.html)

[Nice article](https://raesene.github.io/blog/2019/06/01/docker-capabilities-and-no-new-privs/)

[Here is how runc calls `NO_NEW_PRIVS` right before booting your application binary](https://github.com/opencontainers/runc/blob/3047d61ff95fd56addae9deaa44c78242dfa116f/libcontainer/setns_init_linux.go#L68)