package main

import (
	"fmt"
	"path"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
)

func mixer() {
	fmt.Println("mixing...")
	//oldGit, _ := git.PlainOpen(*oldGitPathBase + "/" + path.Base(*subpath))
	newGit, _ := git.PlainOpen(*outputGitPathBase + "/" + path.Base(*subpath))

	newGit.CreateRemote(&config.RemoteConfig{
		Name: "old",
		URL:  "file://" + *oldGitPathBase + "/" + path.Base(*subpath),
	})
	remote, _ := newGit.Remote("old")
	err := remote.Fetch(&git.FetchOptions{
		RefSpecs: []config.RefSpec{
			"refs/heads/master:refs/remotes/old/master",
		},
		//ReferenceName: "refs/heads/master",
		//SingleBranch:  true,
	})
	fmt.Println(err)
}
