package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"path"

	"github.com/go-git/go-git/v5"
	"github.com/google/go-github/v29/github"
	"github.com/sirupsen/logrus"
)

func pushToBranch(ctx context.Context, l logrus.FieldLogger, repo *git.Repository, client *github.Client, name string, branch string, msg string) error {
	l = l.WithField("file.name", name)
	l.Debug("creating commit for modified file")

	changed, err := isChanged(repo, name)
	if err != nil {
		return err
	}

	if !changed {
		l.Debug("file is not changed -- not commiting")
		return nil
	}

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
		l.Debug("file does not exist yet; creating")

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
	} else {
		remoteContents, err := fc.GetContent()
		if err != nil {
			return err
		}

		if string(contents) != remoteContents {
			l.Debug("file exists and is changed; updating")

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
		}

		return err
	}

	return nil
}

func publishVersion(
	ctx context.Context,
	client *github.Client,
	repo *git.Repository,
	v UpdateVersionSpec,
	latest *TYPO3Version,
	workflowName string,
) error {
	branch := fmt.Sprintf("update-typo-%d-to-%s", v.Major, latest.Version)
	dockerfileName := path.Join(v.Destination, "Dockerfile")

	l := logrus.WithFields(logrus.Fields{"version.branch": v.Destination, "branch.name": branch})

	refs, _, err := client.Git.GetRefs(ctx, Owner, Repo, "heads/master")
	if err != nil {
		return err
	}

	master := refs[0]

	l.WithField("master.sha", *master.Object.SHA).Debug("found master")

	branchRefs, _, _ := client.Git.GetRefs(ctx, Owner, Repo, fmt.Sprintf("heads/%s", branch))

	if len(branchRefs) > 0 {
		l.Info("branch already exists -- force-resetting")

		_, _, err = client.Git.UpdateRef(
			ctx,
			Owner,
			Repo,
			&github.Reference{
				Ref:    strptr(fmt.Sprintf("heads/%s", branch)),
				Object: master.Object,
			},
			true,
		)
	} else {
		l.Info("branch does not exist -- creating")

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

	if err := pushToBranch(ctx, l, repo, client, dockerfileName, branch, fmt.Sprintf("Bump TYPO3 %d to version %s", v.Major, latest.Version)); err != nil {
		return err
	}

	if err := pushToBranch(ctx, l, repo, client, workflowName, branch, fmt.Sprintf("Update Github workflow for TYPO3 %s", v.Destination)); err != nil {
		return err
	}

	t := true

	prs, _, err := client.PullRequests.List(ctx, Owner, Repo, &github.PullRequestListOptions{
		State: "open",
		Head:  branch,
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
