package parseversion

import (
	"github.com/Masterminds/semver/v3"
	log "github.com/sirupsen/logrus"
)

func ToSemver(version string) (*semver.Version, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		log.Errorf("Unable to parse version: %v. Error message: %v\n", version, err)
		return nil, err
	}

	return v, nil
}
