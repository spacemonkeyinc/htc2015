package stdio

import (
	"errors"
	"io"
	"net"
	"os"
	"time"
)

type StdioAddr struct{}

func (StdioAddr) Network() string { return "stdio" }
func (StdioAddr) String() string  { return "stdio" }

type StdioConn struct {
	Output io.WriteCloser
	Input  io.ReadCloser
}

func NewStdioConn() *StdioConn {
	return &StdioConn{
		Output: os.Stdout,
		Input:  os.Stdin}

}

var _ net.Conn = (*StdioConn)(nil)

func (c *StdioConn) Close() error {
	c.Input.Close()
	c.Output.Close()
	return nil
}

func (c *StdioConn) Read(p []byte) (n int, err error) {
	return c.Input.Read(p)
}

func (c *StdioConn) Write(p []byte) (n int, err error) {
	return c.Output.Write(p)
}

func (c *StdioConn) LocalAddr() net.Addr  { return StdioAddr{} }
func (c *StdioConn) RemoteAddr() net.Addr { return StdioAddr{} }

func (c *StdioConn) SetDeadline(t time.Time) error {
	return errors.New("not supported")
}

func (c *StdioConn) SetReadDeadline(t time.Time) error {
	return errors.New("not supported")
}

func (c *StdioConn) SetWriteDeadline(t time.Time) error {
	return errors.New("not supported")
}
