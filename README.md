<!-- markdownlint-disable MD033 MD041 -->

```text
         ╔═══╗           █████ █   █ █████  ███  █   █ █████
    ╔════╩═══╩════╗        █   ██  █   █   █   █ █  █  █    
    ║  ─────────  ║        █   █ █ █   █   █████ ███   ████ 
    ╠═════════════╣        █   █  ██   █   █   █ █  █  █    
    ║ q w e r t y ║      █████ █   █   █   █   █ █   █ █████
    ║   a s d f   ║
    ╠═════════════╣      ━━━━━━━━━━━━━ ISSUE FORGE ━━━━━━━━━━━━━
    ║ [  space  ] ║      Ollama · GitHub CLI · Bubble Tea
    ╚═════════════╝      paste anything → structured issue,
    ═════════════════    one command. MIT · local AI · no cloud.
```

<p align="center">
  <a href="https://github.com/tinydarkforge/Intake/actions/workflows/ci.yml"><img alt="CI" src="https://github.com/tinydarkforge/Intake/actions/workflows/ci.yml/badge.svg"></a>
  <a href="https://github.com/tinydarkforge/Intake/releases/latest"><img alt="release" src="https://img.shields.io/github/v/release/tinydarkforge/Intake?style=flat-square&labelColor=0a0a0a&color=00cc66"></a>
  <a href="LICENSE"><img alt="license" src="https://img.shields.io/badge/license-MIT-00cc66.svg?style=flat-square&labelColor=0a0a0a"></a>
  <img alt="go" src="https://img.shields.io/badge/go-1.26%2B-00cc66.svg?style=flat-square&labelColor=0a0a0a">
  <a href="https://goreportcard.com/report/github.com/tinydarkforge/intake"><img alt="go report" src="https://goreportcard.com/badge/github.com/tinydarkforge/intake"></a>
  <a href="SECURITY.md"><img alt="security" src="https://img.shields.io/badge/security-policy-00cc66.svg?style=flat-square&labelColor=0a0a0a"></a>
</p>

> **Intake** is a terminal UI for GitHub issues powered by a local AI agent. Paste anything — a Slack message, an error log, meeting notes, a stack trace — and intake drafts a structured GitHub issue using a local Ollama model. No cloud. No account. No telemetry.

> **Status:** Early release (`v0.1.0`). MIT-licensed. macOS and Linux. Report vulnerabilities via [SECURITY.md](SECURITY.md).

---

## ░▒▓█ TL;DR

```bash
intake
```

Two-pane TUI: browse your filesystem left, manage GitHub issues right. Press `c` to create an issue — paste anything, the AI agent handles the rest.

---

## ░▒▓█ What it does today

intake wraps two things: a **Norton Commander-style filesystem + issues browser** and an **agentic issue drafter** powered by Ollama.

- **Two-pane layout** — local filesystem left, GitHub issues right
- **Agentic issue creation** — paste raw context; agent asks 1–3 questions if needed, then drafts from your template
- **Template-aware** — reads `.github/ISSUE_TEMPLATES/*.md` frontmatter; auto-applies labels on create
- **Full issue CRUD** — create, edit title+body, comment, assign yourself, close, reopen
- **Streaming drafts** — token-by-token output while the model thinks
- **No cloud dependency** — Ollama runs locally; `gh` handles GitHub auth

intake does not ship its own AI model. Every draft originates from whatever Ollama model you choose.

---

## ░▒▓█ Positioning

intake is **not** a project management tool, a GitHub web app replacement, or a SaaS product. It is a **local TUI gate** for the create-and-triage loop: go from messy context to a structured GitHub issue without leaving the terminal.

| Alternative | When to pick it instead of intake |
|-------------|----------------------------------|
| **GitHub web UI** | You prefer a browser, or need labels/milestones/projects UI |
| **`gh issue create`** | You already have a clean title and body |
| **Linear / Jira** | You need a managed tracker with workflows, not raw GitHub issues |

**intake's niche:** local AI, no account, no telemetry, instant context-to-issue from the terminal.

---

## ░▒▓█ Requirements

