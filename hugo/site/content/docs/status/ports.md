---
weight: 8
title: "Ports"
---

# Status Ports

Display count of listening TCP/UDP ports and details about what is listening.

## Usage

```bash
allbctl status ports
```

## Output

```
Listening Ports:

  TCP Ports: 18
  UDP Ports: 6
  Total:     24

  Details:
    tcp:22
    tcp:80
    tcp:443
    tcp:3000
    tcp:5432
    tcp:6379
    udp:53
    udp:5353
    ...
```

## Information Displayed

- **TCP Ports** — count of listening TCP ports
- **UDP Ports** — count of listening UDP ports
- **Total** — combined count
- **Details** — list of all `protocol:port` combinations currently bound

## Detection

Ports are detected using `ss -tulnp` (Linux) or `netstat -an` (macOS/Windows). Only **listening** ports are shown — not established connections.

## Integration

Port count is shown inline in `allbctl status`:

```
Ports:     24 listening (TCP: 18, UDP: 6)
```

Run `allbctl status ports` to see the full port list with details.
