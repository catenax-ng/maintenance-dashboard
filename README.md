# Maintenance dashboard app

  - version 1.0

## chart/**

Helm chart to deploy the app (e.g. in k8s-cluster-stack as an argo app or app set)

## src/python/app.py

Script to check versions configured in github repository against the latest configured in projects in newreleases.io
The github queries are specific to catenax-ng/k8s-cluster-stack for now. They handle
  - argo-cd kustomize deployment
  - argocd-vault-plugin init container environment variable AVP_VERSION
  - kubernetes-reflector argocd application set (for core cluster)
  - helm chart versions if no dependencies are defined
  - dependent helm chart versions (each require config file in folder "config/deployment")
Queries for newreleases.io handle all preconfigured projects with one exception.
Kube-prometheus-stack shares the repository "prometheus-community" with other helm charts. 

## config

### github_repo.yaml

Configure the github repository in "<owner>/<repository>" format

### deployment

Configuration files for each component in a yaml format

  - name: the name of the software/helm chart
  - path: path of the kustomize/chart/argo app in the repository
  - project: newreleases.io project name (full name)
  - prefix: anything that precedes the semantic version of the software in newreleases.io
    (i.e. "v", or "kube-prometheus-stack-", "" empty string if none)

examples:

ArgoCD using kustomize, prefix: "v"

```yaml
name: "argo-cd"
path: "apps/argocd/base/kustomization.yaml"
project: "argoproj/argo-cd"
prefix: "v"
```

Cert-manager using helm chart, no prefix ("")

```yaml
name: "cert-manager"
path: "apps/certmanager/Chart.yaml"
project: "helm/cert-manager/cert-manager"
prefix: ""
```

Kube-prometheus-stack using helm chart, repository contains multiple helm charts, and they use prefixes
(i.e. "kube-prometheus-stack-") 

```yaml
name: "kube-prometheus-stack"
path: "apps/kube-prometheus-stack/Chart.yaml"
project: "prometheus-community/helm-charts"
prefix: "kube-prometheus-stack-"
```

## Docker

```bash
cd src/python

# Build
docker build -t ghcr.io/catenax-ng/maintenance-dashboard/maintenance-dashboard-app .

# Push
docker push ghcr.io/catenax-ng/maintenance-dashboard/maintenance-dashboard-app

# Run
docker run -u 1000:1000 -p 5000:5000 -d -e NEWRELEASES_API_KEY=$NEWRELEASES_API_KEY -e GITHUB_TOKEN=$GITHUB_TOKEN ghcr.io/catenax-ng/maintenance-dashboard/aintenance-dashboard-app
```
