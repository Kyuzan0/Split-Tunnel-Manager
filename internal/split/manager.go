package split

import (
	"errors"
	"fmt"
	"net"
	"os/exec"
	"strings"
	"syscall"

	"split-tunnel-manager/internal/domain"
)

// Manager applies/removes bypass CIDR routes on Windows (stub).
type Manager struct {
	applied []domain.BypassRule
}

func NewManager() *Manager { return &Manager{} }

func getDefaultGateway() (string, error) {
	cmd := exec.Command("powershell", "-NoProfile", "-Command", "(Get-NetRoute -DestinationPrefix '0.0.0.0/0' | Sort-Object RouteMetric | Select-Object -First 1).NextHop")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	gw := strings.TrimSpace(string(out))
	if gw == "" {
		return "", errors.New("default gateway not found")
	}
	return gw, nil
}

func (m *Manager) Apply(rules []domain.BypassRule) error {
	gw, err := getDefaultGateway()
	if err != nil {
		return fmt.Errorf("failed to get default gateway: %w", err)
	}

	// Hapus semua rule yang pernah diaplikasikan di sesi ini sebelumnya
	_ = m.RemoveAll()

	for _, r := range rules {
		ip, ipNet, err := net.ParseCIDR(r.CIDR)
		if err != nil {
			return fmt.Errorf("%w: %s", domain.ErrInvalidCIDR, r.CIDR)
		}

		mask := net.IP(ipNet.Mask).String()
		network := ip.Mask(ipNet.Mask).String()

		cmd := exec.Command("route", "add", network, "mask", mask, gw, "metric", "1")
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		if err := cmd.Run(); err != nil {
			// Best effort, continue
			fmt.Printf("Warning: failed to add route for %s: %v\n", r.CIDR, err)
		}
	}
	m.applied = append([]domain.BypassRule(nil), rules...)
	return nil
}

func (m *Manager) RemoveAll() error {
	for _, r := range m.applied {
		ip, ipNet, err := net.ParseCIDR(r.CIDR)
		if err != nil {
			continue
		}

		mask := net.IP(ipNet.Mask).String()
		network := ip.Mask(ipNet.Mask).String()

		cmd := exec.Command("route", "delete", network, "mask", mask)
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		_ = cmd.Run()
	}
	m.applied = nil
	return nil
}

func (m *Manager) RemoveRule(cidr string) error {
	ip, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return err
	}

	mask := net.IP(ipNet.Mask).String()
	network := ip.Mask(ipNet.Mask).String()

	cmd := exec.Command("route", "delete", network, "mask", mask)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd.Run()
}

func (m *Manager) Active() []domain.BypassRule {
	return append([]domain.BypassRule(nil), m.applied...)
}
