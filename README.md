Scripts for migrating SVN repositories to Git.

## Quick start

- Any SVN layout: `./export-all [--nostdlayout] [-n|--dry-run] <SVN_URL> <GIT_REMOTE_URL> <DIR_NAME>`
- Batch from `svn-urls.txt`: `./transfer <github-username> <github-token-or-password>`

## Scripts

- `export-all` – Clone an SVN repo (default `git svn --stdlayout`, or plain `git svn clone` with `--nostdlayout`), add the provided Git remote, push `master`, and move the result into `~/repositories/github/<DIR_NAME>`. Supports dry-run via `-n/--dry-run`.
- `transfer` – Loops over URLs in `svn-urls.txt`, recreates GitHub repos via `recreate_repo`, and calls `export-all` for each. Args: `<username> <password-or-token>`.
- `recreate_repo` – Deletes and recreates a GitHub repo using `delete_repo` and `create_repo`.
- `create_repo` / `delete_repo` – Ruby helpers that prefer `GITHUB_TOKEN` but also accept `<username> <password-or-token>` as args.
- `github-import` – For each subdirectory, recreates a GitHub repo with a prefix and pushes `main`. Args: `<username> <password> <prefix> <org-or-user>`.
- `update-gitsvn` – Run `git svn fetch` + `git svn rebase`, then push. Enable force push via `FORCE_PUSH=1`.
- `update-branches` – Track missing remote branches (skips `origin/HEAD`), then `git pull --all` and `git push --all`.
- `update-all` – Call `update-dir` for every repo under a base directory.
- `update-dir` – Invoke `update-gitsvn` and `update-svndcommit` for a single repo.
- `update-svndcommit` – Pull Git changes and `git svn dcommit` back to SVN.
- `fixit` – Cleans stale git reference directories under each repo listed in `~/files` (destructive; review before use).
- `export-all` and other scripts expect `git svn` to be installed and reachable.

## Conventions

- Destinations: staging clones go to `~/repositories/staging`; final Git copies are placed under `~/repositories/github`.
- Authentication: prefer `GITHUB_TOKEN` in the environment; fall back to username/password arguments where supported.
- Dry runs: `export-all` accepts `-n/--dry-run` to print actions without executing.
- Safety: most scripts use `set -euo pipefail`; destructive operations (like `rm -rf` during export) run only after successful steps when not in dry-run.
