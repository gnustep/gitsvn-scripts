# gs-svn-import

This directory contains a program implemented in Go intended to perform the
conversion of GNUstep's Subversion repository into many little repositories.

That is all that is currently implemented. Shortly, a few more necessary
features will be added.

It will also attempt to use `git-replace` refs, or instead of them grafts, in
order to link the old conversion with the new, more nicely done conversion.

It will also perform an export into a remote repository. This could be a public
hosting service, for example a Gitlab instance, or Github.

-- ivucica
