#!/usr/bin/env bash
# intake — guided installer.
# Installs Ollama, GitHub CLI, the intake binary, then walks through first-run config.
# Safe to re-run: every step checks before installing.

set -euo pipefail

# ---------- styling ----------
if [[ -t 1 ]] && command -v tput >/dev/null 2>&1 && [[ "${NO_COLOR:-}" != "1" ]]; then
  BOLD="$(tput bold)"; DIM="$(tput dim)"; RESET="$(tput sgr0)"
  GREEN="$(tput setaf 2)"; YELLOW="$(tput setaf 3)"; BLUE="$(tput setaf 4)"; RED="$(tput setaf 1)"
else
  BOLD=""; DIM=""; RESET=""; GREEN=""; YELLOW=""; BLUE=""; RED=""
fi

step()  { printf "\n${BOLD}${BLUE}==>${RESET} ${BOLD}%s${RESET}\n" "$*"; }
ok()    { printf "    ${GREEN}✓${RESET} %s\n" "$*"; }
warn()  { printf "    ${YELLOW}!${RESET} %s\n" "$*"; }
err()   { printf "    ${RED}✗${RESET} %s\n" "$*" >&2; }
ask()   { printf "${BOLD}?${RESET} %s " "$*"; }

trap 'err "install failed on line $LINENO. Re-run the script — every step is idempotent."' ERR

# ---------- banner ----------
cat <<'EOF'

         ╔═══╗           █████ █   █ █████  ███  █   █ █████
    ╔════╩═══╩════╗        █   ██  █   █   █   █ █  █  █
    ║  ─────────  ║        █   █ █ █   █   █████ ███   ████
    ╠═════════════╣        █   █  ██   █   █   █ █  █  █
    ║ q w e r t y ║      █████ █   █   █   █   █ █   █ █████
    ║   a s d f   ║
    ╠═════════════╣      ━━━━━━━━━━━━━ ISSUE FORGE ━━━━━━━━━━━━━
    ║ [  space  ] ║      guided install — no dev experience required
    ╚═════════════╝
    ═════════════════

This script will:
  1. Check / install Ollama (local AI runtime)
  2. Check / install GitHub CLI and sign you in
  3. Download the intake binary onto your computer
  4. Pull a default AI model (llama3, ~4.7 GB)
  5. Ask which GitHub repo intake should target

You can stop with Ctrl+C at any time and re-run later.
EOF

ask "Continue? [Y/n]"
read -r CONFIRM
case "${CONFIRM:-y}" in
  y|Y|yes|YES) ;;
  *) echo "Cancelled."; exit 0 ;;
esac

# ---------- platform detection ----------
step "Detecting platform"
UNAME_S="$(uname -s)"; UNAME_M="$(uname -m)"
case "$UNAME_S" in
  Darwin) OS="macos" ;;
  Linux)  OS="linux" ;;
  *) err "unsupported OS: $UNAME_S (intake supports macOS and Linux)"; exit 1 ;;
esac
case "$UNAME_M" in
  x86_64|amd64) ARCH="x86_64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  *) err "unsupported architecture: $UNAME_M"; exit 1 ;;
esac
ok "platform: $OS/$ARCH"

# pick install dir — prefer /usr/local/bin if writable, else ~/.local/bin
if [[ -w /usr/local/bin ]]; then
  BIN_DIR="/usr/local/bin"; SUDO=""
elif command -v sudo >/dev/null 2>&1; then
  BIN_DIR="/usr/local/bin"; SUDO="sudo"
else
  BIN_DIR="$HOME/.local/bin"; SUDO=""
  mkdir -p "$BIN_DIR"
fi
ok "binary destination: $BIN_DIR"

# ---------- Ollama ----------
step "Checking Ollama"
if command -v ollama >/dev/null 2>&1; then
  ok "ollama already installed ($(ollama --version 2>/dev/null | head -1))"
else
  warn "ollama not found — installing from ollama.ai"
  if [[ "$OS" == "macos" ]] && command -v brew >/dev/null 2>&1; then
    brew install ollama
  else
    curl -fsSL https://ollama.ai/install.sh | sh
  fi
  ok "ollama installed"
fi

# start daemon if not running
if ! curl -fsS --max-time 2 http://localhost:11434/api/tags >/dev/null 2>&1; then
  warn "ollama daemon not running — starting in background"
  if [[ "$OS" == "macos" ]] && [[ -d "/Applications/Ollama.app" ]]; then
    open -a Ollama || true
  else
    nohup ollama serve >/tmp/ollama.log 2>&1 &
  fi
  for i in 1 2 3 4 5 6 7 8 9 10; do
    sleep 1
    if curl -fsS --max-time 2 http://localhost:11434/api/tags >/dev/null 2>&1; then
      ok "ollama daemon up"
      break
    fi
    if [[ $i -eq 10 ]]; then
      err "ollama daemon did not start. Run 'ollama serve' manually, then re-run this script."
      exit 1
    fi
  done
else
  ok "ollama daemon running"
fi

# ---------- model ----------
step "Pulling default AI model (llama3)"
if ollama list 2>/dev/null | awk 'NR>1 {print $1}' | grep -q '^llama3'; then
  ok "llama3 already pulled"
else
  warn "downloading ~4.7 GB — this can take a while on slow connections"
  ollama pull llama3
  ok "llama3 ready"
fi

# ---------- gh ----------
step "Checking GitHub CLI"
if command -v gh >/dev/null 2>&1; then
  ok "gh already installed ($(gh --version | head -1))"
