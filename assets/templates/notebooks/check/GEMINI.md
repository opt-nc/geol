# Gemini Project Memory: Geol Quarto Stack Dashboard

This document records the architectural decisions, design patterns, and operational knowledge established during the development of the `geol-check-report`.

## 🏗️ Architecture & Stack
- **Engine:** [Quarto](https://quarto.org/) (v1.9+) rendering a Python-based `.qmd` notebook.
- **Data Source:** [geol](https://github.com/opt-nc/geol) CLI tool consuming the [endoflife.date](https://endoflife.date/) API.
- **Orchestration:** [Taskfile](https://taskfile.dev/) (`Taskfile.yml`) managing build dependencies and execution via [uv](https://docs.astral.sh/uv/).
- **Visualization:** [Plotly](https://plotly.com/python/) for interactive charts (Sunburst, Gantt) and [DataTables](https://datatables.net/) for the stack status list.

## 🎨 UI/UX Design System
- **Theme:** Dual-mode support with **Dark Mode (`darkly`)** as the strict default.
- **Icons:** [Font Awesome 6](https://fontawesome.com/) injected via direct HTML/CSS overrides to prevent stripping by Quarto.
- **Colors:**
  - 🟢 **Healthy:** `#50fa7b` (Green)
  - 🟠 **Warning:** `#ffb86c` (Orange)
  - 🔴 **Critical/EOL:** `#ff5555` (Red)
  - 🔵 **Info/Unknown:** `#8be9fd` (Cyan)

## 🧩 Key Features & Logic
### 1. LTS Detection Strategy
- **Logic:** We do *not* rely on the user's YAML name. We map the normalized software name to its `id_eol` slug.
- **Verification:** We compare the user's version against the official `geol product extended` data.
- **Matching:** A version is LTS if it exactly matches an LTS cycle (e.g., `24.04`) or is a sub-version of one (e.g., `24.04.1`).

### 2. Version Lag Analysis (`packaging` library)
- We use `packaging.version` (PEP 440) for robust comparison.
- **Major Lag:** Detected if the major version segments differ.
- **Minor Lag:** Detected if major versions match but minor/patch segments differ.

### 3. "Executive One-Pager" Print Mode
- **Strategy:** We do **not** render to PDF via LaTeX.
- **Implementation:** A custom `@media print` CSS block is injected.
- **Behavior:** When the user presses `Ctrl+P` (Browser Print):
  - Background turns white (ink-saving).
  - Sidebar, ToC, and interactive buttons are hidden.
  - Tables expand to full width.
  - Charts are optimized for static display.

## ⚠️ Known Constraints & Fixes
- **Font Awesome:** Quarto's Markdown processor strips `<i>` tags in headers. **Fix:** Use `<span>` tags and inject CSS with `!important` font-family rules.
- **Geol Logo:** GitHub raw URLs for JPEGs can break due to MIME types. **Fix:** Use `?raw=true` or embed as Base64 (current implementation uses `?raw=true`).
- **Interactive Tables:** DataTables must be initialized via explicit JS injection in the Python block to work with the dark theme.

## 🔄 operational Workflow
- **Setup:** `task check-tools` verifies `uv`, `quarto`, and `geol` availability.
- **Build:** `task build` (default) renders the HTML report.
- **Clean:** `task clean` removes all artifacts.
