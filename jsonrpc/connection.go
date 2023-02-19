package jsonrpc

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Connection interface {
	Close()
	Receive() (*Message, error)
	Send(msg Message) error
}

func NewStdioConnection() Connection {
	return &stdioConnection{
		connectionBase: connectionBase{
			mtx: sync.Mutex{},
			rd:  bufio.NewReader(os.Stdin),
			wr:  bufio.NewWriter(os.Stdout),
		},
		in:  os.Stdin,
		out: os.Stdout,
	}
}

func NewNetConnection(c net.Conn) Connection {
	return &netConnection{
		connectionBase: connectionBase{
			mtx: sync.Mutex{},
			rd:  bufio.NewReader(c),
			wr:  bufio.NewWriter(c),
		},
		c: c,
	}
}

func Listen(c Connection, handler Handler) {
	for {
		msg, err := c.Receive()
		if err != nil {
			handler.OnError(err)
		} else if msg != nil {
			if msg.request() {
				handler.OnRequest(msg.Method, msg.ID, msg.Params)
			} else if msg.notification() {
				handler.OnNotification(msg.Method, msg.Params)
			} else if msg.result() {
				handler.OnResult(msg.ID, msg.Result, msg.Error)
			} else {
				id := msg.ID

				if id == "" {
					id = NullRequestID
				}

				c.Send(Message{
					ID:    id,
					Error: &ErrInvalidRequest,
				})
			}
		} else {
			// Connection was closed.
			break
		}
	}
}

type connectionBase struct {
	mtx sync.Mutex
	rd  *bufio.Reader
	wr  *bufio.Writer
}

func (c *connectionBase) Receive() (*Message, error) {
	clen := -1

	hdr, err := c.readHeaderLine()
	if err != nil {
		return nil, err
	}

	for len(hdr) > 0 {
		idx := strings.IndexByte(hdr, ':')
		if idx >= 0 {
			name := strings.TrimSpace(hdr[0:idx])
			val := strings.TrimSpace(hdr[idx+1:])

			if name == contentLengthHeaderName {
				tmp, err := strconv.ParseInt(val, 10, 64)
				if err != nil {
					return nil, fmt.Errorf("invalid '%s' header value '%s'", contentLengthHeaderName, val)
				}

				clen = int(tmp)
			}
		} else {
			return nil, fmt.Errorf("received invalid header '%s'", hdr)
		}
	}

	if clen < 0 {
		return nil, fmt.Errorf("no '%s' header received", contentLengthHeaderName)
	}

	raw := make([]byte, clen)
	_, err = io.ReadFull(c.rd, raw)
	if err != nil {
		return nil, err
	}

	msg := &Message{}

	if _ = json.Unmarshal(raw, msg); err != nil {
		return nil, ErrParseError
	}

	return msg, nil
}

func (c *connectionBase) Send(msg Message) error {
	raw, err := json.Marshal(&msg)
	if err != nil {
		return err
	}

	c.mtx.Lock()

	c.wr.WriteString(fmt.Sprintf("%s: %s\r\n", contentTypeHeaderName, contentTypeHeaderValue))
	c.wr.WriteString(fmt.Sprintf("%s: %d\r\n", contentLengthHeaderName, len(raw)))
	c.wr.WriteString("\r\n")
	c.wr.Write(raw)
	err = c.wr.Flush()

	c.mtx.Unlock()

	return err
}

func (c *connectionBase) readHeaderLine() (string, error) {
	res, err := c.rd.ReadString('\n')
	if err != nil {
		return "", err
	}

	for !strings.HasSuffix(res, "\r\n") {
		txt, err := c.rd.ReadString('\n')
		if err != nil {
			return "", err
		}

		res += txt
	}

	return strings.TrimSuffix(res, "\r\n"), nil
}

type stdioConnection struct {
	connectionBase
	in  *os.File
	out *os.File
}

func (c *stdioConnection) Close() {
	c.in.Close()
}

func (c *stdioConnection) Receive() (*Message, error) {
	msg, err := c.connectionBase.Receive()

	if err != nil && err == os.ErrClosed {
		err = nil
	}

	return msg, err
}

type netConnection struct {
	connectionBase
	c net.Conn
}

func (c *netConnection) Close() {
	c.c.Close()
}

func (c *netConnection) Receive() (*Message, error) {
	msg, err := c.connectionBase.Receive()

	if err != nil && err == net.ErrClosed {
		err = nil
	}

	return msg, err
}

const (
	contentLengthHeaderName string = "Content-Length"
	contentTypeHeaderName   string = "Content-Type"
	contentTypeHeaderValue  string = "application/vscode-jsonrpc; charset=utf-8"
)
