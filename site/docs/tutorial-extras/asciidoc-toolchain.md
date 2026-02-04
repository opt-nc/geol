---
sidebar_position: 2
---

# `asciidoc` Toolchain

`asciidoc` lets you write human-readable documentation and convert it to HTML, PDF, or slides using tools like `asciidoctor`.
 
See https://pandoc.org/ and https://asciidoctor.org/ for installation and options.

## Notes & workflow examples

If you have a Markdown report (`geol-report.md`) you can:

Convert to asciidoc with `pandoc`:
```bash
pandoc geol-report.md -f markdown -t asciidoc -o geol-report.adoc
```

Process asciidoc (.adoc) files with `asciidoctor`:
```bash
asciidoctor -a toc -a toclevels=4 geol-report.adoc
```
```bash
asciidoctor-pdf -a toc -a toclevels=4 geol-report.adoc
```
