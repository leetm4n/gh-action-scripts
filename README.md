# gh-action-scripts

Opinionated, reusable GitHub Actions CI fragments for use across repos.

## Assumptions

All consuming repos must follow these conventions:

- **[mise](https://mise.jdx.dev/)** — tool version management (`.mise.toml` in repo root)
- **[Taskfile](https://taskfile.dev/)** — task runner (`Taskfile.yml` in repo root)

## Available Workflows

### `ci.yaml` — Reusable CI Pipeline

A standard CI workflow covering setup, codegen, formatting, lint, and tests.

**Reference it from your repo:**

```yaml
# .github/workflows/ci.yaml
name: CI

on:
  push:
    branches: [main]
  pull_request:

jobs:
  ci:
    uses: peteresztari/gh-action-scripts/.github/workflows/ci.yaml@main
    with:
      upload-artifacts: true   # optional, default: false
```

**Inputs:**

| Input              | Type    | Default | Description                                      |
|--------------------|---------|---------|--------------------------------------------------|
| `upload-artifacts` | boolean | `false` | Upload JUnit test results to `reports/junit.xml` |

**Required Taskfile tasks:**

| Task        | Description                                  |
|-------------|----------------------------------------------|
| `setup`     | Install project dependencies                 |
| `generate`  | Run code generation                          |
| `format`    | Format code (must be idempotent)             |
| `lint`      | Run linters                                  |
| `test:ci`   | Run tests, output JUnit XML to `reports/junit.xml` |
| `vuln:ci`   | Run vulnerability scan                       |

The workflow enforces that `generate` and `format` produce no uncommitted changes.

## Default Taskfiles

Pre-built task implementations you can include in your repo's `Taskfile.yml`:

### Go — `taskfiles/go.yml`

Requires these tools in `.mise.toml`: `golangci-lint`, `gotestsum`, `osv-scanner`, `gofumpt`, `task`.

> **Note:** Unlike Taskfile, mise does not support remote config inheritance. Copy `mise/go.toml` from this repo as your `.mise.toml` starting point, then pin versions as needed.

```yaml
# Taskfile.yml
includes:
  go:
    taskfile: https://raw.githubusercontent.com/peteresztari/gh-action-scripts/main/taskfiles/go.yml
    optional: true

tasks:
  setup:
    cmds:
      - task: go:setup
      - npm install  # add repo-specific steps
```

Or use the tasks directly without renaming by flattening includes:

```yaml
includes:
  default:
    taskfile: https://raw.githubusercontent.com/peteresztari/gh-action-scripts/main/taskfiles/go.yml
```

This exposes `task setup`, `task generate`, `task format`, `task lint`, `task test:ci`, `task vuln:ci` directly. Override any task in your own `Taskfile.yml` to replace the default.
