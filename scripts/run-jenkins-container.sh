#!/usr/bin/env bash

set -ou pipefail

#  --name jenkins \
# --mount type=bind,source=$(pwd)/data,target=/var/jenkins_home \
podman run \
  -d \
  --privileged \
  --user 1000:1000 \
  -p 8080:8080 -p 50000:50000 \
  -e JENKINS_ADMIN_ID=ngood \
  -e JENKINS_ADMIN_PASSWORD=p4ssw0rd001 \
  --mount type=bind,source=/var/data/jenkins,target=/var/jenkins_home \
  localhost/atkci
