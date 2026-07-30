---
sidebar_position: 3
---

# ✅ Check Products

Use the `geol check` command to identify products that are approaching or past their end-of-life (EOL) date.

**geol** reads your `.geol.yaml` file and reports support status, warnings, and EOL dates for the products you monitor.

## ❓ Get Help

Display help and available options for the `check` command.

```bash
geol help check
```

## ⚙️ Initialize a Configuration File

:::note

Run the following command from the directory where you want to create your `.geol.yaml` file.

:::

```bash
geol check init
```

Edit the generated `.geol.yaml` file to define the products you want to monitor.

### Example Configuration

```yaml
stack:
  - name: ubuntu
    version: "25.10"
    id_eol: ubuntu

  - name: java temurin
    version: "21"
    id_eol: eclipse-temurin

  - name: non-existent-product
    version: "1.0"
    id_eol: non-existent-product
    skip: true
```

:::info

Products marked with `skip: true` are ignored during the analysis.

:::

## 🔍 Check Your Stack

Run the following command to analyze all products defined in your configuration file:

```bash
geol check
```

**geol** displays the support status, warnings, and end-of-life information for each product.

## 🚨 Use Strict Mode

Enable strict mode to return a non-zero exit code when at least one product has reached its end-of-life date.

```bash
geol check --strict
```

To display the command exit status:

```bash
echo $?
```

:::warning

A non-zero exit code indicates that at least one product is no longer supported.

:::

## 📅 Check a Stack at a Specific Date

By default, **geol** evaluates products using today's date.

Use the `-d` or `--date` option to analyze your stack at a specific point in time.

```bash
geol check --date 2028-01-01
```

or

```bash
geol check -d 2028-01-01
```

:::info

The specified date can be in the future or in the past.

:::