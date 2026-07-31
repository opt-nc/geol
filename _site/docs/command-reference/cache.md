---
sidebar_position: 12
---

# 💾 cache

Manage the local **geol** cache.

## 🖥️ Usage

```bash
geol cache [command] [options]
```

## 📄 Description

The `cache` command manages the local cache stored by **geol**.

The cache file contains the list of products and aliases retrieved from the endoflife.date API. It is stored in the user's configuration directory.

Typical location:

```text
geol/products.json
```

Use this command to:

- Refresh the local cache
- Remove the cache file
- View cache information

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
geol cache refresh --log-level debug
```

## 📋 Available Subcommands

| Subcommand | Description |
|------------|-------------|
| `refresh` | Download the latest products and aliases from endoflife.date |
| `status` | Display information about the local cache |
| `clear` | Delete the local cache file |

## 🔄 Refresh the Cache

Download the latest list of products and aliases from the endoflife.date API and store it locally.

### Usage

```bash
geol cache refresh
```

### Example

```bash
geol cache refresh
```

:::tip

Run this command if you want to make sure you are using the latest product data.

:::

## 📊 Display Cache Status

Display information about the local cache file.

### Usage

```bash
geol cache status
```

### Example

```bash
geol cache status
```

This command can be used to verify whether a cache file exists and check its current status.

## 🗑️ Clear the Cache

Delete the locally cached products file.

### Usage

```bash
geol cache clear
```

### Example

```bash
geol cache clear
```

:::warning

After clearing the cache, **geol** may need to download product information again the next time cache data is required.

:::

## 💡 Common Use Cases

Use this command to:

- Refresh local product data
- Troubleshoot cache-related issues
- Remove outdated cache files
- Verify cache availability