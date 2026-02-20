# DuckDB Export Discover Notebook

This Quarto notebook connects to a DuckDB database, discovers the schema, and generates a visual representation of the tables and their relationships. It also provides an interactive table view of the data.

## Requirements

*   [Quarto CLI](https://quarto.org/docs/get-started/)
*   [Python 3.10+](https://www.python.org/downloads/)
*   [Graphviz](https://graphviz.org/download/)
*   A DuckDB database file (e.g., `geol.duckdb`)

## Installation

1.  **Install Quarto, Python, and Graphviz:**

    Follow the instructions on the respective websites to install Quarto and Python for your operating system.

    To install Graphviz on macOS or Linux using Homebrew, run:
    ```bash
    brew install graphviz
    ```
    On other systems, please refer to the [Graphviz download page](https://graphviz.org/download/).

2.  **Create a virtual environment:**

    It is recommended to use a virtual environment to manage the Python dependencies for this project.
    ```bash
    python3 -m venv .venv
    ```

3.  **Activate the virtual environment:**

    On macOS and Linux:
    ```bash
    source .venv/bin/activate
    ```
    On Windows:
    ```bash
    .venv\Scripts\activate
    ```

4.  **Install the required Python packages:**

    With the virtual environment activated, install the necessary packages using pip:
    ```bash
    pip install duckdb plotly pandas itables graphviz networkx pyyaml
    ```

## Usage

1.  **Activate the virtual environment** (if not already activated):
    ```bash
    source .venv/bin/activate
    ```

2.  **Render the Quarto notebook:**
    ```bash
    quarto render duckdb-export-discover.qmd
    ```
    This will execute the notebook and create an HTML file with the output.

    By default, the notebook will look for a database file named `geol.duckdb` in the same directory. You can specify a different database file by setting the `QUARTO_PARAM_DATABASE` environment variable:
    ```bash
    QUARTO_PARAM_DATABASE=path/to/your/database.duckdb quarto render duckdb-export-discover.qmd
    ```

## Output

The output of the notebook is an HTML file named `duckdb-export-discover.html`. This file contains:

*   An interactive table of all tables in the database.
*   A schema diagram showing the tables and their inferred relationships.
*   Interactive tables for each table in the database.
