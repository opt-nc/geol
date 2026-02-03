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
➜  ~ brew install --cask opt-nc/homebrew-tap/geol    # attempts to install the 'geol' cask via Homebrew
# Output (example error):
Error: Download failed on Cask 'geol' with message: Download failed:
Homebrew-installed `curl` is not installed
for: https://github.com/opt-nc/geol/releases/download/v0.3.0/geol_Linux_x86_64.tar.gz  # the archive couldn't be downloaded because the expected Homebrew-provided curl is missing

# Suggested solution: install curl via Homebrew, then retry
brew install curl    # install 'curl' via Homebrew
brew install --cask opt-nc/homebrew-tap/geol    # retry installing the 'geol' cask
```

Ou installez avec :

```bash
# Go
go install github.com/opt-nc/geol@latest

# Script Shell
curl -fsSL https://raw.githubusercontent.com/opt-nc/geol/main/install.sh | bash
```

Vous pouvez taper cette commande dans l'invite de commandes, le Terminal, ou tout autre terminal intégré de votre éditeur de code.

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