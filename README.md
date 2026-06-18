# go_n8n_notification_windows

A Windows desktop alert notifier driven by an **n8n** workflow. The Go application polls an n8n HTTP endpoint and fires a Windows toast notification whenever an alert becomes active.

---

## Architecture

```
n8n (alert_state Data Table)
        │
        ▼
Alert State Reader workflow  ←── GET /webhook/alert-state
        │
        ▼
go_n8n_notification_windows.exe   (polls every 30 s)
        │
        ▼
Windows Toast Notification
```

---

## n8n Workflows

Two workflows are created and published automatically via n8n MCP.

| Workflow | ID | URL |
|---|---|---|
| Alert State Reader | `kGoaFkn2zKPJ8JUH` | https://n8n.hugojava.dev/workflow/kGoaFkn2zKPJ8JUH |
| Alert Admin | `IEgGHHrMHZDG7wCU` | https://n8n.hugojava.dev/workflow/IEgGHHrMHZDG7wCU |

State is persisted in an **n8n Data Table** (id: `DaFORuYuKLNjfOy3`, name: `alert_state`).

### Persistence tradeoffs

| Option | Tradeoff |
|---|---|
| **n8n Data Table** (chosen) | Zero external dependencies, built-in to n8n, survives restarts. Ideal for single-key state. Limited to n8n's own storage backend. |
| PostgreSQL | Full SQL power, external transactions. Requires running Postgres and credentials. Best if other services also read the state. |
| SQLite | No server needed, file-based. Requires file-system access from n8n and careful path config in containers. |

---

## Alert State API

### GET `/webhook/alert-state`

Returns the current alert state.

```bash
curl https://n8n.hugojava.dev/webhook/alert-state
```

**Response 200**

```json
{
  "enabled": true,
  "priority": "HIGH",
  "title": "Production Incident",
  "message": "Database connection failures detected",
  "sound": "critical"
}
```

### POST `/webhook/alert-admin`

Manages alert state. Requires a JSON body with an `action` field.

#### Enable alert

```bash
curl -X POST https://n8n.hugojava.dev/webhook/alert-admin \
  -H "Content-Type: application/json" \
  -d '{
    "action": "enable",
    "priority": "HIGH",
    "title": "Production Incident",
    "message": "Database connection failures detected",
    "sound": "critical"
  }'
```

#### Disable alert

```bash
curl -X POST https://n8n.hugojava.dev/webhook/alert-admin \
  -H "Content-Type: application/json" \
  -d '{"action": "disable"}'
```

#### Update fields only (keeps enabled state)

```bash
curl -X POST https://n8n.hugojava.dev/webhook/alert-admin \
  -H "Content-Type: application/json" \
  -d '{
    "action": "update",
    "priority": "NORMAL",
    "title": "Maintenance Window",
    "message": "Scheduled maintenance in 15 minutes",
    "sound": "default"
  }'
```

**Valid actions:** `enable` | `disable` | `update`

**Valid sounds:** `default` | `critical` | `alarm` | `call` | `silent`

**Valid priorities:** `LOW` | `NORMAL` | `HIGH` | `CRITICAL` (free-form string displayed as-is)

**Response 200**

```json
{"success": true, "message": "Alert state updated"}
```

**Response 400**

```json
{"error": "Missing required field: action", "valid_actions": ["enable", "disable", "update"]}
```

---

## Go Application

### Requirements

- Windows 10 or later (for toast notifications)
- Network access to the n8n instance

### Install / Download

Download the latest release ZIP from the [Releases](../../releases) page and extract `go_n8n_notification_windows.exe`.

### Usage

```
go_n8n_notification_windows.exe [flags]

Flags:
  -endpoint string    n8n alert-state endpoint URL
                      (default from N8N_ALERT_ENDPOINT env, or https://n8n.hugojava.dev/webhook/alert-state)
  -interval duration  Polling interval (default 30s)
                      (default from POLL_INTERVAL env)
  -appid string       Windows App User Model ID for toast notifications
                      (default from APP_ID env, or "go_n8n_notification_windows")
```

### Environment variables

| Variable | Default | Description |
|---|---|---|
| `N8N_ALERT_ENDPOINT` | `https://n8n.hugojava.dev/webhook/alert-state` | Endpoint URL |
| `POLL_INTERVAL` | `30s` | Go duration string (e.g. `1m`, `10s`) |
| `APP_ID` | `go_n8n_notification_windows` | Windows toast App ID |

### Notification behaviour

The app fires a toast notification **once** when an alert transitions from disabled to enabled. It will not spam repeat notifications while the alert remains active. After the alert is disabled, the next enable transition triggers a new notification.

---

## Build from source

```bash
# Run tests
go test ./...

# Cross-compile for Windows
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 \
  go build -ldflags="-s -w" -o dist/go_n8n_notification_windows.exe ./cmd/notifier
```

---

## CI/CD

GitHub Actions workflow (`.github/workflows/ci.yml`):

1. **Lint & Test** — `gofmt` verification, `go vet`, unit tests with race detector and coverage
2. **Build Windows amd64** — cross-compile and package as ZIP artifact
3. **Publish GitHub Release** — triggered on `v*` tags; attaches the ZIP to a GitHub Release automatically

To release:

```bash
git tag v1.0.0
git push origin v1.0.0
```
