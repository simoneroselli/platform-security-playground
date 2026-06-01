#!/bin/bash

# Install k3d cluster
k3d cluster create sandbox \
  --servers 2 \
  --agents 2 \
  --port "8080:80@loadbalancer" \
  --api-port 6550


