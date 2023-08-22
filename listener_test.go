package telnetter

import (
	"net"
	"testing"
	"time"
)

func TestListener_Accept(t *testing.T) {
	listener, err := Listen("127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		conn, err := listener.Accept()
		if err != nil {
			t.Error(err)
		}
		defer conn.Close()
	}()

	client, err := net.Dial("tcp", listener.listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	time.Sleep(100 * time.Millisecond) // Give some time for the connection to be established
}

func TestListener_SetTimeout(t *testing.T) {
	listener, err := Listen("127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	listener.SetTimeout(5 * time.Second)
	if listener.timeout != 5*time.Second {
		t.Errorf("expected timeout to be set to 5 seconds, got %v", listener.timeout)
	}
}
