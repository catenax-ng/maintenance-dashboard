package data

import "github.com/Masterminds/semver/v3"

type AppVersionInfo struct {
	NewReleasesName    string          // the name of the package on NewReleases.io
	CurrentVersion     *semver.Version // the version of the app currently deployed
	LatestMajorVersion *semver.Version // compared to CurrentVersion
	LatestMinorVersion *semver.Version // compared to CurrentVersion
	LatestPatchVersion *semver.Version // compared to CurrentVersion
}
