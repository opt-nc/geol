---
sidebar_position: 10
---

# 📋 product extended

Display extended lifecycle information for a product.

## 🖥️ Usage

```bash
geol product extended [product]
```

## 📄 Description

The `extended` subcommand displays detailed release information for a product tracked by endoflife.date.

Unlike `product describe`, which provides a summary, `product extended` shows multiple product releases and their lifecycle information.

By default, the latest 10 releases are displayed.

This command is particularly useful when planning upgrades, migrations, or lifecycle reviews.

## 💡 Example

Display extended lifecycle information for Golang:

```bash
geol product extended golang
```


## 📸 Example Output

```bash
geol product extended golang
```

<!--![Output of the geol product extended command](./img/product-extended-golang.png)-->

<!--*Output of the `geol product extended golang` command.*-->

## 📊 Available Information

Depending on the selected product, the output may include:

- Release version
- Release date
- End-of-life date
- Support status
- Lifecycle phase

The information is retrieved from the endoflife.date dataset.

## 🔍 Why Use Extended Output?

Use `product extended` when you need more than a simple product summary.

Typical use cases include:

- Reviewing release history
- Comparing supported versions
- Planning migrations
- Identifying versions approaching end-of-life
- Preparing upgrade roadmaps

:::tip

Use `product describe` when you only need a quick overview of a product.

Use `product extended` when you need detailed lifecycle information across multiple releases.

:::

## ✅ Common Use Cases

Use this command to:

- Review release history
- Analyze lifecycle trends
- Compare product versions
- Identify supported releases
- Prepare upgrade and migration plans

## 📚 Related Commands

- `product describe`
- `check`
- `list products`