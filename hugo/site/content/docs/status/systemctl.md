---
weight: 10
title: "Systemctl"
---

# Status Systemctl

Display count of running systemd system and user services, including any failed services.

## Usage

```bash
allbctl status systemctl
```

## Output

```
Systemd Services:

  System Services:
    Running: 87

  User Services:
    Running: 12
```

With failed services:
```
Systemd Services:

  System Services:
    Running: 37 (1 failed)

  User Services:
    Running: 25
```

## Information Displayed

| Field | Source |
|-------|--------|
| System Services Running | `systemctl list-units --state=running` |
| System Services Failed | `systemctl list-units --state=failed` |
| User Services Running | `systemctl --user list-units --state=running` |

## Notes

- Only available on Linux systems with systemd
- On macOS or Windows, this section is omitted from `allbctl status`
- Failed services are shown in parentheses if count > 0
- Does not list individual service names — just counts

## Integration

Service counts are shown inline in `allbctl status`:

```
Services:  87 running (system), 12 running (user)
```

Run `allbctl status systemctl` to view the breakdown in isolation.
