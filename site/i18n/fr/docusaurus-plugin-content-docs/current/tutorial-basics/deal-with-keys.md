---
sidebar_position: 1
---

# Connaître les clés `endoflife.date` avec geol

Il existe de nombreuses clés utilisées par endoflife.date ; elles associent aux produits des métadonnées telles que les versions, les dates de sortie et les dates de fin de vie.

## Produits

Un produit est défini par l'API endoflife.date ; `geol` utilise cette API et ne définit pas les produits lui-même (par exemple : Windows, Ubuntu, iPhone...).

Affiche les commandes disponibles pour un produit spécifique :
```shell
geol help product
```

Liste les produits disponibles :
```shell
geol list products
```
Pour choisir le nombre de cycles de versions à afficher, utilisez `extended` avec l'option `-n`. Par exemple (avec ubuntu) :
```shell
geol product extended ubuntu -n20
```

## Catégories

Une catégorie regroupe des produits liés (par exemple : bibliothèques, runtimes).
```shell
geol list categories
```

## Tags

Un tag est un mot-clé court utilisé pour étiqueter et filtrer les produits.
```shell
geol list tags
```