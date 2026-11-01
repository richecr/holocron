# Holocron (`hc`)

**Holocron** is a cross-platform script manager for your terminal. It keeps your
frequently used commands in a single TOML catalog and runs them with the right
shell on Linux, macOS and Windows — so `./mvnw clean package`, `docker compose
up -d` or `rmdir /s /q dist` all live under the same name, no matter the OS.

```console
$ hc docker:up
[INFO] Starting core services...
[INFO] Done.
```

## Features

- **One catalog, every OS** — each script can define a generic command and/or
  per-platform commands (`linux`, `darwin`, `windows`), resolved automatically
  at run time.
- **Run it directly** — `hc docker:up` is a shortcut for `hc run docker:up`.
- **Live streaming** — `stdout`, `stderr` and `stdin` are wired directly to
  your terminal, so long builds stream line by line and interactive scripts
  work.
- **Argument forwarding** — pass extra flags to any script:
  `hc build:java -- --debug`.
- **Tags & listing** — organize scripts with tags, filter them
  (`hc list --tag docker`) or consume the catalog as JSON (`hc list --json`).
- **Env vars & working dir** — each script can set its own environment
  variables and working directory in the TOML catalog.
- **Chained scripts** — a script can declare `pre_run` dependencies that run
  before it (cycle-safe): `hc deploy` runs `clean` first.
- **Search** — find scripts anywhere in names, descriptions, tags and
  commands: `hc search mvnw`.
- **Import & export** — back up, version in git, or sync catalogs between
  machines: `hc export > catalog.toml` / `hc import catalog.toml`.
- **Dry run** — preview the resolved command without executing it:
  `hc --dry build:java`.
- **Health check** — `hc doctor` validates shells, platforms and empty
  commands in the catalog.
- **Proper exit codes** — if the script fails with exit code `3`, `hc` exits
  with `3`.
- **Shell completion** — tab-completion for commands, scripts, `--tag` and
  `--platform` values.
- **Star Wars flavored UI** — Jedi cyan, Sith red and Holocron gold.

## Installation

### Pre-built binaries

