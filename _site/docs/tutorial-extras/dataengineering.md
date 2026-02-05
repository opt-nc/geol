---
sidebar_position: 3
---

 # Data Engineering with `geol`
Data Engineering with `geol explains how to collect, normalize, and analyze product metadata (versions, release dates, end-of-life) for downstream analytics. It covers ingestion, storage, querying and exporting workflows (CSV, DuckDB) to build repeatable ETL pipelines and reports.

Requirements: install DuckDB — you can install via Homebrew with the following command:
```bash
brew install duckdb
```
See https://duckdb.org/ for more installation options and details.

## Export to DuckDB

You can export product information to a DuckDB database using the following command:
```bash
geol export
```
This command produces the file `geol.duckdb` containing the exported product information.

## Basic DuckDB examples

Show DuckDB help:
```bash
duckdb -help
```
Open the database in the DuckDB CLI (opens `geol.duckdb`):
```bash
duckdb geol.duckdb
```
Inside the DuckDB CLI you can run SQL. For example, list tags:
```sql
from tags;
```
Count the number of products:
```sql
select count(*) from products;
```