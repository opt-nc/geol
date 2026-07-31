---
sidebar_position: 11
---

# ✅ check

Analyze the lifecycle status of the software products defined in a stack file.

## 🖥️ Usage

```bash
geol check [command] [options]
```

## 📄 Description

The `check` command analyzes each software component listed in a stack configuration file and retrieves lifecycle information from endoflife.date.

The command generates a report showing:

- Supported products
- Products approaching end-of-life
- End-of-life products
- Lifecycle dates and support status

By default, **geol** reads the `.geol.yaml` file located in the current directory.

This command is useful for:

- Technology lifecycle monitoring
- Upgrade planning
- Compliance checks
- CI/CD validation
- Risk assessment

## ⚙️ Global Options

| Option | Description |
|----------|-------------|
| `-d, --date` | Reference date for lifecycle calculations |
| `-f, --file` | Stack file to analyze |
| `--json` | Output results in JSON format |
| `-l, --log-level` | Logging level (`debug`, `info`, `warn`, `error`) |
| `-s, --strict` | Exit with an error if any product is EOL |

## 💡 Examples

Analyze the default stack file:

```bash
geol check
```

Analyze a specific stack file:

```bash
geol check --file stack.yaml
```

Generate JSON output:

```bash
geol check --json
```

Analyze the stack at a specific date:

```bash
geol check --date 2028-01-01
```

Enable strict mode:

```bash
geol check --strict
```

## 📁 Stack File

By default, **geol** reads:

```text
.geol.yaml
```

A typical configuration file contains a list of products and versions to monitor.

Example:

```yaml
stack:
  - name: ubuntu
    version: "24.04"
    id_eol: ubuntu

  - name: eclipse temurin
    version: "21"
    id_eol: eclipse-temurin
```

:::tip

Use `geol check init` to generate a valid configuration template.

:::

## 🛠️ Available Subcommands

| Subcommand | Description |
|------------|-------------|
| `init` | Generate a template configuration file |

### Generate a Template File

Create a sample `.geol.yaml` file:

```bash
geol check init
```

This file can then be customized to match your software stack.

## 🚨 Use Strict Mode

Strict mode is particularly useful in CI/CD pipelines.

```bash
geol check --strict
```

When enabled, **geol** returns a non-zero exit code if at least one product has reached its end-of-life date.

This allows automated workflows to detect unsupported software and fail deployment checks when necessary.

## 📊 Export Results as JSON

Generate machine-readable output:

```bash
geol check --json
```

This format is useful when integrating **geol** with:

- CI/CD pipelines
- Monitoring systems
- Reporting tools
- Automated scripts

## 📅 Check a Specific Date

By default, lifecycle calculations use the current date.

Use the `--date` option to analyze the stack at a specific point in time:

```bash
geol check --date 2028-01-01
```

This is useful for:

- Migration planning
- Future impact analysis
- Historical validation

## ✅ Common Use Cases

Use this command to:

- Identify unsupported software
- Track upcoming end-of-life dates
- Prepare migration projects
- Validate software stacks before deployment
- Generate lifecycle compliance reports