# Advanced Product Lifecycle Comparison

An interactive, high-fidelity dashboard for comparing the lifecycles, maintenance concurrency, and support trends of products documented on [endoflife.date](https://endoflife.date).

## 🚀 Quick Start

1. **Install dependencies**:
   ```bash
   # Create venv and install python packages
   uv venv
   source .venv/bin/activate
   uv pip install pandas duckdb altair IPython
   ```

2. **Generate the report**:
   ```bash
   task
   ```
   *This will automatically generate the `geol.duckdb` database (if missing) and render the Quarto report.*

3. **Open the report**:
   Open `product_lifecycle_comparison.html` in your browser.

## 🛠 Required Software

To run this project, you need the following tools installed on your system:

- **[Quarto](https://quarto.org/)**: The underlying publishing system used to render the dashboard.
- **[uv](https://docs.astral.sh/uv/)**: An extremely fast Python package installer and resolver.
- **[go-task (Task)](https://taskfile.dev/)**: A task runner / build tool used to automate the workflow.
- **[geol](https://github.com/opt-nc/geol)**: (Optional) Used to refresh the underlying data from `endoflife.date`.

## 📋 Automation Tasks

The project uses a `Taskfile.yml` to simplify common operations:

| Command | Description |
|---------|-------------|
| `task` | Alias for `task build` (Default). |
| `task build` | Renders the Quarto notebook into a self-contained HTML file. |
| `task geol` | Generates the `geol.duckdb` database if it doesn't already exist. |
| `task export` | Force-refreshes the database by re-running `geol export --force`. |
| `task clean` | Removes generated HTML files and the DuckDB database. |

## 📊 Dashboard Features

- **KPI Dashboard**: Real-time stats on active versions and upcoming EOL dates (color-coded by urgency).
- **Maintenance Concurrency**: Analysis of how many versions are supported simultaneously over time.
- **Dual-Axis Gantt Chart**: Side-by-side chronological comparison of all version lifecycles.
- **Support Trends**: Visualizations showing if version support is getting longer or shorter over time.
- **Migration Grace Period**: Calculation of "overlap" time between a new release and its predecessor's EOL.

## ⚙️ Configuration

The products being compared are currently defined inside `product_lifecycle_comparison.qmd`:

```python
product1_id = "postgresql"
product2_id = "mariadb"
```

You can change these IDs to any valid product ID from [endoflife.date](https://endoflife.date) to compare different technologies.
