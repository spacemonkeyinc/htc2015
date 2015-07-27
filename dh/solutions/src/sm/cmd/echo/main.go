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
	l, err := net.Listen("tcp", os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	defer l.Close()
	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		conn := sssl.Server(c, dh.Private())
		go func() {
			defer conn.Close()
			io.Copy(conn, conn)
		}()
	}
}
