---
weight: 12
title: "Functions"
description: ""
icon: "article"
date: "2024-01-30T00:15:10+07:00"
lastmod: "2024-01-30T00:15:10+07:00"
toc: true
---

One stop page for `check-function` writing.

## Examples

Use separate PRs for fixes and features. Function checks parsed commit types.

```js
commits.exists(c, c.type == "fix") && commits.exists(c, c.type == "feat")
```

Page contains draft flag and not shown. Function checks changed files and its content.

```js
changes.filter(c,
    c.path.endsWith(".md")
    && c.chunks.exists(ch, ch.content.contains("draft:"))
)
```

Too big file change, more than 250 added lines. Using change stats.

```js
changes.filter(change, change.path.endsWith(".go") && change.stats().additions > 250)
```

Commit author email.

```js
commits.filter(commit, !commit.author.email.endsWith("@gmail.com"))
```

## Values

Top level variables accessible from CEL function.

{{< table >}}

| Name      | Meaning                                                                       | Type        |
|-----------|-------------------------------------------------------------------------------|-------------|
| `commits` | Diff commits. Allows to review commit types, author, messages and other data. | Commit list |
| `changes` | Diff file changes. Allows to review files content, new or deleted lines.      | Change list |
| `now`     | Current time                                                                  | timestamp   |

{{< /table >}}

## Types

### Language types

`bool`, `int`, `double`, `string`, `duration`, `timestamp`

### Containers

`list`, `map`

### App types

#### Commit

Git commit data

```go
// Commit data
type Commit struct {
// Message is the commit message, contains arbitrary text.
message string
// Author is the original author of the commit.
author Signature
// Committer is the one performing the commit, might be different from Author.
committer Signature
// ParentsCount is count of the parent commits of the commit.
parentsCount int
// Commit values
type string
title string
context string
issue string
// ...
}
```

{{< alert context="info" >}}
{{< markdownify >}}

[Commit values]({{< ref "/docs/configuration/commit#values" >}}) accessible by their names.

{{< /markdownify >}}
{{< /alert >}}

#### Signature

Commit author and committer data

```go
// Signature is used to identify who and when created a commit or tag.
type Signature struct {
// Name represents a person name. It is an arbitrary string.
name string
// Email is an email, but it cannot be assumed to be well-formed.
email string
// When is the timestamp of the signature.
when time.Time
}
```

#### Change

```go
// Change represents the necessary steps to transform one file into another.
type Change struct {
// Path is file path after change
// empty if file deleted
path string
// PrevPath is file path before change
// empty if file created
// diff from Path if file renamed
prevPath string
// Chunks returns a slice of ordered changes to transform "from" File into
// "to" File. If the file is a binary one, Chunks will be empty.
chunks []Chunk
}
```

#### Stat

Result of calling `сhange.stats()`

```go
type Stat struct {
additions, deletions int
}
```

#### Chunk

```go
// Chunk represents a portion of a file transformation into another.
type Chunk struct {
// Content contains the portion of the file.
type Operation
// Type contains the Operation to do with this Chunk.
content string
}

// Operation: Equal, Add, Delete
```

{{<alert text="In progress. Chunk type field adaptation to script." />}}

## Operators

{{< table >}}

| Operator                                                   | Meaning                   | Example                             |
|------------------------------------------------------------|---------------------------|-------------------------------------|
| `&&`, `\|\|`, `!`, `==`, `<`, `>`, `-`, `+`, `/`, `%`, `*` | Common logic operators    | `x / 2`, `a <= b`                   |
| `?:`                                                       | Conditional evaluation    | `y != 0 ? x/y : 0`                  |
| `[]`                                                       | Get child by index of key | `someList[0]`, `someMap['keyName']` |

{{< /table >}}

## Functions

### Common for most types

{{< table >}}

| Function | Applied to                               | Example                              |
|----------|------------------------------------------|--------------------------------------|
| `size`   | `string`, `list`, `map`                  | `size("test123")`, `size([1, 2, 3])` |
| `string` | `int`, `double`, `duration`, `timestamp` | `string(123)`                        |
| `int`    | `string`, `double`, `timestamp`          | `int("123")`                         |
| `double` | `int`, `string`                          | `double(123)`                        |

{{< /table >}}

### List and map tests

{{< table >}}

| Function           | Predicate matches                           | Example                           |
|--------------------|---------------------------------------------|-----------------------------------|
| `all(x, p)`        | **all** `list` elements or `map` keys       | `[3, 6, 9].all(n, (n % 3) == 0)`  |
| `exists(x, p)`     | **any** `list` element or `map` key         | `[3, 3, 3].exists(n, n == 3)`     |
| `exists_one(x, p)` | **exactly one** `list` element or `map` key | `[3, 6, 9].exists_one(n, n == 3)` |

{{< /table >}}

### List transformations

{{< table >}}

| Function       | Returns                                             | Example                              |
|----------------|-----------------------------------------------------|--------------------------------------|
| `filter(x, p)` | Sub-list of elements                                | `[1, 2, 3].filter(i, i % 2 > 0)`     |
| `map(x, t)`    | New list with transformed elements                  | `[1, 2, 3].map(n, n * n)`            |
| `map(x, p, t)` | New sub-list with filtered and transformed elements | `[1, 2, 3].map(n, n % 2 > 0, n * n)` |

{{< /table >}}

### Strings

{{< table >}}

| Function     | Returns   | Meaning                                                  | Example                                                  |
|--------------|-----------|----------------------------------------------------------|----------------------------------------------------------|
| `startsWith` | bool      | Has prefix                                               | `"test123".startsWith("test")`                           |
| `endsWith`   | bool      | Has suffix                                               | `"123test".endsWith("test")`                             |
| `contains`   | bool      | Has sub-string                                           | `"123test".contains("test")`                             |
| `matches`    | bool      | Has regexp matches                                       | `"123test".matches("^d+w+$")`, `"123test".matches("w+")` |
| `duration`   | duration  | Convert to duration. Supported suffixes: h, m, s, ms, ns | `duration("10h20m30s40us50ns")`                          |
| `timestamp`  | timestamp | Convert to time. Format RFC3339                          | `timestamp("1972-01-01T10:00:20.021-05:00")`             |

{{< /table >}}

### Timestamp and duration

{{< table >}}

| Function          | Returns     | Applied to               |
|-------------------|-------------|--------------------------|
| `getMilliseconds` | 0-999       | `timestamp`, `duration`  |
| `getSeconds`      | 0-59        | `timestamp`, `duration`  |
| `getMinutes`      | 0-59        | `timestamp`, `duration`  |
| `getHours`        | 0-23        | `timestamp`, `duration`  |
| `getDayOfWeek`    | 0-6         | `timestamp`              |
| `getDayOfMonth`   | 0-31        | `timestamp`              |
| `getMonth`        | 0-11        | `timestamp`              |
| `getFullYear`     | 0-10000     | `timestamp`              |
| `getDayOfYear`    | 0-365       | `timestamp`              |
| `-`, `+`          | `timestamp` | `timestamp +/- duration` |
| `-`, `+`          | `duration`  | `duration +/- duration`  |

{{< /table >}}
