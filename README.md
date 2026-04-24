# gity

Agentic GitHub issue creator powered by Ollama. Paste anything — a Slack message, an error log, meeting notes — and gity drafts a structured GitHub issue using a local AI model.

Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) + [Lipgloss](https://github.com/charmbracelet/lipgloss).

```
╔════════════[ ~/code/myproject ]═════╦════════════[ Issues (open) — 12 ]══╗
║ ../ ◂dir                            ║▶#4617  BE: Beta client rollout  BE  ║
║ .github/                            ║ #4616  QA: Production monitoring QA  ║
║ app/                                ║ #4593  Epic: Quote → Project    epic ║
║ services/                           ║ #4582  FE: Dark mode toggle          ║
║ ui/                                 ║ #4571  BE: Auth timeout fix          ║
║ go.mod                       1.2K  ║ #4560  QA: Mobile perf regression    ║
║ main.go                      3.4K  ║                                      ║
╚═════════════════════════════════════╩══════════════════════════════════════╝
  ↑/↓ nav   enter open   o toggle   r refresh   tab switch   c create   s settings   q quit
```

---

## Requirements

| Tool | Purpose | Min version |
|------|---------|-------------|
| [Ollama](https://ollama.ai) | Runs the AI model locally | 0.1.x |
| [GitHub CLI (`gh`)](https://cli.github.com) | Creates / lists issues | 2.x |
| [Go](https://go.dev/dl/) | Build from source | 1.21+ |

---

## Install

Follow the steps below to install each dependency, then build `gity`.

---

## 1 — Install Ollama

**macOS / Linux:**
```bash
curl -fsSL https://ollama.ai/install.sh | sh
```

**macOS (Homebrew):**
```bash
brew install ollama
```

**Windows:** download the installer from [ollama.ai](https://ollama.ai).

Start the daemon (runs in the background):
```bash
ollama serve
```

Pull the default model:
```bash
ollama pull llama3
```

> Other models work too — `mistral`, `llama3.1`, `gemma2`, etc. You can change the model later in the Settings screen.

Verify it's running:
```bash
curl http://localhost:11434/api/tags
```

---

## 2 — Install GitHub CLI

**macOS:**
```bash
brew install gh
```

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

**Windows:** `winget install --id GitHub.cli`

Authenticate:
```bash
gh auth login
```

Follow the prompts — select GitHub.com → HTTPS → Login with a web browser.

Verify:
```bash
gh auth status
```

---

## 3 — Install Go

**macOS:**
```bash
brew install go
```

**Linux / Windows:** download from [go.dev/dl](https://go.dev/dl/) and follow the installer.

Verify:
```bash
go version   # should print go1.21 or higher
```

---

## 4 — Install gity

Clone and build:
```bash
git clone https://github.com/tinydarkforge/gity.git
cd gity
go build -o gity .
```

Optionally move the binary somewhere on your PATH:
```bash
mv gity /usr/local/bin/gity       # macOS / Linux
# or on macOS:
mv gity /opt/homebrew/bin/gity
```

---

## 5 — Configure

### Option A — Settings screen (recommended)

Run `gity`, press `s` to open Settings, fill in your repo and model, press `ctrl+s` to save.

Config is stored at `~/.config/gity/config.json`.

### Option B — Environment variables

```bash
export GITY_REPO="owner/repo"          # e.g. acme/backend
export GITY_MODEL="llama3"             # any model you have pulled
export OLLAMA_HOST="http://localhost:11434"
export GITY_TIMEOUT=120               # seconds before ollama call times out
export GITY_MAX_TURNS=6               # max Q&A rounds before forcing a draft
export GITY_TEMPLATE_DIR=".github/ISSUE_TEMPLATES"
export GITY_DEBUG=1                   # dump raw model output to stderr
```

### Option C — CLI flags

```bash
gity -repo owner/repo -model mistral -ollama-host http://192.168.1.10:11434
```

Flags override env vars, which override the config file.

---

## 6 — Run

```bash
# make sure ollama is running first
ollama serve &

gity
# or with explicit repo
gity -repo owner/repo
```

---

## Usage

### Main view — Norton Commander layout

gity opens in a two-pane layout inspired by Norton Commander:

- **Left pane** — local filesystem browser, rooted at your current directory
- **Right pane** — GitHub issues list for your configured repo

Press `Tab` to switch between panes. The active pane has a brighter border and title.

#### Global keys (always available)

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
| `enter` | enter directory (or `..` to go up) |

#### Right pane — issues list

| Key | Action |
|-----|--------|
| `↑` / `k` | move up |
| `↓` / `j` | move down |
| `enter` | open issue detail |
| `o` | toggle open / closed issues |
| `r` | refresh list |

#### Shortcut bar

The bottom bar shows context-sensitive shortcuts for the current screen. In the main view:

`↑/↓` nav  `enter` open  `o` toggle open/closed  `r` refresh  `tab` switch pane  `c` create  `s` settings  `q` quit

---

### Create issue (`c`)

**Step 1 — Choose template**

Pick from your `.github/ISSUE_TEMPLATES` with `j/k`, confirm with `enter`.

**Step 2 — Title** *(optional)*

Type a short title, or leave blank — the agent will infer one from your context.

**Step 3 — Context (paste anything)**

This is the power step. Paste whatever you have:

- A Slack message: *"hey the login redirect is broken on staging, users get 403 after SSO"*
- An error log or stack trace
- A PR description you want converted to a task
- Rough meeting notes: *"discussed: need dark mode, auth timeout issue, mobile perf"*
- A one-liner: *"fix the thing that breaks when you click save twice"*

Press `ctrl+s` when ready.

> **Tip:** pastes longer than ~150 characters skip the Q&A round entirely — the agent drafts immediately from your context.

**Step 4 — Agent Q&A** *(short pastes only)*

For brief inputs the agent may ask 1–3 clarifying questions. Type your answer and press `enter`. Type `skip` to force an immediate draft with whatever info is available.

**Step 5 — Preview**

Review the generated title and body.

| Key | Action |
|-----|--------|
| `y` | create the GitHub issue |
| `r` | regenerate (asks agent to try again) |
| `esc` | cancel and start over |

---

### Issue detail (`enter` from list)

| Key | Action |
|-----|--------|
| `↑` / `k` | scroll up |
| `↓` / `j` | scroll down |
| `e` | **edit title and body** |
| `c` | add a comment |
| `a` | assign yourself |
| `x` | close issue |
| `o` | reopen issue |
| `esc` | back to list |

**Editing an issue**

Press `e` to open the edit view. The current title and body are pre-filled.

| Key | Action |
|-----|--------|
| `tab` | switch between title and body fields |
| `ctrl+s` | save changes to GitHub |
| `esc` | cancel, back to read view |

Changes are saved with `gh issue edit` and the issue reloads automatically.

---

### Settings (`s`)

| Field | Description |
|-------|-------------|
| GitHub Repo | `owner/repo` — where issues are created |
| Ollama Model | any model you have pulled (`llama3`, `mistral`, `gemma2`…) |
| Ollama Host | default `http://localhost:11434` — change for remote Ollama |
| Timeout (sec) | how long to wait for the model before giving up |
| Max Turns | maximum Q&A rounds (default 6) |
| Sound | `true` / `false` — audio feedback on macOS/Linux |
| Debug | `true` dumps raw model output to stderr |

Navigate fields with `tab` / `shift+tab`, edit inline, press `ctrl+s` to save.

---

## Issue templates

gity reads templates from `.github/ISSUE_TEMPLATES/*.md`. Each file needs YAML frontmatter:

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

The `labels` array is automatically applied when the issue is created. Any repo already using GitHub issue templates works out of the box.

---

## Troubleshooting

**`ollama is not running — start it with: ollama serve`**
```bash
ollama serve
```

**`gh: command not found`**
Install the GitHub CLI (step 2 above) and run `gh auth login`.

**Model produces garbled JSON**
Try a larger model: `ollama pull llama3.1` then set it in Settings. Smaller models (< 7B) sometimes struggle with strict JSON output.

**Templates not showing**
Make sure `GITY_TEMPLATE_DIR` points to a directory containing `*.md` files with YAML frontmatter. Default is `.github/ISSUE_TEMPLATES` relative to where you run gity.

**Sound not working on Linux**
Install `pulseaudio-utils` (`paplay`) or `alsa-utils` (`aplay`). Or toggle Sound off in Settings.

---

## Environment variable reference

| Variable | Default | Description |
|----------|---------|-------------|
| `GITY_REPO` | *(unset)* | Target GitHub repo, e.g. `owner/repo` |
| `GITY_MODEL` | `llama3` | Ollama model |
| `OLLAMA_HOST` | `http://localhost:11434` | Ollama daemon URL |
| `GITY_TIMEOUT` | `120` | Request timeout in seconds |
| `GITY_MAX_TURNS` | `6` | Max agent Q&A rounds |
| `GITY_TEMPLATE_DIR` | `.github/ISSUE_TEMPLATES` | Template directory |
| `GITY_DEBUG` | *(unset)* | Set to `1` to log raw model output |
| `NO_COLOR` | *(unset)* | Set to any value to disable ANSI colors |
| `GITY_CONFIG` | `~/.config/gity/config.json` | Override config file path |
