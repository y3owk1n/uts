# Contributing to uts

Thanks for your interest in contributing! `uts` is designed to be approachable and modular. We welcome contributions of all kinds — code, documentation, bug reports, or feature suggestions.

---

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Making Changes](#making-changes)
- [Commit Messages](#commit-messages)
- [Pull Requests](#pull-requests)
- [Testing](#testing)
- [Code Style](#code-style)
- [Good First Contributions](#good-first-contributions)
- [Reporting Bugs](#reporting-bugs)
- [Feature Requests](#feature-requests)

---

## Code of Conduct

This project follows our [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you agree to uphold it. Please report unacceptable behavior via [GitHub Issues](https://github.com/y3owk1n/uts/issues) or by contacting [@y3owk1n](https://github.com/y3owk1n) directly.

---

## Getting Started

1. **Search existing issues** — Check if someone is already working on the same thing or if there's a related discussion.
2. **Open an issue first** for non-trivial changes — This avoids wasted effort and lets us align on the approach before you write code.
3. **Small, focused PRs** are preferred over large, sweeping changes.

---

## Development Setup

See our [Development Guide](docs/DEVELOPMENT.md) for full instructions on setting up your dev environment using either Nix/Devbox or manual configurations.

---

## Making Changes

1. **Fork** the repository and clone your fork.
2. **Create a branch** from `main`:

    ```bash
    git checkout -b feat/my-feature
    ```

3. **Make your changes** following our styling and architecture guidelines.
4. **Add or update tests** for any new or changed functionality.
5. **Run the pre-commit checklist**:

    ```bash
    just fmt            # Format code using golangci-lint
    just lint           # Run golangci-lint checking rules
    just test           # Run unit & integration tests
    just build          # Verify the application builds successfully
    ```

6. **Commit** using conventional commit messages.
7. **Push** and open a pull request.

---

## Commit Messages

We use [Conventional Commits](https://www.conventionalcommits.org/) to power automated releases via [Release Please](https://github.com/googleapis/release-please).

**Format:**

```
<type>(<optional scope>): <subject>

<optional body>
<optional footer>
```

**Types:**

| Type       | When to use                            |
| ---------- | -------------------------------------- |
| `feat`     | New feature                            |
| `fix`      | Bug fix                                |
| `docs`     | Documentation only                     |
| `style`    | Formatting, no logic change            |
| `refactor` | Code restructuring, no behavior change |
| `perf`     | Performance improvement                |
| `test`     | Adding or updating tests               |
| `chore`    | Build, CI, dependencies, tooling       |

**Examples:**

```
feat(video): add Support for AV1 conversion
fix(image): fallback safely if optipng is missing
docs: correct example usage in CLI guide
```

---

## Pull Requests

- **Title** should follow the same conventional commit format (e.g., `feat(audio): add FLAC output option`).
- **Description** should explain _what_ changed and _why_. Include before/after details if the output UI changes.
- **Keep PRs focused** — One logical change per PR.
- **Link related issues** (e.g., `Closes #123`).
- All CI checks (lint, test, build) must pass before a pull request can be merged.

---

## Testing

`uts` separates tests into unit and integration tests:

| Type                | File pattern         | Command                 | Build tag / Behavior                                                                             |
| ------------------- | -------------------- | ----------------------- | ------------------------------------------------------------------------------------------------ |
| Unit tests          | `*_test.go`          | `just test-unit`        | Runs fast without third-party CLI dependency                                                     |
| Integration tests   | `*_test.go` with tag | `just test-integration` | Uses `//go:build integration` tag. Tests actual integration with tools (like ffmpeg/imagemagick) |
| Combined test suite | `*_test.go`          | `just test`             | Runs both unit and integration tests                                                             |

**Guidelines:**

- All utility, config-mapping, and core logic paths should have unit tests.
- Use **table-driven tests** where possible.
- Run `just test-coverage-html` to check coverage locally and view it in your browser.

---

## Code Style

All Go code must follow standard idiomatic style guidelines:

- Format code with `just fmt` (runs `golangci-lint fmt` and `golangci-lint run --fix`).
- Lint code with `just lint` (runs `golangci-lint run`).
- Add doc comments for all exported packages, structs, and functions.
- Keep terminal layouts matching the Charmbracelet styling established in the `internal/ui/` package.

---

## Good First Contributions

Not sure where to start? Check these options out:

- 🐛 Bug fixes — Check open issues.
- 📝 Documentation updates or improvements.
- 🧪 Adding tests for areas with low test coverage.
- ⚡ Performance optimizations or faster fallback detection.

---

## Reporting Bugs

Open a [GitHub Issue](https://github.com/y3owk1n/uts/issues/new) detailing:

1. **OS details** (macOS, Linux, etc.) and `uts` version (`uts --version`).
2. **Steps to reproduce** — A minimal reproducible CLI example.
3. **Expected vs actual behavior**.
4. **Verbosity output** — Run with the `-v` or `--verbose` flag and attach the logs.

---

## Feature Requests

Open a [GitHub Issue](https://github.com/y3owk1n/uts/issues/new) or start a discussion describing:

- **What** tool or format you'd like supported.
- **Why** it fits the `uts` single-predictable-API philosophy.
- **How** the options would map to the categories and quality presets.

---

Thank you for helping make `uts` better! 🙏
