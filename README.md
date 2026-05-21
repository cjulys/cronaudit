# cronaudit

Parse and validate cron expressions across multiple systems, outputting a unified schedule report.

## Installation

```bash
go install github.com/youruser/cronaudit@latest
```

Or build from source:

```bash
git clone https://github.com/youruser/cronaudit.git && cd cronaudit && go build ./...
```

## Usage

```bash
# Audit a single cron expression
cronaudit parse "0 */6 * * *"

# Audit all crontabs on the local system
cronaudit audit --system

# Validate expressions from a file and output a unified report
cronaudit audit --file crontabs.txt --output report.json
```

Example output:

```
Expression       Next Run              Description
───────────────  ────────────────────  ─────────────────────────
0 */6 * * *      2024-01-15 18:00:00   Every 6 hours
30 2 * * 1       2024-01-22 02:30:00   At 02:30 on Monday
0 0 1 * *        2024-02-01 00:00:00   At midnight on day 1
```

## Flags

| Flag | Description |
|------|-------------|
| `--system` | Scan system-level crontabs (`/etc/cron*`) |
| `--file` | Read expressions from a file |
| `--output` | Write report to file (`json`, `csv`, `text`) |
| `--tz` | Evaluate schedules in a specific timezone |

## Supported Systems

- Standard POSIX cron
- Vixie cron (Linux)
- macOS `launchd` plists
- AWS EventBridge scheduler expressions

## Contributing

Pull requests are welcome. Please open an issue first to discuss any significant changes.

## License

MIT © youruser