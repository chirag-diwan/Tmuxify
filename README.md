# Tmuxify

![Demo](./assets/demo.gif)

## Motivation

I always wanted a utility that could intelligently set up my tmux sessions. Different projects often require different layouts—opening specific windows, launching editors, running commands, and arranging everything exactly how I like it.

Tmuxify solves this by using a project-specific `.tmuxify.toml` file that describes how a tmux session should be created. Simply run tmuxify , select the directory and tmuxify will then do the rest for you.

## `.tmuxify.toml`

Each project contains a single `.tmuxify.toml` file that defines the tmux session. It consists of a `session` table and a list of `window` entries.

```toml
[session]
name = "tmuxify" # Session name (REQUIRED)
main = 0         # Window to select after creating the session (0-indexed)

# Every field in a window is OPTIONAL.

[[window]]
name = "nvim"
cmds = ["nvim ."]

[[window]]
name = "cmd"
cmds = ["ls"]
```

## Global Configuration

Tmuxify looks for a global configuration file at:

```text
~/.tmuxify-conf.toml
```

All options are optional. If an option is omitted, Tmuxify falls back to its built-in defaults.

Tmuxify also expects your home directory to be available through the `HOME` environment variable. If it cannot determine your home directory, it will use its internal defaults where possible (or return an error if required information is unavailable).

Example configuration:

```toml
roots = ["Projects", ".config", "Streaming"]
# Directories inside your home directory that Tmuxify will search.

ignore = [".git", "node_modules", ".cache", ".bun", ".cargo", ".wrangler"]
# Directory names to ignore anywhere in the search tree.

max_depth = 4
# Maximum directory depth to search.
```

## Installation

```bash
git clone https://github.com/chirag-diwan/Tmuxify.git
cd Tmuxify

go build -o out/tmuxify
mv out/tmuxify ~/.local/bin/
```

Make sure `~/.local/bin` is included in your `PATH`.

## Reporting Issues

If you encounter a bug , please open an issue and include the following information:

* What you expected to happen.
* What actually happened.
* Steps to reproduce the issue.
* Your operating system.
* Your Go version (if relevant).
