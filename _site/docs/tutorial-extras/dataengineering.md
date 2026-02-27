---
sidebar_position: 2
---

 # Data Engineering with `geol`

`geol` offers features that make it possible and easy for datascientists, dataengineers, or devops engineers
to analyze EOL and products at scale comfortably with data tools, hence making it possible to produce
`csv`, `duckdb` or any other portable format.

With these features, it is very easy to load these data in any `ETL` or reporting tools.

## Export to DuckDB

`geol` natively supports `duckdb` export: you can export product information to a DuckDB database using the following command:

```bash
geol export
file geol.duckdb
```

This command produces the file `geol.duckdb` containing the exported product information.

You are now ready to play with the `duckdb` file.

## Requirement

To play with the `duckdb` file, you need to have `duckdb` setup and ready (recommanded install option below):

```bash
brew install duckdb
```

See [DuckDB Installation](https://duckdb.org/install/?platform=linux&environment=cli) for more installation options.


## Basic DuckDB examples

Show DuckDB help:

```shell
duckdb -help
```

Run a single SQL command (The `tags` table can be replaced with `products`, `categories`, ...) without
opening the interactive CLI:

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
