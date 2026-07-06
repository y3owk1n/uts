# Configuration

`uts` can be configured via environment variables. All configuration is optional — the defaults are designed to work well out of the box.

---

## Color Customization

`uts` uses adaptive colors that automatically adjust for light and dark terminals. You can fully customize the color palette using environment variables.

### Palette Colors

Each color role can be overridden with a hex value (e.g. `#FF5733`) or ANSI color name:

| Environment Variable | Role | Default (Light / Dark) |
| --- | --- | --- |
| `UTS_COLOR_PRIMARY` | Logo, titles, active elements | `#6f4d8c` / `#c9a0e9` |
| `UTS_COLOR_TEXT` | Body text | `#2a2738` / `#e0def4` |
| `UTS_COLOR_MUTED` | Dimmed / secondary text | `#5a5672` / `#9a96b5` |
| `UTS_COLOR_SUBTLE` | Hints and subtle labels | `#9a96b5` / `#5a5672` |
| `UTS_COLOR_BORDER` | Panel borders and separators | `#2a2738` / `#3a364d` |
| `UTS_COLOR_ACCENT` | Info highlights and links | `#4068a0` / `#80b8e8` |
| `UTS_COLOR_SUCCESS` | Success messages | `#5a9b65` / `#abe9b3` |
| `UTS_COLOR_WARNING` | Warning messages | `#b89556` / `#f9e2af` |
| `UTS_COLOR_ERROR` | Error messages | `#b86080` / `#f28fad` |

### Light / Dark Variants

Each color can be split into separate light and dark terminal variants:

```bash
# Override just the dark-mode color
export UTS_COLOR_PRIMARY_DARK="#ff6600"

# Override just the light-mode color
export UTS_COLOR_PRIMARY_LIGHT="#0066ff"

# Override both (overrides the base color entirely)
export UTS_COLOR_PRIMARY="#336699"
```

Priority: `_LIGHT` / `_DARK` variants override the base color for their respective terminal mode.

### Disabling / Forcing Color

`uts` respects the [NO_COLOR](https://no-color.org) standard:

```bash
# Disable all color output
export NO_COLOR=1

# Force color output (e.g. when piping to a file)
export FORCE_COLOR=1
```

### Example: Custom Theme

```bash
# Rose Pine-inspired palette for dark terminals
export UTS_COLOR_PRIMARY="#eb6f92"
export UTS_COLOR_ACCENT="#31748f"
export UTS_COLOR_SUCCESS="#9ccfd8"
export UTS_COLOR_WARNING="#f6c177"
export UTS_COLOR_ERROR="#eb6f92"
export UTS_COLOR_BORDER="#26233a"
export UTS_COLOR_TEXT="#e0def4"
export UTS_COLOR_MUTED="#6e6a86"
export UTS_COLOR_SUBTLE="#908caa"

uts info myfile.mp4
```

> [!TIP]
> Add your color exports to your shell profile (`~/.bashrc`, `~/.zshrc`, etc.) to persist them across sessions.
