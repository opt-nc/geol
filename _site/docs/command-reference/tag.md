---
sidebar_position: 7
---

# 🏷️ tag

Display all products associated with a tag.

## 🖥️ Usage

```bash
geol tag [tag] [options]
```

## 📄 Description

The `tag` command displays all products associated with a specific tag.

The tag must exist in the local **geol** cache. Results are displayed in a tree structure to make product discovery easier.

Tags provide an alternative way to explore products beyond categories. They can be used to group products by vendor, technology, ecosystem, organization, or other shared characteristics.

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
geol tag canonical --log-level debug
```

## 💡 Examples

Display products associated with the `os` tag:

```bash
geol tag os
```

## 🌳 Example Output

The command displays matching products in a tree structure.

```text
canonical
├── ubuntu
├── microk8s
├── juju
└── lxd
```

:::info

The actual output depends on the content of your local cache and may change over time as new products are added to endoflife.date.

:::

## 🔍 Discover Available Tags

To list all available tags:

```bash
geol list tags
```

This command can help you identify valid tags before using `geol tag`.

## ✅ Common Use Cases

Use this command to:

- Discover products associated with a tag
- Explore technology ecosystems
- Browse products from a vendor
- Identify related products
- Find products before checking lifecycle information

## 📚 Related Commands

- `list`
- `category`
- `product`