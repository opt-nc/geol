# ❔ About

[GitHub Codespaces](https://github.com/features/codespaces) provides a complete, cloud-based development environment that lets you work on `geol` from anywhere, without needing to install anything locally. This is perfect for trying out `geol`, contributing to the project, or just playing with it in a comfortable way.

The `geol` repository is configured to automatically install `geol` when you create a Codespace, so you can start using it immediately!

# 🚀 Quickstart

## Launch a Codespace with geol pre-installed

1. Navigate to the [geol repository](https://github.com/opt-nc/geol)
2. Click the green **Code** button
3. Select the **Codespaces** tab
4. Click **Create codespace on main** (or your desired branch)

GitHub will create a cloud-based development environment and automatically install `geol` using the installation script.

## Verify geol is ready

Once your Codespace has finished building (wait for the post-create command to complete), verify that `geol` is installed:

```sh
geol version
```

## Start using geol

```sh
geol help
```

That's it! No manual installation required. 🎉

# 🛠️ How it works

The repository includes a `.devcontainer/devcontainer.json` configuration file that:

1. Sets up a Go development environment
2. Automatically runs the `install.sh` script during Codespace creation
3. Adds `geol` to your PATH

This means `geol` is ready to use as soon as your Codespace is created!

# 💡 Tips

- Your Codespace persists your changes, so you can close and reopen it later
- Codespaces are free for personal accounts (up to 60 hours/month for free tier)
- You can customize your Codespace environment by editing `.devcontainer/devcontainer.json`
- The Codespace comes with Go, Git, and GitHub CLI pre-installed

# 📑 Related resources

- [GitHub Codespaces documentation](https://docs.github.com/en/codespaces)
- [GitHub Codespaces features](https://github.com/features/codespaces)
- [Dev Containers documentation](https://docs.github.com/en/codespaces/setting-up-your-project-for-codespaces/adding-a-dev-container-configuration)
- [geol installation guide](README.md#-quickstart)
