package main

import (
	"fmt"
	"io"
	"math/big"
	"os"

	"sm/sssl"
	"sm/stdio"
)

func main() {
	conn := sssl.Client(stdio.NewStdioConn(), big.NewInt(123))
	fmt.Fprintf(os.Stderr, "***** Client's message: hi\n")
	_, err := conn.Write([]byte("hi"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "***** Client failed writing message\n")
		os.Exit(1)
	}
	var p [2]byte
	_, err = io.ReadFull(conn, p[:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "***** Client failed reading server's response\n")
		os.Exit(1)
	}
	if string(p[:]) != "hi" {
		fmt.Fprintf(os.Stderr, "***** Client received %#v but expected %#v\n",
			string(p[:]), "hi")
		os.Exit(1)
	}
}
