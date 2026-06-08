# SchedNext

> A lightweight agent runtime for cloud, edge, and embedded Linux environments.

SchedNext is a lightweight runtime for deploying, scheduling, and observing agents.

It is designed for environments where traditional scheduling platforms may be too heavy or unavailable, including edge devices, cloud VMs, custom Linux distributions, containers, and air-gapped systems.

SchedNext combines:

- Agent execution
- Persistent runtime state
- Filesystem-based observability
- Runtime control via CLI

without requiring databases, web services, or large infrastructure dependencies.

---

## ✨ Features

- ⏱ Cron-based scheduling
- 🔁 Persistent agent state
- 🔒 Runtime locking
- 🧾 Atomic config updates
- 🖥 CLI runtime control
- 📊 Filesystem observability
- 🧠 Dynamic runtime state projection
- 📦 Small static binaries
- ⚡ Low memory footprint

---

## 📊 Observability via StateLens

SchedNext exposes live runtime state through a virtual filesystem.

Inspired by Linux `/proc` and `/sys`, runtime state can be inspected using standard Unix tools.

No web server
No metrics agent
No SDK

Just files

Examples:

```bash
cat /statelens/cpu/summary

cat /statelens/mem/summary

cat /statelens/schednext/all

cat /statelens/schednext/sensor
```

---

## StateLens Layout

```
/statelens
├── cpu
│   └── summary
├── mem
│   └── summary
├── net
│   ├── interfaces
│   └── routes
└── schednext
    ├── all
    ├── heartbeat
    └── sensor
```

## Examples

```
$ cat /statelens/cpu/summary
 18:42:31 up 3 days,  2 users,  load average: 0.15, 0.09, 0.05

$ cat /statelens/mem/summary
              total        used        free      shared  buff/cache   available
Mem:           3.8G        512M        2.4G         12M        896M        3.1G

$ cat /statelens/schednext/all
 heartbeat | idle | enabled=true | exit=0
 sensor | running | enabled=true | exit=0
```

---

## 🎯 Philosophy

SchedNext prioritizes:

- **Simplicity** over complexity
- **Predictable behavior** over magic
- **Filesystem primitives** over infrastructure dependencies
- **Low resource usage** over feature bloat

No database  
No heavy dependencies  
No mandatory log files  

---

## 🧾 Logging Strategy

SchedNext intentionally **does not write logs to disk**.

### Why?

- Disk logging increases I/O
- Log files grow unbounded
- Log buffering consumes RAM
- Minimal OS / edge devices suffer from log overhead

### Instead

- Logs go to **stdout/stderr**
- Users may redirect or pipe logs externally
- Future versions may support pluggable log sinks

Job execution results are persisted via:

- `lastRunAt`
- `lastExitCode`

inside the configuration file.

---

## 🏗 Architecture

SchedNext consists of three runtime components:

| Component | Purpose |
|--------|---------|
| `schednext-agent` | Executes and schedules agents |
| `schednext` | Runtime control |
| `statelens` | Filesystem observability |

The CLI communicates with the agent via **IPC (Unix socket)**.

---

## ⚙ Configuration

SchedNext uses a JSON configuration file:

`schednext.config`

Example:

```json
{
  "jobs": [
    {
      "id": "sensor",
      "binary": "sensor-agent",
      "cron": "*/2 * * * *",
      "enabled": false,
      "lastRunAt": "2025-12-20T12:45:00Z",
      "lastExitCode": 2025,
      "lockUntil": "2025-12-20T12:45:00Z",
      "maxRuntimeSeconds": 2
    }
  ]
}
```

---

## 📘 Job Fields

| Field | Type | Description |
|------|------|-------------|
| `id` | string | Unique job identifier |
| `binary` | string | Executable name or relative path |
| `cron` | string | Cron expression (5-field format) |
| `enabled` | bool | Enables/disables job |
| `lastRunAt` | datetime | Last execution timestamp |
| `lastExitCode` | int | Exit code of last run |
| `lockUntil` | datetime | Prevents overlapping runs |
| `maxRuntimeSeconds` | int | Maximum allowed runtime |

---

## ⏱ Cron Format

SchedNext supports standard 5-field cron expressions:

```
MINUTE HOUR DOM MONTH DOW
```

Examples:

```
*/2 * * * *    → Every 2 minutes
0 * * * *      → Every hour
0 0 * * *      → Daily at midnight
```

---

## 🚀 Running the Agent

Example:

```bash
./schednext-agent
```

Where `/opt` contains schednext runtime directory:

```
/opt/schednext-runtime/schednext.config
```

The agent:

- Locates runtime directory
- Evaluates cron schedules
- Applies runtime locks
- Executes binaries
- Streams runtime state
- Persists job state

---

## 🖥 CLI Usage

Basic commands:

```bash
schednext start <job-id>
schednext stop <job-id>
```

Examples:

```bash
schednext start sensor
schednext stop sensor
```

---

## 🔧 Build

Build binaries from source:

```bash
go build -o schednext-agent ./cmd/agent
go build -o schednext ./cmd/cli
```

Cross-compile examples:

```bash
GOOS=linux GOARCH=amd64 go build -o schednext ./cmd/cli
GOOS=linux GOARCH=amd64 go build -o schednext-agent ./cmd/agent
GOOS=linux GOARCH=arm64 go build -o schednext ./cmd/cli
GOOS=linux GOARCH=arm64 go build -o schednext-agent ./cmd/agent
```

---

## 🧪 Use Cases

SchedNext is designed for environments where lightweight deployment, runtime observability, and minimal operational dependencies are important.

Examples:

- Cloud VM automation
- Edge devices
- Robotics runtimes
- IoT gateways
- Custom Linux distributions
- Air-gapped systems
- Containers without cron or systemd

---

## 🛣 Roadmap

Planned improvements:

- [ ] Event-driven job triggers
- [ ] Webhook-based execution
- [ ] Remote node monitoring
- [ ] Multi-node fleet management
- [ ] Lightweight web dashboard
- [ ] SchedNext OS
- [ ] ARM-first edge deployment tooling

---

## 🤝 Contributing

Contributions are welcome.

Guidelines:

- Keep changes minimal & focused
- Preserve low-resource design goals
- Avoid heavy dependencies

---

## ⚠ Project Status

Early-stage / evolving design.

Config structure and APIs may change.

---

## 📄 License

MIT License
