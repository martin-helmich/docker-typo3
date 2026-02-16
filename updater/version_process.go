package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sort"
	"text/template"

	"github.com/Masterminds/semver"
)

func processVersion(spec *UpdateVersionSpec) (string, *TYPO3Version, error) {
	constraint, err := semver.NewConstraint(spec.Constraint)
	if err != nil {
		return "", nil, err
	}

	u := fmt.Sprintf("https://get.typo3.org/v1/api/major/%d/release/", spec.Major)
	resp, err := http.Get(u)
	if err != nil {
		return "", nil, err
	}

	j, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", nil, err
	}

	versions := make(TYPO3VersionList, 0)
	if err := json.Unmarshal(j, &versions); err != nil {
		return "", nil, err
	}

	matching := make(TYPO3VersionList, 0, len(versions))

	for _, ver := range versions {
		v, err := semver.NewVersion(ver.Version)
		if err != nil {
			fmt.Printf("WARNING: version %s could not be parsed by semver\n", ver.Version)
			continue
		}

		if ver.ELTS {
			continue
		}

		if matches := constraint.Check(v); matches {
			matching = append(matching, ver)
		}
	}

	if len(matching) == 0 {
		return "", nil, fmt.Errorf("no TYPO3 version (non-ELTS) matching constraint %s", spec.Constraint)
	}

	sort.Sort(sort.Reverse(matching))

	latest := matching[0]

	templateContents, err := ioutil.ReadFile(spec.Template)
	if err != nil {
		return "", nil, err
	}

	t, err := template.New("").Parse(string(templateContents))
	if err != nil {
		return "", nil, err
	}

	target := path.Join(spec.Destination, "Dockerfile")

	if _, err = os.Stat(spec.Destination); os.IsNotExist(err) {
		if err := os.MkdirAll(spec.Destination, 0755); err != nil {
			return "", nil, err
		}
	}

	out, err := os.Create(target)
	if err != nil {
		return "", nil, err
	}

	defer out.Close()

	err = t.Execute(out, TemplateData{
		Major:    spec.Major,
		Latest:   latest.Version,
		Checksum: latest.TarPackage.SHA256,
		Values:   spec.TemplateData,
	})

	if err != nil {
		return "", nil, err
	}

	if spec.Major >= 12 {
		entrypointSrc := path.Join(path.Dir(spec.Template), "docker-entrypoint.sh")
		entrypointDst := path.Join(spec.Destination, "docker-entrypoint.sh")

		entrypointContents, err := ioutil.ReadFile(entrypointSrc)
		if err != nil {
			return "", nil, fmt.Errorf("failed to read entrypoint script: %s", err)
		}

		if err := ioutil.WriteFile(entrypointDst, entrypointContents, 0755); err != nil {
			return "", nil, fmt.Errorf("failed to write entrypoint script: %s", err)
		}
	}

	return target, &latest, nil
}

func strptr(s string) *string {
	return &s
}
