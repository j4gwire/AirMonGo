package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--help" || os.Args[1] == "-h") {
		printHelp()
		return
	}

	printBanner()

	osType := runtime.GOOS

	switch osType {
	case "windows":
		scanWindows()
	case "linux":
		scanLinux()
	case "darwin":
		scanMac()
	default:
		fmt.Println("Unsupported OS:", osType)
	}
}

func printBanner() {
	fmt.Println("AirMonGo - Wi-Fi Scanner v1.2")
	fmt.Println(strings.Repeat("-", 40))
}

func printHelp() {
	fmt.Println("Usage: AirMonGo [options]")
	fmt.Println("Options:")
	fmt.Println("  --help, -h  Show help information")
}

// Works with both dBm (negative) and percentage values
func getSignalQuality(signalStrength int) string {
	if signalStrength < 0 {
		// dBm scale
		if signalStrength >= -50 {
			return "Excellent"
		} else if signalStrength >= -60 {
			return "Good"
		} else if signalStrength >= -70 {
			return "Fair"
		} else {
			return "Poor"
		}
	}
	// Percentage scale
	if signalStrength >= 75 {
		return "Excellent"
	} else if signalStrength >= 50 {
		return "Good"
	} else if signalStrength >= 30 {
		return "Fair"
	} else {
		return "Poor"
	}
}

func scanWindows() {
	cmd := exec.Command("netsh", "wlan", "show", "networks", "mode=bssid")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error executing netsh:", err)
		if stderr.Len() > 0 {
			fmt.Println("Details:", stderr.String())
		}
		return
	}

	fmt.Println("Available Wi-Fi Networks (Windows):")
	lines := strings.Split(out.String(), "\n")
	var ssid, bssid, signal, auth, encryption, channel string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "SSID") {
			if ssid != "" {
				printWiFiDetails(ssid, bssid, signal, auth, encryption, channel)
			}
			ssid = strings.TrimPrefix(trimmed, "SSID")
			ssid = strings.TrimSpace(strings.TrimPrefix(ssid, ":"))
		}
		if strings.HasPrefix(trimmed, "BSSID") {
			bssid = strings.TrimPrefix(trimmed, "BSSID")
			bssid = strings.TrimSpace(strings.TrimPrefix(bssid, ":"))
		}
		if strings.HasPrefix(trimmed, "Signal") {
			signal = strings.TrimPrefix(trimmed, "Signal")
			signal = strings.TrimSpace(strings.TrimPrefix(signal, ":"))
		}
		if strings.HasPrefix(trimmed, "Authentication") {
			auth = strings.TrimPrefix(trimmed, "Authentication")
			auth = strings.TrimSpace(strings.TrimPrefix(auth, ":"))
		}
		if strings.HasPrefix(trimmed, "Encryption") {
			encryption = strings.TrimPrefix(trimmed, "Encryption")
			encryption = strings.TrimSpace(strings.TrimPrefix(encryption, ":"))
		}
		if strings.HasPrefix(trimmed, "Channel") {
			channel = strings.TrimPrefix(trimmed, "Channel")
			channel = strings.TrimSpace(strings.TrimPrefix(channel, ":"))
		}
	}
	if ssid != "" {
		printWiFiDetails(ssid, bssid, signal, auth, encryption, channel)
	}
}

func printWiFiDetails(ssid, bssid, signal, auth, encryption, channel string) {
	signalStrength := 0
	fmt.Sscanf(signal, "%d", &signalStrength)
	signalQuality := getSignalQuality(signalStrength)
	fmt.Printf("\n----------------------------------------\n")
	fmt.Printf("SSID       : %-30s\n", ssid)
	fmt.Printf("BSSID      : %-30s\n", bssid)
	fmt.Printf("Signal     : %-10s%%  (%s)\n", signal, signalQuality)
	fmt.Printf("Auth       : %-20s\n", auth)
	fmt.Printf("Encryption : %-20s\n", encryption)
	fmt.Printf("Channel    : %-10s\n", channel)
	fmt.Println("----------------------------------------")
}

