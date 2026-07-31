---
sidebar_position: 6
---

# 📂 category

Display all products associated with a category.

## 🖥️ Usage

```bash
geol category [category] [options]
```

## 📄 Description

Use this command to display all products associated with a category.

The category must exist in the local cache. Results are displayed in a tree structure to make product discovery easier.

This command is useful for exploring products by technology domain, such as operating systems, cloud platforms, databases, runtimes, or programming languages.

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
geol category os --log-level debug
```

## 💡 Examples

Display products in the operating systems category:

```bash
geol category os
```

Display products in the cloud category:

```bash
geol category cloud
```

## 🌳 Example Output

The command displays products in a tree structure.

```text
os
├── ubuntu
├── debian
├── fedora
├── centos
└── windows
```

:::info

The actual output depends on the content of your local cache and may change over time as new products are added to endoflife.date.

:::

## ✅ Common Use Cases

Use this command to:

- Discover products within a category
- Explore available technologies
- Identify products monitored by endoflife.date
- Find products before using other commands such as `product`

## 📚 Related Commands

- `list`
- `product`
- `tag`
```