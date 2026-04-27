# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] — 2026-04-24

### Added

- Norton Commander-style two-pane TUI built on Bubble Tea + Lipgloss
- Local filesystem browser (left pane)
- GitHub issues list with open / closed toggle (right pane)
- Agentic issue creation powered by local Ollama models
- Template-aware drafting — reads `.github/ISSUE_TEMPLATES/*.md` frontmatter
- Auto-applies template labels on issue creation
- Full issue CRUD: create, edit, comment, assign, close, reopen
- Multi-step create flow: template → title → context → agent Q&A → preview
- Configurable via settings screen, env vars (`INTAKE_*`), or CLI flags
- Config persisted at `~/.config/intake/config.json`
- Optional audio feedback on macOS / Linux
- Debug mode dumping raw model output to stderr
- GitHub Actions CI: `go build` + `go test` + `go vet` on macOS and Ubuntu
- MIT License
- Security disclosure policy
- Contributor guidelines and Code of Conduct

[Unreleased]: https://github.com/tinydarkforge/Intake/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/tinydarkforge/Intake/releases/tag/v0.1.0
