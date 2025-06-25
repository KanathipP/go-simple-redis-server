package main

import (
	"fmt"
	"testing"
)

func TestProtocol(t *testing.T) {
	raw := "*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n"

	cmd, err := parseCommnand(raw)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(cmd)
}
