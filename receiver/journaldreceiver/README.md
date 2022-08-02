## `Journald Receiver`

| Status                   |           |
| ------------------------ |-----------|
| Stability                | [alpha]   |
| Supported pipeline types | logs      |
| Distributions            | [contrib] |

Parses Journald events from systemd journal.
Journald receiver is dependent on `journalctl` binary to be present and must be in the $PATH of the agent.

## Configuration

| Field                  | Default          | Description                                                                                                        |
| ---                    | ---              | ---                                                                                                                |
| `directory`            | /run/log/journal or /run/journal | A directory containing journal files to read entries from.     |
| `files`                |                  | A list of journal files to read entries from                  |
| `start_at`              | `end`              | At startup, where to start reading logs from the file. Options are beginning or end          |
| `units`        | `[ssh, kubelet, docker, containerd]` | A list of units to read entries from          |
| `prioriry`             | `info`           | Filter output by message priorities or priority ranges        |

### Example Configurations
```yaml
receivers:
  journald:
    directory: /run/log/journal
    units:
      - ssh
      - kubelet
      - docker
      - containerd
    priority: info
```

[alpha]: https://github.com/open-telemetry/opentelemetry-collector#alpha
[contrib]: https://github.com/open-telemetry/opentelemetry-collector-releases/tree/main/distributions/otelcol-contrib
