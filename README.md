# WinOrchestra

Control Windows windows from the command line.
(*Insert a pun about this because I didn't come up with one*)

Based on my old project [WindowManager](https://github.com/FanaticExplorer/WindowManager) (Python) (which I made as a beginner) and rewritten in Go for speed, compact binaries (and just to prove myself that I can do it).

## Install

Download the latest executable from [Releases](https://github.com/FanaticExplorer/WinOrchestra/releases) and place it anywhere where you want. 
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

| Flag | Description |
|---|---|
| `-t, --title` | Partial window title (case‑insensitive) |
| `-p, --process` | Partial `.exe` name (case‑insensitive) |
| `--class` | Partial window class (case‑insensitive) |
| `--pid` | Exact process ID |

Multiple filters are ANDed — a window must match all of them. `focus`, `minimize`, and `close` require at least one filter; `list` works with or without.

Use `--all` to apply an action to every matching window instead of just the first.

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

# Close all Discord windows
winorchestra close -p "discord" --all

# List only windows of a specific class
winorchestra list --class "MozillaWindowClass"
```

### jq recipes

Since `list` outputs JSON, you can pipe it into something like [jq](https://jqlang.github.io/jq/):

```bash
# Show only minimized windows
winorchestra list | jq '.[] | select(.minimized)'

# Show the currently focused window
winorchestra list | jq '.[] | select(.focused)'

# Check if Notepad is running
winorchestra list -p "notepad" | jq 'length > 0'
```

### JSON output example

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

## Support

I'm a student building tools like this in my free time — if it's useful to you, a [donation](https://donate.fanaticexplorer.dev/) helps me keep going.

## License

MIT — see [LICENSE](LICENSE).
