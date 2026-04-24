# Security Policy

## Supported Versions

Only the latest tagged release receives security fixes.

## Reporting a Vulnerability

Do **not** open a public issue for security reports.

Report privately via GitHub's "Report a vulnerability" button in the
[Security tab](https://github.com/tinydarkforge/gity/security/advisories/new),
or email `security@tinydarkforge.dev`.

Please include:

- Affected version / commit
- Reproduction steps
- Impact assessment
- Suggested fix, if known

You will receive an acknowledgement within 72 hours. Fixes are coordinated
under a 90-day disclosure window unless the issue is actively exploited.

## Scope

In scope:

- `gity` binary and source
- Issue template parsing, Ollama prompt injection, `gh` invocation paths
- Config file handling

Out of scope:

- Vulnerabilities in upstream dependencies (Ollama, `gh` CLI, Bubble Tea) —
  report those to their respective maintainers.
- Local-only issues that require attacker already having shell access.
