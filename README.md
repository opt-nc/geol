# `geol`

![GitHub release (latest by date)](https://img.shields.io/github/v/release/opt-nc/geol)
![GitHub Workflow Status](https://github.com/opt-nc/geol/actions/workflows/test-release.yml/badge.svg)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/opt-nc/geol)
![GitHub License](https://img.shields.io/github/license/opt-nc/geol)
![GitHub Repo stars](https://img.shields.io/github/stars/opt-nc/geol)
[![Powered By: GoReleaser](https://img.shields.io/badge/powered%20by-goreleaser-green.svg)](https://github.com/goreleaser)
[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg)](https://conventionalcommits.org)
[![GoReportCard](https://goreportcard.com/badge/github.com//opt-nc/geol)](https://goreportcard.com/report/github.com/opt-nc/geol)
[![GoDoc](https://godoc.org/github.com/opt-nc/geol?status.svg)](https://pkg.go.dev/github.com/opt-nc/geol)
[![lint-workflow](https://github.com/opt-nc/geol/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/opt-nc/geol/actions/workflows/golangci-lint.yml)
[![OpenSSF Best Practices](https://www.bestpractices.dev/projects/11239/badge)](https://www.bestpractices.dev/projects/11239)
[![📊 OSS Insights](https://img.shields.io/badge/OSS%20Insights-%F0%9F%93%8A-blue)](https://ossinsight.io/analyze/opt-nc/geol#overview)

# ❔ About

`geol` is a Go based alternative to the famous [`hugovk/norwegianblue`](https://github.com/hugovk/norwegianblue) python
based `cli`.

Its ambitions are to : 

1. Deliver a at least UX as good as `hugovk/norwegianblue` aka. `geol` does,
2. First class terminal based UX
3. Make delivery easier and safer
4. Enhance with new custom awaited features
5. Create innovative ways to manager EOLs for more security, for all

# 🧑‍🤝‍🧑 Core team and roles

| Person                                     | 🎯 Product Manager | 👩‍💻 Lead Dev | 📊 Data Scientist | 🛠️ Data Engineer | 🎤 Storyteller | 🧪 End-user & Tester |
|--------------------------------------------|--------------------|---------------|-----------------|------------------|----------------|----------------------|
| [@adriens](https://github.com/adriens)     |         ✅         |               |       ✅        |       ✅         |     ✅         |          ✅          |
| [@supervinh](https://github.com/supervinh) |                    |       ✅      |                 |                  |                |          ✅          |
| [@mbarre](https://github.com/mbarre)       |                    |               |                 |                  |                |           ✅         |
| [@Draks898](https://github.com/Draks898)   |                    |               |                 |                  |                |           ✅         |


# 📑 Resources

- [endoflife API](https://endoflife.date/docs/api/v1/)
- [⌛ Manage EoLs like a boss with endoflife.date 🛑](https://dev.to/adriens/manage-eols-like-a-boss-with-endoflifedate-2ikf)
- [🍃 How Long Your Phone Will Be Supported by the Manufacturer 📅](https://dev.to/adriens/how-long-your-phone-will-be-supported-by-the-manufacturer-3elf)
- [♾️ Efficient stack management with eol on GitHub 🧑‍🤝‍🧑](https://dev.to/optnc/efficient-stack-management-with-eol-on-github-24g8)
- [🔬 Gitlab 15.8 analysis w. endoflife.date, grype, (x)eol 🐋](https://dev.to/optnc/gitlab-158-analysis-w-endoflifedate-grype-xeol-24b0)
- [endoflife.date Series' Articles](https://dev.to/adriens/series/21232)


# 🚀 QuickStart

To install:

- **Manually**: Go to the [releases](https://github.com/opt-nc/geol/releases) page and download the version corresponding to your operating system
- **Automatic**: Install via `brew` see [homebrew-tap](https://github.com/opt-nc/homebrew-tap)

1. Install `geol`

With `brew`:

```sh
brew install --cask opt-nc/homebrew-tap/geol
```

Or simply with `go`:

```sh
go install github.com/opt-nc/geol@latest
```
Or with shell script:

```sh
curl -fsSL https://raw.githubusercontent.com/opt-nc/geol/main/install.sh | bash
```

2. Update:

```sh
brew update && brew upgrade --cask
geol version
```

## Autocompletion

If the autocompletion is not working yet, you need to update your fpath.

For `zsh` you need to add the following lines to your `.zshrc` file : 

```sh
# === Homebrew completions ===
# Add Homebrew completions to fpath
fpath=(/home/linuxbrew/.linuxbrew/share/zsh/site-functions $fpath)

# Activate Homebrew completions
autoload -Uz compinit
compinit

# (Optional) Disable the warning of duplicate completions
zstyle ':completion:*' verbose yes
```

Then restart your terminal or source your `.zshrc` file : 

```sh
source ~/.zshrc
```

## 🧑‍💻 Documentation

You can access the help either by using the `help` subcommand or with the man page.

```sh
geol help
man geol
```

## 🤓 Build it yourself

To build it youself : 

```sh
go build -o geol ./cmd/geol
```
## 🍿 See it in action

<a href= "https://www.youtube.com/watch?v=vhFXWGqB_-g"><img width="1280" height="720" alt="geol-unboxing-video-thumbnail" src="https://github.com/user-attachments/assets/427e704a-d3ac-4be9-8bf8-7ece98d302cc" /></a>



## 🙏 Acknowledgments

- [endoflife.date](https://endoflife.date) for providing the API

## 💡 More around `endoflife.date` ecosystem

`endoflife` is an ever-growing ecosystem, below some resources to be aware of them : 

- Official [`endoflife.date` Known users](https://github.com/endoflife-date/endoflife.date/wiki/Known-Users)
- [`hugovk/norwegianblue`](https://github.com/hugovk/norwegianblue)
- [`xeol`](https://github.com/xeol-io/xeol) _"A scanner for end-of-life (EOL) software and dependencies in container images, filesystems, and SBOMs"_
