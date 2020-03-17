package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
	"path"
	"strings"

	"github.com/google/go-github/v29/github"
)

func publishVersion(ctx context.Context, client *github.Client, v UpdateVersionSpec, latest *TYPO3Version) error {
	contents, err := ioutil.ReadFile(path.Join(v.Destination, "Dockerfile"))
	if err != nil {
		return err
	}

	r := rand.Intn(9000 + 1000)
	//branch := fmt.Sprintf("update-typo-%d-to-%s-%d", v.Major, latest.Version, r)
	branch := fmt.Sprintf("update-typo-%d-to-%s", v.Major, latest.Version, r)
	name := strings.Replace(path.Join(v.Destination, "Dockerfile"), "../", "", -1)

	refs, _, err := client.Git.GetRefs(ctx, Owner, Repo, "heads/master")
	if err != nil {
		return err
	}

	master := refs[0]

	fmt.Printf("master is at %s\n", *master.Object.SHA)

	branchRefs, _, err := client.Git.GetRefs(ctx, Owner, Repo, fmt.Sprintf("heads/%s", branch))
	if err != nil {
		return err
	}

	if len(branchRefs) > 0 {
		fmt.Printf("branch %s already exists\n", branch)
	} else {
		fmt.Printf("creating branch %s\n", branch)

		_, _, err = client.Git.CreateRef(
			ctx,
			Owner,
			Repo,
			&github.Reference{
				Ref:    strptr(fmt.Sprintf("heads/%s", branch)),
				Object: master.Object,
			},
		)

		if err != nil {
			return err
		}
	}

	_, _, err = client.Repositories.CreateFile(
		ctx,
		Owner,
		Repo,
		name,
		&github.RepositoryContentFileOptions{
			Message: strptr(fmt.Sprintf("Bump TYPO3 %d to version %s", v.Major, latest.Version)),
			Content: contents,
			Branch:  &branch,
			Author: &github.CommitAuthor{
				Name:  strptr("TYPO3 Docker Update Bot"),
				Email: strptr("martin@helmich.me"),
			},
			Committer: &github.CommitAuthor{
				Name:  strptr("TYPO3 Docker Update Bot"),
				Email: strptr("martin@helmich.me"),
			},
		},
	)

	if err != nil {
		return err
	}

	t := true

	prs, _, err := client.PullRequests.List(ctx, Owner, Repo, &github.PullRequestListOptions{
		State: "open",
		Head:  fmt.Sprintf("%s:%s", Owner, branch),
	})

	if err != nil {
		return nil
	}

	if len(prs) > 0 {
		fmt.Printf("PR already exists: %d", prs[0].ID)
		return nil
	}

	pr, _, err := client.PullRequests.Create(
		ctx,
		Owner,
		Repo,
		&github.NewPullRequest{
			Title:               strptr(fmt.Sprintf("Bump TYPO3 %d to version %s", v.Major, latest.Version)),
			Head:                &branch,
			Base:                strptr("master"),
			MaintainerCanModify: &t,
		},
	)

	fmt.Printf("Created PR %d\n", *pr.ID)

	return nil
}
