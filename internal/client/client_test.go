package client_test

import (
	"log"
	"net"
	"sync"
	"testing"

	"github.com/jcamiloguz/goftp/internal/client"
	"github.com/jcamiloguz/goftp/internal/server"
)

func TestNewClient(t *testing.T) {
	// Simulate a client connection
	_, clientConn := net.Pipe()

	actions := make(chan *client.Action)

	// Test Cases
	tables := []struct {
		conn     net.Conn
		username string
		actions  chan *client.Action
	}{
		{nil, "", nil},
		{clientConn, "", nil},
		{clientConn, "", nil},
		{clientConn, "TestName", nil},
		{clientConn, "TestName", actions},
		{clientConn, "TestName", actions},
	}

	for _, item := range tables {
		// fmt.Printf("Testing with conn: %v, username: %s, role: %d, register: %v, deregister: %v\n", item.conn, item.username, item.role, item.register, item.deregister)
		client, err := client.NewClient(item.conn, item.username, item.actions)

		if item.conn == nil || item.username == "" || item.actions == nil {
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

func TestHandle(t *testing.T) {

	serverMock, connMock := net.Pipe()
	defer serverMock.Close()
	defer connMock.Close()

	s, err := server.NewServer(&server.Config{
		Host:      "localhost",
		Port:      "3090",
		NChannels: 3,
	})
	if err != nil {
		log.Fatal(err)
	}

	clientTest, _ := client.NewClient(connMock, "TestName", s.Actions)
	// Test Cases
	tables := []struct {
		msg    []byte
		expect client.ACTIONID
	}{
		{[]byte("register\n"), client.REG},
		{[]byte("ERR\n"), client.ERR},
		{[]byte("out\n"), client.OUT},
		{[]byte("publish\n"), client.PUB},
		{[]byte("subscribe\n"), client.SUB},
		{[]byte("unsubscribe\n"), client.UNSUB},
		{[]byte("\n"), client.ERR},
	}

	for _, item := range tables {

		defer serverMock.Close()
		wg := sync.WaitGroup{}
		go func() {
			serverMock.Write(item.msg)
		}()
		wg.Add(1)
		go func() {
			select {
			case action := <-s.Actions:
				if action.Id != item.expect {
					t.Errorf("expected %d, got %d", item.expect, action.Id)
				}
				if action.Client.Username != "TestName" {
					t.Errorf("expected %s, got %s", "TestName", action.Client.Username)
				}
				wg.Done()
			}
		}()

		go clientTest.Read()
		wg.Wait()

	}

}
