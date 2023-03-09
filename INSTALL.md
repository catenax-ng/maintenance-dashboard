# Installation

The Maintenance Dashboard repository can be found [here](https://catenax-ng.github.io/maintenance-dashboard/index.yaml)

Use the latest release of the "maintenance-dashboard" chart.

Supply the required configuration properties (see chapter [Configuration](#configuration)) in a values.yaml file or
override the settings directly.

### Deployment using Helm

Add the Maintenance Dashboard repository:

```(shell)
    helm repo add maintenance-dashboard https://catenax-ng.github.io/maintenance-dashboard/index.yaml
```

Then install the Helm chart into your cluster:

```(shell)
    helm install -f your-values.yaml maintenance-dashboard maintenance-dashboard/maintenance-dashboard
```

Or create a new Helm chart and use the Maintenance Dashboard as a dependency.

```(yaml)
    dependencies:
      - name: maintenance-dashboard
        repository: https://catenax-ng.github.io/maintenance-dashboard
        version: 1.0.x
```

Then provide your configuration as the values.yaml of that chart.

Create a new application in ArgoCD and point it to your repository / Helm chart folder.

## Configuration

All the configurable variables are available in the [values.yaml](charts/maintenance-dashboard/values.yaml) file and comments can be found for each of them on what they do. Some remarks that are specific to this application:

- `newReleasesApiKey` is a token to access the [NewReleases.io](newreleases.io) account where the repositories are followed.
- `inCluster` specifies if the resources that are scanned for maintenance are inside the cluster where the Maintenance Dashboard is running.
- `kubeConfig` if the resources that are scanned for maintenance are outside of the cluster where the Maintenance Dashboard is running then a configuration needs to be provided to access the Kubernetes cluster.
- `refreshIntervalSeconds` determines how often the scan is running.
