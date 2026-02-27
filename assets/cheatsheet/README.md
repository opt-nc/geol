# ❔ About `geol` CLI Cheatsheet Build

This directory contains source code to build cheatsheet.

It uses :

- The `Taskfile.yml` to define build tasks
- The `cheatsheet_geol.tex` LaTeX source file

**All files in the `dist/` directory are generated from this source : do not edit them manually.**

## Prerequisites

To build the cheatsheet, you need to have a LaTeX distribution installed, specifically with the `xelatex` command available. On Debian-based systems, you can install a full TeX Live distribution (which includes xelatex) with:

```sh
sudo apt-get update && sudo apt-get install texlive-full
```

Alternatively, you can install just the `xelatex` package with:

```sh
sudo apt-get update && sudo apt-get install texlive-xetex
```

You also need to have the `texlive-fonts-extra` package installed, which provides the `fontawesome5` package. On Debian-based systems, you can install it with:

```sh
sudo apt-get update && sudo apt-get install texlive-fonts-extra
```

## How to Build

Run:

```sh
task build
```

This will generate the cheatsheet in various formats in the `dist/` directory.
