package main

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"sync"
	"time"

	"sm/dh"
	"sm/sssl"
	"sm/stdio"
)

type WriteCloser struct {
	io.Writer
	io.Closer
}

type ReadCloser struct {
	io.Reader
	io.Closer
}

func main() {
	underlying := stdio.NewStdioConn()
	conn := sssl.Client(underlying, dh.Private())
	err := conn.Handshake()
	if err != nil {
		os.Exit(1)
	}

	r := rand.New(rand.NewSource(time.Now().Unix()))

	buf_out := &bytes.Buffer{}
	underlying.Output = WriteCloser{
		Writer: buf_out,
		Closer: ioutil.NopCloser(nil)}
	buf_in := bufio.NewReaderSize(os.Stdin, 1024*1024)
	underlying.Input = ReadCloser{
		Reader: buf_in,
		Closer: os.Stdin}

	var expected int64
	bytes_remaining := 63 * 1024 * 1024
	var buffer [65536]byte
	for bytes_remaining > 0 {
		amount := r.Intn(len(buffer))
		if amount > bytes_remaining {
			amount = bytes_remaining
		}
		bytes_remaining -= amount
		// just write out different sizes of zeros
		buf_slice := buffer[:amount]
		_, err = conn.Write(buf_slice)
		if err != nil {
			panic(err)
		}
		expected += int64(amount)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		_, err := buf_out.WriteTo(os.Stdout)
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		defer wg.Done()
		var total int64
		var buffer [65536]byte
		var zeros [65536]byte
		for total < expected {
			amount := int64(len(buffer))
			if expected-total < amount {
				amount = expected - total
			}
			buf_slice := buffer[:amount]
			n, err := io.ReadFull(conn, buf_slice)
			if err != nil {
				panic(err)
			}
			if !bytes.Equal(buf_slice[:n], zeros[:amount]) {
				panic("not equal")
			}
			total += amount
		}
	}()

	wg.Wait()
}
