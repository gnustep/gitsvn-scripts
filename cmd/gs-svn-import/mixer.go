package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
)

func mixer(ctx context.Context, matches GitMatches) error {
	fmt.Println("mixing...")

	oldGitPath := *oldGitPathBase + "/" + path.Base(*subpath)
	newGitPath := *outputGitPathBase + "/" + path.Base(*subpath)

	oldGit, err := git.PlainOpen(oldGitPath)
	if err != nil {
		return fmt.Errorf("failed to open oldGit: %s", err)
	}
	newGit, err := git.PlainOpen(newGitPath)
	if err != nil {
		return fmt.Errorf("failed to open newGit: %s", err)
	}

	newGit.CreateRemote(&config.RemoteConfig{
		Name: "old",
		URL:  "file://" + oldGitPath,
	})

	/*
		        // The following block gives us "unknown channel NAK"
			remote, _ := newGit.Remote("old")
			err := remote.Fetch(&git.FetchOptions{
				RefSpecs: []config.RefSpec{
					"refs/heads/master:refs/remotes/old/old",
				},
				//ReferenceName: "refs/heads/master",
				//SingleBranch:  true,
			})

			fmt.Println(err)
	*/
	cmd := exec.CommandContext(ctx, "git", "fetch", "old", "refs/heads/*:refs/remotes/old/*")
	cmd.Dir = newGitPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Run()

	if err := addOldHeadToNewGit(ctx, oldGit, newGitPath); err != nil {
		// not adding extra context as context is already available in the error as returned
		return err
	}

	if err := addReplaceRefs(ctx, matches, newGitPath); err != nil {
		return err
	}
	return nil
}

func addOldHeadToNewGit(ctx context.Context, oldGit *git.Repository, newGitPath string) error {
	oldBranch, err := oldGit.Reference("refs/heads/master", true)
	if err != nil {
		return fmt.Errorf("failed to find oldGit's master branch: %s", err)
	}

	f, err := os.Create(newGitPath + "/.git/refs/heads/old")
	if err != nil {
		return fmt.Errorf("failed to create 'old' branch in newGit: %s", err)
	}
	defer f.Close()

	if _, err := f.WriteString(oldBranch.Hash().String()); err != nil {
		return fmt.Errorf("failed to write 'old' branch hash in newGit: %s", err)
	}

	return nil
}

func addReplaceRefs(ctx context.Context, matches GitMatches, newGitPath string) error {
	err := os.MkdirAll(newGitPath+"/.git/refs/replace", os.ModeDir|0755)
	if err != nil {
		return fmt.Errorf("could not make replacerefs dir: %s", err)
	}

	for rev, match := range matches {
		fmt.Printf("writing rev %d: %+v\n", rev, match)

		f, err := os.Create(newGitPath + "/.git/refs/replace/" + match.OldGitHash.String())
		if err != nil {
			return fmt.Errorf("failed to create replace ref for r%d (old hash: %s, new hash: %s): %s", match.SubversionRev, match.OldGitHash, match.NewGitHash, err)
		}
		defer f.Close()

		f.WriteString(match.NewGitHash.String())
	}
	return nil
}
