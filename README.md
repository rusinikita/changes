![linter](https://img.shields.io/badge/golangci--lint-blue?logo=verified_user)
![coverage](https://img.shields.io/badge/test_coverage-89%25-green)
![mutation](https://img.shields.io/badge/mutation_coverage-90%25-green)

`Changes` helps automate code review and enforce team agreements.

It validates git diff with simple configuration and [CEL scripting language](https://github.com/google/cel-spec).

[Visit project GH pages site for detail information](https://rusinikita.github.io/changes/)

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

```shell
Changes
-------
docs/content/docs/ADR/24.02-output-message.md
- draft page will not be shown

Commits
-------
12341324: test
- type: value must be one of [fix ref feat]
```

### Independent of platform

You can use it locally, as a CI step and in combination with messages posting tools.

![An github PR comment](docs/assets/images/github-message.png)

[See working GitHub actions configuration](https://rusinikita.github.io/changes/docs/mvp-usage/#github)

## Alternatives

[Danger](https://danger.systems/js) and [conventional-changelog](https://github.com/conventional-changelog) - great
tools for JS and Ruby community.

`Changes` brings advantages for other languages users:

1. yaml/toml configuration familiarity and simplicity
2. One tool to rule. Set message format and use parsed values in diff validation and changelog generation
3. Well-designed configuration and scripting DSL
4. Go code maintainability and execution speed

> [!IMPORTANT]
> ##### Project in a beta test
> Tutorials, footer values parsing and changelog generation in development.
> Please provide your feedback, needs, and ideas in project issues and forum.
