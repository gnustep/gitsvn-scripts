package main

import (
	"fmt"
	"gopkg.in/src-d/go-git.v4"
	"os"
	"path"

	"github.com/golang/glog"
)

func oldGit() (*git.Repository, error) {
	oldGitPath := *oldGitPathBase + "/" + *subpath
	if err := os.MkdirAll(path.Dir(oldGitPath), os.ModeDir|0755); err != nil {
		return nil, fmt.Errorf("could not make a directory to host old git path %s: %s", oldGitPath, err)
	}
	if _, err := os.Stat(oldGitPath); os.IsNotExist(err) {
		opts := &git.CloneOptions{
			URL:               "https://github.com/gnustep/" + path.Base(*subpath),
			RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		}
		glog.Infof("cloning old repo from %s to %s\n", opts.URL, oldGitPath)

		oldGit, err := git.PlainClone(*oldGitPathBase+"/"+*subpath, false, opts)
		if err != nil {
			glog.Errorf("failed to git clone the old repo: %s", err)
			return nil, err
		}
		return oldGit, nil
	} else if err != nil {
		glog.Errorf("error checking if oldgit path %s exists: %s", oldGitPath, err)
		return nil, err
	} else {
		glog.Info("skipping oldgit clone as the oldgit repo exists locally")
		return git.PlainOpen(oldGitPath)
	}
}
