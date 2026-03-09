# gh-action-scripts

Opinionated, reusable GitHub Actions CI fragments for use across repos.

## Assumptions

All consuming repos must follow these conventions:

- **[mise](https://mise.jdx.dev/)** — tool version management (`mise.toml` in repo root)
- **[Taskfile](https://taskfile.dev/)** — task runner (`Taskfile.yml` in repo root)

## Reusable Workflows

All reusable workflows live in `.github/workflows/shared/` and are triggered via `workflow_call`.

### `shared/ci.yaml` — CI Pipeline

A standard CI workflow covering setup, codegen, formatting, lint, build, tests, and vulnerability scanning.

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
    uses: peteresztari/gh-action-scripts/.github/workflows/shared/ci.yaml@main
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
| `build`     | Build the project                            |
| `test:ci`   | Run tests, output JUnit XML to `reports/junit.xml` |
| `vuln:ci`   | Run vulnerability scan                       |

The workflow enforces that `generate` and `format` produce no uncommitted changes.

### `shared/upgrade-tools.yaml` — Tool Version Upgrades

Runs `mise upgrade`, and if `mise.toml` changed, creates a PR with the updated versions.

**Reference it from your repo:**

```yaml
# .github/workflows/upgrade-tools.yaml
name: Upgrade Tools

on:
  schedule:
    - cron: "0 6 * * 1"
  workflow_dispatch:

jobs:
  upgrade:
    uses: peteresztari/gh-action-scripts/.github/workflows/shared/upgrade-tools.yaml@main
```

Requires `contents: write` and `pull-requests: write` permissions (set in the reusable workflow).

## Go Defaults

This repo provides ready-to-use `mise.toml` and `Taskfile.yml` for Go projects. Copy them into your repo root — no extension or inclusion needed.

```sh
# From your repo root
curl -O https://raw.githubusercontent.com/peteresztari/gh-action-scripts/main/mise.toml
curl -O https://raw.githubusercontent.com/peteresztari/gh-action-scripts/main/Taskfile.yml
```

**`mise.toml`** — pinned tool versions:

| Tool             | Description          |
|------------------|----------------------|
| `go`             | Go compiler          |
| `golangci-lint`  | Linter               |
| `gotestsum`      | Test runner (JUnit)  |
| `osv-scanner`    | Vulnerability scanner|
| `gofumpt`        | Formatter            |
| `task`           | Taskfile runner      |

**`Taskfile.yml`** — all required CI tasks plus a local `ci` task:

| Task        | Description                                  |
|-------------|----------------------------------------------|
| `setup`     | `go mod download` + `go mod verify`          |
| `build`     | `go build ./...`                             |
| `generate`  | `go generate ./...`                          |
| `format`    | `gofumpt -l -w .`                            |
| `lint`      | `golangci-lint run ./...`                    |
| `test:ci`   | `gotestsum` with JUnit XML to `reports/`     |
| `vuln:ci`   | `osv-scanner scan -r .`                      |
| `ci`        | Runs the full pipeline locally               |

Adjust tool versions in `mise.toml` and task commands in `Taskfile.yml` to fit your project.