func scanLinux() {
	cmd := exec.Command("nmcli", "-t", "-f", "SSID,SIGNAL,SECURITY,BSSID,CHAN", "dev", "wifi")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error executing nmcli:", err)
		if stderr.Len() > 0 {
			fmt.Println("Details:", stderr.String())
		}
		return
	}

	fmt.Println("Available Wi-Fi Networks (Linux):")
	lines := strings.Split(out.String(), "\n")
	for _, line := range lines {
		if line != "" {
			fields := strings.Split(line, ":")
			if len(fields) >= 5 {
				ssid := fields[0]
				signal := fields[1]
				security := fields[2]
				bssid := fields[3]
				channel := fields[4]
				printWiFiDetailsLinux(ssid, bssid, signal, security, channel)
			}
		}
	}
}

func printWiFiDetailsLinux(ssid, bssid, signal, security, channel string) {
	signalStrength := 0
	fmt.Sscanf(signal, "%d", &signalStrength)
	signalQuality := getSignalQuality(signalStrength)
	fmt.Printf("\n----------------------------------------\n")
	fmt.Printf("SSID       : %-30s\n", ssid)
	fmt.Printf("Signal     : %-10s%%  (%s)\n", signal, signalQuality)
	fmt.Printf("Security   : %-20s\n", security)
	fmt.Printf("BSSID      : %-30s\n", bssid)
	fmt.Printf("Channel    : %-10s\n", channel)
	fmt.Println("----------------------------------------")
}

// Fixed: Using wdutil instead of deprecated airport command
func scanMac() {
	// Try wdutil first (works on macOS 14.4+)
	cmd := exec.Command("wdutil", "info")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	
	if err == nil && out.Len() > 0 {
		parseMacWdutil(out.String())
		return
	}

	// Fallback to system_profiler
	cmd = exec.Command("system_profiler", "SPAirPortDataType")
	out.Reset()
	stderr.Reset()
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	
	if err != nil {
		fmt.Println("Error executing system_profiler:", err)
		if stderr.Len() > 0 {
			fmt.Println("Details:", stderr.String())
		}
		return
	}

	parseMacSystemProfiler(out.String())
}

func parseMacWdutil(output string) {
	fmt.Println("Available Wi-Fi Networks (macOS):")
	lines := strings.Split(output, "\n")
	
	var ssid, bssid, signal, security, channel string
	inNetworkSection := false
	
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		
		if strings.Contains(trimmed, "SSID") && strings.Contains(trimmed, ":") {
			if ssid != "" && inNetworkSection {
				printWiFiDetailsMac(ssid, bssid, signal, security, channel)
				ssid, bssid, signal, security, channel = "", "", "", "", ""
			}
			parts := strings.SplitN(trimmed, ":", 2)
			if len(parts) == 2 {
				ssid = strings.TrimSpace(parts[1])
				inNetworkSection = true
			}
		} else if strings.Contains(trimmed, "BSSID") && strings.Contains(trimmed, ":") {
			parts := strings.SplitN(trimmed, ":", 2)
			if len(parts) == 2 {
				bssid = strings.TrimSpace(parts[1])
			}
		} else if strings.Contains(trimmed, "RSSI") && strings.Contains(trimmed, ":") {
			parts := strings.SplitN(trimmed, ":", 2)
			if len(parts) == 2 {
				signal = strings.TrimSpace(strings.Fields(parts[1])[0])
			}
		} else if strings.Contains(trimmed, "Security") && strings.Contains(trimmed, ":") {
			parts := strings.SplitN(trimmed, ":", 2)
			if len(parts) == 2 {
				security = strings.TrimSpace(parts[1])
			}
		} else if strings.Contains(trimmed, "Channel") && strings.Contains(trimmed, ":") {
			parts := strings.SplitN(trimmed, ":", 2)
			if len(parts) == 2 {
				channel = strings.TrimSpace(strings.Fields(parts[1])[0])
			}
		}
	}
	
	if ssid != "" && inNetworkSection {
		printWiFiDetailsMac(ssid, bssid, signal, security, channel)
	}
}

