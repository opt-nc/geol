---
sidebar_position: 1
---

# `asciidoc` Toolchain

`asciidoc` lets you write human-readable documentation and convert it to HTML, PDF, or slides using tools like `asciidoctor`.

### Requirements

- `asciidoctor` — install with Homebrew using the following command:
```bash
brew install asciidoc
```
See https://asciidoctor.org/ for installation and options.
Note: Asciidoctor often produces more polished rendering than Pandoc.

## Notes & workflow examples

If you have a Markdown report (`geol-report.md`) you can:

Convert to asciidoc with `pandoc`:
```bash
pandoc geol-report.md -f markdown -t asciidoc -o geol-report.adoc
```
Generate HTML from the `.adoc` input:
```bash
asciidoctor -a toc -a toclevels=4 geol-report.adoc
```
Generate PDF from the `.adoc` input:
```bash
asciidoctor-pdf -a toc -a toclevels=4 geol-report.adoc
```
