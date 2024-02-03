---
weight: 1
title: "Overview"
description: ""
icon: "help"
date: "2023-11-19T21:16:54+02:00"
lastmod: "2023-11-19T21:16:54+02:00"
toc: true
---

## What

`Changes` helps automate code review and enforce team agreements.

It validates git diff with simple configuration and [CEL scripting language](https://github.com/google/cel-spec).

```yaml
commit:
  # Commit message format validation. Example - `TEAM-123: fix something`
  subject: '(issue): (type)? (title)'
  values:
    issue:
      regexp: TEAM-\d+
check-functions:
  - message: use separate PRs for fixes and features
    # Function checks parsed commit types
    func: 'commits.exists(c, c.type == "fix") && commits.exists(c, c.type == "feat")'
  - message: page contains draft flag and will not be shown
    # Function checks changed files and its content
    func: changes.filter(c,
      c.path.endsWith(".md") && c.chunks.exists(ch, ch.content.contains("draft:"))
      )
```

### Output

{{< alert text="In progress. Ugly CLI output. Make it better and add screenshot." />}}

### Independent of platform

{{< alert text="In progress. No warnings and markdown output. Screenshot and example setup for PR message posting" />}}

You can use it locally, as a CI step and in combination with messages posting tools.

## Alternatives

[Danger](https://danger.systems/js) and [conventional-changelog](https://github.com/conventional-changelog) - great
tools for JS and Ruby community.

`Changes` brings advantages for other languages users:

1. yaml/toml configuration familiarity and simplicity
2. One tool to rule. Set message format and use parsed values in diff validation and changelog generation
3. Well-designed configuration and scripting DSL
4. Best quality Go code maintainability and execution speed

## Current state

{{< alert context="warning" >}}
{{< markdownify >}}

##### Project in a beta test

Tutorials, footer values parsing and changelog generation in development.

Please provide your feedback, needs, and ideas in project issues and forum.

{{< /markdownify >}}
{{< /alert >}}

See [architectural decision records]({{< ref "/docs/adr" >}}) for more information.
