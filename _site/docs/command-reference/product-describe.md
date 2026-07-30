---
sidebar_position: 9
---

# 📄 product describe

Display a summary of a product.

## 🖥️ Usage

```bash
geol product describe [product]
```

## 📄 Description

The `describe` subcommand displays a summary of a product tracked by endoflife.date.

It provides detailed information about the selected product, including lifecycle metadata and additional product information when available.

This command is useful when you need a quick overview of a product before exploring its release history or lifecycle details.

## 💡 Example

Display information about Ubuntu:

```bash
geol product describe ubuntu
```

## 📸 Example Output

```bash
geol product describe ubuntu
```

<!--![Output of the geol product describe command](./img/product-describe-ubuntu.png)-->

<!--*Output of the `geol product describe ubuntu` command.* -->

## 📦 Product Information

Depending on the product and the data available from endoflife.date, the output may include:

- Product description
- Product name
- Available versions
- Release dates
- End-of-life dates
- Version detection command
- Additional metadata

:::info

The available information depends on the selected product and the data provided by endoflife.date.

:::

## 🔍 Find Available Products

Use the following command to list available products:

```bash
geol list products
```

Once you have identified a product, use `describe` to display its details.

## ✅ Common Use Cases

Use this command to:

- Explore a product tracked by endoflife.date
- Review lifecycle information
- Identify supported versions
- Prepare upgrades and migrations
- Retrieve product metadata

## 📚 Related Commands

- `product extended`
- `list products`
- `category`
- `tag`