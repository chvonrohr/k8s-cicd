#!/bin/bash

cat > sync-start-excludes.txt << EOF
.dockerignore
.gitlab-ci.yml
cleanup.sh
run-local-kubernetes.sh
build/ci/*.yml
build/package/*.Dockerfile
deployments/*.yaml
deployments/*/*.yaml
EOF
rsync -av --progress --dry-run \
    --exclude-from=sync-start-excludes.txt
    project-solution project-start

rm sync-start-excludes.txt