# Security Policy

## Supported Versions

Only the **latest release** receives security fixes. We do not back-port patches to older versions.

| Version        | Supported |
| -------------- | --------- |
| Latest release | ✅ Yes    |
| Older releases | ❌ No     |

---

## Reporting a Vulnerability

**Please do not open a public GitHub issue for security vulnerabilities.**

Instead, report them privately:

1. **GitHub Security Advisories (preferred)** — Go to [Security → Report a vulnerability](https://github.com/y3owk1n/uts/security/advisories/new) on this repository. This creates a private advisory only visible to maintainers.
2. **Direct contact** — Reach out to [@y3owk1n](https://github.com/y3owk1n) via GitHub if you cannot use the advisory flow.

Please include:

- A description of the vulnerability and its potential impact.
- Steps to reproduce or a proof-of-concept.
- The version(s) affected.
- Any suggested fix, if you have one.

### What to Expect

- **Acknowledgment** within **48 hours** of your report.
- A fix or mitigation plan within **7 days** for confirmed vulnerabilities.
- Credit in the release notes (unless you prefer to remain anonymous).

---

## Security Model

Understanding `uts`'s security posture helps frame what constitutes a vulnerability.

### Local-First Architecture

`uts` is a CLI wrapper designed to run entirely on your local machine. It:

- Does **not** spin up any local server or daemon.
- Does **not** make any outbound network connections.
- Does **not** send telemetry, analytics, or crash reports.
- Does **not** contact update servers or phone home.

All operations are executed synchronously or asynchronously in the foreground of your terminal session.

### Execution of External Binaries

`uts` delegates heavy lifting (like image processing, audio/video encoding, and PDF compilation) to third-party tools installed on your system (e.g., `ffmpeg`, `imagemagick`, `ghostscript`, `pngquant`).

- `uts` invokes these tools using standard Go system execution (`os/exec`).
- Ensure you source these third-party tools from trusted distributors (such as official Homebrew taps or official Nix packages).
- If a security vulnerability is found in the underlying encoder/tool itself (e.g., a buffer overflow in `ffmpeg`), that issue should be reported directly to that project's security team.

### Path Sanitization and Shell Injection

Since `uts` passes file paths and options as arguments to external CLI commands, we take path sanitization and argument safety seriously to prevent command or shell injection. If you find a way to execute arbitrary commands by supplying crafted filenames or flag options to `uts`, please report it immediately as a vulnerability.

### Dependencies

`uts` keeps a minimal dependency footprint, depending primarily on standard library packages and trusted libraries like `cobra` and Charmbracelet's styling tools (`lipgloss`, `bubbles`, `log`). Go modules are managed via `go.sum` for integrity verification.

---

## Scope

The following are **in scope** for security reports:

- Shell injection or execution of arbitrary code via malicious filenames, flags, or inputs.
- Privilege escalation or unexpected directory traversal during compression/extraction.
- Dependency vulnerabilities with a realistic exploit path.

The following are **out of scope**:

- Issues requiring the attacker to already have local write/execution access on your machine.
- Vulnerabilities within the third-party binary dependencies (e.g., `ffmpeg`, `imagemagick`, `ghostscript`) themselves.
- Denial-of-service/crashes when processing extremely large or corrupted media inputs.

---

Thank you for helping keep `uts` and its users safe.
