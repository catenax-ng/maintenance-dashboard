package latestversions

import (
	"context"
	"os"
	"sort"

	log "github.com/sirupsen/logrus"

	"github.com/Masterminds/semver/v3"
	"github.com/catenax-ng/maintenance-dashboard/internal/data"
	"github.com/catenax-ng/maintenance-dashboard/internal/parseversion"
	"newreleases.io/newreleases"
)

type ProjectInfo struct {
	ID          string
	ReleaseName string
}

var (
	NewReleasesApiKey = os.Getenv("NEWRELEASES_API_KEY")
)

// Get stable releases for an app from NewReleases.io, parse them to semver and compare to current version
func GetForApp(appVersionInfo data.AppVersionInfo) *data.AppVersionInfo {
	client := newreleases.NewClient(NewReleasesApiKey, nil)
	ctx := context.Background()
	vs := []*semver.Version{}

	for i := 1; i < 10; i++ {
		releases, lastPage, err := client.Releases.ListByProjectName(ctx, "github", appVersionInfo.NewReleasesName, i)
		if err != nil {
			log.Panic(err)
		}

		for _, release := range releases {
			parsedVersion, err := parseversion.ToSemver(release.Version)
			if err != nil {
				log.Warningf("Skipping invalid version: %v", release.Version)
			} else if parsedVersion.Prerelease() == "" {
				vs = append(vs, parsedVersion)
			}
		}

		if i >= lastPage {
			break
		}
	}

	sort.Sort(sort.Reverse(semver.Collection(vs)))
	latestMajorVersion := vs[0]
	var latestMinorVersion *semver.Version
	var latestPatchVersion *semver.Version

	for i := 0; i < len(vs); i++ {
		if latestMinorVersion == nil && vs[i].Major() == appVersionInfo.CurrentVersion.Major() {
			latestMinorVersion = vs[i]
		}

		if latestPatchVersion == nil && vs[i].Major() == appVersionInfo.CurrentVersion.Major() && vs[i].Minor() == appVersionInfo.CurrentVersion.Minor() {
			latestPatchVersion = vs[i]
		}
	}

	return &data.AppVersionInfo{
		NewReleasesName:    appVersionInfo.NewReleasesName,
		CurrentVersion:     appVersionInfo.CurrentVersion,
		ResourceName:       appVersionInfo.ResourceName,
		LatestMajorVersion: latestMajorVersion,
		LatestMinorVersion: latestMinorVersion,
		LatestPatchVersion: latestPatchVersion,
	}
}
