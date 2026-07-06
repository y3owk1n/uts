# Development Guide

This guide covers building, testing, linting, and contributing to the `uts` codebase.

---

## Table of Contents

- [Quick Start](#quick-start)
- [Prerequisites & Development Setup](#prerequisites--development-setup)
    - [Option A: Devbox (Recommended)](#option-a-devbox-recommended)
    - [Option B: Manual Setup](#option-b-manual-setup)
- [Project Architecture](#project-architecture)
- [Justfile Tasks (Build/Test/Lint)](#justfile-tasks-buildtestlint)
- [Testing Standards](#testing-standards)
    - [Unit Tests](#unit-tests)
    - [Integration Tests](#integration-tests)
    - [Code Coverage](#code-coverage)
- [Coding & Styling Standards](#coding--styling-standards)
- [Release Process](#release-process)

---

## Quick Start

Get `uts` building and running locally in under 5 minutes:

```bash
# 1. Clone the repository
git clone https://github.com/y3owk1n/uts.git
cd uts

# 2. Enter development shell
devbox shell  # Or use direnv automatically

# 3. Build the binary
just build

# 4. Run tests
just test

# 5. Execute local binary
./bin/uts --version
```

---

## Prerequisites & Development Setup

### Option A: Devbox (Recommended)

`uts` uses [Devbox](https://www.jetify.com/devbox) to provide a predictable, isolated development shell with all required toolchains pre-configured.

1. Install Devbox:
    ```bash
    curl -fsSL https://get.jetify.com/devbox | bash
    ```
2. Enter the development shell:
    ```bash
    devbox shell
    ```

_Tip:_ If you use `direnv`, the devbox environment will automatically load whenever you `cd` into the project directory via the provided `.envrc` config.

The Devbox environment automatically installs:

- **Go 1.26.4**
- **just** (command runner)
- **golangci-lint**
- formatting tools (`gofumpt`, `golines`)
- Go language server tools (`gopls`)

### Option B: Manual Setup

If you prefer installing tools manually, ensure you have:

- **Go 1.26+**
- **Just** command runner (`brew install just`)
- **golangci-lint** (`brew install golangci-lint`)

---

## Project Architecture

```
uts/
├── main.go               # Application entrypoint
├── cmd/                  # CLI Command definition (spf13/cobra)
│   ├── root.go           # Base command, persistent flags
│   ├── video.go          # Video subcommands (compress, convert)
│   ├── image.go          # Image subcommands (compress, convert)
│   ├── audio.go          # Audio subcommands (compress, convert)
│   ├── pdf.go            # PDF subcommands (compress, convert)
│   ├── archive.go        # Archive subcommands (compress, extract, list)
│   ├── convert.go        # Top-level convert shortcuts
│   ├── info.go           # File inspection command
│   └── genman/           # Auto-generates man documentation pages
├── internal/             # Core business logic packages
│   ├── compress/         # Media compression logic (invoking CLI tools)
│   ├── convert/          # Media format conversion logic
│   ├── archive/          # Listing/extracting zip, tar, zstd, etc.
│   ├── info/             # File-type detection & recommendations
│   ├── ui/               # Terminal rendering, spinners, colors (lipgloss)
│   ├── util/             # File, sizing, and preset conversion helpers
│   └── core/errors/      # Standard domain errors and exit codes
├── docs/                 # Project markdown guides
└── nix/                  # Declarative Nix building structures
```

---

## Justfile Tasks (Build/Test/Lint)

The project tasks are defined in the root `Justfile`. Here are the key commands:

| Task Command              | Description                                                                          |
| ------------------------- | ------------------------------------------------------------------------------------ |
| `just build`              | Compiles the binary into `bin/uts` using CGO-disabled flags.                         |
| `just test`               | Runs both unit and integration test suites.                                          |
| `just test-unit`          | Runs standard unit tests (`go test ./... -v`).                                       |
| `just test-integration`   | Runs tests tagged with the `integration` build tag.                                  |
| `just test-race`          | Runs unit & integration tests with Go's race detector active.                        |
| `just test-coverage`      | Generates a coverage file (`coverage.txt`) for unit tests.                           |
| `just test-coverage-all`  | Generates a coverage file containing integration tests as well.                      |
| `just test-coverage-html` | Generates and opens an HTML representation of the coverage in your browser.          |
| `just fmt`                | Standardizes formatting using `golangci-lint fmt` and fixes autofixable lint errors. |
| `just lint`               | Evaluates rules in `.golangci.yml`.                                                  |
| `just genman`             | Builds and runs the manpage generator to compile manual pages.                       |
| `just clean`              | Removes `bin/`, `build/`, and coverage text/html files.                              |

---

## Testing Standards

### Unit Tests

Unit tests are placed in `*_test.go` files beside the code they test.

- Must execute quickly.
- Do not rely on external command-line binaries (like `ffmpeg` or `imagemagick`).
- Code under test should mock out or decouple system execution logic where appropriate.
- We prefer **table-driven tests** (using structs and loops) to cover multiple boundary conditions cleanly.

Example:

```go
func TestHumanSize(t *testing.T) {
	tests := []struct {
		bytes    int64
		expected string
	}{
		{0, "0 B"},
		{1024, "1.00 KB"},
		{1048576, "1.00 MB"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := HumanSize(tt.bytes); got != tt.expected {
				t.Errorf("HumanSize(%d) = %q, want %q", tt.bytes, got, tt.expected)
			}
		})
	}
}
```

### Integration Tests

Integration tests verify that `uts` compiles the right command arguments and successfully invokes third-party binaries (e.g. compressing a real PNG or generating a video output).

- Must be marked with the build constraint:
    ```go
    //go:build integration
    ```
- Make sure to check that the necessary tool is installed on the host before running target commands, or skip gracefully if missing.
- Run integration tests using `just test-integration` or `just test`.

### Code Coverage

Maintain code coverage above 80% for utility, formatting, error mapping, and option validation packages.

- Check current coverage values using:
    ```bash
    just test-coverage-html
    ```

---

## Coding & Styling Standards

1. **Idiomatic Go**:
    - Follow formatting principles checked by `golangci-lint` (rules configured in `.golangci.yml`).
    - Use meaningful, explicit variable naming. Avoid short acronyms unless local context is extremely narrow.
2. **Terminal UI consistency**:
    - `uts` uses Charmbracelet's `lipgloss` library to format terminal text, lists, and output tables.
    - When printing messages, warnings, or panels, use the styling facades defined under `internal/ui/facade.go` (e.g. `ui.Banner()`, `ui.Panel()`, `ui.Message()`) rather than standard printing statements (`fmt.Printf`).
3. **Graceful Error Handling**:
    - Never bubble raw panic blocks up to the user.
    - Map standard library and operating system execution errors to the domain-specific errors in `internal/core/errors/errors.go` to provide high-quality terminal debugging advice.

---

## Release Process

`uts` utilizes Google's `release-please` to automate version increments, changelog compilation, and GitHub release publishing.

1. When code changes are merged into `main`, `release-please` reads conventional commit messages since the last release.
2. If it detects `feat` or `fix` commits, it creates/updates a Release PR.
3. Merging the Release PR tags the commit, creates a new GitHub Release, and triggers the `publish-artifacts.yml` GitHub action.
4. The publishing workflow builds pre-compiled binaries for Darwin and Linux (arm64 and amd64 architectures), bundles them with the auto-generated manpages, and uploads them to the release.
