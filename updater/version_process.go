package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sort"
	"text/template"

	"github.com/Masterminds/semver"
)

func processVersion(spec *UpdateVersionSpec) (bool, *TYPO3Version, error) {
	constraint, err := semver.NewConstraint(spec.Constraint)
	if err != nil {
		return false, nil, err
	}

	u := fmt.Sprintf("https://get.typo3.org/v1/api/major/%d/release/", spec.Major)
	resp, err := http.Get(u)
	if err != nil {
		return false, nil, err
	}

	j, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, nil, err
	}

	versions := make(TYPO3VersionList, 0)
	if err := json.Unmarshal(j, &versions); err != nil {
		return false, nil, err
	}

	matching := make(TYPO3VersionList, 0, len(versions))

	for _, ver := range versions {
		v, err := semver.NewVersion(ver.Version)
		if err != nil {
			fmt.Printf("WARNING: version %s could not be parsed by semver\n", ver.Version)
			continue
		}

		if matches := constraint.Check(v); matches {
			matching = append(matching, ver)
		}
	}

	if len(matching) == 0 {
		return false, nil, fmt.Errorf("no TYPO3 version matching constraint %s", spec.Constraint)
	}

	sort.Sort(sort.Reverse(matching))

	latest := matching[0]

	templateContents, err := ioutil.ReadFile(spec.Template)
	if err != nil {
		return false, nil, err
	}

	t, err := template.New("").Parse(string(templateContents))
	if err != nil {
		return false, nil, err
	}

	var currentHash = ""

	target := path.Join(spec.Destination, "Dockerfile")

	if _, err = os.Stat(spec.Destination); os.IsNotExist(err) {
		if err := os.MkdirAll(spec.Destination, 0755); err != nil {
			return false, nil, err
		}
	}

	h := md5.New()
	existing, err := os.Open(target)
	if err != nil {
		if os.IsNotExist(err) {
			currentHash = ""
		} else {
			return false, nil, err
		}
	} else {
		if _, err := io.Copy(h, existing); err != nil {
			return false, nil, err
		}

		currentHash = fmt.Sprintf("%x", h.Sum(nil))
	}

	out, err := os.Create(target)
	if err != nil {
		return false, nil, err
	}

	defer out.Close()

	err = t.Execute(out, TemplateData{
		Major:    spec.Major,
		Latest:   latest.Version,
		Checksum: latest.TarPackage.SHA256,
		Values:   spec.TemplateData,
	})

	if err != nil {
		return false, nil, err
	}

	h2 := md5.New()
	written, err := os.Open(target)
	if err != nil {
		return false, nil, err
	}

	if _, err := io.Copy(h2, written); err != nil {
		return false, nil, err
	}

	newHash := fmt.Sprintf("%x", h2.Sum(nil))

	return currentHash != newHash, &latest, nil
}

func strptr(s string) *string {
	return &s
}
