---
sidebar_position: 3
---

# Learn the `check` command
Quickly check your stack for end-of-life (EOL) statuses. The `geol check` command scans products from your `.geol.yaml` and reports support status and EOL dates.

## Get help for `check`

Use `geol help check` to display help and available options for the `check` command.

```bash
geol help check
```

## Initialize a check file

Run the command to create a template `.geol.yaml` in the current directory:

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
...
```

## Statuses and warnings

Run the check to view statuses and warnings:

```shell
geol check
```

Use this flag to make `geol check` return a non-zero exit code when any product is past its EOL.
```bash
--strict
```