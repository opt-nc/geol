# ❔ About

To be able to reproduce compilation/delivery/installation issues (or just for fun), we provide a
set of ready-to-use [`distrobox`](https://distrobox.it/) machines.

# 🚀 Quickstart

Build the boxes: 

```sh
distrobox-assemble create
```

List them:

```sh
distrobox-list
```

Enter one of them (wait a little bit):

```sh
distrobox enter arch-latest
```

then

```sh
fastfetch
```

or even:

```sh
distrobox enter ubuntu-24-arm64
```

Once done, better cleanup:

```sh
distrobox assemble rm
```
# 📑 Related resources

- https://distrobox.it/
- [🍿Run any version of any OS on any architecture anywhere w.distrobox](https://www.youtube.com/clip/UgkxM65ozHrNb9lYSm6g-h6eBRCGjJ_3okIa)
