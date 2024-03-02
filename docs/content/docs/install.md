---
weight: 999
title: "Installation"
description: "Local and CI setup"
icon: "download"
date: "2024-01-14T18:38:29+08:00"
lastmod: "2024-01-14T18:38:29+08:00"
toc: true
---

## Go users

```shell
go install github.com/rusinikita/changes/cmd/changes@latest
```

## Homebrew users

```shell
brew tap rusinikita/tap
brew install rusinikita/tap/devex
```

## Binaries

###### One line installer
```shell
curl -sSfL https://raw.githubusercontent.com/rusinikita/changes/master/install.sh | sh -s -- -b ~/bin latest
```

- `-b` flag - folder to place binary. Leave `~/bin` or replace with any folder under `$PATH` works.
- `latest` argument - version to install. Look at [releases](https://github.com/rusinikita/changes/releases) for specific version.


###### Verify installation
```shell
changes help
```

#### Windows

On Windows, you can run the installation commands with Git Bash, which comes with [Git for Windows](https://git-scm.com/download/win).

#### No curl

Some OS or containers don't have `curl`. For example, [alpine linux](https://github.com/alpinelinux/docker-alpine).

Use `wget -O- -nv` instead of `curl -sSfL`.
```shell
wget -O- -nv https://raw.githubusercontent.com/rusinikita/changes/master/install.sh | sh -s -- -b ~/bin latest
```

## Project folder

- Create .changes.yaml (or toml) file inside your repository
- Copy rules from example

## GitHub actions

- Use `actions/checkout` and `actions/setup-go` to create environment
- Call `go run github.com/rusinikita/changes/cmd/changes check`

```yaml
git-check:
  name: git-checks
  runs-on: ubuntu-latest
  steps:
    - uses: actions/checkout@v4
      with:
        # skip gh merge commit from diff
        ref: ${{ github.event.pull_request.head.sha }}
        fetch-depth: 0

    - uses: actions/setup-go@v4
      with:
        go-version: '1.21'

    - name: git-changes check
      run: |
        go run github.com/rusinikita/changes/cmd/changes check --output=git-check.md

    - name: post result comment
      if: always() # post even if previous step failed
      uses: thollander/actions-comment-pull-request@v2
      with:
        filePath: git-check.md
        comment_tag: git-check
        mode: recreate
        create_if_not_exists: true
```

## Docker

Create dockerfile with this content.

```dockerfile
FROM golang:1.21

WORKDIR /app

RUN go install github.com/rusinikita/changes/cmd/changes@570e84e
LABEL authors="github.com/rusinikita"

CMD ["changes", "check"]
```

Build and run image with `-v .:/app`

## GitLab runner or other software

Use [binary](#binaries) or [docker](#docker) setup for your runner/worker.