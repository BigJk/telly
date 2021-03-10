package telly

import (
	"net"
	"time"
)

const (
	READ_STATE_NORMAL  = 1
	READ_STATE_COMMAND = 2
	READ_STATE_SUBNEG  = 3

	TN_INTERPRET_AS_COMMAND = 255
	TN_ARE_YOU_THERE        = 246
	TN_WILL                 = 251
	TN_WONT                 = 252
	TN_DO                   = 253
	TN_DONT                 = 254
	TN_SUBNEGOTIATION_START = 250
	TN_SUBNEGOTIATION_END   = 240
)

// A Listener is a telnet server listener.
type Listener struct {
	listener net.Listener
	timeout  time.Duration
}

// Listen starts a new telnet listener that can accept connections.
func Listen(bind string) (*Listener, error) {
	listener, err := net.Listen("tcp4", bind)
	if err != nil {
		return nil, err
	}

	return &Listener{
		listener: listener,
		timeout:  0,
	}, nil
}

// SetTimeout set's the duration that the server will wait for data before closing
// the connection. If you don't set a timeout it's easy to get overflowed by idle
// connections.
func (l *Listener) SetTimeout(dur time.Duration) {
	l.timeout = dur
}

// Accept waits for and returns the next connection to the listener.
func (l *Listener) Accept() (*Conn, error) {
	conn, err := l.listener.Accept()
	if err != nil {
		return nil, err
	}

	telCon := &Conn{
		conn: conn,
	}

	go func() {
		state := READ_STATE_NORMAL
		buf := make([]byte, 2048)
		curMsg := ""

		for {
			if int64(l.timeout) > 0 {
				if err := conn.SetReadDeadline(time.Now().Add(l.timeout)); err != nil {
					break
				}
			}

			read, err := conn.Read(buf)
			if err != nil {
				break
			}

			state = READ_STATE_NORMAL
			for i := 0; i < read; i++ {
				switch state {
				case READ_STATE_NORMAL:
					if buf[i] == TN_INTERPRET_AS_COMMAND {
						state = READ_STATE_COMMAND
					} else if buf[i] == '\n' {
						telCon.Lock()
						if telCon.msgHandler != nil {
							telCon.msgHandler(telCon, curMsg)
						}
						telCon.Unlock()
						curMsg = ""
					} else if buf[i] == '\x08' {
						if len(curMsg) > 0 {
							curMsg = curMsg[:len(curMsg)-1]
						}
					} else {
						curMsg += string(buf[i])
					}
				case READ_STATE_COMMAND:
					if buf[i] == TN_SUBNEGOTIATION_START {
						state = READ_STATE_SUBNEG
					} else {
						switch buf[i] {
						case TN_WILL:
							fallthrough
						case TN_WONT:
							fallthrough
						case TN_DO:
							fallthrough
						case TN_DONT:
							state = READ_STATE_COMMAND
						default:
							state = READ_STATE_NORMAL
						}
					}
				case READ_STATE_SUBNEG:
					if buf[i] == TN_SUBNEGOTIATION_END {
						state = READ_STATE_NORMAL
					}
				}
			}
		}

		_ = conn.Close()

		telCon.Lock()
		if telCon.disHandler != nil {
			telCon.disHandler(telCon)
		}
		telCon.Unlock()
	}()

	return telCon, err
}
