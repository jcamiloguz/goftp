package client_test

import (
	"net"
	"testing"

	"github.com/jcamiloguz/goftp/internal/client"
)

func TestNewClient(t *testing.T) {
	// Simulate a client connection
	_, clientConn := net.Pipe()

	register := make(chan *client.Client)
	deregister := make(chan *client.Client)

	// Test Cases
	tables := []struct {
		conn       net.Conn
		username   string
		register   chan<- *client.Client
		deregister chan<- *client.Client
	}{
		{nil, "", nil, nil},
		{clientConn, "", nil, nil},
		{clientConn, "", nil, nil},
		{clientConn, "jcamiloguz", nil, nil},
		{clientConn, "jcamiloguz", nil, deregister},
		{clientConn, "jcamiloguz", register, deregister},
	}

	for _, item := range tables {
		// fmt.Printf("Testing with conn: %v, username: %s, role: %d, register: %v, deregister: %v\n", item.conn, item.username, item.role, item.register, item.deregister)
		client, err := client.NewClient(item.conn, item.username, item.register, item.deregister)

		if item.conn == nil || item.username == "" || item.register == nil || item.deregister == nil {
			if err == nil {
				t.Errorf("expected error, got nil")
			}
		} else {
			if err != nil {
				t.Errorf("expected nil, got %v", err)
			}
			if client.Username != item.username {
				t.Errorf("expected %s, got %s", item.username, client.Username)
			}
			if client.Connection != item.conn {
				t.Errorf("expected %v, got %v", item.conn, client.Connection)
			}
			if client.StartConnTime.IsZero() {
				t.Errorf("expected non-zero time, got zero")
			}
			if !client.EndConnTime.IsZero() {
				t.Errorf("expected zero time, got non-zero")
			}
		}
	}

}
