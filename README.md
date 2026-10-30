
<p align="center">
  <img alt="header" src="https://shieldcn.dev/header/surface.svg?title=ENSOR&amp;subtitle=Censor+secrets+from+code&amp;mode=dark" />
<a href="https://www.producthunt.com/products/ensor-2/reviews/new?utm_source=badge-product_review&utm_medium=badge&utm_source=badge-ensor&#0045;2" target="_blank"><img src="https://api.producthunt.com/widgets/embed-image/v1/product_review.svg?product_id=1361958&theme=light" alt="ensor - cli&#0032;tool&#0032;to&#0032;detect&#0032;and&#0032;redact&#0032;secrets | Product Hunt" style="width: 250px; height: 54px;" width="250" height="54" /></a>
</p>

**ensor** is a fast, interactive CLI tool for detecting and redacting secrets from codebase.
By default, `ensor` extracts discovered secrets into a `.env` file and substitutes them with customizable placeholders across your project, but also includes a versatile converter for `.env`, JSON, YAML, and TOML.

---

## ✨ Features

> scan for detecting secrets
> redact secrets from codebase
> put secrets into a `.env` file
> convert between different secret formats

## 🎬 Demo


Watch `ensor` in action:

https://github.com/user-attachments/assets/e817c8d7-f4eb-4eee-9cf1-f26c14f055c0

[![ensor](https://img.youtube.com/vi/3d1xLYkeVsA/hqdefault.jpg)](http://www.youtube.com/watch?v=3d1xLYkeVsA "ensor cli")

## 📦 Installation

### Homebrew

```bash
brew trust --cask pratyay360/tap/ensor
brew install pratyay360/tap/ensor
```

### mise

```bash
mise use github:pratyay360/ensor@latest
```

### curl (automated install script)

```bash
curl -sSL https://raw.githubusercontent.com/Pratyay360/ensor/main/install.sh | sh
```

### goblin

```bash
curl -sf http://goblin.run/github.com/pratyay360/ensor | PREFIX=~/.local/bin sh 
```

### Go Install

```bash
go install github.com/Pratyay360/ensor@latest
```

### Scoop
```bash
scoop bucket add pratyay360 https://github.com/pratyay360/scoop-bucket
scoop install pratyay360/ensor
```

### Winget (not yet available)
```bash
winget install pratyay360.ensor
```

### Aqua registry

```bash
aqua install Pratyay360/ensor
```



### Prebuilt Binaries & Packages

Download ready-to-run binaries or `.deb`, `.rpm`, and `.apk` packages from the 
[Releases](https://github.com/Pratyay360/ensor/releases) page for Linux, macOS, Windows, and FreeBSD.

---

## 🚀 Usage

### 1. Scan and Censor Secrets

Scan the current directory:

```bash
ensor
```

interactive mode feels slow and boring try.

```bash
ensor --yolo
```

for automatic redaction of secrets from the entire codebase.

That's it. `ensor` walks all the files and directories in codebase and detects + redacts the secrets from there.

By default, `ensor` skips scanning any file or directory that starts with a trailing dot, like `.env`, `.env.json`, or `.git/*` — so your actual secret files stay untouched.

> Run `ensor --help` to see all available options.

## Convert env files

`ensor convert` converts between any format of env files. Convert `.env`, YAML, TOML, JSON, and more — interchangeably.

## Live demos

**ensor in action**

[![asciicast](https://asciinema.org/a/1270301.svg)](https://asciinema.org/a/1270301)

**ensor convert**


[![asciicast](https://asciinema.org/a/1264777.svg)](https://asciinema.org/a/1264777)



**result ensor**

[![asciicast](https://asciinema.org/a/1270303.svg)](https://asciinema.org/a/1270303)

---

![chart](https://shieldcn.dev/chart/github/stars/Pratyay360/ensor.svg?theme=emerald&bg=transparent)


[![releaseszn](https://releaseszn.com/badge/ensor.svg)](https://releaseszn.com/p/ensor)
