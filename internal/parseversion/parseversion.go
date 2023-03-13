package parseversion

import (
	"github.com/Masterminds/semver/v3"
	log "github.com/sirupsen/logrus"
)

// Tries to parse the provided string to SemVer. Returns nil when failed.
func ToSemver(version string) (*semver.Version, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		log.Errorf("Unable to parse version: %v. Error message: %v\n", version, err)
		return nil, err
	}

	return v, nil
}
