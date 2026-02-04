---
sidebar_position: 2
---

# Reporting products with `geol`

`geol` fetches metadata for each product (versions, release dates, and end-of-life dates) and lets you generate summaries, version lists, or exportable reports.

## Get a product overview

Displays detailed metadata for a product (description, version command,...); `ubuntu` is an example — replace it with any product name.

 ```bash
 geol product describe ubuntu
 ```
This output is a summary of the information available on endoflife.date — for Ubuntu see https://endoflife.date/ubuntu.