package main

import (
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type TnetClient struct {
	conn         net.Conn
	connAddress  string
	retryTimeout time.Duration
	inStream     io.ReadCloser
	outStream    io.Writer
}

func (tc *TnetClient) Connect() error {
	dialer := &net.Dialer{
		Timeout:  tc.retryTimeout,
		Resolver: net.DefaultResolver,
	}
	conn, err := dialer.Dial("tcp", tc.connAddress)
	if err != nil {
		return err
	}
	tc.conn = conn
	return nil
}

func (tc *TnetClient) Send() error {
	if _, err := io.Copy(tc.conn, tc.inStream); err != nil {
		return err
	}
	return nil
}

func (tc *TnetClient) Receive() error {
	if _, err := io.Copy(tc.outStream, tc.conn); err != nil {
		return err
	}
	return nil
}

func (tc *TnetClient) Close() error {
	if tc.conn == nil {
		return nil
	}

	err := tc.conn.Close()
	if err != nil {
		tc.conn = nil
		return err
	}

	return nil
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &TnetClient{
		connAddress:  address,
		retryTimeout: timeout,
		inStream:     in,
		outStream:    out,
	}
}
