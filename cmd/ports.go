package cmd

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

var PortsCmd = &cobra.Command{
	Use:   "ports",
	Short: "Display listening ports",
	Long:  `Display count of listening TCP/UDP ports and details about what's listening.`,
	Run: func(cmd *cobra.Command, args []string) {
		PrintPortsInfo()
	},
}

type PortInfo struct {
	TCPPorts int
	UDPPorts int
	Ports    []string
}

func PrintPortsInfo() {
	fmt.Println("Listening Ports:")
	fmt.Println()

	info := gatherPortsInfo()

	fmt.Printf("  TCP Ports: %d\n", info.TCPPorts)
	fmt.Printf("  UDP Ports: %d\n", info.UDPPorts)
	fmt.Printf("  Total:     %d\n", info.TCPPorts+info.UDPPorts)

	if len(info.Ports) > 0 {
		fmt.Println()
		fmt.Printf("  Details:\n")
		for _, port := range info.Ports {
			fmt.Printf("    %s\n", port)
		}
	}
}

func gatherPortsInfo() *PortInfo {
	info := &PortInfo{
		Ports: []string{},
	}

	osType := runtime.GOOS

	switch osType {
	case "linux":
		// Use ss (preferred) or netstat
		if exists("ss") {
			cmd := exec.Command("ss", "-tulpn")
			out, err := cmd.Output()
			if err == nil {
				lines := strings.Split(string(out), "\n")
				for _, line := range lines {
					line = strings.TrimSpace(line)
					if line == "" || strings.HasPrefix(line, "Netid") || strings.HasPrefix(line, "State") {
						continue
					}

					parts := strings.Fields(line)
					if len(parts) >= 5 {
						proto := parts[0]
						localAddr := parts[4]

						// Count by protocol
						if strings.HasPrefix(proto, "tcp") {
							info.TCPPorts++
						} else if strings.HasPrefix(proto, "udp") {
							info.UDPPorts++
						}

						// Extract port
						portParts := strings.Split(localAddr, ":")
						if len(portParts) >= 2 {
							port := portParts[len(portParts)-1]
							process := ""
							if len(parts) >= 7 {
								process = parts[6]
							}
							if process != "" {
								info.Ports = append(info.Ports, fmt.Sprintf("%s:%s (%s)", proto, port, process))
							} else {
								info.Ports = append(info.Ports, fmt.Sprintf("%s:%s", proto, port))
							}
						}
					}
				}
			}
		} else if exists("netstat") {
			cmd := exec.Command("netstat", "-tulpn")
			out, err := cmd.Output()
			if err == nil {
				lines := strings.Split(string(out), "\n")
				for _, line := range lines {
					line = strings.TrimSpace(line)
					if line == "" || strings.HasPrefix(line, "Active") || strings.HasPrefix(line, "Proto") {
						continue
					}

					parts := strings.Fields(line)
					if len(parts) >= 4 {
						proto := parts[0]
						localAddr := parts[3]

						// Count by protocol
						if strings.HasPrefix(proto, "tcp") {
							info.TCPPorts++
						} else if strings.HasPrefix(proto, "udp") {
							info.UDPPorts++
						}

						// Extract port
						portParts := strings.Split(localAddr, ":")
						if len(portParts) >= 2 {
							port := portParts[len(portParts)-1]
							process := ""
							if len(parts) >= 7 {
								process = parts[6]
							}
							if process != "" {
								info.Ports = append(info.Ports, fmt.Sprintf("%s:%s (%s)", proto, port, process))
							} else {
								info.Ports = append(info.Ports, fmt.Sprintf("%s:%s", proto, port))
							}
						}
					}
				}
			}
		}
	case "darwin":
		// Use lsof on macOS
		if exists("lsof") {
			cmd := exec.Command("lsof", "-iTCP", "-sTCP:LISTEN", "-P", "-n")
			out, err := cmd.Output()
			if err == nil {
				lines := strings.Split(string(out), "\n")
				for _, line := range lines {
					line = strings.TrimSpace(line)
					if line == "" || strings.HasPrefix(line, "COMMAND") {
						continue
					}

					parts := strings.Fields(line)
					if len(parts) >= 9 {
						info.TCPPorts++
						command := parts[0]
						portInfo := parts[8]
						info.Ports = append(info.Ports, fmt.Sprintf("tcp:%s (%s)", portInfo, command))
					}
				}
			}

			// UDP ports
			cmd = exec.Command("lsof", "-iUDP", "-P", "-n")
			out, err = cmd.Output()
			if err == nil {
				lines := strings.Split(string(out), "\n")
				for _, line := range lines {
					line = strings.TrimSpace(line)
					if line == "" || strings.HasPrefix(line, "COMMAND") {
						continue
					}

					parts := strings.Fields(line)
					if len(parts) >= 9 {
						info.UDPPorts++
						command := parts[0]
						portInfo := parts[8]
						info.Ports = append(info.Ports, fmt.Sprintf("udp:%s (%s)", portInfo, command))
					}
				}
			}
		}
	case "windows":
		// Use netstat on Windows
		if exists("netstat") {
			cmd := exec.Command("netstat", "-ano")
			out, err := cmd.Output()
			if err == nil {
				lines := strings.Split(string(out), "\n")
				for _, line := range lines {
					line = strings.TrimSpace(line)
					if line == "" || strings.HasPrefix(line, "Active") || strings.HasPrefix(line, "Proto") {
						continue
					}

					parts := strings.Fields(line)
					if len(parts) >= 4 {
						proto := parts[0]
						localAddr := parts[1]
						state := ""
						if len(parts) >= 4 {
							state = parts[3]
						}

						if state == "LISTENING" || proto == "UDP" {
							// Count by protocol
							if strings.HasPrefix(proto, "TCP") {
								info.TCPPorts++
							} else if strings.HasPrefix(proto, "UDP") {
								info.UDPPorts++
							}

							info.Ports = append(info.Ports, fmt.Sprintf("%s:%s", proto, localAddr))
						}
					}
				}
			}
		}
	}

	return info
}
