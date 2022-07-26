package channel_test

import (
	"testing"

	ch "github.com/jcamiloguz/goftp/internal/channel"
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
		_, err := ch.NewChannel(table.idChannel)
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

func TestNewFile(t *testing.T) {

	// Test cases
	tables := []struct {
		name string
		size int
	}{
		{"", 0},
		{"test", 0},
		{"test", 1},
	}

	// Run test cases
	for _, table := range tables {
		_, err := ch.NewFile(table.name, table.size)
		if table.name == "" || table.size < 1 {
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

// func TestBroadcast(t *testing.T) {
// 	chann:= ch.NewChannel(0)
// 	_,publisherConn := net.Pipe()
// 	publisher := cl.NewClient( )
// 	fileTestContent := []byte("test")
// 	// Test cases
// 	tables:=[]struct {
// 		nSubscribers int
// 		file 			 *ch.File
// 	}{

// 	}

// 	}

// }
