> skip

## Versions - Kustomize

project-start/deployments/kustomization
```yaml
resources:
  - deployments/database/service.yaml
  - deployments/database/deployment.yaml
  - deployments/database/pvc.yaml
  - deployments/frontend/deployment.yaml
  - deployments/frontend/service.yaml
  - deployments/ingress.yaml
  - deployments/crawler/deployment.yaml
  - deployments/crawler/pvc.yaml
  - deployments/scheduler/cronjob.yaml
  - deployments/backend/deployment.yaml
  - deployments/backend/service.yaml
  - deployments/namespace.yaml
  - deployments/queue/service.yaml
  - deployments/queue/deployment.yaml
```
