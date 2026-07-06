# Installation Guide

This guide covers installation methods for `uts` and explains how to configure its third-party system dependencies on macOS and Linux.

---

## Table of Contents

- [Quick Start](#quick-start)
- [System Dependencies (Crucial)](#system-dependencies-crucial)
    - [macOS Installation](#macos-installation)
    - [Linux (Debian/Ubuntu) Installation](#linux-debianubuntu-installation)
    - [Nix Installation](#nix-installation)
- [Installation Methods](#installation-methods)
    - [Method 1: Homebrew (Recommended)](#method-1-homebrew-recommended)
    - [Method 2: Nix Flake](#method-2-nix-flake)
    - [Method 3: Go Install](#method-3-go-install)
    - [Method 4: GitHub Release Binaries](#method-4-github-release-binaries)
- [Verification](#verification)
- [Troubleshooting](#troubleshooting)

---

## Quick Start

```bash
# 1. Install uts via Homebrew
brew install y3owk1n/tap/uts

# 2. Install recommended dependencies (macOS example)
brew install ffmpeg imagemagick pngquant jpegoptim cwebp gifsicle poppler ghostscript

# 3. Verify it works
uts info --version
```

---

## System Dependencies (Crucial)

`uts` acts as a high-level command orchestrator. It does not perform encoding directly; instead, it detects and executes specialized command-line utilities already installed on your system.

To get the most out of `uts`, you should install the recommended tools for the categories you intend to use.

| Category    | Utility Tools Used                                                                          | Purpose                                                                    |
| ----------- | ------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------- |
| **Image**   | `imagemagick`, `pngquant`, `jpegoptim`, `cwebp`, `gifsicle`, `optipng`, `sips` (macOS only) | Compressing PNG/JPG/WebP/GIF and converting between various image formats. |
| **Video**   | `ffmpeg`                                                                                    | Compressing and converting video files (using `libx265`).                  |
| **Audio**   | `ffmpeg`                                                                                    | Compressing and converting audio files (using `aac` / `mp3` encoders).     |
| **PDF**     | `ghostscript`, `poppler` (specifically `pdftoppm`), `imagemagick`                           | Compressing PDFs, converting PDFs to images, or stitching images to PDFs.  |
| **Archive** | `tar`, `zstd`, `xz`, `brotli`, `zip` / `unzip`                                              | Compressing, listing, and extracting archive formats.                      |

### macOS Installation

We recommend using [Homebrew](https://brew.sh/) to install all required dependencies:

```bash
# Install everything at once:
brew install ffmpeg imagemagick pngquant jpegoptim cwebp gifsicle optipng poppler ghostscript zstd xz brotli
```

> [!NOTE]
> macOS includes `sips` (Scriptable Image Processing System) by default, which `uts` will leverage as a fast fallback for basic image conversions.

### Linux (Debian/Ubuntu) Installation

Use your system packager (`apt`) to install dependencies:

```bash
sudo apt-get update
sudo apt-get install -y \
  ffmpeg \
  imagemagick \
  pngquant \
  jpegoptim \
  webp \
  gifsicle \
  optipng \
  poppler-utils \
  ghostscript \
  zstd \
  xz-utils \
  brotli
```

### Nix Installation

If you are using Nix, you can install dependencies using `nix-env` or declare them in your Nix shell:

```bash
nix-env -iA \
  nixpkgs.ffmpeg \
  nixpkgs.imagemagick \
  nixpkgs.pngquant \
  nixpkgs.jpegoptim \
  nixpkgs.webp \
  nixpkgs.gifsicle \
  nixpkgs.poppler_utils \
  nixpkgs.ghostscript \
  nixpkgs.zstd \
  nixpkgs.xz \
  nixpkgs.brotli
```

---

## Installation Methods

### Method 1: Homebrew (Recommended)

`uts` is available through Homebrew via the custom tap:

```bash
brew install y3owk1n/tap/uts
```

### Method 2: Nix Flake

`uts` is available as a Nix flake with support for installing as a Nix overlay, direct package reference, or declaring it through Home Manager.

> [!NOTE]
> `pkgs.uts` uses the published release binary (packed as a zip file) and `pkgs.uts-source` builds the Go module from source.

#### Add Flake Input

First, add `uts` to your system's `flake.nix` inputs:

```nix
# flake.nix
{
  inputs = {
    # ... other inputs
    uts.url = "github:y3owk1n/uts";
  };
}
```

#### Option 1: Home Manager Module (User-Level Declarative)

`uts` exposes a Home Manager module that manages user packages and setup.

1. **Import the Module & Configure**:
   Pass the inputs to your configuration and import the default Home Manager module:

    ```nix
    # home.nix or configuration.nix
    { inputs, pkgs, ... }: {
      imports = [
        inputs.uts.homeManagerModules.default
      ];

      # Enable the program
      programs.uts.enable = true;

      # Optional: use build-from-source package instead of pre-compiled release binary
      # programs.uts.package = pkgs.uts-source;
    }
    ```

#### Option 2: Using overlays or packages (System-Level or Declarative nix-darwin/NixOS)

If you prefer managing the packages manually without using the Home Manager module, you can apply the flake overlay:

```nix
# flake.nix
{
  outputs = { self, nixpkgs, uts, ... }: {
    darwinConfigurations.your-hostname = nix-darwin.lib.darwinSystem {
      modules = [
        # Apply the default overlay
        {
          nixpkgs.overlays = [ uts.overlays.default ];
        }
        # Add packages to system configuration
        {
          environment.systemPackages = [
            pkgs.uts         # Pre-compiled release binary
            # pkgs.uts-source # Build from Go source
          ];
        }
      ];
    };
  };
}
```

Or for standalone NixOS/home-manager configurations:

```nix
{
  home.packages = [
    uts.packages.${system}.default # Pre-compiled release
    # uts.packages.${system}.source  # Built from source
  ];
}
```

---

### Method 3: Go Install

If you have a Go development toolchain set up (Go 1.26+), install directly via Go:

```bash
go install github.com/y3owk1n/uts@latest
```

Ensure your Go bin directory (usually `~/go/bin`) is included in your system's `PATH`:

```bash
export PATH="$HOME/go/bin:$PATH"
```

### Method 4: GitHub Release Binaries

1. Navigate to the [Releases page](https://github.com/y3owk1n/uts/releases).
2. Download the package matching your OS and Architecture:
    - `uts-darwin-arm64.zip` (Apple Silicon Mac)
    - `uts-darwin-amd64.zip` (Intel Mac)
    - `uts-linux-amd64.tar.gz` (64-bit Linux)
    - `uts-linux-arm64.tar.gz` (ARM Linux)
3. Unpack the archive and move the `uts` binary into a directory in your `PATH` (e.g., `/usr/local/bin` or `~/.local/bin`).
4. Make the file executable:
    ```bash
    chmod +x /usr/local/bin/uts
    ```

---

## Verification

To verify that `uts` was installed correctly and can see your environment's utility tools, run the following validation:

```bash
# Check version
uts --version

# Run dry-run on info to see if dependencies are accessible
uts info somefile.png
```

---

## Troubleshooting

### "uts: command not found"

If you installed via `go install`, make sure your shell profile (`~/.zshrc`, `~/.bashrc`, or `~/.config/fish/config.fish`) exports the Go binary path:

```bash
# Zsh / Bash
export PATH="$HOME/go/bin:$PATH"
```

### "exec: \"ffmpeg\": executable file not found in $PATH"

This error indicates `uts` tried to run compression or conversion on a media file but could not find the required utility tool.

- Refer back to [System Dependencies](#system-dependencies-crucial) to install the missing tool for that category.
- Verify the tool is available in your current terminal session by running:
    ```bash
    which ffmpeg
    which convert # ImageMagick
    ```

### Nix Flake build fails on macOS

If you encounter errors building the flake locally, ensure you are running a supported Nix version with flake support enabled. Check your `~/.config/nix/nix.conf` or `/etc/nix/nix.conf` contains:

```conf
experimental-features = nix-command flakes
```

If issues persist, please open an issue in the [uts Github repository](https://github.com/y3owk1n/uts/issues).
