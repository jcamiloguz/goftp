package channel_test

import (
	"testing"

	"github.com/jcamiloguz/goftp/internal/channel"
)

func TestNewChannel(t *testing.T) {
	// Test cases
	tables := []struct {
		idChannel int
	}{
		{0},
		{1},
	}

	// Run test cases
	for _, table := range tables {
		_, err := channel.NewChannel(table.idChannel)
		if table.idChannel < 1 {
			if err == nil {
				t.Errorf("expected error, got nil")
			}
		} else {
			if err != nil {
				t.Errorf("expected nil, got %v", err)
			}
		}
	}
}

func TestBroadcast(t *testing.T) {
	testChannel, _ := channel.NewChannel(1)
	// Test cases
	tables := []struct {
		channel  *channel.Channel
		username string
		content  []byte
	}{
		{
			channel:  testChannel,
			username: "jcamiloguz",
			content:  []byte("Hello World!"),
		},
		{
			channel:  testChannel,
			username: "",
			content:  []byte("Hello World!"),
		},
		{
			channel:  testChannel,
			username: "jcamiloguz",
			content:  nil,
		},
	}

	// Run test cases
	for _, item := range tables {
		err := item.channel.Broadcast(item.username, item.content)
		if item.username == "" || item.content == nil {
			if err == nil {
				t.Errorf("expected error, got nil")
			}
		} else {
			if err != nil {
				t.Errorf("expected nil, got %v", err)
			}
		}
	}

}
