---
weight: 11
title: "Commit"
description: ""
icon: "article"
date: "2024-01-27T18:34:45+08:00"
lastmod: "2024-01-27T18:34:45+08:00"
toc: true
---

Commit messages validation configuration.

{{< alert context="info" >}}
{{< markdownify >}}

##### Optional

`changes check` uses conventional commit by default.

Empty `changes.yaml` or no file - **valid**.
{{< /markdownify >}}
{{< /alert >}}

```yaml

commit:
  # Commit message format for subject validation  
  subject: '(issue): (type)? (title)'
  # values extracted from commit message
  values:
    issue:
      regexp: ^TEAM-\d+$

```

## Subject format

`commit.subject` - simplified regexp for subject validation.

### Examples

- `(!)?(type)((context))?: (title)` - conventional commit with optional context and `!`
  - `!feat: something`
  - `fix(ci): something`
- `(issue): (title)` - simple Jira/Basecamp/GitHub issue, `:` and title
  - `CMD-123: something`
  - `#123: something`
- `(issue): (type)? (title)` - previous format with optional `type` keyword
  - `#123: something`
  - `CMD-123: Fix something`

### Notation

1. `()` brackets - declaration for a value symbols group. Name inside brackets - value name.
2. `?` question symbol - group becomes optional.
3. `(([(^_approved_^)]))` - symbols inside brackets except letters treated as not value symbols, other brackets too.

[Format decision record](docs/ADR/23.12-format-template.md).

## Values

Values validation customization. You can add new values or override existed.

### Goals

1. Split giant regexp into parts for simplicity.
2. Commit value names used in custom validation functions and changelog generator.

### Defaults

```yaml

values:
  type:
    allowed: [ feat, fix ]
  title:
    regexp: "^[\w]+$"
    max-len: 20
  context:
    regexp: "^[a-z-_]+$"
    max-len: 10
  issue:
    regexp: "^#\d+$"
    max-len: 5
```

### Structure

- `values.{{value-name}}.alloved` - list of allowed strings for this value
- `values.{{value-name}}.max-len` - length limit
- `values.{{value-name}}.regexp` - value regular expression, `^` and `$` required to match full string
