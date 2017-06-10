package main

import (
	"bufio"
	"context"
	"fmt"
	"path"
	"strconv"
	"strings"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

type SubversionRevision uint64

type GitMatch struct {
	SubversionRev SubversionRevision
	OldGitHash    plumbing.Hash
	NewGitHash    plumbing.Hash
}

type GitMatches map[SubversionRevision]*GitMatch

func matcher(ctx context.Context) (GitMatches, error) {
	// map from subversion revision to a match entry
	matches := make(GitMatches)

	// in this function, we're abusing closures because .ForEach() API in
	// go-git does not support passing any context.

	oldGit, _ := git.PlainOpen(*oldGitPathBase + "/" + path.Base(*subpath))
	oldHeadRef, _ := oldGit.Head()
	oldHead, _ := oldGit.CommitObject(oldHeadRef.Hash())
	oldHashesVisited := make(map[plumbing.Hash]bool)
	var matchOld func(commit *object.Commit) error
	matchOld = func(commit *object.Commit) error {
		if _, ok := oldHashesVisited[commit.Hash]; ok {
			return nil
		}
		oldHashesVisited[commit.Hash] = true

		rev, err := revisionFromGitCommitMessage(commit.Message)
		if err != nil {
			return err
		}
		if _, ok := matches[rev]; !ok {
			matches[rev] = &GitMatch{
				SubversionRev: rev,
				OldGitHash:    commit.Hash,
			}
		} else {
			matches[rev].OldGitHash = commit.Hash
		}

		commit.Parents().ForEach(matchOld)
		return nil
	}
	matchOld(oldHead)

	newGit, err := git.PlainOpen(*outputGitPathBase + "/" + path.Base(*subpath))
	if err != nil {
		return nil, fmt.Errorf("could not open newgit repo: %s", err)
	}
	newHeadRef, err := newGit.Head()
	if err != nil {
		return nil, fmt.Errorf("could not get newgit repo's head: %s", err)
	}
	newHead, err := newGit.CommitObject(newHeadRef.Hash())
	if err != nil {
		return nil, fmt.Errorf("could not get newgit repo's head commit object: %s", err)
	}
	newHashesVisited := make(map[plumbing.Hash]bool)
	var matchNew func(commit *object.Commit) error
	matchNew = func(commit *object.Commit) error {
		if _, ok := newHashesVisited[commit.Hash]; ok {
			return nil
		}
		newHashesVisited[commit.Hash] = true

		rev, err := revisionFromGitCommitMessage(commit.Message)
		if err != nil {
			return err
		}
		if _, ok := matches[rev]; !ok {
			matches[rev] = &GitMatch{
				SubversionRev: rev,
				NewGitHash:    commit.Hash,
			}
		} else {
			matches[rev].NewGitHash = commit.Hash
		}

		commit.Parents().ForEach(matchNew)
		return nil
	}
	matchNew(newHead)

	return matches, nil
}

func revisionFromGitCommitMessage(message string) (SubversionRevision, error) {
	msgLines := bufio.NewReader(strings.NewReader(message))
	for msgLine, _, err := msgLines.ReadLine(); err == nil; msgLine, _, err = msgLines.ReadLine() {
		msgLineField := strings.Split(string(msgLine), ": ")
		if msgLineField[0] == "git-svn-id" {
			return revisionFromGitSvnId(msgLineField[1])
		}
	}
	return 0, fmt.Errorf("no git svn id found")
}

func revisionFromGitSvnId(gitSvnId string) (SubversionRevision, error) {
	pathRevAndRepoId := strings.Split(gitSvnId, " ")
	pathRev := pathRevAndRepoId[0]
	pathAndRev := strings.Split(pathRev, "@")
	if len(pathAndRev) != 2 {
		return 0, fmt.Errorf("misformatted git-svn-id: %s", gitSvnId)
	}
	rev := pathAndRev[1]
	revInt, err := strconv.Atoi(rev)
	return SubversionRevision(revInt), err
}
