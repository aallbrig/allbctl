package cmd

import (
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
	"strings"

	"github.com/shirou/gopsutil/v4/net"
	"github.com/spf13/cobra"
)

var NetworkCmd = &cobra.Command{
	Use:   "network",
	Short: "Display network interface information",
	Long: `Display network interface information including IP addresses, router, connection type, VPN status, DNS, and connectivity.

This is the same output shown in the 'Network:' section of 'allbctl status'.`,
	Run: func(cmd *cobra.Command, args []string) {
		PrintNetworkInfo()
	},
}

// NetworkDetails holds comprehensive network information
type NetworkDetails struct {
	Interfaces     []InterfaceInfo
	VPNActive      bool
	VPNInterface   *InterfaceInfo
	PrimaryIface   *InterfaceInfo
	DefaultGateway string
	DNSServers     []string
	VPNDNSServers  []string
	PublicIP       string
	InternetOK     bool
	WiFiDetails    *WiFiInfo
}

// InterfaceInfo holds interface details
type InterfaceInfo struct {
	Name    string
	IP      string
	Status  string
	IsVPN   bool
	Gateway string
}

// WiFiInfo holds WiFi-specific details
type WiFiInfo struct {
	SSID      string
	Frequency string
	Standard  string
	Signal    string
	Quality   string
	Speed     string
}

// PrintNetworkInfo outputs comprehensive network information
func PrintNetworkInfo() {
	details := gatherNetworkDetails()

	// Primary Interface
	if details.PrimaryIface != nil {
		fmt.Printf("Network:\n")
		fmt.Printf("  Primary Interface: %s (%s)\n", details.PrimaryIface.Name, details.PrimaryIface.IP)

		// WiFi details if available
		if details.WiFiDetails != nil {
			if details.WiFiDetails.SSID != "" {
				fmt.Printf("    WiFi: %s", details.WiFiDetails.SSID)
				if details.WiFiDetails.Frequency != "" {
					fmt.Printf(" @ %s", details.WiFiDetails.Frequency)
				}
				if details.WiFiDetails.Standard != "" {
					fmt.Printf(" (%s)", details.WiFiDetails.Standard)
				}
				fmt.Println()
			}

			if details.WiFiDetails.Speed != "" || details.WiFiDetails.Signal != "" {
				fmt.Print("    ")
				if details.WiFiDetails.Speed != "" {
					fmt.Printf("Speed: %s", details.WiFiDetails.Speed)
				}
				if details.WiFiDetails.Signal != "" {
					if details.WiFiDetails.Speed != "" {
						fmt.Print(" | ")
					}
					fmt.Printf("Signal: %s", details.WiFiDetails.Signal)
					if details.WiFiDetails.Quality != "" {
						fmt.Printf(" (%s)", details.WiFiDetails.Quality)
					}
				}
				fmt.Println()
			}
		}

		if details.PrimaryIface.Gateway != "" {
			fmt.Printf("    Gateway: %s\n", details.PrimaryIface.Gateway)
		}
		fmt.Println()
	}

	// VPN Status
	if details.VPNActive && details.VPNInterface != nil {
		fmt.Printf("  VPN Active: ✓ %s (%s)\n", details.VPNInterface.Name, details.VPNInterface.IP)
		if details.VPNInterface.Gateway != "" {
			fmt.Printf("    Gateway: %s\n", details.VPNInterface.Gateway)
		}
		fmt.Printf("    Status: Traffic routed via VPN\n")
		fmt.Println()
	}

	// DNS
	if len(details.DNSServers) > 0 || len(details.VPNDNSServers) > 0 {
		fmt.Printf("  DNS:\n")
		if len(details.DNSServers) > 0 {
			fmt.Printf("    System: %s\n", strings.Join(details.DNSServers, ", "))
		}
		if len(details.VPNDNSServers) > 0 {
			fmt.Printf("    VPN: %s\n", strings.Join(details.VPNDNSServers, ", "))
		}
		fmt.Println()
	}

	// Connectivity
	fmt.Printf("  Connectivity:\n")
	if details.PublicIP != "" {
		fmt.Printf("    Public IP: %s\n", details.PublicIP)
	}
	if details.InternetOK {
		fmt.Printf("    Internet: ✓ Connected\n")
	} else {
		fmt.Printf("    Internet: ✗ No connection\n")
	}

	// Other interfaces
	otherIfaces := []InterfaceInfo{}
	for _, iface := range details.Interfaces {
		if details.PrimaryIface != nil && iface.Name == details.PrimaryIface.Name {
			continue
		}
		if details.VPNInterface != nil && iface.Name == details.VPNInterface.Name {
			continue
		}
		otherIfaces = append(otherIfaces, iface)
	}

	if len(otherIfaces) > 0 {
		fmt.Println()
		fmt.Printf("  Other Interfaces:\n")
		for _, iface := range otherIfaces {
			if iface.Status == "DOWN" {
				fmt.Printf("    %s: %s\n", iface.Name, iface.Status)
			} else {
				fmt.Printf("    %s: %s\n", iface.Name, iface.IP)
			}
		}
	}
}

