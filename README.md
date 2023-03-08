# Maintenance dashboard app

The purpose of this application is to gather version information about the applications deployed to a kubernetes cluster and then compare it to the latest available versions received from [NewReleases.io](newreleases.io). The combined result is exposed as a prometheus metric endpoint.
