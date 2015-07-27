package main

import (
	"fmt"
	"io"
	"net"
	"os"

	"sm/dh"
	"sm/sssl"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <address>\n", os.Args[0])
		os.Exit(1)
	}
	c, err := net.Dial("tcp", os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	conn := sssl.Client(c, dh.Private())
	defer conn.Close()

	go func() {
		defer conn.Close()
		io.Copy(os.Stdout, conn)
	}()
	io.Copy(conn, os.Stdin)
}