// gatherNetworkDetails collects comprehensive network information
func gatherNetworkDetails() *NetworkDetails {
	details := &NetworkDetails{
		Interfaces: []InterfaceInfo{},
	}

	// Get interfaces
	netIfaces, err := net.Interfaces()
	if err == nil {
		for _, iface := range netIfaces {
			info := InterfaceInfo{
				Name:   iface.Name,
				Status: "UP",
			}

			// Check if VPN interface
			info.IsVPN = isVPNInterface(iface.Name)

			// Get IP address
			if len(iface.Addrs) > 0 {
				for _, addr := range iface.Addrs {
					if !strings.Contains(addr.Addr, ":") { // Skip IPv6
						info.IP = addr.Addr
						break
					}
				}
			}

			if info.IP == "" {
				info.Status = "DOWN"
			}

			details.Interfaces = append(details.Interfaces, info)
		}
	}

	// Detect primary and VPN interfaces
	detectPrimaryInterface(details)
	detectVPNInterface(details)

	// Get gateway information
	details.DefaultGateway = getRouterIP()
	if details.PrimaryIface != nil {
		details.PrimaryIface.Gateway = details.DefaultGateway
	}
	if details.VPNInterface != nil {
		details.VPNInterface.Gateway = getVPNGateway(details.VPNInterface.Name)
	}

	// Get WiFi details for primary interface
	if details.PrimaryIface != nil {
		details.WiFiDetails = getWiFiDetails(details.PrimaryIface.Name)
	}

	// Get DNS servers
	details.DNSServers, details.VPNDNSServers = getDNSServers()

	// Check internet connectivity
	details.InternetOK = checkInternetConnectivity()

	// Get public IP
	if details.InternetOK {
		details.PublicIP = getPublicIP()
	}

	return details
}

// isVPNInterface checks if an interface is a VPN tunnel
func isVPNInterface(name string) bool {
	vpnPrefixes := []string{"tun", "tap", "wg", "ppp", "vpn"}
	for _, prefix := range vpnPrefixes {
		if strings.HasPrefix(name, prefix) {
			return true
		}
	}
	return false
}

// detectPrimaryInterface finds the primary network interface
func detectPrimaryInterface(details *NetworkDetails) {
	if runtime.GOOS != "linux" {
		// For non-Linux, just pick first non-loopback, non-VPN interface
		for i := range details.Interfaces {
			iface := &details.Interfaces[i]
			if iface.Name != "lo" && !iface.IsVPN && iface.Status == "UP" {
				details.PrimaryIface = iface
				return
			}
		}
		return
	}

	// On Linux, use ip route to find default interface
	cmd := exec.Command("sh", "-c", "ip route | grep default | awk '{print $5}' | head -n1")
	out, err := cmd.Output()
	if err == nil && len(out) > 0 {
		primaryName := strings.TrimSpace(string(out))
		for i := range details.Interfaces {
			if details.Interfaces[i].Name == primaryName {
				details.PrimaryIface = &details.Interfaces[i]
				return
			}
		}
	}
}

// detectVPNInterface finds active VPN interfaces
func detectVPNInterface(details *NetworkDetails) {
	for i := range details.Interfaces {
		iface := &details.Interfaces[i]
		if iface.IsVPN && iface.Status == "UP" {
			details.VPNInterface = iface
			details.VPNActive = true
			return
		}
	}
}

// getVPNGateway gets the gateway for a VPN interface
func getVPNGateway(ifaceName string) string {
	if runtime.GOOS != "linux" {
		return ""
	}

	cmd := exec.Command("sh", "-c", fmt.Sprintf("ip route | grep %s | grep -v default | awk '{print $1}' | head -n1", ifaceName))
	out, err := cmd.Output()
	if err == nil && len(out) > 0 {
		return strings.TrimSpace(string(out))
	}

	return ""
}