| Tool | Purpose | Min version |
|------|---------|-------------|
| [Ollama](https://ollama.ai) | Runs the AI model locally | 0.1.x |
| [GitHub CLI (`gh`)](https://cli.github.com) | Creates / lists issues | 2.x |
| [Go](https://go.dev/dl/) | Build from source only | 1.26+ |

**Platforms:** macOS (Intel + Apple Silicon), Linux (x86\_64 + arm64). Pre-built binaries on the [releases page](https://github.com/tinydarkforge/Intake/releases/latest). No Windows binary — build from source if needed.

---

## ░▒▓█ Install

### Pre-built binary (recommended)

Grab the tarball for your platform from the [latest release](https://github.com/tinydarkforge/Intake/releases/latest), extract, and move the binary onto your `$PATH`:

```bash
# example — macOS Apple Silicon
tar xz intake < intake_*_macos_arm64.tar.gz
mv intake /usr/local/bin/
```

Builds are provided for `macos_arm64`, `macos_x86_64`, `linux_arm64`, and `linux_x86_64`.

### go install

```bash
go install github.com/tinydarkforge/intake@latest
```

### Build from source

```bash
git clone https://github.com/tinydarkforge/intake.git
cd intake
go build -o intake .
mv intake /usr/local/bin/
```

---

## ░▒▓█ Dependencies

### Ollama

**macOS / Linux:**
```bash
curl -fsSL https://ollama.ai/install.sh | sh
# or: brew install ollama
```

Start the daemon:
```bash
ollama serve
```

Pull the default model:
```bash
ollama pull llama3
```

> Other models work — `mistral`, `llama3.1`, `gemma2`. Change the model in the Settings screen.

Verify:
```bash
curl http://localhost:11434/api/tags
```

### GitHub CLI

**macOS:** `brew install gh`

**Linux (Debian/Ubuntu):**
```bash
(type -p wget >/dev/null || (sudo apt update && sudo apt-get install wget -y)) \
  && sudo mkdir -p -m 755 /etc/apt/keyrings \
  && wget -qO- https://cli.github.com/packages/githubcli-archive-keyring.gpg \
     | sudo tee /etc/apt/keyrings/githubcli-archive-keyring.gpg > /dev/null \
  && echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" \
     | sudo tee /etc/apt/sources.list.d/github-cli.list > /dev/null \
  && sudo apt update && sudo apt install gh -y
```

Authenticate:
```bash
gh auth login
```

---

## ░▒▓█ Configure

### Option A — Settings screen (recommended)

Run `intake`, press `s`, fill in repo and model, press `ctrl+s` to save.

Config is stored at `~/.config/intake/config.json` (or `$INTAKE_CONFIG`).

### Option B — Environment variables

```bash
export INTAKE_REPO="owner/repo"
export INTAKE_MODEL="llama3"
export OLLAMA_HOST="http://localhost:11434"
export INTAKE_TIMEOUT=120
export INTAKE_MAX_TURNS=6
export INTAKE_TEMPLATE_DIR=".github/ISSUE_TEMPLATES"
export INTAKE_DEBUG=1
```

### Option C — CLI flags

```bash
intake -repo owner/repo -model mistral -ollama-host http://192.168.1.10:11434
intake -no-sound -debug
```

Flags override env vars, which override the config file.

---

## ░▒▓█ Run

```bash
# make sure Ollama is running first
ollama serve &

intake
# or with explicit repo
intake -repo owner/repo
```

---

## ░▒▓█ Usage

### Main view — Norton Commander layout

intake opens in a two-pane layout:

- **Left pane** — local filesystem browser, rooted at your working directory
- **Right pane** — GitHub issues list for your configured repo

Press `Tab` to switch panes. The active pane has a brighter border and title.

#### Global keys

| Key | Action |
|-----|--------|
| `Tab` | switch active pane |
| `c` | open Create issue screen |
| `s` | open Settings screen |
| `?` | show help |
| `q` | quit |
| `esc` | go back / cancel |

#### Left pane — filesystem

| Key | Action |
|-----|--------|
| `↑` / `k` | move up |
| `↓` / `j` | move down |
| `enter` | enter directory |

#### Right pane — issues list

| Key | Action |
|-----|--------|
| `↑` / `k` | move up |
| `↓` / `j` | move down |
| `enter` | open issue detail |
| `o` | toggle open / closed |
| `r` | refresh list |

---

### Create issue (`c`)

**Step 1 — Choose template**

Pick from `.github/ISSUE_TEMPLATES` with `j/k`, confirm with `enter`.

**Step 2 — Title** *(optional)*

Type a short title, or leave blank — the agent infers one from your context.

**Step 3 — Context (paste anything)**

This is the power step. Paste whatever you have:

- A Slack message: *"login redirect broken on staging, users get 403 after SSO"*
- An error log or stack trace
- A PR description you want converted to a task
- Meeting notes: *"discussed: dark mode, auth timeout issue, mobile perf"*
- A one-liner: *"fix the thing that breaks when you click save twice"*

Press `ctrl+s` when ready.

> **Tip:** pastes longer than ~150 characters skip the Q&A round entirely — the agent drafts immediately from your context.

**Step 4 — Agent Q&A** *(short pastes only)*

For brief inputs the agent may ask 1–3 clarifying questions. Answer and press `enter`. Type `skip` to force an immediate draft.

**Step 5 — Preview**

| Key | Action |
|-----|--------|
| `y` | create the GitHub issue |
| `r` | regenerate |
| `esc` | cancel |

---

### Issue detail (`enter`)

| Key | Action |
|-----|--------|
| `↑` / `k` | scroll up |
| `↓` / `j` | scroll down |
| `e` | edit title and body |
| `c` | add a comment |
| `a` | assign yourself |
| `x` | close issue |
| `o` | reopen issue |
| `esc` | back to list |

**Editing:** press `e`, use `tab` to switch between title and body, `ctrl+s` to save. Changes are applied with `gh issue edit` and the issue reloads automatically.

---

### Settings (`s`)

| Field | Description |
|-------|-------------|
| GitHub Repo | `owner/repo` — where issues are created |
| Ollama Model | any model you have pulled (`llama3`, `mistral`, `gemma2`…) |
| Ollama Host | default `http://localhost:11434` — change for remote Ollama |
| Timeout (sec) | how long to wait for the model before giving up |
| Max Turns | max Q&A rounds (default 6) |
| Sound | audio feedback on macOS/Linux |
| Debug | dumps raw model output to stderr |

Navigate with `tab` / `shift+tab`, press `ctrl+s` to save.

---

## ░▒▓█ Issue templates

intake reads templates from `.github/ISSUE_TEMPLATES/*.md`. Each file needs YAML frontmatter:

```markdown
---
name: "🐞 Bug Report"
about: "Report an issue"
title: "🐞 [Brief title]"
labels: ["bug"]
assignees: []
---

## Summary
...
```

Labels are auto-applied on create. Any repo already using GitHub issue templates works out of the box.

---

## ░▒▓█ Troubleshooting

**`ollama is not running — start it with: ollama serve`**
```bash
ollama serve
```

**`gh: command not found`**
Install the GitHub CLI (above) and run `gh auth login`.

**Model produces garbled JSON**
Try a larger model: `ollama pull llama3.1`, then set it in Settings. Small models (< 7B) sometimes struggle with strict JSON output.

**Templates not showing**
Ensure `INTAKE_TEMPLATE_DIR` points to a directory with `*.md` files that have YAML frontmatter. Default is `.github/ISSUE_TEMPLATES` relative to where you run intake.

**Sound not working on Linux**
Install `pulseaudio-utils` (`paplay`) or `alsa-utils` (`aplay`). Or toggle Sound off in Settings.

---

## ░▒▓█ Env reference

| Variable | Default | Description |
|----------|---------|-------------|
| `INTAKE_REPO` | *(unset)* | Target GitHub repo, e.g. `owner/repo` |
| `INTAKE_MODEL` | `llama3` | Ollama model |
| `OLLAMA_HOST` | `http://localhost:11434` | Ollama daemon URL |
| `INTAKE_TIMEOUT` | `120` | Request timeout in seconds |
| `INTAKE_MAX_TURNS` | `6` | Max agent Q&A rounds |
| `INTAKE_TEMPLATE_DIR` | `.github/ISSUE_TEMPLATES` | Template directory |
| `INTAKE_DEBUG` | *(unset)* | Set to `1` to log raw model output |
| `NO_COLOR` | *(unset)* | Disable ANSI colors |
| `INTAKE_CONFIG` | `~/.config/intake/config.json` | Override config file path |
