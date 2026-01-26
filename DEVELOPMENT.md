# Development Guide

## Table of Contents

- [Development Guide](#development-guide)
  - [Table of Contents](#table-of-contents)
  - [Contribution Workflow](#contribution-workflow)
    - [Conventional Commits](#conventional-commits)
      - [Format](#format)
      - [Types](#types)
      - [Scope](#scope)
      - [Examples](#examples)
    - [Pull Request Squashing](#pull-request-squashing)
    - [Review Checklist](#review-checklist)


## Contribution Workflow

### Conventional Commits

All pull request titles **must** follow the [Conventional Commits](https://www.conventionalcommits.org/) specification to maintain a clear and organized commit history.

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

### Pull Request Squashing

All pull requests **must be squashed** before merging to `main`. This keeps the commit history clean and aligns with the conventional commit history.

- Use GitHub's **"Squash and merge"** option when merging
- Ensure the squashed commit message follows the Conventional Commits format
- One logical change = one squashed commit

### Review Checklist

Before opening a PR:

- [ ] Code follows Go idioms and project conventions
- [ ] Tests added/updated for new functionality
- [ ] No Hungarian notation or inconsistent naming
- [ ] Error messages are wrapped with context using `fmt.Errorf(...%w...)`
- [ ] Godoc comments added for exported types/functions
- [ ] PR title follows Conventional Commits format