---
sidebar_position: 1
---

# `asciidoc` Toolchain

`AsciiDoc` lets you write human-readable documentation and convert it to HTML, PDF, or slides using tools like `Asciidoctor`.
Here are two common commands to convert an `AsciiDoc` report to HTML (with a table of contents) and to PDF:

```bash
asciidoctor -a toc -a toclevels=4 geol-report.adoc  # convert geol-report.adoc to HTML with a table of contents
asciidoctor-pdf -a toc -a toclevels=4 geol-report.adoc  # convert geol-report.adoc to PDF
```
 
See https://pandoc.org/ and https://asciidoctor.org/ for installation and options.

## Notes & workflow examples

If you have a Markdown report (`geol-report.md`) you can:

 - Convert to AsciiDoc and process with `Asciidoctor`:

```bash
pandoc geol-report.md -f markdown -t asciidoc -o geol-report.adoc
asciidoctor -a toc -a toclevels=4 geol-report.adoc
asciidoctor-pdf -a toc -a toclevels=4 geol-report.adoc
```

- Or produce EPUB/PDF directly with `pandoc`:

```bash
pandoc geol-report.md -o geol-report.epub
pandoc geol-report.md -o geol-report.pdf --pdf-engine=xelatex
```
