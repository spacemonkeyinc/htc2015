package main

import (
	"io"

	"sm/dh"
	"sm/sssl"
	"sm/stdio"
)

func main() {
	conn := sssl.Server(stdio.NewStdioConn(), dh.Private())
	io.Copy(conn, conn)
}
