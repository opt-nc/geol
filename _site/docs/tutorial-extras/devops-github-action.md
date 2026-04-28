---
sidebar_position: 3
---

# DevOps with `geol` as a GitHub Action

Integrating `geol` into your GitHub Actions workflow allows you to automatically monitor the end-of-life (EOL) status of your stack and fail the build if any product is no longer supported.

## See it in action

Watch this video to see how `geol` integrates into GitHub Actions to provide clear EOL reporting and automated checks:

<div style={{position: 'relative', paddingBottom: '56.25%', height: 0, marginBottom: '2rem'}}>
  <iframe
    src="https://www.youtube.com/embed/0havqKL-Suo"
    title="geol as a GitHub Action"
    frameBorder="0"
    allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share"
    style={{position: 'absolute', top: 0, left: 0, width: '100%', height: '100%'}}
    allowFullScreen
  />
</div>

## Why use `geol` in CI/CD?

- **Security**: Ensure no unsupported (and potentially unpatched) software is used in your environment.
- **Compliance**: Maintain an up-to-date inventory of your stack's lifecycle.
- **Automation**: Get alerted immediately when a product reaches its end-of-life.

## Official GitHub Action

The easiest way to integrate `geol` is to use the official [geol-action](https://github.com/opt-nc/geol-action). It installs the binary and makes it available in your workflow path.

### Workflow Example

Create a file named `.github/workflows/geol-check.yml` in your repository. 

:::tip
It is highly recommended to include a **schedule** trigger (like the Monday morning example below). Since EOL dates are external events, your build should check for them even if you haven't pushed any new code.
:::

```yaml
name: Check EOL

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]
  schedule:
    - cron: '0 0 * * 1' # Run every Monday at midnight to catch new EOLs

jobs:
  check-eol:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Install geol
        uses: opt-nc/geol-action@v1
        with:
          version: 'v2.12.1' # Optional: specify a version or use 'latest'

      - name: Check stack EOL
        run: |
          geol check --strict
```

## How it works

1. **Installation**: The `opt-nc/geol-action` step downloads the specified version of `geol` and adds it to the `PATH`.
2. **Execution**: You can then call `geol` directly in any subsequent `run` step.
3. **Strict Mode**: Adding the `--strict` flag to `geol check` is essential for CI/CD. It ensures that the command returns a non-zero exit code if any product is past its EOL date, which effectively "fails" the GitHub Action and alerts your team.
4. **Schedule**: By running this weekly, you ensure that "stable" projects are still monitored for underlying software obsolescence.

## Go further

### Automate Issue Generation

You can go beyond failing the build by automatically opening a GitHub Issue when your stack reaches EOL. This ensures the task is tracked in your backlog.

You can use a combination of `geol check` (capturing output to a file) and an action like `peter-evans/create-issue-from-file`:

```yaml
      - name: Check stack EOL and save report
        id: geol_check
        continue-on-error: true
        run: |
          geol check --strict > eol-report.txt

      - name: Create Issue on EOL failure
        if: steps.geol_check.outcome == 'failure'
        uses: peter-evans/create-issue-from-file@v5
        with:
          title: "Critical: End-of-Life software detected in stack"
          content-filepath: eol-report.txt
          labels: |
            security
            obsolescence
```

For more details on how to configure your products and the `.geol.yaml` format, see the [Learn the check command](../tutorial-basics/check-command) tutorial.
