package main

import (
	"bufio"
	"context"
	"errors"
	"io"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"
)

var errCannotScan = errors.New("can't scan data from connection")

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
	FinishChan() chan os.Signal
}

type TnetClient struct {
	ctx          context.Context
	conn         net.Conn
	connAddress  string
	retryTimeout time.Duration
	inStream     io.ReadCloser
	outStream    io.Writer
	cancel       context.CancelFunc
	lastError    error
	finish       chan os.Signal
	wg           sync.WaitGroup
}

func (tc *TnetClient) Connect() error {
	signal.Notify(tc.finish, os.Interrupt)
	dialer := &net.Dialer{
		Timeout:  tc.retryTimeout,
		Resolver: net.DefaultResolver,
	}
	tc.ctx, tc.cancel = context.WithCancel(context.Background())
	tc.conn, tc.lastError = dialer.DialContext(tc.ctx, "tcp", tc.connAddress)
	if tc.lastError != nil {
		return tc.lastError
	}
	return nil
}

func (tc *TnetClient) Send() error {
	tc.wg.Add(1)
	go func() {
		defer tc.wg.Done()
		stdin := stdinScan(tc.inStream)
		// scanner := bufio.NewScanner(tc.inStream)
		sb := strings.Builder{}
		count := 0
	FORLOOP:
		for {
			select {
			case <-tc.ctx.Done():
				break FORLOOP
			case data, ok := <-stdin:
				if !ok {
					// // fmt.Println("got eof")
					tc.finish <- os.Interrupt
					return
				}
				sb.WriteString(data)
				sb.WriteString("\n")
				_, err := tc.conn.Write([]byte(sb.String()))
				sb.Reset()
				if err != nil {
					if count >= 1 {
						tc.finish <- os.Interrupt
						break FORLOOP
					}
					count++
				}
			}
		}
	}()
	return nil
}

func (tc *TnetClient) Receive() error {
	// tc.wg.Add(1)
	// +- ok
	go func() {
		// defer tc.wg.Done()
		scanner := bufio.NewScanner(tc.conn)
	FOR:
		for /*scanner.Scan()*/ {
			select {
			case <-tc.ctx.Done():
				break FOR
			default:
				//fmt.Println("Start cycle")
				if !scanner.Scan() {
					tc.lastError = errCannotScan
					break FOR
				}
				wbytes, err := tc.outStream.Write(append(scanner.Bytes(), '\n'))
				if err != nil {
					//return err
					break FOR
				}
				if wbytes == 0 {
					break FOR
				}
				// fmt.Println("Start cycle")
			} // - select
		}
	}()
	// +- ok
	return nil
}

func (tc *TnetClient) Close() error {
	for s := range tc.finish {
		if s == os.Interrupt {
			close(tc.finish)
			// // fmt.Println("Log: context cancel function call")
			tc.cancel()
			break
		}
	}
	// // fmt.Println("Log: w8ing for gorutines")
	tc.wg.Wait()
	// // fmt.Println("Log: goruties are closed")
	tc.conn.Close()
	return nil
}

func (tc *TnetClient) FinishChan() chan os.Signal {
	return tc.finish
}

func stdinScan(in io.ReadCloser) chan string {
	out := make(chan string)
	go func() {
		scanner := bufio.NewScanner(in)
		for scanner.Scan() {
			out <- scanner.Text()
		}
		close(out)
	}()
	return out
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &TnetClient{
		connAddress:  address,
		retryTimeout: timeout,
		inStream:     in,
		outStream:    out,
		finish:       make(chan os.Signal, 1),
	}
}
