package netcheck

import (
	"bytes"
	"encoding/json"
	"net"
	"os/exec"
	"strings"
	"syscall"
)

type Adapter struct {
	Name         string
	Desc         string // Windows InterfaceDescription (e.g. Phantom TAP-Windows Adapter V9)
	HardwareAddr string
	IPs          []string
}

// psAdapter maps to PowerShell Get-NetAdapter output
type psAdapter struct {
	Name                 string `json:"Name"`
	InterfaceDescription string `json:"InterfaceDescription"`
}

func GetActiveAdapters() ([]Adapter, error) {
	// Ambil metadata dari PowerShell agar nama deskriptif (Avira/Phantom) terbaca
	descMap := make(map[string]string)
	cmd := exec.Command("powershell", "-NoProfile", "-Command", "Get-NetAdapter | Select-Object Name, InterfaceDescription | ConvertTo-Json")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err == nil {
		var ps []psAdapter
		if err := json.Unmarshal(out.Bytes(), &ps); err == nil {
			for _, p := range ps {
				descMap[p.Name] = p.InterfaceDescription
			}
		} else {
			// If it's a single object, json array unmarshal fails
			var single psAdapter
			if err := json.Unmarshal(out.Bytes(), &single); err == nil {
				descMap[single.Name] = single.InterfaceDescription
			}
		}
	}

	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var adapters []Adapter
	for _, iface := range interfaces {
		addrs, _ := iface.Addrs()
		
		var ips []string
		for _, addr := range addrs {
			ips = append(ips, addr.String())
		}

		desc := descMap[iface.Name]
		if desc == "" {
			desc = "Standard Adapter"
		}

		// Filter out irrelevant adapters
		lowerName := strings.ToLower(iface.Name)
		lowerDesc := strings.ToLower(desc)
		if strings.Contains(lowerName, "bluetooth") || strings.Contains(lowerDesc, "bluetooth") ||
			strings.Contains(lowerName, "loopback") || strings.Contains(lowerDesc, "loopback") ||
			strings.Contains(lowerName, "vmware") || strings.Contains(lowerDesc, "vmware") {
			continue
		}

		if iface.Flags&net.FlagUp == 0 || len(ips) == 0 {
			desc += " (Offline/Disconnected)"
		}

		adapters = append(adapters, Adapter{
			Name:         iface.Name,
			Desc:         desc,
			HardwareAddr: iface.HardwareAddr.String(),
			IPs:          ips,
		})
	}
	return adapters, nil
}
