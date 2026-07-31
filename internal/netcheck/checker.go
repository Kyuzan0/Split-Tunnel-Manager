package netcheck

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Checker probes public IP and LAN reachability.
type Checker struct {
	PublicIPURL string
	Timeout     time.Duration
	client      *http.Client
}

func NewChecker(publicIPURL string, timeout time.Duration) *Checker {
	if timeout == 0 {
		timeout = 8 * time.Second
	}
	if publicIPURL == "" {
		publicIPURL = "https://api.ipify.org"
	}
	return &Checker{
		PublicIPURL: publicIPURL,
		Timeout:     timeout,
		client:      &http.Client{Timeout: timeout},
	}
}

func (c *Checker) PublicIP(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.PublicIPURL, nil)
	if err != nil {
		return "", err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("netcheck: status %d", resp.StatusCode)
	}
	b, err := io.ReadAll(io.LimitReader(resp.Body, 64))
	if err != nil {
		return "", err
	}
	return string(b), nil
}
