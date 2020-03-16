package main

import (
	"encoding/json"

	"github.com/Masterminds/semver"
)

type TYPO3VersionList []TYPO3Version

type TYPO3VersionChecksums struct {
	MD5    string `json:"md5sum"`
	SHA1   string `json:"sha1sum"`
	SHA256 string `json:"sha256sum"`
}

type TYPO3Version struct {
	Version         string
	SemanticVersion *semver.Version
	TarPackage      TYPO3VersionChecksums
}

func (t *TYPO3Version) UnmarshalJSON(b []byte) error {
	temp := struct {
		Version    string                `json:"version"`
		TarPackage TYPO3VersionChecksums `json:"tar_package"`
	}{}

	if err := json.Unmarshal(b, &temp); err != nil {
		return err
	}

	v, err := semver.NewVersion(temp.Version)
	if err != nil {
		return err
	}

	t.Version = temp.Version
	t.SemanticVersion = v
	t.TarPackage = temp.TarPackage

	return nil
}

func (l TYPO3VersionList) Len() int {
	return len(l)
}

func (l TYPO3VersionList) Less(i, j int) bool {
	return l[i].SemanticVersion.LessThan(l[j].SemanticVersion)
}

func (l TYPO3VersionList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}
