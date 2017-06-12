package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/golang/glog"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

func export(ctx context.Context, skippingOldIsFine bool) error {
	glog.Info("exporting new git repository to github")
	newGitPath := *outputGitPathBase + "/" + *subpath
	newGit, err := git.PlainOpen(newGitPath)
	if err != nil {
		return fmt.Errorf("failed to open newGit: %s", err)
	}

	remoteCfg := &config.RemoteConfig{
		Name: "github",
		URL:  "https://github.com/gnustep/" + strings.Replace(*subpath, "/", "-", -1),
	}
	if _, err := newGit.Remote("github"); err != nil {
		_, err = newGit.CreateRemote(remoteCfg)
		if err != nil {
			return fmt.Errorf("could not create remote config for %s: %s", *subpath, err)
		}
	}

	cmd := exec.CommandContext(ctx, "git", "push", "-u", "github", "old")
	cmd.Dir = newGitPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		if skippingOldIsFine {
			glog.Warningf("could not push old branch to remote repo for %s (but continuing because the relevant option was passed): %s", *subpath, err)
		} else {
			return fmt.Errorf("could not push old branch to remote repo for %s: %s", *subpath, err)
		}
	}

	cmd = exec.CommandContext(ctx, "git", "push", "-u", "-f", "github", "master")
	cmd.Dir = newGitPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("could not push replace master branch on the remote repo %s: %s", remoteCfg.URL, err)
	}

	// push other branches (forcibly)
	refs, err := newGit.References()
	if err != nil {
		return fmt.Errorf("could not get refs in local new repo %s: %s", newGitPath, err)
	}
	refNames := []string{}
	refs.ForEach(func(ref *plumbing.Reference) error {
		if strings.HasPrefix(ref.Name().String(), "refs/heads/") {
			if ref.Name().String() == "refs/heads/master" || ref.Name().String() == "refs/heads/old" {
				// no need to push master or old
				return nil
			}

			refNames = append(refNames, ref.Name().String()) //path.Base(ref.Name().String()))
		}
		return nil
	})

	args := append([]string{"push", "-u", "-f", "github"}, refNames...)

	cmd = exec.CommandContext(ctx, "git", args...)
	cmd.Dir = newGitPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("could not replace other branches onto the remote repo %s: %s", remoteCfg.URL, err)
	}

	// push tags
	cmd = exec.CommandContext(ctx, "git", "push", "--tags", "github")
	cmd.Dir = newGitPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("could not push tags onto remote repo %s: %s", remoteCfg.URL, err)
	}

	// push replace refs
	cmd = exec.CommandContext(ctx, "git", "push", "github", "refs/replace/*")
	cmd.Dir = newGitPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("could not push replace refs onto remote repo %s: %s", remoteCfg.URL, err)
	}

	return nil
}
