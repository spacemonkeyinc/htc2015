package main

import (
	"io"
	"os"

	"sm/dh"
	"sm/sssl"
	"sm/stdio"
)

func main() {
	conn := sssl.Client(stdio.NewStdioConn(), dh.Private())
	_, err := conn.Write([]byte("hello, world!"))
	if err != nil {
		os.Exit(1)
	}
	var p [13]byte
	_, err = io.ReadFull(conn, p[:])
	if err != nil {
		os.Exit(1)
	}
	if string(p[:]) != "hello, world!" {
		os.Exit(1)
	}
}
