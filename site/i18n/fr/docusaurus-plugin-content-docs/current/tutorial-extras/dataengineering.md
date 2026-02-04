---
sidebar_position: 3
---


# Data Engineering with `geol`

Data Engineering with `geol` explique comment collecter, normaliser et analyser les métadonnées produits (versions, dates de sortie, dates de fin de vie) pour des analyses en aval. Il couvre le stockage, les requêtes et les flux d'export (CSV, DuckDB) pour construire des pipelines ETL reproductibles et des rapports.

Il est possible d'exporter des informations vers une base DuckDB grâce à la commande suivante :

```bash
geol export
```

Prérequis : installez DuckDB en suivant les instructions sur https://duckdb.org/.