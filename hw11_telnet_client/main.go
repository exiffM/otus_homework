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

	// inChan := stdinScan(os.Stdin)

	// Receiver goroutine
	go func() {
		errChan := make(chan error)
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case errChan <- client.Receive():
				err := <-errChan
				if err != nil {
					fmt.Println(err.Error())
					return
				}
			}
		}
	}()

	count := 0
	// Sender goroutine
	go func() {
		defer func() {
			wg.Done()
			cancel()
		}()

		stdinBuff := bufio.NewReader(os.Stdin)

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
			data, err := stdinBuff.ReadString('\n')
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			_, err = inBuff.WriteString(data)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			err = client.Send()
			if err != nil {
				if count >= 1 {
					fmt.Println("Resend attempts are over")
					return
				}
				count++
			}
		}
	}()

	wg.Wait()
}