Download the archive for your platform from the
[GitHub Releases](https://github.com/richecr/holocron/releases) page, extract
it and put `hc` on your `PATH`:

```sh
# Linux (amd64) example
curl -LO https://github.com/richecr/holocron/releases/download/v0.3.0/hc_0.3.0_linux_amd64.tar.gz
tar -xzf hc_0.3.0_linux_amd64.tar.gz
sudo mv hc_0.3.0_linux_amd64/hc /usr/local/bin/hc
```

| Asset | Platform |
|---|---|
| `hc_<version>_linux_amd64.tar.gz` | Linux x86_64 |
| `hc_<version>_linux_arm64.tar.gz` | Linux ARM64 |
| `hc_<version>_darwin_amd64.tar.gz` | macOS Intel |
| `hc_<version>_darwin_arm64.tar.gz` | macOS Apple Silicon |
| `hc_<version>_windows_amd64.zip` | Windows x86_64 (`hc.exe`) |
| `hc_<version>_windows_arm64.zip` | Windows ARM64 (`hc.exe`) |

Each archive also contains a `SHA256SUMS.txt` file to verify downloads.

### From source

Requires [Go](https://go.dev/dl/) 1.27+:

```sh
git clone https://github.com/richecr/holocron.git
cd holocron
make build        # produces ./hc  (or: go build -o hc .)
sudo mv hc /usr/local/bin/
```

### Add `hc` to your `PATH`

The examples above drop the binary in `/usr/local/bin`, which is already on
your `PATH` — that is all it takes to run plain `hc` instead of
`<path_to_binary>/hc`. If you prefer a different location, put the binary in
any directory and add that directory to your `PATH`:

```sh
# Example: user-writable directory, no sudo needed (Linux/macOS)
mkdir -p ~/.local/bin
mv hc ~/.local/bin/
```

Then make your shell look there:

**bash** — add to `~/.bashrc`:

```sh
export PATH="$HOME/.local/bin:$PATH"
```

**zsh** — add to `~/.zshrc`:

```sh
export PATH="$HOME/.local/bin:$PATH"
```

**fish** — one command, no config editing needed:

```sh
fish_add_path ~/.local/bin
```

**PowerShell (Windows)** — add for your user account (no admin needed):

```powershell
[Environment]::SetEnvironmentVariable(
  "Path",
  [Environment]::GetEnvironmentVariable("Path", "User") + ";$env:USERPROFILE\.local\bin",
  "User"
)
```

Alternatively on Windows: *Settings → System → About → Advanced system
settings → Environment Variables*, then append
`C:\Users\<you>\.local\bin` to your user `PATH`. Keep the file named
`hc.exe` — Windows finds it as plain `hc`.

Open a new terminal (or `source ~/.bashrc` / `source ~/.zshrc`) and verify:

```sh
hc --help
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
NAME       PLATFORMS               TAGS      DESCRIPTION
clean       linux   darwin   win             Remove build output
docker:up   any                     docker   Start core services

# Run a script (extra args are appended to the command)
hc docker:up              # direct execution
hc run docker:up          # same thing, explicit "run"
hc --dry docker:up        # preview the resolved command first
```

The first `hc` command creates the configuration file automatically
(see [Configuration](#configuration)).

## Commands

| Command | What it does |
|---|---|
| [`hc <script> [args...]`](#hc-run) | Run a script (shortcut for `hc run`) |
| [`hc run <script> [args...]`](#hc-run) | Run a script from the catalog |
| [`hc list`](#hc-list) | List the catalog as a table or JSON |
| [`hc add`](#hc-add) | Add or update a script |
| [`hc show <name>`](#hc-show) | Show the details of a script |
| [`hc remove <name>`](#hc-remove) | Remove a script (alias: `hc rm`) |
| [`hc search <term>`](#hc-search) | Search scripts by name, description, tag or command |
| [`hc doctor`](#hc-doctor) | Check the catalog for problems |
| [`hc init`](#hc-init) | Create a starter catalog with example scripts |
| [`hc config`](#hc-config) | Print or edit the config file path |
| [`hc export`](#hc-export) | Export the catalog as TOML |
| [`hc import`](#hc-import) | Import scripts from a TOML file |
| [`hc completion`](#hc-completion) | Generate shell completion scripts |

Every command also accepts `--help` (or `-h`) for its own usage text, and the
root command accepts `--version` to print the binary version.

**How invocation works:**

- The `run` keyword is optional: `hc build:java` and `hc run build:java` are
  equivalent.
- Built-in commands take precedence over scripts of the same name — a script
  named `add` would be shadowed by `hc add`.
- Flags go **before** the script name; everything after it is forwarded to the
  script (`hc docker:up --verbose` passes `--verbose` to the script).

### `hc run`

```sh
hc run <script> [args...]
hc <script> [args...]
```

Runs a script from the catalog, resolving the command for the current OS and
executing it in the configured shell.

| Flag | Description |
|---|---|
| `--dry` | Print the resolved command (including forwarded args) instead of running it |
| `--shell <shell>` | Run with this shell for one call (`sh`, `bash`, `zsh`, `powershell`, `cmd`), overriding `settings.default_shell` |

```sh
$ hc --dry build:java -- --debug
./mvnw clean package --debug

$ hc build:java -- --debug        # actually run it
```

Behavior:

- Extra arguments after the script name are appended to the resolved command —
  the `--` separator is optional: `hc test -- -v` and `hc test -v` are
  equivalent.
- Output is streamed live and the script inherits your terminal's stdin, so
  interactive scripts work.
- The script's exit code is propagated: a script failing with `3` makes
  `hc` exit with `3`.
- If `confirm_before_run` is enabled in the config, Holocron shows the
  resolved command and asks `Run <command> ? [y/N]` before executing.
- Typos get suggestions:
  `hc buld:java` → `script not found: buld:java (did you mean build:java?)`.
- Scripts with `env` / `dir` run with those environment variables and in that
  working directory (see [Configuration](#configuration)).
- Scripts with `pre_run` run their dependencies first, in order, each at most
  once: `hc deploy` with `pre_run = ["clean", "test"]` prints
  `pre: clean`, `pre: test` and then runs `deploy`. A failing dependency
  aborts the chain (its exit code is propagated). Unknown references and
  cycles are rejected before anything runs.

### `hc list`

```sh
hc list [--tag <tag>] [--group tag] [--json]
```

Renders the catalog as a table with the columns `NAME`, `PLATFORMS`, `TAGS`
and `DESCRIPTION`. Scripts with a generic command show the `any` badge;
platform-specific scripts show one badge per OS.

| Flag | Description |
|---|---|
| `--tag <tag>` | Only show scripts carrying this tag |
| `--group tag` | Group rows by tag (a script with several tags appears under each; untagged scripts under `(untagged)`) |
| `--json` | Print the catalog as a JSON array (for scripting, e.g. `\| jq`) |

```sh
$ hc list --tag docker
NAME       PLATFORMS  TAGS      DESCRIPTION
docker:up   any        docker   Start core services

$ hc list --group tag
build
NAME        PLATFORMS   TAGS     DESCRIPTION
build:java   any         java     Build Java

docker
NAME       PLATFORMS   TAGS      DESCRIPTION
docker:up   any         docker   Start core services

$ hc list --json
[
  {
    "name": "docker:up",
    "description": "Start core services",
    "tags": [
      "docker"
    ],
    "command": "docker compose up -d"
  }
]

$ hc list --json | jq -r '.[].name'
clean
docker:up
```

### `hc add`

```sh
hc add <name> "<command>" [--desc "..."] [--platform <os>]... [--tags a,b] [--pre x,y]
```

Adds a script to the catalog, or updates an existing one.

| Flag | Description |
|---|---|
| `--desc "..."` | Set the script description |
| `--platform <os>` | Store the command for this OS (`linux`, `darwin`, `windows`). Repeat the flag for multiple OSes |
| `--tags a,b` | Comma-separated list of tags (replaces existing tags) |
| `--pre x,y` | Comma-separated `pre_run` scripts, in order (replaces existing; pass `--pre ""` to clear). A warning is printed if a referenced script does not exist yet |

- Without `--platform`, the command is stored as the generic fallback
  (`command`), valid on every OS.
- With one or more `--platform` flags, the command is stored under
  `platforms["<os>"]` instead.
- Calling `hc add` again for an existing script **merges**: platform commands
  are added or updated per OS (e.g. one call for `linux`/`darwin` and another
  for `windows`), while description, tags and `pre_run` only change when
  `--desc` / `--tags` / `--pre` are passed. To replace a script entirely,
  edit it with `hc config edit`.

```sh
$ hc add docker:up "docker compose up -d" --tags docker
Saved docker:up  any

$ hc add deploy "./ship.sh" --pre test,clean
Saved deploy  any
```

The badges show where the command was stored: `any` for the generic command,
or one badge per OS.

### `hc show`

```sh
hc show <name>
```

Shows the details of one script: description, tags, the generic and
per-platform commands, and the command that would be resolved on the current
OS.

```sh
$ hc show clean
clean

  Description   Remove build output
  Tags

Platforms
   linux   rm -rf dist/
   darwin  rm -rf dist/
   win     rmdir /s /q dist

  On this OS (linux)  rm -rf dist/
```

Like `hc run`, misspelling the name prints a "did you mean" suggestion.

### `hc remove`

```sh
hc remove <name>    # or: hc rm <name>
```

Removes a script from the catalog.

```sh
$ hc remove docker:up
Removed docker:up
```

### `hc search`

```sh
hc search <term> [--json]
```

Case-insensitive search across script **names, descriptions, tags, commands
(generic and per-platform) and `dir`**. Results are shown in the same table
as `hc list`; `--json` prints them as a JSON array.

```sh
$ hc search mvnw
NAME        PLATFORMS   TAGS            DESCRIPTION
build:java   any         java   build   Build the Java project

$ hc search docker
No scripts match docker     # when nothing matches
```

### `hc doctor`

```sh
hc doctor
```

Validates the catalog and reports problems:

| Check | Severity |
|---|---|
| Invalid shell in `settings.default_shell` (accepted: `sh`, `bash`, `zsh`, `powershell`, `cmd`) | error |
| Script with no `command` and no `platforms` | error |
| Empty per-platform command | error |
| `pre_run` referencing an unknown script | error |
| `pre_run` cycle (e.g. `a -> b -> a`) | error |
| Unknown platform key (e.g. `windwos`) | warning |

Warnings do not fail the check; errors exit with code `1`.

```sh
$ hc doctor
OK 2 scripts, no issues found

$ hc doctor
error scripts."broken": no command and no platforms defined
warning scripts."typo-os".platforms: unknown platform "windwos" (expected linux, darwin or windows)
1 error(s), 1 warning(s)
```

### `hc init`

```sh
hc init [--yes]
```

Creates a starter catalog: if the config file does not exist yet it is
created; if it exists but is empty, three harmless example scripts are added
(never overwriting scripts you already have):

- `hello` — generic command, sanity check
- `clean` — per-platform command
- `deploy` — demonstrates `pre_run` chaining (`clean` runs first)

Asks for confirmation before seeding; pass `--yes` to skip the prompt.

```sh
$ hc init --yes
Initialized 3 example script(s) in ~/.config/holocron/config.toml

Next steps:
  hc list                 see the catalog
  hc hello                run the first example
  hc add my:script "echo hi"   add your own
```

### `hc config`

```sh
hc config          # print the absolute path of config.toml
hc config path     # same as above
hc config edit     # open config.toml in your editor
```

`config edit` uses `$VISUAL`, then `$EDITOR`, falling back to `notepad` on
Windows and `nano` elsewhere.

### `hc export`

```sh
hc export [--tag <tag>]
```

Prints the whole catalog (settings and scripts) as TOML to stdout. Pipe it to
a file to back it up or version it in git. With `--tag`, only the scripts
carrying that tag are exported.

```sh
hc export > holocron.toml
hc export --tag docker > docker-scripts.toml
```

### `hc import`

```sh
hc import <file>    # or: hc import -   (read from stdin)
```

Imports scripts from a TOML file into the catalog — typically one previously
produced by `hc export`.

- Imported scripts **replace** scripts with the same name; scripts that exist
  only in the current catalog are kept.
- The `[settings]` section of the file is ignored — import never touches your
  local settings.
- The incoming scripts are validated first; errors (e.g. a script with no
  command) abort the import before anything is saved.

```sh
# Sync the catalog from a git repo
git -C ~/dotfiles pull && hc import ~/dotfiles/holocron.toml

# Or pipe it
curl -s https://example.com/team.toml | hc import -
```

### `hc completion`

```sh
hc completion bash
hc completion zsh
hc completion fish
hc completion powershell
```

Generates a completion script for the given shell. Install it in your shell's
completion directory to get tab-completion for:

- all built-in commands and catalog script names (`hc <TAB>`),
- `--tag` values on `hc list`,
- `--platform` values on `hc add`.

For example, with bash:

```sh
hc completion bash | sudo tee /etc/bash_completion.d/hc > /dev/null
```

## Configuration

The catalog lives in a single TOML file:

| OS | Path |
|---|---|
| Linux | `~/.config/holocron/config.toml` |
| macOS | `~/Library/Application Support/holocron/config.toml` |
| Windows | `%APPDATA%\holocron\config.toml` |

> Paths follow `os.UserConfigDir()` (e.g. `$XDG_CONFIG_HOME` is respected on
> Linux).

The file is created automatically with defaults on first run.

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
# Optional per-script working directory, resolved relative to where you run hc.
dir = "./backend"

[scripts."build:java".platforms]
linux = "./mvnw clean package"
darwin = "./mvnw clean package"
windows = ".\\mvnw.cmd clean package"

# Optional per-script environment: added on top of your environment.
[scripts."build:java".env]
MAVEN_OPTS = "-Xmx1g"

# Optional: scripts to run before this one (in list order, each at most once).
# hc deploy runs test, then clean, then deploy. Cycles are rejected.
[scripts."deploy"]
description = "Run tests, then ship"
command = "./ship.sh"
pre_run = ["test", "clean"]

[scripts."test"]
command = "./mvnw test"

[scripts."clean"]
command = "rm -rf dist/"
```

### Command resolution

When a script runs, Holocron picks the command in this order:

1. `platforms[<current OS>]` — the platform-specific command, if non-empty.
2. `command` — the generic fallback.
3. Error — if neither matches the current platform.

### Shell selection

The shell used to execute the command is chosen in this order:

1. `settings.default_shell[<current OS>]`, if set.
2. OS default: `/bin/sh -c` on Linux/macOS,
   `powershell.exe -NoProfile -Command` on Windows.

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
