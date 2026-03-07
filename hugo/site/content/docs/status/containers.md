---
weight: 5
title: "Containers"
---

# Status Containers

Display information about container runtimes (Docker, Podman) and virtualization status.

## Usage

```bash
allbctl status containers
```

## Output

Shows running containers, local images, and virtualization status:

```
Containers/Virtualization:

  Docker:
    Running Containers: 3
    Images: 5
    Image List:
      - alpine/curl:latest
      - node:18-alpine
      - localstack/localstack:latest
      - public.ecr.aws/lambda/nodejs:20-x86_64
      - postgres:15-alpine

  Virtualization: None detected (bare metal)
```

If Docker is not running:
```
Containers/Virtualization:

  Docker: not running

  Virtualization: None detected (bare metal)
```

With virtual machines detected:
```
Containers/Virtualization:

  Docker:
    Running Containers: 0
    Images: 2

  Virtualization:
    VirtualBox: 2 VMs detected
    Vagrant: 2 environments
```

## Detected Runtimes

### Containers
- **Docker** — detects running containers, local images, and Docker Compose projects
- **Podman** — detects running containers and images (if installed)

### Virtualization
- **VirtualBox** — detects VMs via `VBoxManage list vms`
- **Vagrant** — detects Vagrant environments via `vagrant global-status`
- **Native** — reports "bare metal" if no hypervisor is detected

## Integration

The containers summary is included in the `allbctl status` output under the **Containers/Virtualization:** section.
