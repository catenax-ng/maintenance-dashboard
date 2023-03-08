package latestversions

import (
	"context"
	"os"
	"sort"
	"strings"

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
	NewReleasesApiKey = os.Getenv("NEWREKEASES_API_KEY")
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
				log.Warningf("Skipping invalid version: %v\n", release.Version)
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
		LatestMajorVersion: latestMajorVersion,
		LatestMinorVersion: latestMinorVersion,
		LatestPatchVersion: latestPatchVersion,
	}
}

// Get all the projects added to the newreleases account and parse them to ProjectInfo struct
func GetAllProjects() []ProjectInfo {
	client := newreleases.NewClient(NewReleasesApiKey, nil)
	ctx := context.Background()
	var pi []ProjectInfo
	o := &newreleases.ProjectListOptions{
		Page: 1,
	}
	for {
		projects, lastPage, err := client.Projects.List(ctx, *o)
		if err != nil {
			log.Panic(err)
		}

		for _, proj := range projects {
			releaseNameParts := strings.Split(proj.Name, "/")
			releaseName := releaseNameParts[len(releaseNameParts)-1]
			pi = append(pi, ProjectInfo{ID: proj.ID, ReleaseName: releaseName})
		}

		if o.Page >= lastPage {
			break
		}
		o.Page++
	}

	return pi
}

// Get max 10 pages of releases for an app and parse versions to semVer
func GetLatestVersionsForApp(projectID string) semver.Collection {
	client := newreleases.NewClient(NewReleasesApiKey, nil)
	ctx := context.Background()
	vs := []*semver.Version{}

	for i := 1; i < 10; i++ {
		releases, lastPage, err := client.Releases.ListByProjectID(ctx, projectID, i)
		if err != nil {
			log.Panic(err)
		}

		for _, release := range releases {
			parsedVersion, err := parseversion.ToSemver(release.Version)
			if err != nil {
				log.Warningf("Skipping invalid version: %v\n", release.Version)
			} else {
				vs = append(vs, parsedVersion)
			}
		}

		if i >= lastPage {
			break
		}
	}

	return vs
}
