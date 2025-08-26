![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/bladeacer/mnemosync?style=for-the-badge&logo=go)
![GitHub License](https://img.shields.io/github/license/bladeacer/mnemosync?style=for-the-badge)

# mnemosync

A CLI tool that lets you add folders to backup manually to a target Git repository.

The name is inspired by the Greek Goddess of memory Mnemosyne.

## Installation guide

This installation guide assumes you know how to create and set up a Git repository.

**Always backup your files before using mmsync**.

```bash
go install github.com/bladeacer/mnemosync
```

Ensure that you can access Go binaries in your $PATH.

```bash
mmsync
```

## Project status
WIP. See [this GitHub project](https://github.com/users/bladeacer/projects/3) for
the progress tracker.

This is my first project using the Go programming language, but I hope it will
be useful.

## Planned features
- Check if required binaries are available before calling the tool
  - Required binaries: `git, rsync, tar, zip`

- Help command line flag
___
- CRUD target directories which user wishes to backup e.g.

- Rsync to mirror said target directories to a `~/.mnemosync/folders`
  - Either manually triggered or we integrate `cron`
- Wrapper for user to manually copy the files and push them in their Git repository

- Wrapper to let user set default commit message format

## License

This Golang CLI app, "mnemosync" is released under the GNU General Public
License version 3 (GPLv3) License.

### License Notice

```
This file is part of mnemosync. mnemosync is a CLI tool that lets you add
folders to backup manually to a target Git repository. 

Copyright (c) 2025 bladeacer

mnemosync is free software: you can redistribute it and/or modify it under the
terms of the GNU General Public License as published by the Free Software
Foundation, either version 3 of the License, or (at your option) any later version.

mnemosync is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY;
without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR
PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with mnemosync.
If not, see <https://www.gnu.org/licenses/>. 
```

### License file

You can find the [license file here](./LICENSE).

## Credits

This CLI was made possible by [Cobra CLI](https://github.com/spf13/cobra).

## Planned CLI spec

TBC. View currently available options by running mnemosync without any flags or arguments.

```bash
## Init and config
# Init the app with helpers to get the user to set path config and all
mmsync init 
mmsync get-config-path
mmsync set-config-path
mmsync set-repo-path
mmsync get-repo-path

## CRUD directories to mmsync before staging
mmsync add <target_path> -a <optional_alias>
mmsync list
mmsync change <target_path-or-alias> <new-target_path-or-alias>
mmsync rm <target_path-or-alias>

## Find a mmsync path or alias that has been added
mmsync search <query-by-path-or-alias>

## Add warning for user to confirm if they wish to delete all directories they added
mmsync clear

## Location of mmsync related files, defaults to ~/.config/.mmsync/
mmsync locate
mmsync set-location <custom_path>

# Backup related
## Technical info: staging is just a temp directory over at ~/.config/.mmsync/staging.
## Inherits from the mmsync path setting
## You can use . to include all directories and aliases

# rsyncs all added target mmsync directories or aliases to staging 
mmsync stage <target_path-or-alias> 
# rsyncs unstages added target mmsync directories or aliases to staging 
mmsync unstage <target_path-or-alias> 
mmsync get-staging # get status of staging

mmsync get-hist # get staging history
## get staging history limit in days before it is cleared. Defaults to 7 days and a max of 1024 MB
mmsync get-hist-limit 
mmsync set-hist-limit -d <number_of_days> -s <max_size_in_mb>
mmsync clear-hist # clears staging history

mmsync set-archiver tar|zip
mmsync get-archiver # gets archive tool used, defaults to tar

# Git related

## Configure commit messages
mmsync get-commit-fmt # Defaults to mnemosync archive ISO timestamp
mmsync set-commit-fmt <custom_format>

## checks if anything in staging, if yes it compresses writes the archive file over to be pushed
## if not, warns the user that staging is empty
## When pushing, write folder and filenames affected to viewable local db as part of staging history
## Also does the needed git commit and push on behalf of the user.
mmsync push 

## Respecting .gitignore
## mmsync respects gitignore in the directories you provide by default
## Returns true or 1 by default
mmsync get-ignore 
mmsync set-ignore 0|1

# Misc
mmsync version
mmsync help
mmsync status # status?
mmsync status --verbose # status with verbose flag
mmsync logs # gets runtime logs opened in user's $EDITOR
```
