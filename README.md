# ensor

Environment variables (secrets) manager + censoring codebase from leaks.

a small cli tool for managing and censoring secrets.

by default ensor is using dotenv file but you are free to use any file format
you want .

A small video demo of the tool.
[![ensor](https://img.youtube.com/vi/3d1xLYkeVsA/hqdefault.jpg)](http://www.youtube.com/watch?v=3d1xLYkeVsA "ensor cli")

## manual

`ensor` that's it. Rest is handled by the tool itself. It will walk all the files
and directories and redact the secrets from the files.

by default, `ensor` skips scanning any files or directories that starts
with a trailing dot like `.env` `.env.json` `.git/*` .


how to install ?

binaries are available in the release section

```bash
mise use github:pratyay360/ensor@latest
```


```bash
curl -sSL https://raw.githubusercontent.com/Pratyay360/ensor/main/install.sh | sh
```

# ensor live demo

[![asciicast](https://asciinema.org/a/1264776.svg)](https://asciinema.org/a/1264776)

## ensor convert
ensor can convert between any format of `env` files. Convert `.env`, YAML, TOML, JSON, etc interchangeably.

[![asciicast](https://asciinema.org/a/1264777.svg)](https://asciinema.org/a/1264777)

License: apache-2.0

thanks to amazing:
[betterleaks](https://github.com/betterleaks/betterleaks)
[huh](https://github.com/charmbracelet/huh)
[lipgloss](https://github.com/charmbracelet/lipgloss)
[bubbletea](https://github.com/charmbracelet/bubbletea)
