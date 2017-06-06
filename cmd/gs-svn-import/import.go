package main // import "github.com/gnustep/gitsvn-scripts/cmd/gs-svn-import"

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
)

var (
	canonicalSubversionURLBase = flag.String("canonical_subversion_url_base", "svn+ssh://svn.gna.org/svn/gnustep", "base URL hosting GNUstep repositories which is to be used in git-svn-id in commit descriptions")
	actualSubversionURLBase    = flag.String("actual_subversion_url_base", "svn+ssh://svn.gna.org/svn/gnustep", "actual URL from which the GNUstep repositories will be fetched; this could be a locally rsync'ed copy")

	authorsFilePath = flag.String("authors_file_path", os.Getenv("GOPATH")+"/authors.txt", "path to the authors.txt file")

	outputGitPathBase = flag.String("output_git_path_base", os.Getenv("GOPATH")+"/gs-svn/git", "base path at which output git repos will be placed / updated")

	oldGitPathBase      = flag.String("old_git_path_base", os.Getenv("GOPATH")+"/gs-svn/oldgit", "base path for old git repositories for which to generate replace refs")
	generateReplaceRefs = flag.Bool("generate_replace_refs", true, "whether replace refs should be generated for this repository")

	stdLayout = flag.Bool("stdlayout", true, "whether the Subversion repository contains a 'standard' layout; i.e. trunk/branches/tags structure")
	subpath   = flag.String("subpath", "", "subpath to lib or app to convert, sans the 'trunk' part. basename will be used as the output git repo name. example: libs/gui => gui, apps/gorm => gorm.")
)

type SvnCloner struct {
	CanonicalSubversionURLBase string
	ActualSubversionURLBase    string
	OutputGitPathBase          string
	StdLayout                  bool
	Subpath                    string
	AuthorsFilePath            string
}

func (opts SvnCloner) Clone(ctx context.Context) {

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

	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Run()
}

func main() {
	flag.Parse()

	if *subpath == "" {
		fmt.Println("specify --subpath")
		return
	}

	c := SvnCloner{
		ActualSubversionURLBase:    *actualSubversionURLBase,
		CanonicalSubversionURLBase: *canonicalSubversionURLBase,
		OutputGitPathBase:          *outputGitPathBase,
		StdLayout:                  *stdLayout,
		Subpath:                    *subpath,
		AuthorsFilePath:            *authorsFilePath,
	}
	c.Clone(context.TODO())
}
