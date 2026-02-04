---
sidebar_position: 2
---

# Chaîne d'outils `asciidoc`

`asciidoc` permet d'écrire une documentation lisible par l'humain et de la convertir en HTML, PDF, ou diaporama en utilisant des outils comme `asciidoctor`.
 
Voir https://pandoc.org/ et https://asciidoctor.org/ pour l'installation et les options.

## Notes & exemples de workflow

Si vous disposez d'un rapport Markdown (`geol-report.md`), vous pouvez :

Convertir en asciidoc avec `pandoc` :
```bash
pandoc geol-report.md -f markdown -t asciidoc -o geol-report.adoc
```

Traiter les fichiers asciidoc (.adoc) avec `asciidoctor` :
```bash
asciidoctor -a toc -a toclevels=4 geol-report.adoc
```
```bash
asciidoctor-pdf -a toc -a toclevels=4 geol-report.adoc
```