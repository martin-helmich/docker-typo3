package main

import (
	"github.com/go-git/go-git/v5"
	"github.com/sirupsen/logrus"
)

func mustIsChanged(repo *git.Repository, file string) bool {
	c, err := isChanged(repo, file)
	if err != nil {
		panic(err)
	}

	return c
}

func isChanged(repo *git.Repository, file string) (bool, error) {
	l := logrus.WithField("file.name", file)
	l.Debug("testing if file is changed")

	tree, err := repo.Worktree()
	if err != nil {
		return false, err
	}

	status, err := tree.Status()
	if err != nil {
		return false, err
	}

	// See [1] as to why this is necessary.
	//   [1]: https://github.com/src-d/go-git/issues/1300
	fileStatus, ok := status[file]
	if !ok {
		l.Debug("file seems unmodified")
		return false, nil
	}

	changed := fileStatus.Worktree != git.Unmodified
	l.
		WithField("git.status", string(fileStatus.Worktree)).
		WithField("file.modified", changed).
		Debug("check complete")

	return changed, nil
}
