package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"
)

type ContainerInfo struct {
	ID      string   `json:"id"`
	Names   []string `json:"names"`
	Image   string   `json:"image"`
	State   string   `json:"state"`
	Status  string   `json:"status"`
	Created int64    `json:"created"`
}

type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	// Unix Domain Socket HTTP Transport for /var/run/docker.sock
	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return net.Dial("unix", "/var/run/docker.sock")
		},
	}
	return &Client{
		httpClient: &http.Client{
			Transport: transport,
			Timeout:   10 * time.Second,
		},
	}
}

func (c *Client) ListContainers() ([]ContainerInfo, error) {
	resp, err := c.httpClient.Get("http://localhost/v1.41/containers/json?all=true")
	if err != nil {
		return nil, fmt.Errorf("docker socket unreadable or docker inactive: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("docker API returned status: %s", resp.Status)
	}

	var containers []ContainerInfo
	if err := json.NewDecoder(resp.Body).Decode(&containers); err != nil {
		return nil, err
	}
	return containers, nil
}

func (c *Client) ActionContainer(id string, action string) error {
	allowedActions := map[string]bool{"start": true, "stop": true, "restart": true}
	if !allowedActions[action] {
		return fmt.Errorf("unsupported container action: %s", action)
	}

	url := fmt.Sprintf("http://localhost/v1.41/containers/%s/%s", id, action)
	resp, err := c.httpClient.Post(url, "application/json", nil)
	if err != nil {
		return fmt.Errorf("failed to %s container %s: %w", action, id, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("docker %s container failed with status: %s", action, resp.Status)
	}
	return nil
}