// getWiFiDetails gets WiFi-specific information
func getWiFiDetails(ifaceName string) *WiFiInfo {
	if runtime.GOOS != "linux" {
		return nil
	}

	info := &WiFiInfo{}

	// Get WiFi info using iwconfig
	cmd := exec.Command("iwconfig", ifaceName)
	out, err := cmd.Output()
	if err != nil {
		return nil
	}

	output := string(out)

	// Parse SSID
	if strings.Contains(output, "ESSID:") {
		re := regexp.MustCompile(`ESSID:"([^"]+)"`)
		if matches := re.FindStringSubmatch(output); len(matches) > 1 {
			info.SSID = matches[1]
		}
	}

	// Parse frequency
	if strings.Contains(output, "Frequency:") {
		re := regexp.MustCompile(`Frequency:([0-9.]+)\s*GHz`)
		if matches := re.FindStringSubmatch(output); len(matches) > 1 {
			info.Frequency = matches[1] + " GHz"
		}
	}

	// Parse signal strength
	if strings.Contains(output, "Signal level=") {
		re := regexp.MustCompile(`Signal level=(-?[0-9]+)\s*dBm`)
		if matches := re.FindStringSubmatch(output); len(matches) > 1 {
			signalDBm := matches[1]
			info.Signal = signalDBm + " dBm"

			// Determine quality
			signal := 0
			fmt.Sscanf(signalDBm, "%d", &signal)
			if signal >= -50 {
				info.Quality = "Excellent"
			} else if signal >= -60 {
				info.Quality = "Good"
			} else if signal >= -70 {
				info.Quality = "Fair"
			} else {
				info.Quality = "Poor"
			}
		}
	}

	// Get link speed using iw
	cmd = exec.Command("iw", "dev", ifaceName, "link")
	out, err = cmd.Output()
	if err == nil {
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			if strings.Contains(line, "tx bitrate:") {
				re := regexp.MustCompile(`tx bitrate:\s*([0-9.]+)\s*MBit/s`)
				if matches := re.FindStringSubmatch(line); len(matches) > 1 {
					info.Speed = matches[1] + " Mbps"
				}
			}
		}
	}

	// Determine WiFi standard from frequency
	if info.Frequency != "" {
		freq := 0.0
		fmt.Sscanf(info.Frequency, "%f", &freq)
		if freq >= 5.0 {
			info.Standard = "802.11ac/ax"
		} else if freq >= 2.4 {
			info.Standard = "802.11n"
		}
	}

	if info.SSID == "" {
		return nil
	}

	return info
}

// getDNSServers gets system and VPN DNS servers
func getDNSServers() ([]string, []string) {
	systemDNS := []string{}
	vpnDNS := []string{}

	if runtime.GOOS == "linux" {
		// Try resolvectl first (systemd-resolved)
		cmd := exec.Command("resolvectl", "status")
		out, err := cmd.Output()
		if err == nil {
			lines := strings.Split(string(out), "\n")
			inVPNSection := false

			for _, line := range lines {
				line = strings.TrimSpace(line)

				// Check if entering VPN interface section
				if strings.Contains(line, "tun") || strings.Contains(line, "wg") || strings.Contains(line, "vpn") {
					inVPNSection = true
				}

				// Check if leaving interface section
				if strings.HasPrefix(line, "Link ") && !strings.Contains(line, "tun") && !strings.Contains(line, "wg") {
					inVPNSection = false
				}

				// Parse DNS servers
				if strings.HasPrefix(line, "DNS Servers:") || strings.HasPrefix(line, "Current DNS Server:") {
					parts := strings.Split(line, ":")
					if len(parts) >= 2 {
						dns := strings.TrimSpace(parts[1])
						if inVPNSection {
							vpnDNS = append(vpnDNS, dns)
						} else {
							systemDNS = append(systemDNS, dns)
						}
					}
				}
			}
		} else {
			// Fallback to /etc/resolv.conf
			cmd = exec.Command("cat", "/etc/resolv.conf")
			out, err = cmd.Output()
			if err == nil {
				lines := strings.Split(string(out), "\n")
				for _, line := range lines {
					if strings.HasPrefix(line, "nameserver") {
						parts := strings.Fields(line)
						if len(parts) >= 2 {
							systemDNS = append(systemDNS, parts[1])
						}
					}
				}
			}
		}
	} else if runtime.GOOS == "darwin" {
		cmd := exec.Command("scutil", "--dns")
		out, err := cmd.Output()
		if err == nil {
			lines := strings.Split(string(out), "\n")
			for _, line := range lines {
				if strings.Contains(line, "nameserver[") {
					parts := strings.Split(line, ":")
					if len(parts) >= 2 {
						dns := strings.TrimSpace(parts[1])
						systemDNS = append(systemDNS, dns)
					}
				}
			}
		}
	} else if runtime.GOOS == "windows" {
		cmd := exec.Command("ipconfig", "/all")
		out, err := cmd.Output()
		if err == nil {
			lines := strings.Split(string(out), "\n")
			for _, line := range lines {
				if strings.Contains(line, "DNS Servers") {
					parts := strings.Split(line, ":")
					if len(parts) >= 2 {
						dns := strings.TrimSpace(parts[1])
						systemDNS = append(systemDNS, dns)
					}
				}
			}
		}
	}

	// Deduplicate
	systemDNS = uniqueStrings(systemDNS)
	vpnDNS = uniqueStrings(vpnDNS)

	return systemDNS, vpnDNS
}

// checkInternetConnectivity checks if internet is accessible
func checkInternetConnectivity() bool {
	cmd := exec.Command("ping", "-c", "1", "-W", "2", "8.8.8.8")
	err := cmd.Run()
	return err == nil
}

// getPublicIP gets the public IP address
func getPublicIP() string {
	cmd := exec.Command("curl", "-s", "--max-time", "3", "ifconfig.me")
	out, err := cmd.Output()
	if err == nil && len(out) > 0 {
		ip := strings.TrimSpace(string(out))
		// Validate it looks like an IP
		if strings.Count(ip, ".") == 3 {
			return ip
		}
	}
	return ""
}

// uniqueStrings removes duplicates from string slice
func uniqueStrings(input []string) []string {
	seen := make(map[string]bool)
	result := []string{}

	for _, item := range input {
		if item != "" && !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	return result
}
