package telnetter

import (
	"net"
	"testing"
)

func TestConn_ReadWrite(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	go func() {
		conn, err := ln.Accept()
		if err != nil {
			t.Error(err)
			return
		}
		defer conn.Close()

		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			t.Error(err)
			return
		}
		if string(buf[:n]) != "hello" {
			t.Errorf("expected 'hello', got '%s'", buf[:n])
		}
	}()

	client, err := net.Dial("tcp", ln.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	tConn := &Conn{rw: client, conn: client}
	err = tConn.WriteString("hello")
	if err != nil {
		t.Errorf("failed to write to connection: %s", err)
	}
}

func TestConn_SetMessageHandler(t *testing.T) {
	c := &Conn{}
	handler := func(conn *Conn, msg string) {}
	c.SetMessageHandler(handler)

	if c.msgHandler == nil {
		t.Errorf("failed to set message handler")
	}
}

func TestConn_SetDisconnectHandler(t *testing.T) {
	c := &Conn{}
	handler := func(conn *Conn) {}
	c.SetDisconnectHandler(handler)

	if c.disHandler == nil {
		t.Errorf("failed to set disconnect handler")
	}
}
