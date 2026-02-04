---
sidebar_position: 3
---

# Apprendre la commande `check`
Vérifiez rapidement si les composants de votre environnement logiciel sont en fin de vie (EOL) ou s'en approchent. La commande `geol check` analyse les produits du `.geol.yaml` et affiche leur statut et leurs dates de fin de vie.

## Aide pour `check`

Utilisez `geol help check` pour afficher l'aide et les options disponibles de la commande `check`.

```bash
geol help check
```

## Initialiser un fichier de check

Exécutez la commande pour créer un modèle `.geol.yaml` dans le répertoire courant :

```shell
geol check init
```

Éditez le `.geol.yaml` généré pour lister les produits que vous souhaitez surveiller.

Exemple minimal de `.geol.yaml` (créé par `geol check init`) :

```yaml
stack:
  - name: ubuntu
    version: "25.10"
    id_eol: ubuntu

  - name: java temurin
    version: "21"
    id_eol: eclipse-temurin
...
```

## Statuts et avertissements

Lancez la vérification pour afficher les statuts et les avertissements :

```shell
geol check
```

Utilisez ce `flag` pour que `geol check` renvoie un code d'erreur si au moins un produit est en fin de vie.
```bash
--strict
```