---
sidebar_position: 1
---

# Chaîne d'outils `asciidoc`

`AsciiDoc` permet d'écrire une documentation lisible par l'humain et de la convertir en HTML, PDF, ou diaporama en utilisant des outils comme `Asciidoctor`.
Voici deux commandes courantes pour convertir un rapport `AsciiDoc` en HTML (avec table des matières) et en PDF :

```bash
asciidoctor -a toc -a toclevels=4 geol-report.adoc  # convertir geol-report.adoc en HTML avec table des matières
asciidoctor-pdf -a toc -a toclevels=4 geol-report.adoc  # convertir en PDF
```
 
Voir https://pandoc.org/ et https://asciidoctor.org/ pour l'installation et les options.

## Notes & exemples de workflow

Si vous disposez d'un rapport Markdown (`geol-report.md`), vous pouvez :

 - Convertir en AsciiDoc et traiter avec `Asciidoctor` :

```bash
pandoc geol-report.md -f markdown -t asciidoc -o geol-report.adoc
asciidoctor -a toc -a toclevels=4 geol-report.adoc
asciidoctor-pdf -a toc -a toclevels=4 geol-report.adoc
```

- Ou produire EPUB/PDF directement avec `pandoc` :

```bash
pandoc geol-report.md -o geol-report.epub
pandoc geol-report.md -o geol-report.pdf --pdf-engine=xelatex
```
