---
sidebar_position: 2
---

# 📊 Export for Analytics

**geol** can export lifecycle data to portable database formats that can be used with analytics, reporting, and ETL tools.

Supported export formats include:

- DuckDB
- SQLite

The exported data can be queried using SQL and integrated with a wide range of analytics platforms.

## 🦆 Export to DuckDB

Export the complete endoflife.date dataset to a DuckDB database.

```bash
geol export
```

This command creates a `geol.duckdb` file containing structured lifecycle data from endoflife.date.

:::tip

The exported database contains multiple tables, including `products`, `categories`, and `tags`.

:::

## 🗄️ Export to SQLite

Export the complete endoflife.date dataset to a SQLite database.

```bash
geol export sqlite
```

This command creates a portable SQLite database that can be queried using standard SQL tools.

:::tip

SQLite is well suited for lightweight applications, embedded systems, and portable analytics workflows.

:::

## ⚙️ Install Required Tools

### DuckDB

Install DuckDB with Homebrew:

```bash
brew install duckdb
```

:::info

See the official installation guide:

https://duckdb.org/install/

:::

### SQLite

Check whether SQLite is available:

```bash
sqlite3 --version
```

Install SQLite if needed:

```bash
brew install sqlite
```

:::info

See the official download page:

https://www.sqlite.org/download.html

:::

## 🔍 Query Exported Data

### DuckDB Examples

Display help:

```bash
duckdb -help
```

Run a query without opening the interactive shell:

```bash
duckdb geol.duckdb -c "select * from tags;"
```

:::tip

This is useful for scripts and quick data checks.

:::

Open the interactive DuckDB shell:

```bash
duckdb geol.duckdb
```

List available tags:

```sql
from tags;
```

Count products:

```sql
select count(*) from products;
```

### SQLite E*amples

Run a query without openin* the interactive shell:

```bash
s*lite3 geol.db "SELECT * FROM produ*ts LIMIT 10;"
```

:::tip

Use `-h*ader` and `-column` for improved r*adability:

```bash
sqlite3 -heade* -column geol.db "SELECT * FROM pr*ducts LIMIT 10;"
```

:::

Open th* SQLite shell:

```bash
sqlite3 ge*l.db
```

List tables:

```sql
.ta*les
```

Display a table schema:

*``sql
.schema products
```

Count *roducts:

```sql
SELECT COUNT(*) FROM products;
```

Exit SQLite:

```sql
.quit
```

## 📚 Generate Database Documentation

SchemaCrawler can generate documentation and diagrams from the exported database.

Project page:

https://github.com/schemacrawler/SchemaCrawler

### Install SchemaCrawler

```bash
brew tap schemacrawler/homebrew-tap
brew install --formula schemacrawler
```

### Install Graphviz

```bash
brew install graphviz
```

:::note

Graphviz is required to generate PNG schema diagrams.

:::

### Generate a Schema Diagram

```bash
schemacrawler \
  --url="jdbc:duckdb:geol.duckdb" \
  --command=schema \
  --info-level=standard \
  --output-format=png \
  --output-file=geol_duckdb_chart.png
```

:::tip

This command generates a visual diagram showing tables and relationships.

:::

### Generate HTML Documentation

```bash
schemacrawler \
  --url="jdbc:duckdb:geol.duckdb" \
  --command=schema \
  --info-level=standard \
  --output-format=htmlx \
  --output-file=geol_duckdb_doc.html
```

:::info

The generated HTML documentation provides detailed information about tables and columns.

:::

## 🔗 Integration Examples

The exported database can be integrated with many tools:

- Python / Pandas
- Jupyter Notebooks
- Tableau
- Power BI
- Apache Airflow
- dbt

Example with Python:

```python
import duckdb

con = duckdb.connect("geol.duckdb")

df = con.execute(
    "SELECT * FROM products WHERE eol_date < CURRENT_DATE"
).df()

print(f"Found {len(df)} products past EOL")
```

## 📹 Additional Resources

### Exploring DuckDB with Quarto

<iframe
  width="100%"
  height="415"
  src="https://www.youtube.com/embed/G_x2Aven5Yg"
  title="Exploring geol DuckDB export with Quarto"
  frameborder="0"
  allowfullscreen>
</iframe>

### Using SQLite with geol

<iframe
  width="100%"
  height="415"
  src="https://www.youtube.com/embed/SojNndY8vrk"
  title="Using geol with SQLite"
  frameborder="0"
  allowfullscreen>
</iframe>

## 📈 Export to Other Formats

DuckDB supports exporting data to additional formats.

Export to CSV:

```sql
COPY (SELECT * FROM products)
TO 'products.csv'
(HEADER, DELIMITER ',');
```

Export to Parquet:

```sql
COPY products
TO 'products.parquet'
(FORMAT PARQUET);
```

Export to JSON:

```sql
COPY (SELECT * FROM products LIMIT 100)
TO 'products.json'
(FORMAT JSON);
```

## 💡 Use Cases

Typical use cases include:

- Proactive lifecycle monitoring
- Compliance reporting
- Technology inventory management
- Vendor analysis
- Risk assessment
- Dashboard and reporting solutions