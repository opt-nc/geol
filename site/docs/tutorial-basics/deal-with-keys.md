---
sidebar_position: 1
---

# Know `endoflife.date` keys with geol

There are a lot of keys used by endoflife.date; they map products to metadata such as versions, release dates and end-of-life dates.

## Products

A product is defined by the endoflife.date API; `geol` uses that API and does not define products itself (for example: Windows, Ubuntu, iPhone...).

Show available product commands:
```shell
geol help product
```

List available products:
```shell
geol list products
```
To choose how many release cycles to display, use `extended` with the `-n` flag. For example ( with ubuntu ):
```shell
geol product extended ubuntu -n20
```

## Categories

A category groups related products (for example: libraries, runtimes).
```shell
geol list categories
```

## Tags

A tag is a short keyword used to label and filter products.
```shell
geol list tags
```