package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
)

type SvnCloner struct {
	CanonicalSubversionURLBase string
	ActualSubversionURLBase    string
	OutputGitPathBase          string
	StdLayout                  bool
	Subpath                    string
	AuthorsFilePath            string
}

func (opts SvnCloner) Clone(ctx context.Context) error {
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
		opts.OutputGitPathBase + "/" + path.Base(opts.Subpath),
	}...)

	fmt.Printf("git svn clone'ing\n")
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
