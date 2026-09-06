# ensor

Environment variables (secrets) manager + censoring codebase from leaks.

a small cli tool for managing and censoring secrets.

by default ensor is using dotenv file but you are free to use any file format
you want .

A small video demo of the tool.
[![ensor .](https://img.youtube.com/vi/3d1xLYkeVsA/hqdefault.jpg)](http://www.youtube.com/watch?v=3d1xLYkeVsA "ensor cli")

## manual

`ensor` that's it. Rest is handled by the tool itself. It will walk all the files
and directories and redact the secrets from the files.

by default, `ensor` skips scanning any files or directories that starts
with a trailing dot like `.env` `.env.json` `.git/*` .


how to install 

```bash
mise use github:pratyay360/ensor@latest
```

```bash
curl -sSL https://gobinaries.com/pratyay360/ensor@latest | sh
```

License: apache-2.0

thanks to amazing:
[betterleaks](https://github.com/betterleaks/betterleaks)
[huh](https://github.com/charmbracelet/huh)
[lipgloss](https://github.com/charmbracelet/lipgloss)
[bubbletea](https://github.com/charmbracelet/bubbletea)
