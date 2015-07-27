package sssl

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net"
	"os"
	"sync"
	"time"

	"sm/dh"
)

type Conn struct {
	client     bool
	private    *big.Int
	aead       cipher.AEAD
	conn       net.Conn
	mtx        sync.Mutex
	low_nonce  []byte
	high_nonce []byte
	buf        bytes.Buffer
}

var _ net.Conn = (*Conn)(nil)

func Client(conn net.Conn, private_key *big.Int) (c *Conn) {
	c = Server(conn, private_key)
	c.client = true
	return c
}

func Server(conn net.Conn, private_key *big.Int) *Conn {
	return &Conn{
		private: private_key,
		conn:    conn}
}

func (c *Conn) Handshake() (err error) {
	defer func() {
		if err != nil {
			c.Close()
		}
	}()
	c.mtx.Lock()
	defer c.mtx.Unlock()
	if c.aead != nil {
		return nil
	}
	if c.client {
		fmt.Fprintf(os.Stderr, "***** Client's private key: %d\n", c.private)
		fmt.Fprintf(os.Stderr, "***** Client's public key: %d\n",
			dh.Public(c.private))
		_, err := fmt.Fprintf(c.conn, "SimpleSSLv0\n%x\n", dh.Public(c.private))
		if err != nil {
			return err
		}
		server_public := new(big.Int)
		_, err = fmt.Fscanf(c.conn, "OK\n%x\n", server_public)
		if err != nil {
			return err
		}
		if c.client {
			fmt.Fprintf(os.Stderr, "***** Server's public key: %d\n", server_public)
		}
		return c.configure(server_public)
	}
	client_public := new(big.Int)
	_, err = fmt.Fscanf(c.conn, "SimpleSSLv0\n%x\n", client_public)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(c.conn, "OK\n%x\n", dh.Public(c.private))
	if err != nil {
		return err
	}
	return c.configure(client_public)
}

func (c *Conn) configure(other_public *big.Int) error {
	session := dh.SessionId(c.private, other_public)
	if c.client {
		fmt.Fprintf(os.Stderr, "***** Client's session key: %x\n", session[:16])
	}
	block, err := aes.NewCipher(session[:16])
	if err != nil {
		return err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}
	c.aead = aead
	c.low_nonce = make([]byte, aead.NonceSize())
	c.high_nonce = make([]byte, aead.NonceSize())
	for i := range c.high_nonce {
		c.high_nonce[i] = 0xff
	}
	return nil
}

func (c *Conn) incrementNonce() []byte {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	last := c.low_nonce
	next := new(big.Int).Add(new(big.Int).SetBytes(last), big.NewInt(1)).Bytes()
	if len(next) < len(last) {
		next = append(make([]byte, len(last)-len(next)), next...)
	}
	c.low_nonce = next
	return last
}

func (c *Conn) decrementNonce() []byte {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	last := c.high_nonce
	next := new(big.Int).Sub(new(big.Int).SetBytes(last), big.NewInt(1)).Bytes()
	if len(next) < len(last) {
		next = append(make([]byte, len(last)-len(next)), next...)
	}
	c.high_nonce = next
	return last
}

func (c *Conn) nextOutNonce() []byte {
	if c.client {
		return c.decrementNonce()
	}
	return c.incrementNonce()
}

func (c *Conn) nextInNonce() []byte {
	if c.client {
		return c.incrementNonce()
	}
	return c.decrementNonce()
}

func (c *Conn) Read(p []byte) (n int, err error) {
	err = c.Handshake()
	if err != nil {
		return 0, err
	}
	if c.buf.Len() > 0 {
		return c.buf.Read(p)
	}

	var size_buf [4]byte
	_, err = io.ReadFull(c.conn, size_buf[:])
	if err != nil {
		if c.client {
			fmt.Fprintf(os.Stderr,
				"***** Client failed reading header size from server\n")
		}
		return 0, err
	}

	size := binary.BigEndian.Uint32(size_buf[:])
	if c.client {
		fmt.Fprintf(os.Stderr,
			"***** Client read header size from server: %x (%d)\n",
			size_buf[:], size)
	}
	data := make([]byte, size)
	_, err = io.ReadFull(c.conn, data)
	if err != nil {
		if c.client {
			fmt.Fprintf(os.Stderr, "***** Client failed reading data from server\n")
		}
		return 0, err
	}

	if c.client {
		fmt.Fprintf(os.Stderr, "***** Client read %x from server\n", data)
	}
	plaintext, err := c.aead.Open(nil, c.nextInNonce(), data, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "***** Client failed decrypting data\n")
		return 0, err
	}

	_, err = c.buf.Write(plaintext)
	if err != nil {
		return 0, err
	}
	return c.buf.Read(p)
}

func (c *Conn) Write(p []byte) (n int, err error) {
	err = c.Handshake()
	if err != nil {
		return 0, err
	}

	ciphertext := c.aead.Seal(nil, c.nextOutNonce(), p, nil)
	var size [4]byte
	binary.BigEndian.PutUint32(size[:], uint32(len(ciphertext)))
	encrypted := append(size[:], ciphertext...)
	n, err = c.conn.Write(encrypted)
	if err != nil {
		return 0, err
	}
	if n != len(ciphertext)+4 {
		return 0, errors.New("partial ciphertext write")
	}
	if c.client {
		fmt.Fprintf(os.Stderr, "***** Client's encrypted and framed message: %x\n",
			encrypted)
	}
	return len(p), nil
}

func (c *Conn) Close() error         { return c.conn.Close() }
func (c *Conn) LocalAddr() net.Addr  { return c.conn.LocalAddr() }
func (c *Conn) RemoteAddr() net.Addr { return c.conn.RemoteAddr() }

func (c *Conn) SetDeadline(t time.Time) error {
	return c.conn.SetDeadline(t)
}

func (c *Conn) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

func (c *Conn) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}
