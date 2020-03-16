package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/v29/github"
	"gopkg.in/yaml.v2"
)

const Owner = "martin-helmich"
const Repo = "docker-typo3"

var configFile string
var keyFile string
var appID int64
var installationID int64
var noPR bool

type UpdateVersionSpec struct {
	Major        int         `yaml:"major"`
	Constraint   string      `yaml:"constraint"`
	Template     string      `yaml:"template"`
	Destination  string      `yaml:"destination"`
	TemplateData interface{} `yaml:"templateData"`
}

type UpdateSpec struct {
	Versions []UpdateVersionSpec `yaml:"versions"`
}

type TemplateData struct {
	Major    int
	Latest   string
	Checksum string
	Values   interface{}
}

func main() {
	flag.StringVar(&configFile, "config", ".updater.yaml", "path to config")
	flag.StringVar(&keyFile, "key-file", "", "path to key file")
	flag.Int64Var(&appID, "app-id", 0, "application ID")
	flag.Int64Var(&installationID, "installation-id", 0, "installation ID")
	flag.BoolVar(&noPR, "no-pr", false, "do not create pull requests")
	flag.Parse()

	updaterContents, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic(err)
	}

	updateSpec := UpdateSpec{}
	if err := yaml.Unmarshal(updaterContents, &updateSpec); err != nil {
		panic(err)
	}

	itr, err := ghinstallation.NewKeyFromFile(http.DefaultTransport, appID, installationID, keyFile)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	client := github.NewClient(&http.Client{Transport: itr})

	for _, v := range updateSpec.Versions {
		changed, latest, err := processVersion(&v)
		if err != nil {
			panic(err)
		}

		if !changed {
			fmt.Printf("Dockerfile for version %d did not change\n", v.Major)
			continue
		}

		if !noPR {
			if err := publishVersion(ctx, client, v, latest); err != nil {
				fmt.Printf("ERROR while creating PR for version %d: %s\n", v.Major, err.Error())
				continue
			}
		}
	}

}
