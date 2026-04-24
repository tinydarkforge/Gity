# Contributing to gity

Thanks for your interest in improving gity. This document covers the contribution workflow.

## Before you start

- Search [existing issues](https://github.com/tinydarkforge/Gity/issues) to avoid duplicates.
- For non-trivial changes, open an issue first to discuss scope.
- Security issues → see [SECURITY.md](SECURITY.md), do not open a public issue.

## Development setup

```bash
git clone https://github.com/tinydarkforge/Gity.git
cd Gity
go build -o gity .
go test ./...
```

Requirements: Go 1.22+, [Ollama](https://ollama.ai), [GitHub CLI](https://cli.github.com).

Pull a model for local testing:

```bash
ollama pull llama3
ollama serve
```

Run against a scratch repo:

```bash
./gity -repo your-user/sandbox-repo
```

## Pull request workflow

1. Fork and create a branch: `feat/short-description` or `fix/short-description`.
2. Make focused changes. One logical change per PR.
3. Run locally before pushing:
   ```bash
   go vet ./...
   go test ./...
   go build -o /tmp/gity-check .
   ```
4. Commit using [Conventional Commits](https://www.conventionalcommits.org/):
   - `feat: add template caching`
   - `fix: handle empty ollama response`
   - `docs: clarify install instructions`
5. Push and open a PR against `main`. Link the related issue.

## Commit style

- Header ≤72 chars, no trailing period.
- Valid types: `feat fix docs style refactor perf test chore ci security revert`.
- Body wrapped at 100 chars, explains **why** not **what**.

## Code style

- Standard Go formatting: `gofmt -w .`
- Prefer explicit error wrapping: `fmt.Errorf("loading config: %w", err)`.
- No comments explaining what the code does — only why when non-obvious.
- Keep packages focused: `services/` for external calls, `ui/` for Bubble Tea screens, `types/` for shared structs.

## Tests

- Unit tests in `*_test.go` next to code.
- Table-driven tests where practical.
- Mock Ollama / `gh` only at the boundary — prefer passing interfaces.

## Scope of contributions

**Welcome:**

- Bug fixes
- New issue-template features (frontmatter fields, label mapping)
- Additional Ollama model compatibility
- TUI polish (key bindings, themes, accessibility)
- Windows support (untested today)
- Documentation improvements

**Likely rejected without prior discussion:**

- New third-party service integrations (Jira, Linear, etc.) — scope creep
- Breaking changes to config format
- Replacing Bubble Tea with another TUI library

## License

By contributing you agree your contributions are licensed under the [MIT License](LICENSE).
