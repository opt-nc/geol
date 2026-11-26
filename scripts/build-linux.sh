#!/bin/sh
# Builds the Linux release archives for the host architecture.
# Meant to run inside an AlmaLinux 8 container: glibc 2.28 keeps the binary
# runnable on older distributions, and the static libstdc++ removes any
# dependency on the system C++ runtime.
set -eu

GORELEASER_VERSION=2.20.3
SYFT_VERSION=1.57.2

case "$(uname -m)" in
x86_64) GOARCH=amd64; GR_ARCH=x86_64 ;;
aarch64) GOARCH=arm64; GR_ARCH=arm64 ;;
*) echo "unsupported architecture: $(uname -m)" >&2; exit 1 ;;
esac

GO_VERSION=$(awk '/^go /{print $2}' go.mod)

dnf install -y gcc-toolset-13-gcc-c++ gcc-toolset-13-libstdc++-devel git tar gzip >/dev/null

curl -fsSL "https://go.dev/dl/go${GO_VERSION}.linux-${GOARCH}.tar.gz" | tar -xz -C /usr/local
curl -fsSL "https://github.com/goreleaser/goreleaser/releases/download/v${GORELEASER_VERSION}/goreleaser_Linux_${GR_ARCH}.tar.gz" | tar -xz -C /usr/local/bin goreleaser
curl -fsSL "https://github.com/anchore/syft/releases/download/v${SYFT_VERSION}/syft_${SYFT_VERSION}_linux_${GOARCH}.tar.gz" | tar -xz -C /usr/local/bin syft

export PATH=/opt/rh/gcc-toolset-13/root/usr/bin:/usr/local/go/bin:$PATH
git config --global --add safe.directory '*'

# "goreleaser release" has no --single-target flag: restrict the build to the host arch
CONFIG=$(mktemp)
sed "s/^    goos:/    goarch:\n      - ${GOARCH}\n    goos:/" .goreleaser_linux.yaml > "$CONFIG"

goreleaser release --clean --skip=publish --config "$CONFIG" "$@"
