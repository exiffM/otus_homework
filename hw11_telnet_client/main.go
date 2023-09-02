package main

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	args     []string
	duration time.Duration
)

func init() {
	flag.DurationVar(&duration, "timeout", 10*time.Second, "use --timeout=<value> in seconds")
}

// func stdinScan(in io.ReadCloser) chan string {
// 	out := make(chan string)
// 	go func() {
// 		scanner := bufio.NewScanner(in)
// 		for scanner.Scan() {
// 			out <- scanner.Text()
// 		}
// 		close(out)
// 	}()
// 	return out
// }

var errInvalidArgs = errors.New("invalid arguments")

func checkParams() error {
	flag.Parse()
	args = flag.Args()
	if len(args) != 2 {
		return errInvalidArgs
	}
	return nil
}

func main() {
	if err := checkParams(); err != nil {
		fmt.Println("invalid usage of go-telnet\nuse command go-telnet [--timeout=x] <hostname/ipadress> <port>")
		return
	}

	sb := strings.Builder{}
	sb.WriteString(args[0]) // link/ip-address
	sb.WriteString(":")
	sb.WriteString(args[1]) // port

	inBuff := bytes.Buffer{}

	client := NewTelnetClient(sb.String(), duration, io.NopCloser(&inBuff), os.Stdout)
	defer client.Close()

	err := client.Connect()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	ctx, cancel := context.WithCancel(context.Background())

	// Receiver goroutine
	Receiver := func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				err := client.Receive()
				if err != nil {
					return
				}
			}
		}
	}
	go Receiver()

	count := 0
	// Sender goroutine
	Sender := func() {
		defer func() {
			wg.Done()
			cancel()
			client.Close()
		}()
		stdinBuff := bufio.NewScanner(os.Stdout)

		for {
			// select {
			// case <-ctx.Done():
			// 	fmt.Println("Send done by context")
			// 	return
			// case _, ok := <-inChan:
			// 	fmt.Println("Got something from stdin to send")
			// 	if !ok {
			// 		fmt.Println("error in stdin chan(mb closed)")
			// 		return
			// 	}
			// 	err := client.Send()
			// 	fmt.Println("Send something from stdin")
			// 	if err != nil {
			// 		fmt.Println("Send error: " + err.Error())
			// 		if count >= 3 {
			// 			fmt.Println("Send attempts are over")
			// 			return
			// 		}
			// 		count++
			// 	}
			// }
			ok := stdinBuff.Scan()
			if !ok {
				return
			}

			_, err = inBuff.WriteString(stdinBuff.Text() + "\n")
			if err != nil {
				return
			}

			err = client.Send()
			if err != nil {
				if count >= 1 {
					return
				}
				count++
			}
		}
	}
	Sender()

	wg.Wait()
}
