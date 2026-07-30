---
sidebar_position: 8
---

# 📦 product

Display lifecycle information for one or more products.

## 🖥️ Usage

```bash
geol product [command] [options]
```

## 📄 Description

The `product` command provides access to lifecycle information for software and hardware products available in the endoflife.date dataset.

Use this command to retrieve:
 
- Product details
- Release information
- Lifecycle metadata
- End-of-life dates

Depending on your needs, you can display either a summary or detailed release information.

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
geol product describe ubuntu --log-level debug
```

## 📋 Available Subcommands

| Subcommand | Description |
|------------|-------------|
| `describe` | Display a product summary |
| `extended` | Display detailed release information |

## 💡 Examples

Display lifecycle information for one or more products:

```bash
geol product linux ubuntu
```

Display detailed release information:

```bash
geol product extended golang k8s
```

Display a product summary:

```bash
geol product describe nodejs
```

## 📚 Subcommand Documentation

The following pages provide detailed documentation for each subcommand:
 
- 📄 **describe** - Display a product summary
- 📋 **extended** - Display detailed release information

## ✅ Common Use Cases

Use this command to:

- View lifecycle information for a product
- Check the latest supported versions
- Identify end-of-life dates
- Review release history
- Prepare upgrade and migration plans