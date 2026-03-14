# Geol Quarto Stack Dashboard

An automated, professional End-of-Life (EoL) reporting dashboard for your technology stack. This project integrates [geol](https://opt-nc.github.io/geol/) with [Quarto](https://quarto.org/) to produce an interactive, dark-themed HTML analysis based on the [endoflife.date](https://endoflife.date/) API.

## 🚀 Features

- **Dynamic Reporting:** Automatically reads your `.geol.yaml` configuration.
- **Interactive Visualizations:**
  - **Sunburst Drill-down:** Click to explore status by component.
  - **Lifecycle Timeline:** Interactive Plotly timeline of product releases and EOL dates.
  - **Searchable Table:** DataTables-powered list with Export buttons (CSV, Excel, Copy).
- **Professional KPI Dashboard:** At-a-glance health metrics (Total, Healthy, EoL, Untracked).
- **Enterprise Ready:** Full metadata tracking (MD5 hashes, tool versions, dual timestamps).
- **Responsive Design:** Dark mode by default with Font Awesome icons.

## 🛠 Prerequisites

Ensure you have the following tools installed:

- **uv**: `brew install uv`
- **Quarto**: `brew install --cask quarto` (MacOs only, see [Quarto installation](https://quarto.org/docs/get-started/) for other platforms)
- **geol**: `brew install curl && brew install --cask opt-nc/homebrew-tap/geol`
- **Taskfile**: `brew install go-task` (optional, for automation)

Python dependencies are managed via `requirements.txt` and will be installed in a virtual environment automatically when building.

## 📖 Usage

### Quick Start
To build the report using the default `.geol.yaml`:
```bash
task build
```

This command will:
1. Verify that all required tools are installed
2. Create a virtual environment (`.venv`) if it doesn't exist
3. Install Python dependencies from `requirements.txt`
4. Render the Quarto report to HTML

### Manual Dependency Installation
If you prefer to install Python dependencies separately:
```bash
task install-deps
```

### Advanced Usage
Run the report for a custom configuration or add extra context:
```bash
task build GEOL_CONFIG="my-app.yaml" ADDITIONAL_CONTEXT="Internal audit for Q1 2026."
```

### Cleanup
Remove generated artifacts:
```bash
task clean
```

Remove everything including the virtual environment:
```bash
task clean-all
```

## 📄 Output
The resulting dashboard is generated as a self-contained HTML file: **`geol-check-report.html`**.

---
*Built with ❤️ by [geol](https://github.com/opt-nc/geol)*