else
  warn "gh not found — installing"
  if [[ "$OS" == "macos" ]]; then
    if command -v brew >/dev/null 2>&1; then
      brew install gh
    else
      err "Homebrew not found. Install brew (https://brew.sh) or 'gh' manually, then re-run."
      exit 1
    fi
  else
    if command -v apt-get >/dev/null 2>&1; then
      $SUDO mkdir -p -m 755 /etc/apt/keyrings
      curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg \
        | $SUDO tee /etc/apt/keyrings/githubcli-archive-keyring.gpg > /dev/null
      $SUDO chmod go+r /etc/apt/keyrings/githubcli-archive-keyring.gpg
      echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" \
        | $SUDO tee /etc/apt/sources.list.d/github-cli.list > /dev/null
      $SUDO apt-get update
      $SUDO apt-get install -y gh
    elif command -v dnf >/dev/null 2>&1; then
      $SUDO dnf install -y gh
    elif command -v pacman >/dev/null 2>&1; then
      $SUDO pacman -S --noconfirm github-cli
    else
      err "Could not detect a supported package manager. Install gh from https://cli.github.com and re-run."
      exit 1
    fi
  fi
  ok "gh installed"
fi

# auth
if gh auth status >/dev/null 2>&1; then
  ok "gh already signed in"
else
  warn "gh not signed in — launching 'gh auth login'"
  echo "    Choose: GitHub.com → HTTPS → Yes (auth git) → Login with a web browser"
  gh auth login
  ok "gh signed in"
fi

# ---------- intake binary ----------
step "Installing intake"
RELEASE_JSON="$(curl -fsSL https://api.github.com/repos/tinydarkforge/Intake/releases/latest)" || {
  err "Could not reach GitHub API. Check your internet connection."
  exit 1
}
LATEST_TAG="$(printf '%s' "$RELEASE_JSON" | grep -oE '"tag_name":\s*"[^"]+"' | head -1 | sed -E 's/.*"([^"]+)"/\1/')"
[[ -z "$LATEST_TAG" ]] && { err "Could not parse latest release tag."; exit 1; }

# match any asset whose name ends with _<OS>_<ARCH>.tar.gz — survives project renames
ASSET_URL="$(printf '%s' "$RELEASE_JSON" \
  | grep -oE '"browser_download_url":\s*"[^"]+\.tar\.gz"' \
  | sed -E 's/.*"([^"]+)"/\1/' \
  | grep -E "_${OS}_${ARCH}\.tar\.gz$" \
  | head -1)"
if [[ -z "$ASSET_URL" ]]; then
  err "No release asset found for ${OS}/${ARCH} in $LATEST_TAG."
  err "Check the releases page: https://github.com/tinydarkforge/Intake/releases"
  exit 1
fi
TARBALL="$(basename "$ASSET_URL")"
TMPDIR="$(mktemp -d)"
trap 'rm -rf "$TMPDIR"' EXIT

ok "fetching $LATEST_TAG → $TARBALL"
curl -fsSL "$ASSET_URL" -o "$TMPDIR/$TARBALL" || {
  err "download failed: $ASSET_URL"
  exit 1
}
tar -xzf "$TMPDIR/$TARBALL" -C "$TMPDIR"

# binary inside tarball may be 'intake' or legacy 'gity'
BIN_SRC=""
for candidate in intake gity; do
  if [[ -f "$TMPDIR/$candidate" ]]; then BIN_SRC="$TMPDIR/$candidate"; break; fi
done
[[ -z "$BIN_SRC" ]] && { err "No intake/gity binary found inside $TARBALL"; exit 1; }
$SUDO install -m 0755 "$BIN_SRC" "$BIN_DIR/intake"
ok "intake installed → $BIN_DIR/intake"

# PATH check
if ! echo ":$PATH:" | grep -q ":$BIN_DIR:"; then
  warn "$BIN_DIR is not in your PATH"
  warn "add this to your shell profile (~/.zshrc or ~/.bashrc):"
  echo "        export PATH=\"$BIN_DIR:\$PATH\""
fi

# ---------- config ----------
step "Configuring intake"
CONFIG_DIR="${XDG_CONFIG_HOME:-$HOME/.config}/intake"
CONFIG_FILE="$CONFIG_DIR/config.json"
mkdir -p "$CONFIG_DIR"

DEFAULT_REPO=""
if gh repo view --json nameWithOwner -q .nameWithOwner >/dev/null 2>&1; then
  DEFAULT_REPO="$(gh repo view --json nameWithOwner -q .nameWithOwner 2>/dev/null || true)"
fi

while true; do
  if [[ -n "$DEFAULT_REPO" ]]; then
    ask "Which GitHub repo should intake target? [${DEFAULT_REPO}]"
  else
    ask "Which GitHub repo should intake target? (e.g. owner/name)"
  fi
  read -r REPO_INPUT
  REPO="${REPO_INPUT:-$DEFAULT_REPO}"
  if [[ "$REPO" =~ ^[A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+$ ]]; then
    break
  fi
  warn "format must be owner/name — try again"
done

cat > "$CONFIG_FILE" <<JSON
{
  "repo": "$REPO",
  "model": "llama3",
  "ollama_host": "http://localhost:11434",
  "timeout": 120,
  "max_turns": 6,
  "template_dir": ".github/ISSUE_TEMPLATES",
  "sound": true
}
JSON
ok "config written → $CONFIG_FILE"

# ---------- done ----------
cat <<EOF

${GREEN}${BOLD}Done.${RESET}

  Run intake:        ${BOLD}intake${RESET}
  Open settings:     press ${BOLD}s${RESET} inside intake
  Create an issue:   press ${BOLD}c${RESET} inside intake

Tips:
  • Ollama must be running. If you reboot, run: ${BOLD}ollama serve${RESET}
  • Switch model anytime in Settings (e.g. mistral, llama3.1, gemma2)
  • Config lives at: $CONFIG_FILE

Issues / questions: https://github.com/tinydarkforge/Intake/issues
EOF
