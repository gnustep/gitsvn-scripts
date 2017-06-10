# gs-svn-import

This directory contains a program implemented in Go intended to perform the
conversion of GNUstep's Subversion repository into many little repositories.

It attempts to use `git-replace` refs in order to link the old conversion with
the new, more nicely done conversion.

It will also perform an export into a remote repository. This could be a public
hosting service, for example a Gitlab instance, or Github.

NOTE: This export currently not done, and you'll want to do this:

    git remote add github https://github.com/gnustep/REPONAME
    REMOTENAME=github

    # push old branch
    git push -u $REMOTENAME old

    # forcibly push new master branch
    git push -u -f $REMOTENAME master

    # push all other branches (including old and master, which is noop)
    for i in $(ls .git/refs/heads) ; do echo "refs/heads/$i" ; done | xargs git push -u $REMOTENAME

    # push replace refs
    git push $REMOTENAME 'refs/replace/*'

    # later, if needed, pull replace refs:
    git pull $REMOTENAME 'refs/replace/*:refs/replace/*'  # pull replace refs

Or to obtain replace refs after cloning:

    git clone https://github.com/gnustep/REPONAME
    cd REPONAME
    git fetch origin 'refs/replace/*:refs/replace/*'

-- ivucica
