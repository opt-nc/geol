---
sidebar_position: 5
---

# 📋 list

Display cached products, categories, and tags.

## 🖥️ Usage

```bash
geol list [command] [options]
```

## 📄 Description

The `list` command displays the data currently available in the local **geol** cache.

This command can be used to discover:

- Available products
- Available categories
- Available tags

The displayed data comes from the local cache and is maintained by **geol**.

:::tip

If the data appears outdated, refresh the cache using:

```bash
geol cache refresh
```

:::

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
geol list products --log-level debug
```

## 📋 Available Subcommands

| Subcommand | Description |
|------------|-------------|
| `products` | List all cached product names |
| `categories` | List all cached category names |
| `tags` | List all cached tag names |

## 📦 List Products

Display all products currently available in the cache.

### Usage

```bash
geol list products
```

### Example

```bash
geol list products
```

Example output:

```text
ubuntu
debian
nginx
nodejs
postgresql
```

## 📂 List Categories

Display all available categories.

### Usage

```bash
geol list categories
```

### Example

```bash
geol list categories
```

Example output:

```text
os
database
runtime
cloud
language
```

## 🏷️ List Tags

Display all available tags.

### Usage

```bash
geol list tags
```

### Example

```bash
geol list tags
```

Example output:

```text
linux
java
cloud
database
container
```

## 💡 Common Use Cases

Use this command to:

- Discover available products
- Explore categories
- Explore tags
- Find values for the `category` command
- Find values for the `tag` command
- Explore the content of the local cache

## 📚 Related Commands

- `cache`
- `category`
- `tag`
- `product`