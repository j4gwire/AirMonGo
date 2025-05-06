package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// Define the main entry point for the program
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

// Print banner with tool's version
func printBanner() {
	fmt.Println("AirMonGo - Wi-Fi Scanner v1.1")
	fmt.Println(strings.Repeat("-", 40))
}

// Print help information
func printHelp() {
	fmt.Println("Usage: AirMonGo [options]")
	fmt.Println("Options:")
	fmt.Println("  --help, -h  Show help information")
}

// Get the signal quality based on signal strength
func getSignalQuality(signalStrength int) string {
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

// Windows scanner
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
		}
		if strings.HasPrefix(trimmed, "BSSID") {
			bssid = strings.TrimPrefix(trimmed, "BSSID")
		}
		if strings.HasPrefix(trimmed, "Signal") {
			signal = strings.TrimPrefix(trimmed, "Signal")
		}
		if strings.HasPrefix(trimmed, "Authentication") {
			auth = strings.TrimPrefix(trimmed, "Authentication")
		}
		if strings.HasPrefix(trimmed, "Encryption") {
			encryption = strings.TrimPrefix(trimmed, "Encryption")
		}
		if strings.HasPrefix(trimmed, "Channel") {
			channel = strings.TrimPrefix(trimmed, "Channel")
		}
	}
	printWiFiDetails(ssid, bssid, signal, auth, encryption, channel)
}

// Print Wi-Fi details
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

// Linux scanner
func scanLinux() {
	cmd := exec.Command("nmcli", "-t", "-f", "SSID,SIGNAL,SECURITY,BSSID,CHAN", "dev", "wifi")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error:", err)
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

// Print Wi-Fi details for Linux with channel information
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

// macOS scanner
func scanMac() {
	cmd := exec.Command("/System/Library/PrivateFrameworks/Apple80211.framework/Versions/Current/Resources/airport", "-s")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Available Wi-Fi Networks (macOS):")
	lines := strings.Split(out.String(), "\n")
	for _, line := range lines {
		if line != "" {
			fields := strings.Fields(line)
			if len(fields) >= 7 {
				ssid := fields[0]
				signal := fields[1]
				bssid := fields[2]
				channel := fields[4]
				security := fields[6] // assuming the last field is security
				printWiFiDetailsMac(ssid, bssid, signal, security, channel)
			}
		}
	}
}

// Print Wi-Fi details for macOS with channel information
func printWiFiDetailsMac(ssid, bssid, signal, security, channel string) {
	signalStrength := 0
	fmt.Sscanf(signal, "%d", &signalStrength)
	signalQuality := getSignalQuality(signalStrength)
	fmt.Printf("\n----------------------------------------\n")
	fmt.Printf("SSID       : %-30s\n", ssid)
	fmt.Printf("Signal     : %-10s%%  (%s)\n", signal, signalQuality)
	fmt.Printf("BSSID      : %-30s\n", bssid)
	fmt.Printf("Security   : %-20s\n", security)
	fmt.Printf("Channel    : %-10s\n", channel)
	fmt.Println("----------------------------------------")
}
