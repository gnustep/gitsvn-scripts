package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"gopkg.in/src-d/go-git.v4/plumbing"
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

func main() {
	flag.Parse()

	if *subpath == "" {
		fmt.Println("specify --subpath")
		return
	}

	// perform an svn clone
	c := SvnCloner{
		ActualSubversionURLBase:    *actualSubversionURLBase,
		CanonicalSubversionURLBase: *canonicalSubversionURLBase,
		OutputGitPathBase:          *outputGitPathBase,
		StdLayout:                  *stdLayout,
		Subpath:                    *subpath,
		AuthorsFilePath:            *authorsFilePath,
	}
	c.Clone(context.TODO())

	// fetch the oldgit repo
	oldGit, err := oldGit()
	if err != nil {
		return
	}

	branchesIter, err := oldGit.Branches()
	if err != nil {
		fmt.Printf("failed to get branches: %s\n", err)
		return
	}

	// TODO(ivucica): maybe validate that the old repo only has refs/heads/master
	// TODO(ivucica): support branch 'oldimport'?
	err = branchesIter.ForEach(func(r *plumbing.Reference) error {
		fmt.Printf("%s\n", r.Strings())
		return nil
	})
	if err != nil {
		fmt.Printf("failed to iterate over branches: %s\n", err)
		return
	}

	// invoking matcher
	fmt.Printf("matching\n")
	matcher()
}
