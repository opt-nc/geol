---
sidebar_position: 1
---


# Commencez maintenant?

Découvrons **geol en moins de 5 minutes**.

Commencez par installer `geol` avec Homebrew.

Remarque : Homebrew fournit souvent des paquets plus à jour que d'autres sources, donc l'installation via `brew` vous donnera généralement une version plus récente.

### Comment installer?

- `brew` installé sur votre machine. Voir le site officiel <a href="https://brew.sh/" target="_blank" rel="noreferrer noopener">Homebrew</a> pour les instructions d'installation.

```bash
brew install --cask opt-nc/homebrew-tap/geol
```

Remarque : si l'installation via `brew` échoue, c'est peut-être parce que `curl` fourni par Homebrew n'est pas installé et que votre distribution utilise le `curl` du système (apt).

```bash
# Solution : installer curl avec brew puis réessayer
brew install curl
brew install --cask opt-nc/homebrew-tap/geol
```



## Obtenir la version de `geol`

```bash
geol version
```



## En savoir plus sur `geol`

```bash
geol about
```

## Obtenir de l'aide

```bash
geol help
man geol
```
