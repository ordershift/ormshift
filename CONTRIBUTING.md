# Contributing to ormshift

## Table of Contents

- [Contributing to ormshift](#contributing-to-ormshift)
  - [Table of Contents](#table-of-contents)
  - [How to Contribute](#how-to-contribute)
    - [Commit Message Guidelines](#commit-message-guidelines)
      - [Format](#format)
      - [Types](#types)
      - [Scope](#scope)
      - [Examples](#examples)
    - [Merging Pull Requests](#merging-pull-requests)
    - [Pre-Submission Checklist](#pre-submission-checklist)


## How to Contribute

### Commit Message Guidelines

All pull request titles **must** follow the [Conventional Commits](https://www.conventionalcommits.org/) specification. Individual commits within your PR can use any format you prefer; only the PR title matters, as it becomes the commit message on `main` after squashing.

#### Format

```
<type>(<scope>): <subject>
```


#### Types

- **feat**: A new feature
- **fix**: A bug fix
- **docs**: Documentation changes
- **refactor**: Code refactoring without feature changes or bug fixes
- **test**: Adding or updating tests
- **chore**: Build, CI, dependencies, or tooling changes
- **perf**: Performance improvements

#### Scope

Optional but recommended. Use the affected component:
- `migrations` - Migration system changes
- `schema` - Schema package changes
- `builder` - SQL builder changes
- `executor` - Executor changes
- `dialects/postgresql` - PostgreSQL dialect
- `dialects/sqlite` - SQLite dialect
- `dialects/sqlserver` - SQL Server dialect

#### Examples

- `feat(migrations): add transaction support for Up/Down operations`
- `fix(schema): validate table names case-insensitively`
- `docs(README): add foreign key examples`
- `test(dialects/postgresql): add boolean conversion tests`
- `refactor: remove hungarian notation from variable names`

### Merging Pull Requests

All pull requests **must be squashed** before merging to `main`. This ensures that only the PR title (following Conventional Commits) becomes a commit on the main branch.

- Use GitHub's **"Squash and merge"** option when merging
- The PR title becomes the commit message on `main`
- One logical change = one squashed commit

### Pre-Submission Checklist

Before opening a PR:

- [ ] Code follows Go idioms and project conventions
- [ ] Tests added/updated for new functionality
- [ ] No Hungarian notation or inconsistent naming
- [ ] Error messages are wrapped with context using `fmt.Errorf(...%w...)`
- [ ] Godoc comments added for exported types/functions
- [ ] PR title follows Conventional Commits format
