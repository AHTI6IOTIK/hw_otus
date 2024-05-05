package telnet

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type Client interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type client struct {
	address  string
	timeout  time.Duration
	in       io.ReadCloser
	out      io.Writer
	conn     net.Conn
	isClosed bool
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) Client {
	return &client{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (c *client) Connect() error {
	conn, err := net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrConnectionFailEstablished, err) //nolint:errorlint,nolintlint
	}
	log.Println("connected to", c.address)

	c.conn = conn

	return nil
}

func (c *client) Close() error {
	if c.isClosed {
		return nil
	}

	if c.conn == nil {
		return ErrConnectionNotEstablished
	}

	defer func() {
		c.isClosed = true
		c.conn = nil
	}()

	err := c.conn.Close()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrConnectionClose, err) //nolint:errorlint,nolintlint
	}

	return nil
}

func (c *client) Send() error {
	scanner := bufio.NewScanner(c.in)

	for scanner.Scan() {
		_, err := c.conn.Write(append(scanner.Bytes(), byte('\n')))
		if err != nil {
			return err
		}
	}

	if scanner.Err() != nil {
		return scanner.Err()
	}

	return nil
}

func (c *client) Receive() error {
	scanner := bufio.NewScanner(c.conn)

	for scanner.Scan() {
		bbb := append(scanner.Bytes(), byte('\n'))
		_, err := c.out.Write(bbb)
		if err != nil {
			return err
		}
	}

	if scanner.Err() != nil {
		return scanner.Err()
	}

	return nil
}
