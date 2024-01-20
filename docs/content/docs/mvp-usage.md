---
weight: 999
title: "Mvp Usage"
description: ""
icon: "article"
date: "2024-01-14T18:38:29+08:00"
lastmod: "2024-01-14T18:38:29+08:00"
toc: true
---

{{< alert context="info" >}}
Not complete and good documentation. Just reminders for first adopters.
Quick support and communication in project issues.
{{< /alert >}}

## Project folder

- Create .changes.yaml (or toml) file inside your repository
- Setup rules

## GitLab runner

Install `changes` on to runner machine.

```
go install github.com/rusinikita/changes/cmd/changes@570e84e
```

Add cli call in .gitlab-ci.yaml

```
changes check
```

## GitHub

- Use `actions/checkout` and `actions/setup-go` to create environment
- Call `go run github.com/rusinikita/changes/cmd/changes@570e84e check`

```yaml
git-check:
  name: git-checks
  runs-on: ubuntu-latest
  steps:
    - uses: actions/checkout@v4
      with:
        ref: ${{ github.event.pull_request.head.sha }}
        fetch-depth: 0

    - uses: actions/setup-go@v4
      with:
        go-version: '1.21'

    - name: git-changes_check
      run: |
        go run github.com/rusinikita/changes/cmd/changes@570e84e check
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
