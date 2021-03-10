# telly
Minimal Telnet Server I wrote for a MUD backend. 

## Example
```go
package main

import (
	"fmt"
	"github.com/BigJk/telly"
	"time"
)

func main() {
	listener, err := telly.Listen(":5050")
	if err != nil {
		panic(err)
	}
	
	// If no message is received for 60 seconds close connection.
	listener.SetTimeout(time.Second * 60)

	for {
		if conn, err := listener.Accept(); err == nil {
			conn.SetMessageHandler(func(conn *telly.Conn, s string) {
				fmt.Println(s)
				_ = conn.Write(s)
			})

			conn.SetDisconnectHandler(func(conn *telly.Conn) {
				fmt.Println(conn.RemoteAddr(), "disconnected")
			})

			fmt.Println(conn.RemoteAddr(), "connected")
		}
	}
}
```

## Why didn't I use ``github.com/reiver/go-telnet``?

The Project might be RFC conform but seems abondoned and doesn't support two crucial features. There is no clear way to close a connection from the server side and no support for custom timeouts. Without custom timeouts it's easy to overflow a server with idle connections.

## Used References
- https://github.com/Frimkron/mud-pi/blob/master/mudserver.py#L327
- http://pcmicro.com/netfoss/telnet.html
