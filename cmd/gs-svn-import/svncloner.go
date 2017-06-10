package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"gopkg.in/src-d/go-git.v4"
)

type SvnCloner struct {
	CanonicalSubversionURLBase string
	ActualSubversionURLBase    string
	OutputGitPathBase          string
	StdLayout                  bool
	Subpath                    string
	AuthorsFilePath            string
}

func (opts SvnCloner) CopySubversionRemotesToTagsAndHeads(ctx context.Context) error {
	outputGitPath := opts.OutputGitPathBase + "/" + opts.Subpath
	repo, err := git.PlainOpen(outputGitPath)
	if err != nil {
		return fmt.Errorf("could not open gitsvn repo %s: %s", outputGitPath, err)
	}

	refs, err := repo.References()
	if err != nil {
		return fmt.Errorf("could not get gitsvn repo's refs: %s", err)
	}

	for ref, err := refs.Next(); ; ref, err = refs.Next() {
		if ref == nil {
			break
		}
		if err != nil {
			return fmt.Errorf("error listing ref: %s", err)
		}

		if !strings.HasPrefix(ref.Name().String(), "refs/remotes/svn") {
			continue
		}

		if ref.Name().String() == "refs/remotes/svn/trunk" {
			continue
		}
		if strings.HasPrefix(ref.Name().String(), "refs/remotes/svn/tags/") {
			tagName := path.Base(ref.Name().String())
			f, err := os.Create(outputGitPath + "/.git/refs/tags/" + tagName)
			if err != nil {
				return fmt.Errorf("could not create tag %s for repo %s: %s", tagName, outputGitPath, err)
			}
			defer f.Close()

			f.WriteString(ref.Hash().String())
		} else {
			// non-tag
			branchName := path.Base(ref.Name().String())
			if ref.Name().String() != "refs/remotes/svn/"+branchName {

				fmt.Printf("omitting %s while copying svn refs to usual heads/tags: it's probably in a subdir and probably not a branch", ref.Name().String())
				continue
			}
			f, err := os.Create(outputGitPath + "/.git/refs/heads/" + branchName)
			if err != nil {
				return fmt.Errorf("could not create branch %s for repo %s: %s", branchName, outputGitPath, err)
			}
			defer f.Close()

			f.WriteString(ref.Hash().String())
		}
	}
	return nil
}

func (opts SvnCloner) Clone(ctx context.Context) error {

	if err := os.MkdirAll(opts.OutputGitPathBase+"/"+path.Dir(opts.Subpath), os.ModeDir|0755); err != nil {
		return fmt.Errorf("could not make a directory %s to store git svn clone'd repo: %s", opts.OutputGitPathBase+"/"+path.Dir(opts.Subpath))
	}

	args := []string{
		"svn", "clone",
		"--prefix=svn/",
		"--preserve-empty-dirs",
	}

	if opts.ActualSubversionURLBase != "" {
		args = append(args, "--rewrite-root="+opts.CanonicalSubversionURLBase)
	}
	if opts.StdLayout {
		args = append(args, "--stdlayout")
	}
	if opts.AuthorsFilePath != "" {
		args = append(args, "--authors-file="+opts.AuthorsFilePath)
	}

	args = append(args, []string{
		opts.ActualSubversionURLBase + "/" + opts.Subpath,
		opts.OutputGitPathBase + "/" + opts.Subpath,
	}...)

	fmt.Printf("git svn clone'ing\n")
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
