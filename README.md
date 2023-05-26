# Maintenance dashboard app

The purpose of this application is to gather version information about the applications deployed to a kubernetes cluster and then compare it to the latest available versions received from [NewReleases.io](newreleases.io). The combined result is exposed as a prometheus metric endpoint.

## How it works

1. The application is collecting information from the cluster the following way:
   - Gets the `kubelet` version from each node
   - Searches for every service, deployment, statefulset, daemonset in every namespace that is labeled with `maintenance/scan=true`. The version is extracted from the value of the `app.kubernetes.io/version` label. The version format has to follow [SemVer](https://semver.org/) rules. It also reads the value of the `maintenance/releasename` annotation to match the name with the one in [NewReleases.io](newreleases.io).
2. [NewReleases.io](newreleases.io) is queried and the latest versinos are pulled. The application then finds the latest non-prerelease version.
3. The results are exposed at the `/metrics` endpoint as the application is listening on the `:2112` port within the cluster.

## Installation

See in the [INSTALL.md](INSTALL.md).

## Configure apps to be scanned

``` sh

# Example for a service resource.
# 1. Select a k8s resource that has a kind of `service`, `deployment`, `statefulset` or `daemonset`.

# 2. Add the following label to this resource: `maintenance/scan=true`.
kubectl label svc -n [NAMESPACE] [SERVICE_NAME] maintenance/scan=true

# 3. Check if this recommended label exists on the resource: `app.kubernetes.io/version`.
kubectl get svc -n [NAMESPACE] [SERVICE_NAME] -o jsonpath="{.metadata.labels.app\.kubernetes\.io/version}"

# 4. If not, add it with the proper app version in `semver` format.
kubectl label svc -n [NAMESPACE] [SERVICE_NAME] app.kubernetes.io/version=[SEMVER_VERSION]

# 5. Annotate resource with the key `maintenance/releasename` and the value from the name of the project found on NewReleases.io.
kubectl annotate svc -n [NAMESPACE] [SERVICE_NAME] maintenance/releasename=[NEWRELEASES_PROJECT_NAME]

```

When these labels and annotations are set, the maintenance dashboard will pick up and serve metrics about the application.
