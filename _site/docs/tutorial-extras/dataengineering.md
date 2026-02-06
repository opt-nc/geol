---
sidebar_position: 2
---

 # Data Engineering with `geol`
Data Engineering with `geol explains how to collect, normalize, and analyze product metadata (versions, release dates, end-of-life) for downstream analytics. It covers ingestion, storage, querying and exporting workflows (CSV, DuckDB) to build repeatable ETL pipelines and reports.

### Requirements

- `duckdb` — install with Homebrew using the following command:
```bash
brew install duckdb
```
See https://duckdb.org/ for more installation options and details.

## Export to DuckDB

You can export product information to a DuckDB database using the following command:
```bash
geol export
file geol.duckdb
```
This command produces the file `geol.duckdb` containing the exported product information.

## Basic DuckDB examples

Show DuckDB help:
```shell
duckdb -help
```
Run a single SQL command (The `tags` table can be replaced with `products`, `categories`, ...) without opening the interactive CLI:
```bash
duckdb geol.duckdb -c "select * from tags;"
```
Open the database in the DuckDB CLI (opens `geol.duckdb`):
```shell
duckdb geol.duckdb
```
Inside the DuckDB CLI you can run SQL. For example, list tags:
```shell
from tags;
```
Count the number of products:
```shell
select count(*) from products;
```