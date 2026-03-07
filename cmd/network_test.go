package cmd

import (
	"io"
	"os"
	"strings"
	"testing"
)

// TestIsVPNInterface tests VPN interface name detection — pure function, no I/O
func TestIsVPNInterface(t *testing.T) {
	tests := []struct {
		name     string
		iface    string
		expected bool
	}{
		{"tun0 is VPN", "tun0", true},
		{"tun1 is VPN", "tun1", true},
		{"tap0 is VPN", "tap0", true},
		{"wg0 WireGuard is VPN", "wg0", true},
		{"wg1 WireGuard is VPN", "wg1", true},
		{"ppp0 is VPN", "ppp0", true},
		{"vpn0 is VPN", "vpn0", true},
		{"eth0 is not VPN", "eth0", false},
		{"en0 is not VPN", "en0", false},
		{"lo loopback is not VPN", "lo", false},
		{"Ethernet is not VPN", "Ethernet", false},
		{"Wi-Fi is not VPN", "Wi-Fi", false},
		{"empty string is not VPN", "", false},
		{"utun0 macOS VPN is not matched (no utun prefix)", "utun0", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isVPNInterface(tt.iface)
			if got != tt.expected {
				t.Errorf("isVPNInterface(%q) = %v, want %v", tt.iface, got, tt.expected)
			}
		})
	}
}

// TestUniqueStrings tests string deduplication — pure function, no I/O
func TestUniqueStrings(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{"empty slice", []string{}, []string{}},
		{"single element", []string{"a"}, []string{"a"}},
		{"no duplicates", []string{"a", "b", "c"}, []string{"a", "b", "c"}},
		{"with duplicates", []string{"a", "b", "a", "c"}, []string{"a", "b", "c"}},
		{"all duplicates", []string{"a", "a", "a"}, []string{"a"}},
		{"empty strings are filtered", []string{"", "a", "", "b"}, []string{"a", "b"}},
		{"preserves insertion order", []string{"3", "1", "2", "1", "3"}, []string{"3", "1", "2"}},
		{"DNS-like IPs deduplicated", []string{"1.1.1.1", "8.8.8.8", "1.1.1.1"}, []string{"1.1.1.1", "8.8.8.8"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := uniqueStrings(tt.input)
			if len(got) != len(tt.expected) {
				t.Errorf("uniqueStrings(%v) len=%d, want len=%d; got=%v", tt.input, len(got), len(tt.expected), got)
				return
			}
			for i, v := range tt.expected {
				if got[i] != v {
					t.Errorf("uniqueStrings(%v)[%d] = %q, want %q", tt.input, i, got[i], v)
				}
			}
		})
	}
}

// TestDetectVPNInterface tests VPN detection state machine — no I/O, just struct manipulation
func TestDetectVPNInterface(t *testing.T) {
	t.Run("detects active VPN interface", func(t *testing.T) {
		details := &NetworkDetails{
			Interfaces: []InterfaceInfo{
				{Name: "eth0", IP: "192.168.1.10", Status: "UP", IsVPN: false},
				{Name: "tun0", IP: "10.8.0.2", Status: "UP", IsVPN: true},
			},
		}
		detectVPNInterface(details)
		if !details.VPNActive {
			t.Error("expected VPNActive=true")
		}
		if details.VPNInterface == nil || details.VPNInterface.Name != "tun0" {
			t.Errorf("expected VPNInterface.Name=tun0, got %v", details.VPNInterface)
		}
	})

	t.Run("ignores VPN interface that is DOWN", func(t *testing.T) {
		details := &NetworkDetails{
			Interfaces: []InterfaceInfo{
				{Name: "tun0", IP: "", Status: "DOWN", IsVPN: true},
			},
		}
		detectVPNInterface(details)
		if details.VPNActive {
			t.Error("expected VPNActive=false when VPN interface is DOWN")
		}
		if details.VPNInterface != nil {
			t.Errorf("expected VPNInterface=nil, got %v", details.VPNInterface)
		}
	})

	t.Run("no VPN when only normal interfaces present", func(t *testing.T) {
		details := &NetworkDetails{
			Interfaces: []InterfaceInfo{
				{Name: "eth0", IP: "192.168.1.10", Status: "UP", IsVPN: false},
				{Name: "lo", IP: "127.0.0.1", Status: "UP", IsVPN: false},
			},
		}
		detectVPNInterface(details)
		if details.VPNActive {
			t.Error("expected VPNActive=false when no VPN interfaces")
		}
	})

	t.Run("no VPN when no interfaces at all", func(t *testing.T) {
		details := &NetworkDetails{Interfaces: []InterfaceInfo{}}
		detectVPNInterface(details)
		if details.VPNActive {
			t.Error("expected VPNActive=false for empty interface list")
		}
	})
}

// TestCheckInternetConnectivity verifies no panic and returns a bool.
// In a real network environment this should return true.
func TestCheckInternetConnectivity(t *testing.T) {
	result := checkInternetConnectivity()
	t.Logf("checkInternetConnectivity() = %v", result)
	// Result can be true or false depending on environment — we just verify no panic.
}

// TestGatherNetworkDetails verifies the function returns a fully-initialized struct.
func TestGatherNetworkDetails(t *testing.T) {
	details := gatherNetworkDetails()
	if details == nil {
		t.Fatal("gatherNetworkDetails() returned nil")
	}
	if details.Interfaces == nil {
		t.Error("gatherNetworkDetails() Interfaces is nil, expected initialized slice")
	}
	t.Logf("interfaces=%d, primaryIface=%v, vpnActive=%v, internetOK=%v",
		len(details.Interfaces), details.PrimaryIface != nil, details.VPNActive, details.InternetOK)
}

// TestGetPublicIP verifies the function returns either empty string or a valid IPv4 address.
func TestGetPublicIP(t *testing.T) {
	ip := getPublicIP()
	if ip != "" && strings.Count(ip, ".") != 3 {
		t.Errorf("getPublicIP() = %q, want empty string or valid IPv4 (3 dots)", ip)
	}
	t.Logf("getPublicIP() = %q", ip)
}

// TestPrintNetworkInfo verifies PrintNetworkInfo produces output without panicking.
func TestPrintNetworkInfo(t *testing.T) {
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stdout = w

	PrintNetworkInfo()

	w.Close()
	os.Stdout = old

	var sb strings.Builder
	if _, err := io.Copy(&sb, r); err != nil {
		t.Fatalf("failed to read captured output: %v", err)
	}
	output := sb.String()

	if !strings.Contains(output, "Connectivity:") {
		t.Errorf("PrintNetworkInfo() missing 'Connectivity:' section\noutput:\n%s", output)
	}
}
