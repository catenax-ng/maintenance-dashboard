# Maintenance dashboard app

  - version 1.0

TODO:

  - create actions pipeline for building the image?

## maintenance-dashboard-app.py

The script that is checking the deployed versions in config.json agains the latest in newreleases.io

## Dockerfile

Build a docker image from maintenance-dashboard-app.py

## chart/**

A Helm chart to deploy the app (e.g. in k8s-cluster-stack as an argo app or app set)

## config.json

Mapping between github helm charts/kustomizations/deployments and newreleases.io project in json format

```json
{
  "<github repo>":"<owner/repo>"
  "<newreleases.io api url>":"<value>"
  "<charts>":[
    {
      "<app>":"<app name>"
      "<path>":"<path in github repo>"
      "<deplendency>":"<index of dependencies array(starts at 0)>"
      "<project>":"<newreleases.io project name>"
      "<prefix>":"<prefix of version in newreleases.io>"
    }
  ]
  "<kustomizes>":[{"<same structure>": "<as charts>"}]
  "<deployments>":[{"<same structure>": "<as charts>"}]
}
```

## Docker

```bash
# Build
docker build -t ghcr.io/catenax-ng/maintenance-dashboard/maintenance-dashboard-app .

# Push
docker push ghcr.io/catenax-ng/maintenance-dashboard/maintenance-dashboard-app

# Run
docker run -u 1000:1000 -p 8000:8000 -d -e NR_API_KEY=$NR_API_KEY mghcr.io/catenax-ng/maintenance-dashboard/aintenance-dashboard-app
```
