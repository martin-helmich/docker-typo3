package main

import (
	"context"
	"flag"
	"io/ioutil"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/google/go-github/v29/github"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v2"
)

const Owner = "martin-helmich"
const Repo = "docker-typo3"

var configFile string
var noPR bool
var accessToken string
var verbose bool

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

func init() {
	logrus.SetOutput(os.Stderr)
	logrus.SetLevel(logrus.WarnLevel)
}

func main() {
	flag.StringVar(&configFile, "config", ".updater.yaml", "path to config")
	flag.StringVar(&accessToken, "access-token", "", "access token")
	flag.BoolVar(&noPR, "no-pr", false, "do not create pull requests")
	flag.BoolVar(&verbose, "v", false, "more logging")
	flag.Parse()

	if verbose {
		logrus.SetLevel(logrus.DebugLevel)
	}

	ctx := context.Background()

	updaterContents, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic(err)
	}

	updateSpec := UpdateSpec{}
	if err := yaml.Unmarshal(updaterContents, &updateSpec); err != nil {
		panic(err)
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)

	httpClient := oauth2.NewClient(ctx, ts)

	client := github.NewClient(httpClient)

	repo, err := git.PlainOpen(".")
	if err != nil {
		panic(err)
	}

	for _, v := range updateSpec.Versions {
		l := logrus.WithField("version.branch", v.Destination)
		l.Info("processing version")

		df, latest, err := processVersion(&v)
		if err != nil {
			l.WithError(err).Error("error while processing version")
			continue
		}

		wf, err := updateWorkflowFile(&v)
		if err != nil {
			l.WithError(err).Error("error while updating workflow file")
			continue
		}

		if !mustIsChanged(repo, df) && !mustIsChanged(repo, wf) {
			l.Info("skipping -- no files have changed")
			continue
		}

		if noPR {
			l.Info("skipping creating PR")
			continue
		}

		if err := publishVersion(ctx, client, repo, v, latest, wf); err != nil {
			l.WithError(err).Error("error while creating PR")
			continue
		}
	}

}
