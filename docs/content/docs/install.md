---
weight: 2
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
brew install rusinikita/tap/changes
```

## Binaries

###### One line installer
```shell
curl -sSfL https://raw.githubusercontent.com/rusinikita/changes-action/master/changes-binary-install.sh | sh -s -- -b ~/bin latest
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
wget -O- -nv https://raw.githubusercontent.com/rusinikita/changes-action/master/changes-binary-install.sh | sh -s -- -b ~/bin latest
```

## GitHub actions

```yaml

steps:
  - uses: rusinikita/changes-action@v0.7
```

`Changes` has [github-action](https://github.com/rusinikita/changes-action) for fastest execution capable of:
- Binary installation
- Caching binary to skip installation in next runs
- Running `check` command
- Posting `check` result message in pull request

#### Full workflow

Create `.github/wokrflows/changes.yml` file and fill it with following content.

```yaml

name: changes verification
on:
  pull_request:

permissions:
  contents: read
  pull-requests: write # important for message posting, removable

jobs:
  lint:
    name: changes-verification
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          # skip pull request merge commit, important if "no merge commits" validation enabled, removable
          ref: ${{ github.event.pull_request.head.sha }}
          # fetch all branches and commits, important for git diff retrieving
          fetch-depth: 0
      - uses: rusinikita/changes-action@v1
        with:
          version: v0.3.2 # version of changes cli, default latest   
          pr-message: true # enables message posting, default false
          config: .changes.yaml # config file path, default is .changes.[yaml,yml,toml,json]
```

{{< alert context="info" >}}
{{< markdownify >}}

Using `latest` version slows action run a little. 
It enables network request for resolving latest release. 
Resolved version binary cashed until new release detected.

{{< /markdownify >}}
{{< /alert >}}

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
