# SchedNext

> A lightweight, filesystem-based job scheduler for cloud VMs, minimal OS environments, and edge devices.

SchedNext is a minimal job scheduler designed for:

- Cloud virtual machines
- Minimal / custom OS environments
- Edge & ARM devices
- Containers
- Air-gapped systems

It provides cron-based execution, persistent job state, atomic config updates, and a simple CLI for runtime interaction — **without requiring databases or heavy infrastructure.**

---

## ✨ Features

- ⏱ Cron-based scheduling
- 🔁 Persistent job state (`lastRunAt`, `lastExitCode`)
- 🔒 File locking & runtime protection
- 🧾 Atomic config writes
- 🖥 Simple CLI for job control
- 📦 Small static binaries
- 🧠 Designed for minimal OS / low-resource environments

---

## 🎯 Philosophy

SchedNext prioritizes:

- **Simplicity** over complexity
- **Predictable behavior** over magic
- **Filesystem persistence** over databases
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

SchedNext consists of two binaries:

| Binary | Purpose |
|--------|---------|
| `schednext-agent` | Scheduler engine |
| `schednext` | CLI controller |

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
      "id": "test-file",
      "binary": "test_file",
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

Where `/home` contains user directories:

```
/home/<user>/schednext.config
```

The agent:

- Scans user directories
- Evaluates cron schedules
- Applies runtime locks
- Executes binaries
- Persists job state

---

## 🖥 CLI Usage

Basic commands:

```bash
schednext start <user> <job-id>
schednext stop <user> <job-id>
```

Examples:

```bash
schednext start customer-ezeji test-file
schednext stop customer-ezeji test-file
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
GOOS=linux GOARCH=amd64 go build ./cmd/cli
GOOS=linux GOARCH=amd64 go build ./cmd/agent
GOOS=linux GOARCH=arm64 go build ./cmd/cli
GOOS=linux GOARCH=arm64 go build ./cmd/agent
```

---

## 🧪 Example Use Cases

- Cloud VM automation
- Edge device scheduling
- Minimal OS task execution
- Container cron replacement
- Embedded automation runtimes
- Air-gapped systems

---

## ❓ Why Not Use cron / systemd timers / Kubernetes?

SchedNext is **not meant to replace existing schedulers everywhere**.

It targets environments where traditional tools may be unavailable, heavy, or insufficient.

---

### cron

cron is excellent for traditional Linux systems.

**Limitations in some scenarios:**

- No persistent job metadata
- No built-in overlap protection
- Harder runtime/dynamic control
- Often absent in minimal OS / containers

**SchedNext advantages:**

- Persistent job state
- Runtime locking
- Atomic config updates
- CLI-driven control

---

### systemd timers

systemd timers are powerful but:

- systemd may not exist in minimal OS builds
- Higher operational complexity
- Less portable

**SchedNext advantages:**

- No init system dependency
- Portable static binary
- Minimal footprint

---

### Kubernetes CronJobs

K8s CronJobs are ideal inside clusters but:

- Require Kubernetes
- Heavy infrastructure dependency
- Unsuitable for small VMs / edge devices

**SchedNext advantages:**

- No cluster required
- Tiny runtime
- Simple deployment

---

### When SchedNext Makes Sense

- Minimal / custom OS
- Edge / ARM devices
- Small cloud VMs
- Containers without cron/systemd
- Embedded runtimes
- Air-gapped systems

---

## 🛣 Roadmap

Planned improvements:

- [ ] Job execution history
- [ ] Metrics / observability hooks
- [ ] Optional pluggable logging sinks
- [ ] REST / HTTP control API

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
