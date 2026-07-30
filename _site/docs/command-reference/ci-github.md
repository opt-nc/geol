---
sidebar_position: 15
---

# 🚀 ci-github

Generate GitHub Actions workflow files for **geol**.

## 🖥️ Usage

```bash
geol ci-github [command] [options]
```

## 📄 Description

The `ci-github` command generates GitHub Actions workflow files that can be used to automate lifecycle checks with **geol**.

By default, this command generates a ready-to-use GitHub Actions workflow.

This is equivalent to:

```bash
geol ci-github init
```

The generated workflow can be committed directly to a repository and used to monitor product end-of-life (EOL) status in CI/CD pipelines.

## ⚙️ Global Options

| Option | Description |
|----------|-------------|
| `-f, --force` | Overwrite the output file if it already exists |
| `-o, --output` | Path to the generated workflow file |
| `-l, --log-level` | Logging level (`debug`, `info`, `warn`, `error`) |

## 📋 Available Subcommands

| Subcommand | Description |
|------------|-------------|
| `init` | Generate a GitHub Actions workflow file |

## 📄 Generate a Workflow

Create a ready-to-use GitHub Actions workflow:

```bash
geol ci-github init
```

This command generates a workflow file that can be added directly to a GitHub repository.

By default, the file is created as:

```text
.github/workflows/geol-action.yml
```

## 📁 Custom Output File

Generate the workflow in a custom location:

```bash
geol ci-github init --output .github/workflows/eol-check.yml
```

:::tip

Use a custom file name if your repository already contains other workflow definitions.

:::

## 🔄 Overwrite an Existing File

If the output file already exists, use the `--force` option to replace it.

```bash
geol ci-github init --force
```

Example:

```bash
geol ci-github init \
  --output .github/workflows/geol-action.yml \
  --force
```

:::warning

Existing files will be overwritten when using the `--force` option.

:::

## 💡 Examples

Generate the default workflow:

```bash
geol ci-github
```

Generate the workflow explicitly:

```bash
geol ci-github init
```

Generate a workflow in a custom file:

```bash
geol ci-github init \
  --output .github/workflows/eol-check.yml
```

Overwrite an existing workflow:

```bash
geol ci-github init --force
```

## ✅ Common Use Cases

Use this command to:

- Generate a GitHub Actions workflow
- Automate lifecycle checks
- Monitor software end-of-life dates
- Integrate **geol** into CI/CD pipelines
- Quickly bootstrap GitHub-based lifecycle monitoring

## 📚 Related Documentation

For a complete GitHub Actions integration example, see:

- 🚀 Use geol in GitHub Actions