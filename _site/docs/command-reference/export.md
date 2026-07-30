---
sidebar_position: 13
---

# 📤 export

Export lifecycle data to a database format.

## 🖥️ Usage

```bash
geol export [command] [options]
```

## 📄 Description

The `export` command exports all known products and their end-of-life (EOL) metadata to a database file.

By default, **geol** exports data to a DuckDB database.

The exported dataset contains lifecycle information retrieved from endoflife.date and can be queried using standard SQL tools.

Supported export formats:

- DuckDB (default)
- SQLite

## ⚙️ Global Options

| Option | Description |
|----------|-------------|
| `-f, --force` | Overwrite the output file if it already exists |
| `-o, --output` | Output database file path |
| `-l, --log-level` | Logging level (`debug`, `info`, `warn`, `error`) |

## 📋 Available Subcommands

| Subcommand | Description |
|------------|-------------|
| `duckdb` | Export data to a DuckDB database |
| `sqlite` | Export data to a SQLite database |

## 🦆 Export to DuckDB

Export lifecycle data to a DuckDB database.

### Usage

```bash
geol export
```

or

```bash
geol export duckdb
```

### Example

```bash
geol export
```

By default, the database is created as:

```text
geol.duckdb
```

Use a custom output file:

```bash
geol export --output products.duckdb
```

## 🗄️ Export to SQLite

Export lifecycle data to a SQLite database.

### Usage

```bash
geol export sqlite
```

### Example

```bash
geol export sqlite
```

Use a custom output file:

```bash
geol export sqlite --output products.db
```

## 📁 Output Files

Default output files:

| Format | Default File |
|----------|-------------|
| DuckDB | `geol.duckdb` |
| SQLite | Platform dependent or specified with `--output` |

:::info

Use the `--output` option to control the location and name of the generated database file.

:::

## 🔄 Overwrite Existing Files

If the output file already exists, use the `--force` option to overwrite it.

```bash
geol export --force
```

Example with a custom file:

```bash
geol export --output products.duckdb --force
```

:::warning

Existing files will be replaced when using the `--force` option.

:::

## 💡 Examples

Export to DuckDB:

```bash
geol export
```

Export to SQLite:

```bash
geol export sqlite
```

Export to a custom file:

```bash
geol export --output lifecycle.duckdb
```

Overwrite an existing export:

```bash
geol export --force
```

## ✅ Common Use Cases

Use this command to:

- Analyze lifecycle data with SQL tools
- Build dashboards and reports
- Feed ETL pipelines
- Export data for data engineering workflows
- Integrate lifecycle data with BI platforms

## 📊 Related Documentation

See **Export for Analytics** in the **Advanced Usage** section for examples using DuckDB, SQLite, SQL queries, and reporting tools.