---
sidebar_position: 2
---

 # Data Engineering with `geol`

`geol` offers features that make it possible and easy for data scientists, data engineers, or devops engineers
to analyze EOL and products at scale comfortably with data tools, hence making it possible to produce
`csv`, `duckdb` or any other portable format.

With these features, it is very easy to load these data in any `ETL` or reporting tools.

## Export to DuckDB

`geol` natively supports `duckdb` export: you can export product information to a DuckDB database using the following command:

```bash title="Export to DuckDB"
geol export
```

This command produces the file `geol.duckdb` containing the whole [`endoflife.date`](https://endoflife.date/) as single structured database.

:::tip
The exported database contains multiple tables including `products`, `categories`, and `tags` that you can query using standard SQL.
:::

You are now ready to play with the `duckdb` file.

## Export to SQLite

`geol` also supports `sqlite` export: you can export product information to a SQLite database using the following command:

```bash title="Export to SQLite"
geol export sqlite
```

This command produces a SQLite database file containing the whole [`endoflife.date`](https://endoflife.date/) dataset as a single structured database.

:::tip
SQLite is perfect for lightweight applications, embedded systems, and scenarios where you need a portable, self-contained database without a separate server process.
:::

You can query the SQLite database using the `sqlite3` command-line tool or any SQLite-compatible client.

## Requirements

### DuckDB

To play with the `duckdb` file, you need to have `duckdb` setup and ready (recommended install option below):

```bash title="Install DuckDB (brew install)"
brew install duckdb
```

:::info
See [DuckDB Installation](https://duckdb.org/install/?platform=linux&environment=cli) for more installation options on Linux, Windows, and other platforms.
:::

### SQLite

To work with SQLite databases, you need the `sqlite3` command-line tool. It's often pre-installed on most systems:

```bash title="Check if SQLite is installed"
sqlite3 --version
```

If not installed, you can install it:

```bash title="Install SQLite (brew install)"
brew install sqlite
```

:::info
See [SQLite Download](https://www.sqlite.org/download.html) for more installation options on different platforms.
:::


## Basic DuckDB Examples

Show DuckDB help:

```shell title="Display help"
duckdb -help
```

Run a single SQL command (The `tags` table can be replaced with `products`, `categories`, ...) 
without opening the interactive CLI:

```bash title="Run a single query"
duckdb geol.duckdb -c "select * from tags;"
```

:::tip
This is useful for scripting or quick data checks without entering the interactive shell.
:::

Open the database in the DuckDB CLI (opens `geol.duckdb`):

```shell title="Open interactive CLI"
duckdb geol.duckdb
```

Inside the DuckDB CLI you can run SQL. For example, list tags:

```sql title="List all tags"
from tags;
```

Count the number of products:

```sql title="Count products"
select count(*) from products;
```

## Basic SQLite Examples

Run a single SQL command without opening the interactive CLI:

```bash title="Run a single query"
sqlite3 geol.db "SELECT * FROM products LIMIT 10;"
```

:::tip
Use the `-header` and `-column` flags for prettier output: `sqlite3 -header -column geol.db "SELECT * FROM products LIMIT 10;"`
:::

Open the database in the SQLite interactive CLI:

```shell title="Open interactive CLI"
sqlite3 geol.db
```

Inside the SQLite CLI, you can run SQL commands:

```sql title="List all tables"
.tables
```

```sql title="Show table schema"
.schema products
```

```sql title="Count products"
SELECT COUNT(*) FROM products;
```

```sql title="Exit SQLite CLI"
.quit
```

## Understand database structure with `schemacrawler`

[`schemacrawler`](https://github.com/schemacrawler/SchemaCrawler) is

> a free database schema discovery and comprehension tool. [...]
> You can search for database schema objects using regular expressions, and
> output the schema and data in a readable text format. The output serves
> for database documentation, [...]
> SchemaCrawler also generates schema diagrams. You can execute scripts in any standard scripting 
> language against your database.

Let's install it:

```shell title="Install SchemaCrawler (brew install)"
brew tap schemacrawler/homebrew-tap
brew install --formula schemacrawler
```

Also, for diagram, let's install [Graphviz](https://graphviz.org/):

```sh title="Install Graphviz (brew install)"
brew install graphviz
```

:::note
Graphviz is required to generate visual schema diagrams in PNG format.
:::

Now, we are ready to discover the database documentation.

Generate PNG Chart:

```sh title="Generate schema diagram"
schemacrawler \
  --url="jdbc:duckdb:geol.duckdb" \
  --command=schema \
  --info-level=standard \
  --output-format=png \
  --output-file=geol_duckdb_chart.png
```

:::tip
This generates a visual diagram showing all tables and their relationships, perfect for documentation.
:::

Generate HTML Documentation:

```sh title="Generate HTML documentation"
schemacrawler \
  --url="jdbc:duckdb:geol.duckdb" \
  --command=schema \
  --info-level=standard \
  --output-format=htmlx \
  --output-file=geol_duckdb_doc.html
```

:::info
The HTML output provides an interactive, browsable documentation of your database schema with detailed information about each table and column.
:::

## Happy Data Engineering Experience 🚀

Now that you have your EOL data in DuckDB, the possibilities are endless! Here are some exciting things you can explore:


### 🔗 Integration Possibilities

Your `geol.duckdb` database can be easily integrated with:

- **Python/Pandas**: Load data for machine learning or visualization
- **Jupyter Notebooks**: Create interactive EOL analysis reports
- **Tableau/Power BI**: Build executive dashboards
- **Apache Airflow**: Schedule automated EOL monitoring pipelines
- **dbt**: Transform and model your EOL data for analytics

```python title="Quick Python integration example"
import duckdb

# Connect and query
con = duckdb.connect('geol.duckdb')
df = con.execute("SELECT * FROM products WHERE eol_date < CURRENT_DATE").df()
print(f"Found {len(df)} products past EOL")
```

### 📹 Video Tutorials

#### Exploring with Quarto

Watch this hands-on video tutorial that demonstrates how to use **Quarto** to explore and analyze the DuckDB export interactively:

<iframe width="100%" height="415" src="https://www.youtube.com/embed/G_x2Aven5Yg" title="Exploring geol DuckDB export with Quarto" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>

:::tip[Why Quarto?]
Quarto combines the power of markdown, code execution, and beautiful output rendering. It's perfect for creating reproducible EOL analysis reports that mix SQL queries, visualizations, and narrative documentation.
:::

#### endoflife.date as SQLite

Watch this tutorial to learn how to work with endoflife.date data using SQLite with geol:

<iframe width="100%" height="415" src="https://www.youtube.com/embed/SojNndY8vrk" title="endoflife.date as SQLite with geol" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>

### 📈 Export to Other Formats

DuckDB makes it easy to export your data:

```sql title="Export to CSV"
COPY (SELECT * FROM products) TO 'products.csv' (HEADER, DELIMITER ',');
```

```sql title="Export to Parquet for big data tools"
COPY products TO 'products.parquet' (FORMAT PARQUET);
```

```sql title="Export to JSON"
COPY (SELECT * FROM products LIMIT 100) TO 'products.json' (FORMAT JSON);
```

### 🎯 Real-World Use Cases

Consider these practical applications:

- **Proactive Alerts**: Identify products nearing EOL to plan migrations
- **Compliance Reports**: Track security updates and EOL status for audits
- **Cost Analysis**: Correlate EOL dates with upgrade/replacement budgets
- **Vendor Management**: Analyze patterns across different product vendors
- **Risk Assessment**: Identify critical systems running EOL software

:::tip[Get Creative!]
The data is yours to explore. Try joining tables in new ways, create custom views, or build automated reports. DuckDB's SQL capabilities are extensive—don't hesitate to experiment!
:::

:::success[Share Your Insights]
Found an interesting pattern or built a cool integration? Consider sharing your queries and insights with the community. Data engineering is more fun when we learn from each other!
:::

**Happy querying and may your EOL insights be ever enlightening!** ✨

