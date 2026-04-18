# portwatch

Lightweight CLI to monitor and alert on open port changes across hosts.

## Installation

```bash
go install github.com/yourusername/portwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/portwatch.git && cd portwatch && go build -o portwatch .
```

## Usage

Watch a host for port changes and print alerts to stdout:

```bash
portwatch --host 192.168.1.1 --interval 30s
```

Monitor multiple hosts using a config file:

```bash
portwatch --config hosts.yaml
```

Example `hosts.yaml`:

```yaml
hosts:
  - address: 192.168.1.1
    ports: [22, 80, 443]
  - address: 10.0.0.5
    ports: [3306, 5432]
interval: 60s
```

When a port change is detected, portwatch outputs an alert:

```
[ALERT] 192.168.1.1 — port 8080 is now OPEN (detected at 2024-01-15 10:32:01)
[ALERT] 10.0.0.5    — port 3306 is now CLOSED (detected at 2024-01-15 10:32:04)
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--host` | — | Single host to monitor |
| `--config` | — | Path to YAML config file |
| `--interval` | `60s` | Polling interval |
| `--notify` | `stdout` | Alert output (`stdout`, `webhook`) |
| `--webhook` | — | Webhook URL for notifications |

## License

MIT © 2024 portwatch contributors