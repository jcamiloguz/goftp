package client

import (
	"net"
	"time"
)

type ROLE int

const (
	SUB ROLE = iota
	PUB
)

type Client struct {
	Username      string
	Role          ROLE
	Connection    net.Conn
	StartConnTime time.Time
	EndConnTime   time.Time
	Register      chan<- *Client
	Deregister    chan<- *Client
}

func NewClient(conn net.Conn, username string, role ROLE, register chan<- *Client, deregister chan<- *Client) *Client {
	return &Client{
		Username:      username,
		Role:          role,
		Connection:    conn,
		StartConnTime: time.Now(),
		Register:      register,
		Deregister:    deregister,
	}
}
