---
sidebar_position: 1
---

# Getting Started?

Let's discover **geol in less than 5 minutes**.

Get started by **installing with <code>brew</code>**.

Note: Homebrew often provides more up-to-date packages than other sources, so installing via `brew` will typically give you a newer version.

### How to install?

- <code>brew</code> installed on your machine. See the official <a href="https://brew.sh/" target="_blank" rel="noreferrer noopener">Homebrew website</a> for installation instructions.

```bash
brew install --cask opt-nc/homebrew-tap/geol
```
Note: If the `brew` installation fails, it may be because Homebrew's `curl` is not installed and your system is using the distribution `curl` (apt).

```bash
# Solution: install curl with brew then retry
brew install curl
brew install --cask opt-nc/homebrew-tap/geol
```



## Get `geol` version

```bash
geol version
```



## Know a bit more about `geol`

```bash
geol about
```

## Get help

```bash
geol help
man geol
```
