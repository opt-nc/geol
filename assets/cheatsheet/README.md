# ❔ About `geol` CLI Cheatsheet Build

This directory contains source code to build cheatsheet.

It uses :

- The `Taskfile.yml` to define build tasks
- The `cheatsheet_geol.tex` LaTeX source file

**All files in the `dist/` directory are generated from this source : do not edit them manually.**

## How to Build

Run:

```sh
task build
```

This will generate the cheatsheet in various formats in the `dist/` directory.
