# Multistage Builds

This example demonstrates a secure container build using a multistage Dockerfile.

## Security benefits

- Builds the app binary in a dedicated build stage, keeping build tools and source code out of the final image.
- Ships only the compiled artifact in the runtime image, reducing the container attack surface.
- Uses a minimal runtime base image (distroless or similar), which removes shells, package managers, and other unnecessary utilities.
- Limits runtime environment size and privileges, making the container harder to compromise.

## Why this matters

Multistage builds are a best practice for container hardening because they separate the build environment from the runtime environment. The final image contains only what is needed to run the application, which helps reduce vulnerability exposure and improves overall container security.
