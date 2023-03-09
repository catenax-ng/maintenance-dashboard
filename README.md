# Maintenance dashboard app

The purpose of this application is to gather version information about the applications deployed to a kubernetes cluster and then compare it to the latest available versions received from [NewReleases.io](newreleases.io). The combined result is exposed as a prometheus metric endpoint.

## How it works

1. The application is collecting information from the cluster the following way:
   - Gets the `kubelet` version from each node
   - Searches for every service in every namespace that is labeled with `maintenance/scan=true`. The version is extracted from the value of the `app.kubernetes.io/version` label. The version format has to follow [SemVer](https://semver.org/) rules. It also reads the value of the `maintenance/releasename` annotation to match the name with the one in [NewReleases.io](newreleases.io).
2. [NewReleases.io](newreleases.io) is queried and the latest versinos are pulled. The application then calculates 3 versions compared to the current version:
     - Latest major version
     - Latest minor version
     - Latest patch version
3. The results are exposed at the `/metrics` endpoint as the application is listening on the `:2112` port within the cluster.

## Installation

See in the [INSTALL.md](INSTALL.md).

## TODOs

- Improve label and annotation process for components to scan in the cluster