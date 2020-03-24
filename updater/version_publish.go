package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"path"

	"github.com/google/go-github/v29/github"
)

func pushToBranch(ctx context.Context, client *github.Client, name string, branch string, msg string) error {
	contents, err := ioutil.ReadFile(name)
	if err != nil {
		return err
	}

	fc, _, _, err := client.Repositories.GetContents(
		ctx,
		Owner,
		Repo,
		name,
		&github.RepositoryContentGetOptions{
			Ref: fmt.Sprintf("heads/%s", branch),
		},
	)

	if fc == nil {
		_, _, err = client.Repositories.CreateFile(
			ctx,
			Owner,
			Repo,
			name,
			&github.RepositoryContentFileOptions{
				Message: &msg,
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

		return err
	} else if *fc.Content != string(contents) {
		_, _, err = client.Repositories.UpdateFile(
			ctx,
			Owner,
			Repo,
			name,
			&github.RepositoryContentFileOptions{
				Message: &msg,
				SHA:     fc.SHA,
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

		return err
	}

	return nil
}

func publishVersion(ctx context.Context, client *github.Client, v UpdateVersionSpec, latest *TYPO3Version, workflowName string) error {
	branch := fmt.Sprintf("update-typo-%d-to-%s", v.Major, latest.Version)
	dockerfileName := path.Join(v.Destination, "Dockerfile")

	refs, _, err := client.Git.GetRefs(ctx, Owner, Repo, "heads/master")
	if err != nil {
		return err
	}

	master := refs[0]

	fmt.Printf("master is at %s\n", *master.Object.SHA)

	branchRefs, _, _ := client.Git.GetRefs(ctx, Owner, Repo, fmt.Sprintf("heads/%s", branch))

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

	if err := pushToBranch(ctx, client, dockerfileName, branch, fmt.Sprintf("Bump TYPO3 %d to version %s", v.Major, latest.Version)); err != nil {
		return err
	}

	if err := pushToBranch(ctx, client, workflowName, branch, fmt.Sprintf("Update Github workflow for TYPO3 %s", v.Destination)); err != nil {
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
