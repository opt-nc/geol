---
sidebar_position: 1
---

# 🔑 Discover endoflife.date keys

**geol** uses data provided by the endoflife.date API. This page explains the different types of keys available and how to explore them.


:::info About endoflife.date

endoflife.date maintains information about software and hardware products, including versions, release dates, support periods, and end-of-life dates.

:::

## 📦 Products

A product represents a technology tracked by the endoflife.date API, such as Ubuntu, Windows, iPhone, Node.js, or PostgreSQL.


### Show available product commands:
```bash
geol help product
```

### List available products:
```bash
geol list products
```

### Display More Release Cycles

Use the `extended` command together with the `-n` option to control how many release cycles are displayed.

For example, to display the last 20 Ubuntu release cycles:
```shell
geol product extended ubuntu -n20
```

## 🗂️ Categories

A category groups related products together.

Examples include libraries, runtimes, databases, operating systems, and programming languages.

List available categories:

```bash
geol list categories
```

## 🏷️ Tags

Tags are keywords used to classify and filter products.

List available tags:

```bash
geol list tags
```