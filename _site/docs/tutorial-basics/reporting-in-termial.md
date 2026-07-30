---
sidebar_position: 2
---

# 📦 Get Product Information

**geol** retrieves product metadata including available versions, release dates, and end-of-life dates. You can use this information to generate summaries, version lists, or exportable reports.

## 🔍 Get a Product Overview

Display detailed information about a product, such as its description, available versions, release dates, end-of-life dates, and version detection command.

The example below uses `ubuntu`, but you can replace it with any supported product.

 ```bash
 geol product describe ubuntu
 ```

 :::tip Data source
This output is a summary of the information available on endoflife.date — for Ubuntu see https://endoflife.date/ubuntu.
:::