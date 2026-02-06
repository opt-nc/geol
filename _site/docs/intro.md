---
sidebar_position: 1
---

# Getting Started

Let's discover **geol in less than 5 minutes**.

Get started by **installing with <code>brew</code>**.

Note: Homebrew often provides more up-to-date packages than other sources, so installing via `brew` will typically give you a newer version.

### How to install?

- <code>brew</code> installed on your machine. See the official <a href="https://brew.sh/" target="_blank" rel="noreferrer noopener">Homebrew website</a> for installation instructions.

```bash
brew install curl
brew install --cask opt-nc/homebrew-tap/geol
```

If your installed `geol` version is out of date, you can update it (Homebrew) with:

```bash
# Update Homebrew and upgrade the geol cask
brew update && brew upgrade --cask geol
```

## Get `geol` version
```bash
geol version
```

## Know a bit more about `geol`
Displays a brief description and information about `geol`.
```bash
geol about
```

## Get help
Displays quick help and the main commands.
```bash
geol help
```

Shows the full man page: all commands and detailed help (includes the information from `geol help`).
```bash
man geol
```