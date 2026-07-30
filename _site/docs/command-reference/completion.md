---
sidebar_position: 14
---

# ⌨️ completion

Generate shell autocompletion scripts for **geol**.

## 🖥️ Usage

```bash
geol completion [shell] [options]
```

## 📄 Description

The `completion` command generates shell completion scripts for **geol**.

Shell completion improves the command-line experience by providing:

- Command suggestions
- Option completion
- Faster command entry
- Reduced typing errors

Generated scripts can be loaded into your shell to enable tab completion.

## ⚙️ Global Options

The following option is available for this command:

```text
-l, --log-level
```

Supported values:

- `debug`
- `info` (default)
- `warn`
- `error`

Example:

```bash
geol completion bash --log-level debug
```

## 📋 Available Subcommands

| Subcommand | Description |
|------------|-------------|
| `bash` | Generate a Bash completion script |
| `fish` | Generate a Fish completion script |
| `powershell` | Generate a PowerShell completion script |
| `zsh` | Generate a Zsh completion script |

## 🐚 Generate Bash Completion

Generate an autocompletion script for Bash:

```bash
geol completion bash
```

## 🐟 Generate Fish Completion

Generate an autocompletion script for Fish:

```bash
geol completion fish
```

## ⚡ Generate PowerShell Completion

Generate an autocompletion script for PowerShell:

```bash
geol completion powershell
```

## 💎 Generate Zsh Completion

Generate an autocompletion script for Zsh:

```bash
geol completion zsh
```

## 💡 Example

Display a Zsh completion script:

```bash
geol completion zsh
```

The generated output can then be saved and loaded by your shell according to your local configuration.

:::tip

Run the command for your preferred shell and follow the shell-specific instructions provided by the generated output.

:::

## ✅ Common Use Cases

Use this command to:

- Enable shell autocompletion
- Improve CLI productivity
- Reduce typing errors
- Discover commands and options more easily