# WinOrchestra

Control Windows windows from the command line.  

Based on [WindowManager](https://github.com/FanaticExplorer/WindowManager) (Python) and rewritten in Go for speed, compact binaries, and native Win32 access with zero runtime dependencies.

## Install

Download the latest `WinOrchestra.exe` from [Releases](https://github.com/FanaticExplorer/WinOrchestra/releases) and place it anywhere — no installer, no runtime.

Or build from source:

```
go install github.com/FanaticExplorer/WinOrchestra@latest
```

## Usage

```
winorchestra                       Check that it works
winorchestra --help                List available commands
winorchestra [command] --help      Show help for a specific command
```

### Commands

| Command | Description |
|---|---|
| `list` | List windows as JSON |
| `focus` | Restore and bring a window to the foreground |
| `minimize` | Minimize a window |
| `close` | Gracefully close a window (same as clicking ✕) |

### Filters

All commands except `list` require at least one filter:

| Flag | Description |
|---|---|
| `-t, --title` | Partial window title (case‑insensitive) |
| `-p, --process` | Partial `.exe` name (case‑insensitive) |
| `--pid` | Exact process ID |

Multiple filters are ANDed — a window must match all of them.

### Examples

```bash
# List all windows
winorchestra list

# List only Firefox windows
winorchestra list -p "firefox"

# Compact JSON output
winorchestra list --raw

# Focus Chrome
winorchestra focus -p "chrome"

# Minimize a window by title
winorchestra minimize -t "Calculator"

# Close Discord (it will minimize to tray instead of exiting)
winorchestra close -p "discord"
```

### jq recipes

Since `list` outputs JSON, you can pipe it into [jq](https://jqlang.github.io/jq/):

```bash
# Show only minimized windows
winorchestra list | jq '.[] | select(.minimized)'

# Show the currently focused window
winorchestra list | jq '.[] | select(.focused)'

# Check if Notepad is running
winorchestra list -p "notepad" | jq 'length > 0'
```

### JSON output

```json
{
  "pid": 21348,
  "exe": "Discord.exe",
  "class": "Chrome_WidgetWin_1",
  "title": "Friends - Discord",
  "minimized": false,
  "focused": false
}
```

## License

MIT © FanaticExplorer
