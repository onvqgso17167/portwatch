# portwatch

A lightweight CLI daemon that monitors open ports and alerts on unexpected changes.

---

## Installation

```bash
go install github.com/youruser/portwatch@latest
```

Or build from source:

```bash
git clone https://github.com/youruser/portwatch.git && cd portwatch && go build -o portwatch .
```

---

## Usage

Start the daemon with a default scan interval of 30 seconds:

```bash
portwatch start
```

Specify a custom interval and log file:

```bash
portwatch start --interval 60 --log /var/log/portwatch.log
```

Run a one-time snapshot of currently open ports:

```bash
portwatch scan
```

**Example output:**

```
[INFO]  Watching ports... (interval: 30s)
[ALERT] New port detected: 8080 (tcp)
[ALERT] Port closed: 3306 (tcp)
```

Define a whitelist of expected ports in `portwatch.yaml` to suppress known services:

```yaml
whitelist:
  - 22
  - 80
  - 443
```

---

## Configuration

| Flag | Default | Description |
|------|---------|-------------|
| `--interval` | `30` | Scan interval in seconds |
| `--log` | stdout | Path to log file |
| `--config` | `portwatch.yaml` | Path to config file |

---

## License

MIT © 2024 youruser