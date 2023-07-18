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

// var errCannotScan = errors.New("can't scan data from connection")

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
	count := 0
	go func() {
		defer tc.wg.Done()
		stdin := stdinScan(tc.inStream)
		// defer close(stdin)
		sb := strings.Builder{}
		// FORCYCLE:
		for {
			// fmt.Println("In send for cycle")
			select {
			case <-tc.ctx.Done():
				// fmt.Println("Send is done by context cancel func")
				return
			case data, ok := <-stdin:
				sb.WriteString(data)
				sb.WriteString("\n")
				// fmt.Println("Before write in socket")
				_, err := tc.conn.Write([]byte(sb.String()))
				// fmt.Println("After write in socket")
				sb.Reset()
				if err != nil {
					return
				}
				if !ok {
					if count > 2 {
						return
					}
					count++
				}
			}
			// fmt.Println("End send for cycle")
		}
		// fmt.Println("Out of send for cycle")
	}()
	return nil
}

func (tc *TnetClient) Receive() error {
	// tc.wg.Add(1)
	// go func() {
	// defer tc.wg.Done()
	scanner := bufio.NewScanner(tc.conn)
	// FORCYCLE:
	// for {
	// 	// // fmt.Println("in recv for cycle")
	// 	if !scanner.Scan() {
	// 		tc.lastError = errCannotScan
	// 		break
	// 	}
	// 	wbytes, err := tc.outStream.Write(append(scanner.Bytes(), '\n'))
	// 	if err != nil {
	// 		break
	// 	}
	// 	if wbytes == 0 {
	// 		break
	// 	}
	// 	// fmt.Println("Out default branch of recv")
	// }
	for scanner.Scan() {
		wbytes, err := tc.outStream.Write(append(scanner.Bytes(), '\n'))
		if err != nil {
			return err
		}
		if wbytes == 0 {
			break
		}
	}
	// fmt.Println("Out of recv for cycle")
	// }()

	// if tc.lastError != nil {
	// 	return tc.lastError
	// }

	return nil
}

func (tc *TnetClient) Close() error {
	for s := range tc.finish {
		if s == os.Interrupt {
			close(tc.finish)
			// fmt.Println("Log: context cancel function call")
			tc.cancel()
			break
		}
	}
	// fmt.Println("Log: w8ing for gorutines")
	tc.wg.Wait()
	// fmt.Println("Log: goruties are closed")
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
		if scanner.Err() != nil {
			// fmt.Println("stdin chan is closed 1")
			close(out)
		}
		if errors.Is(scanner.Err(), io.EOF) {
			// fmt.Println("stdin chan is closed 2")
			close(out)
		}
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
