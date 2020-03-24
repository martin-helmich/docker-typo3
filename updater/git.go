package main

import (
	"gopkg.in/src-d/go-git.v4"
)

func mustIsChanged(repo *git.Repository, file string) bool {
	c, err := isChanged(repo, file)
	if err != nil {
		panic(err)
	}

	return c
}

func isChanged(repo *git.Repository, file string) (bool, error) {
	tree, err := repo.Worktree()
	if err != nil {
		return false, err
	}

	status, err := tree.Status()
	if err != nil {
		return false, err
	}

	fileStatus := status.File(file)

	return fileStatus.Worktree != git.Unmodified, nil
}
