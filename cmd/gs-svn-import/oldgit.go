package main

import (
	"fmt"
	"gopkg.in/src-d/go-git.v4"
	"os"
	"path"
)

func oldGit() (*git.Repository, error) {
	oldGitPath := *oldGitPathBase + "/" + path.Base(*subpath)
	if _, err := os.Stat(oldGitPath); os.IsNotExist(err) {
		opts := &git.CloneOptions{
			URL:               "https://github.com/gnustep/" + path.Base(*subpath),
			RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		}
		fmt.Printf("cloning old repo from %s to %s\n", opts.URL, oldGitPath)

		oldGit, err := git.PlainClone(*oldGitPathBase+"/"+path.Base(*subpath), false, opts)
		if err != nil {
			fmt.Printf("failed to git clone the old repo: %s\n", err)
			return nil, err
		}
		return oldGit, nil
	} else if err != nil {
		fmt.Printf("error checking if oldgit path %s exists: %s\n", oldGitPath, err)
		return nil, err
	} else {
		fmt.Printf("skipping oldgit clone as the oldgit repo exists locally\n")
		return git.PlainOpen(oldGitPath)
	}
}
