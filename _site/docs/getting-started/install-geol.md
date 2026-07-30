---
sidebar_position: 1
---

# 📥 Install geol

Install **geol** and verify that it is working correctly.

**geol** is a command-line tool that helps you monitor the end-of-life (EOL) status of software products and application stacks.

## ⚡ Install with Homebrew

:::info

Homebrew is the recommended installation method.

:::

Install **geol**:

```bash
brew install curl
brew install --cask opt-nc/homebrew-tap/geol
```

:::warning

Homebrew requires its own version of `curl` to download packages from GitHub.

Skipping the `brew install curl` step may cause the installation to fail.

:::

### macOS Security Notice

:::info

Because **geol** is not signed with an Apple Developer certificate, macOS may block its execution the first time you run it.

:::

Allow **geol** to run from:

```text
System Settings → Privacy & Security
```

or remove the quarantine flag from the terminal:

```bash
xattr -d com.apple.quarantine $(which geol)
```

## 🌐 Install with the Installation Script

Install the latest version using the official installation script:

```bash
curl -fsSL https://raw.githubusercontent.com/opt-nc/geol/main/install.sh | bash
```

## 🐹 Install with Go

Install **geol** using Go:

```bash
go install github.com/opt-nc/geol/v2@latest
```

:::info

This method requires a working Go environment.

:::

## 📦 Manual Installation

You can also download a prebuilt binary from the GitHub releases page.

[Download geol Releases](https://github.com/opt-nc/geol/releases matches your operating system and architecture.

## 🔄 Update geol

If you installed **geol** with Homebrew, update it with:

```bash
brew update && brew upgrade --cask
```

## 🏷️ Check the Installed Version

Verify that **geol** is installed correctly:

```bash
geol version
```

The command should display the installed version.

## ✅ Verify the Installation

Run:

```bash
geol version
```

If the version is displayed correctly, **geol** is ready to use.

Continue with the next page to learn the most common commands and start exploring available products.