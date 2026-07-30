---
sidebar_position: 3
---

# 🚀 Use geol in GitHub Actions

Integrate **geol** into your GitHub Actions workflows to automatically monitor the end-of-life (EOL) status of your software stack.

With strict mode enabled, workflows can fail when unsupported products are detected.

## 🎬 Watch a Demo

See how **geol** integrates with GitHub Actions to provide automated lifecycle checks and EOL reporting.

<iframe
  width="100%"
  height="415"
  src="https://www.youtube.com/embed/0havqKL-Suo"
  title="Use geol in GitHub Actions"
  frameborder="0"
  allowfullscreen>
</iframe>

## 🎯 Why Use geol in CI/CD?

Using **geol** in your CI/CD pipeline helps you:

- 🔒 Detect unsupported software early
- 📋 Track product lifecycle status
- ⚡ Automate EOL monitoring
- 🚨 Alert teams when products reach end-of-life

## ⚙️ Install the GitHub Action

The easiest way to integrate **geol** is to use the official GitHub Action:

https://github.com/opt-nc/geol-action

The action downloads the **geol** binary and makes it available in your workflow.

## 📄 Example Workflow

Create a `.github/workflows/geol-check.yml` file in your repository.

:::tip

It is recommended to use a scheduled workflow in addition to push and pull request triggers.

Since EOL dates change independently from your source code, scheduled checks help detect newly unsupported products even when no code changes have been made.

:::

```yaml
name: Check EOL

on:
  push:
    branches: [ main ]

  pull_request:
    branches: [ main ]

  schedule:
    - cron: '0 0 * * 1'

jobs:
  check-eol:
    runs-on: ubuntu-latest

    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Install geol
        uses: opt-nc/geol-action@v1
        with:
          version: 'v2.12.1'

      - name: Check stack EOL
        run: |
          geol check --strict
```

## 🔍 How It Works

### 1. Install geol

The `opt-nc/geol-action` step downloads **geol** and adds it to the workflow `PATH`.

### 2. Run Lifecycle Checks

Once installed, **geol** can be used in any subsequent workflow step.

```bash
geol check --strict
```

### 3. Enable Strict Mode

The `--strict` option causes `geol check` to return a non-zero exit code when at least one product has reached its end-of-life date.

This allows GitHub Actions to fail the workflow and notify the team.

### 4. Schedule Regular Checks

Scheduled workflows ensure that lifecycle monitoring continues even when no new code is pushed.

## 📝 Create Issues Automatically

Instead of only failing the workflow, you can automatically create a GitHub issue when unsupported software is detected.

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

:::tip

Automatically creating issues helps track remediation work and ensures that EOL findings are visible in your backlog.

:::