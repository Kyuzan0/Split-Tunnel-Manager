package domain

import "time"

// BypassRule excludes traffic from the VPN tunnel (split tunnel v1).
type BypassRule struct {
	ID      string
	CIDR    string // e.g. 192.168.2.0/24
	Label   string
	Enabled bool
	Source  string // "default" | "user"
}

// ConnState is the connection FSM for UI.
type ConnState string

const (
	ConnIdle       ConnState = "idle"
	ConnApplying   ConnState = "applying"
	ConnConnected  ConnState = "connected"
	ConnError      ConnState = "error"
)

// SplitStatus is what the Connect screen binds to.
type SplitStatus struct {
	State        ConnState
	InterfaceName string // Name of the adapter the routes were applied to
	ConnectedAt  *time.Time
	LastError    string
	ActiveBypass []BypassRule
}
