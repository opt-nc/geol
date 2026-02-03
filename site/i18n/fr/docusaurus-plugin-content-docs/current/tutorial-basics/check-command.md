---
sidebar_position: 3
---

# Apprendre la commande `check`
Vérifiez rapidement si les composants de votre environnement logiciel sont en fin de vie (EOL) ou s'en approchent.

## Initialiser un fichier de check

Exécutez la commande pour créer un modèle `.geol.yaml` dans le répertoire courant :

```shell
geol check init
```

Éditez le `.geol.yaml` généré pour lister les produits que vous souhaitez surveiller.

Exemple minimal de `.geol.yaml` (créé par `geol check init`) :

```yaml
stack:
  - name: Go
    version: "1.25"
    id_eol: go
```

## Statuts et avertissements

Lancez la vérification pour afficher les statuts et les avertissements :

```shell
geol check

# Exemple de sortie (ce que `geol check` affiche) :
# Software       │ Version │ EOL Date    │ Status │   Days │ Is Latest │ Latest
# ───────────────┼─────────┼─────────────┼────────┼────────┼───────────┼--------
# maven          │ 3.6     │ 2021-03-30  │ EOL    │ -1756  │ false     │ 3.9
# quarkus        │ 3.12    │ 2024-07-31  │ EOL    │ -537   │ false     │ 3.30
# traefik        │ 2.11    │ 2026-02-01  │ WARN   │ 12     │ false     │ 3.6
# ubuntu         │ 25.10   │ 2026-07-01  │ OK     │ 162    │ true      │ 25.10
# postgresql     │ 14      │ 2026-11-12  │ OK     │ 296    │ false     │ 18
# java temurin   │ 21      │ 2029-12-31  │ OK     │ 1441   │ false     │ 25
# opensearch     │ 2       │             │ OK     │ -      │ false     │ 3
```

> **Légende** — champs affichés dans l'exemple
>
> **Champs de sortie**
> - Software : nom du produit (ex. `ubuntu`).
> - Version : la version rapportée ou configurée.
> - EOL Date : date de fin de vie connue (vide si inconnue).
> - Status : `OK`, `WARN`, `EOL`, ou `UNKNOWN`.
> - Days : jours avant la EOL (positif) ou depuis la EOL (négatif) ; `-` si inconnu.
> - Is Latest : indique si cette version est la dernière connue.
> - Latest : la dernière version connue pour le produit.
>
> **Signification des statuts**
> - OK — supporté et pas proche de la fin de vie.
> - WARN — approche de la fin de vie (moins de 30 jours).
> - EOL — la version est passée en fin de vie ; mise à jour recommandée.
> - UNKNOWN — pas de date de fin de vie disponible pour cette version.