func parseMacSystemProfiler(output string) {
	fmt.Println("Available Wi-Fi Networks (macOS):")
	lines := strings.Split(output, "\n")
	
	var ssid, bssid, signal, security, channel string
	
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		
		if strings.HasPrefix(trimmed, "SSID:") || strings.Contains(trimmed, "Current Network Information:") {
			if ssid != "" {
				printWiFiDetailsMac(ssid, bssid, signal, security, channel)
				ssid, bssid, signal, security, channel = "", "", "", "", ""
			}
			if strings.HasPrefix(trimmed, "SSID:") {
				parts := strings.SplitN(trimmed, ":", 2)
				if len(parts) == 2 {
					ssid = strings.TrimSpace(parts[1])
				}
			}
		} else if strings.Contains(trimmed, "BSSID:") {
			parts := strings.SplitN(trimmed, ":", 2)
			if len(parts) >= 2 {
				bssid = strings.TrimSpace(strings.Join(parts[1:], ":"))
			}
		} else if strings.Contains(trimmed, "Signal / Noise:") {
			parts := strings.SplitN(trimmed, ":", 2)
			if len(parts) == 2 {
				signalPart := strings.TrimSpace(parts[1])
				signalFields := strings.Fields(signalPart)
				if len(signalFields) > 0 {
					signal = signalFields[0]
				}
			}
		} else if strings.Contains(trimmed, "Signal:") && !strings.Contains(trimmed, "Noise") {
			parts := strings.SplitN(trimmed, ":", 2)
			if len(parts) == 2 {
				signal = strings.TrimSpace(strings.Fields(parts[1])[0])
			}
		} else if strings.Contains(trimmed, "Security:") {
			parts := strings.SplitN(trimmed, ":", 2)
			if len(parts) == 2 {
				security = strings.TrimSpace(parts[1])
			}
		} else if strings.Contains(trimmed, "Channel:") {
			parts := strings.SplitN(trimmed, ":", 2)
			if len(parts) == 2 {
				channelPart := strings.TrimSpace(parts[1])
				channelFields := strings.Fields(channelPart)
				if len(channelFields) > 0 {
					channel = channelFields[0]
				}
			}
		}
	}
	
	if ssid != "" {
		printWiFiDetailsMac(ssid, bssid, signal, security, channel)
	}
}

func printWiFiDetailsMac(ssid, bssid, signal, security, channel string) {
	signalStrength := 0
	
	if signal != "" {
		cleanSignal := strings.TrimSpace(signal)
		cleanSignal = strings.Replace(cleanSignal, "dBm", "", -1)
		cleanSignal = strings.Replace(cleanSignal, "%", "", -1)
		cleanSignal = strings.TrimSpace(cleanSignal)
		
		if val, err := strconv.Atoi(cleanSignal); err == nil {
			signalStrength = val
		}
	}
	
	signalQuality := getSignalQuality(signalStrength)
	
	fmt.Printf("\n----------------------------------------\n")
	fmt.Printf("SSID       : %-30s\n", ssid)
	
	if signalStrength < 0 {
		fmt.Printf("Signal     : %-10s dBm  (%s)\n", signal, signalQuality)
	} else {
		fmt.Printf("Signal     : %-10s%%  (%s)\n", signal, signalQuality)
	}
	
	fmt.Printf("BSSID      : %-30s\n", bssid)
	fmt.Printf("Security   : %-20s\n", security)
	fmt.Printf("Channel    : %-10s\n", channel)
	fmt.Println("----------------------------------------")
}
