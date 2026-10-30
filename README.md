# Holocron (`hc`)

**Holocron** is a cross-platform script manager for your terminal. It keeps your
frequently used commands in a single TOML catalog and runs them with the right
shell on Linux, macOS and Windows — so `./mvnw clean package`, `docker compose
up -d` or `rmdir /s /q dist` all live under the same name, no matter the OS.

```console
$ hc run build:java -- --debug
[INFO] Building with debug flags...
[INFO] BUILD SUCCESS
```

## Features

- **One catalog, every OS** — each script can define a generic command and/or
  per-platform commands (`linux`, `darwin`, `windows`), resolved automatically at run time.
- **Live streaming** — `stdout`, `stderr` and `stdin` are wired directly to your
  terminal, so long builds stream line by line and interactive scripts work.
- **Argument forwarding** — pass extra flags to any script: `hc run build:java -- --debug`.
- **Tags & listing** — organize scripts with tags and filter them: `hc list --tag docker`.
- **Proper exit codes** — if the script fails with exit code `3`, `hc` exits with `3`.
- **Star Wars flavored UI** — Jedi cyan, Sith red and Holocron gold.

## Installation

### Pre-built binaries

Download the archive for your platform from the
[GitHub Releases](https://github.com/richecr/holocron/releases) page, extract it
and put `hc` on your `PATH`:

```sh
# Linux (amd64) example
curl -LO https://github.com/richecr/holocron/releases/download/v0.1.0/hc_0.1.0_linux_amd64.tar.gz
tar -xzf hc_0.1.0_linux_amd64.tar.gz
sudo mv hc_0.1.0_linux_amd64/hc /usr/local/bin/hc
```

Available assets:

| Asset | Platform |
|---|---|
| `hc_<version>_linux_amd64.tar.gz` | Linux x86_64 |
| `hc_<version>_linux_arm64.tar.gz` | Linux ARM64 |
| `hc_<version>_darwin_amd64.tar.gz` | macOS Intel |
| `hc_<version>_darwin_arm64.tar.gz` | macOS Apple Silicon |
| `hc_<version>_windows_amd64.zip` | Windows x86_64 (`hc.exe`) |

Each archive also contains a `SHA256SUMS.txt` file to verify downloads.

### From source

Requires [Go](https://go.dev/dl/) 1.27+:

```sh
git clone https://github.com/richecr/holocron.git
cd holocron
make build        # produces ./hc  (or: go build -o hc .)
sudo mv hc /usr/local/bin/
```

## Quick start

```sh
# Add a generic script (same command on every OS)
hc add docker:up "docker compose up -d" --desc "Start core services" --tags docker

# Add a platform-specific script
hc add clean "rm -rf dist/" --platform linux --platform darwin --desc "Remove build output"
hc add clean "rmdir /s /q dist" --platform windows

# List everything in the catalog
hc list
# NAME         PLATFORMS        DESCRIPTION
# clean        linux darwin win Remove build output
# docker:up    any              Start core services

# Run a script (extra args are appended to the command)
hc run docker:up
hc run clean
```

The first `hc` command creates the configuration file automatically (see below).

## Commands

### `hc run <script> [args...]`

Runs a script from the catalog.

- Extra arguments after the script name are appended to the resolved command —
  `--` is optional: `hc run test -- -v` and `hc run test -v` are equivalent.
- Streams output live and inherits your terminal's stdin.
- Propagates the script's exit code.
- If `confirm_before_run` is enabled, shows the resolved command and asks
  `Run <command> ? [y/N]` first.

### `hc list [--tag <tag>]`

Renders the catalog as a table (`NAME`, `PLATFORMS`, `DESCRIPTION`).
Use `--tag` to filter scripts by tag. Scripts with a generic command show the
`any` badge; platform-specific scripts show one badge per OS.

### `hc add <name> "<command>"`

Adds or updates a script in the catalog.

| Flag | Description |
|---|---|
| `--desc "..."` | Script description |
| `--platform <os>` | Store the command for this OS (`linux`, `darwin`, `windows`). Repeatable |
| `--tags a,b` | Comma-separated tags |

Without `--platform`, the command is stored as the generic fallback (`command`).
With one or more `--platform` flags, it is stored under `platforms["<os>"]`.

Calling `hc add` again for an existing script **merges**: platform commands are
added or updated per OS (e.g. one call for `linux`/`darwin` and another for
`windows`), while description and tags only change when `--desc`/`--tags` are
passed. To replace a script entirely, edit it with `hc config edit`.

### `hc config [path|edit]`

- `hc config` or `hc config path` — prints the absolute path of `config.toml`.
- `hc config edit` — opens `config.toml` in `$VISUAL` / `$EDITOR`
  (falls back to `notepad` on Windows, `nano` elsewhere).

## Configuration

The catalog lives in a single TOML file:

| OS | Path |
|---|---|
| Linux | `~/.config/holocron/config.toml` |
| macOS | `~/Library/Application Support/holocron/config.toml` |
| Windows | `%APPDATA%\holocron\config.toml` |

> Paths follow `os.UserConfigDir()` (e.g. `$XDG_CONFIG_HOME` is respected on Linux).

### Full example

```toml
version = "1.0.0"

[settings]
# Ask for confirmation before running a script (hc run).
confirm_before_run = false

[settings.default_shell]
# Optional: override the default shell per OS.
# Accepted values: "sh", "bash", "zsh" (Linux/macOS), "powershell", "cmd" (Windows).
# When empty, the OS default is used: /bin/sh -c on Linux/macOS,
# powershell.exe -NoProfile -Command on Windows.
# linux = "bash"
# darwin = "zsh"
# windows = "powershell"

# ------------------------------------------------------------------
# Scripts
# ------------------------------------------------------------------

# Generic command: identical on every OS.
[scripts."docker:up"]
description = "Start core services with Docker Compose"
command = "docker compose up -d"
tags = ["docker"]

# Per-platform commands: picked automatically by hc at run time.
[scripts."build:java"]
description = "Build the Java project with Maven Wrapper"
tags = ["java", "build"]

[scripts."build:java".platforms]
linux = "./mvnw clean package"
darwin = "./mvnw clean package"
windows = ".\\mvnw.cmd clean package"
```

### Command resolution

When a script runs, Holocron picks the command in this order:

1. `platforms[<current OS>]` — the platform-specific command, if non-empty.
2. `command` — the generic fallback.
3. Error — if neither matches the current platform.

### Shell selection

The shell used to execute the command is chosen in this order:

1. `settings.default_shell[<current OS>]`, if set.
2. OS default: `/bin/sh -c` on Linux/macOS, `powershell.exe -NoProfile -Command` on Windows.

## Development

```sh
make help        # list all targets
make build       # build ./hc
make test        # run tests
make race        # tests + race detector + coverage (same as CI)
make lint        # golangci-lint
make fmt         # gofmt -w
make ci          # everything the CI runs
```

CI runs on every pull request (`.github/workflows/ci.yml`): `gofmt` check,
`go vet`, `golangci-lint` and `go test -race`. Releases
(`.github/workflows/release.yml`) build binaries for Linux, macOS and Windows
and attach them to the GitHub release.

See [PLAN.md](PLAN.md) for the implementation plan.
