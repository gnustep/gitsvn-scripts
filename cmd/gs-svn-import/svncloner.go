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
