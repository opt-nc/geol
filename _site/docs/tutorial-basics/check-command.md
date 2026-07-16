---
sidebar_position: 3
---

# Learn the `check` command

Use the `geol check` command to identify products that are approaching or past their end-of-life (EOL) date.

`geol` reads your `.geol.yaml` file and reports support status, warnings, and EOL dates.
`


## ❓ Get help
Use `geol help check` to display help and available options for the `check` command.
```bash
geol help check
```

## ⚙️ Initialize a configuration file

:::note
Run the commands from the directory where you want to create your `.geol.yaml` file.
:::
```shell
geol check init
```

Edit the generated `.geol.yaml` to list the products you want to monitor.

Minimal example `.geol.yaml` (created by `geol check init`):
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

The example also contains a product marked with `skip: true`.
Skipped products are ignored during the analysis.


## 🔍 Check your stack

Run the check to view statuses and warnings:
```shell
geol check
```

## 🚨 Use strict mode

Use this flag to make `geol check` return a non-zero exit code when any product is past its EOL.
```bash
geol check --strict
```

Run the following command to print the exit status (a non-zero value indicates an error):
```bash
echo $?
```

## 📅 Check a stack at a specific date


By default, `geol check` analyzes your stack using today's date.

Use `-d` or `--date` to evaluate the support status of your products at a specific date.

```bash
geol check --date 2028-01-01
```
or
```bash
geol check -d 2028-01-01
```

:::info
The date can be in the future or in the past.
:::
