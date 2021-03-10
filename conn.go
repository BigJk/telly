package telly

import (
	"net"
	"sync"
)

type MessageHandler func(*Conn, string)
type DisconnectHandler func(*Conn)

// A Conn is a telnet connection.
type Conn struct {
	sync.Mutex
	conn       net.Conn
	msgHandler MessageHandler
	disHandler DisconnectHandler
}

// Close closes the connection.
func (c *Conn) Close() error {
	return c.conn.Close()
}

// RemoteAddr returns the remote network address. The Addr returned is shared by
// all invocations of RemoteAddr, so do not modify it.
func (c *Conn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

// SetMessageHandler set's the handler that will be invoked if one message is fully
// received from the telnet connection.
func (c *Conn) SetMessageHandler(handler MessageHandler) {
	c.Lock()
	defer c.Unlock()

	c.msgHandler = handler
}

// SetDisconnectHandler set's the handler that will be invoked if a connection is
// lost (disconnect happens).
func (c *Conn) SetDisconnectHandler(handler DisconnectHandler) {
	c.Lock()
	defer c.Unlock()

	c.disHandler = handler
}

// Write writes a linebreak terminated message to the connection.
func (c *Conn) Write(message string) error {
	_, err := c.conn.Write([]byte(message + "\n"))
	return err
}
