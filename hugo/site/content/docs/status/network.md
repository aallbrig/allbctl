---
weight: 7
title: "Network"
---

# Status Network

Display network interface information including IP addresses, gateway, DNS, VPN status, and internet connectivity.

## Usage

```bash
allbctl status network
```

## Output

### Wired Connection
```
Network:
  Primary Interface: eth0 (192.168.1.100/24)
    Gateway: 192.168.1.1

  DNS:
    System: 192.168.1.1

  Connectivity:
    Public IP: 203.0.113.10
    Internet: ✓ Connected

  Other Interfaces:
    lo: 127.0.0.1/8
    docker0: 172.17.0.1/16
```

### WiFi Connection
```
Network:
  Primary Interface: wlp0s20f3 (192.168.50.90/24)
    WiFi: MY_SSID @ 5.24 GHz (802.11ac/ax)
    Signal: -56 dBm (Good)
    Gateway: 192.168.50.1

  DNS:
    System: 192.168.50.1

  Connectivity:
    Public IP: 68.54.122.205
    Internet: ✓ Connected

  Other Interfaces:
    lo: 127.0.0.1/8
    enp0s31f6: DOWN
    docker0: 172.17.0.1/16
```

### No Internet
```
Network:
  Primary Interface: eth0 (192.168.1.100/24)
    Gateway: 192.168.1.1

  Connectivity:
    Internet: ✗ No connectivity
```

## Information Displayed

### Primary Interface
- Interface name and IP address with CIDR notation
- **WiFi**: SSID, frequency band, protocol (802.11ac/ax), signal strength in dBm with quality label (Excellent/Good/Fair/Poor)
- Default gateway

### DNS
- System DNS server(s) from `/etc/resolv.conf`

### Connectivity
- Public IP address (detected via external lookup)
- Internet connectivity status (✓ Connected / ✗ No connectivity)

### Other Interfaces
- All non-primary interfaces with their IPs or DOWN status
- Includes loopback, Ethernet, Docker bridges, VPN tunnels, VirtualBox interfaces

## VPN Detection

allbctl detects common VPN interfaces by name pattern:
- `tun*`, `tap*` — OpenVPN / WireGuard / generic
- `wg*` — WireGuard
- `vpn*`, `utun*` — macOS/generic VPN
- `proton*` — ProtonVPN

## Integration

Network information is shown in the `allbctl status` output under the **Network:** section.
