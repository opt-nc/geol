---
sidebar_position: 1
---

# Getting Started

Get up and running with `geol` in less than 5 minutes.

`geol` is a command-line tool that helps you check the End-of-Life (EOL) status of your application stack, so you can keep your dependencies healthy and up to date.

:::info Prerequisite:
You need [Homebrew](https://brew.sh/) installed on your machine.
:::

## 📦 Install `geol`

Homebrew requires its own version of `curl` to download casks from GitHub — even if `curl` is already available on your system.
```bash
brew install curl
```
:::warning
Skipping this step will cause the next command to fail with a *"Homebrew-installed `curl` is not installed"* error.
:::

Installs `geol` on your machine using Homebrew.
```bash
brew install --cask opt-nc/homebrew-tap/geol
```
:::tip
Homebrew usually ships more up-to-date packages than system package managers, which is why we recommend it as the default install method.
:::

## 🔄 Update `geol`
Updates Homebrew and upgrades `geol` to the latest version.
```bash
brew update && brew upgrade --cask geol
```

## 🏷️ Get `geol` version
Displays the installed version of `geol`
```bash
geol version
```

## ℹ️ About `geol`
Displays a brief description and information about `geol`.
```bash
geol about
```

## ❓ Get help
Displays quick help and the main commands.
```bash
geol help
```

## 📖 Open the full manual
Shows all commands and detailed help (includes the information from `geol help`).
```bash
man geol
```