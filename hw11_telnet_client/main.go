package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

var (
	args     []string
	duration time.Duration
)

func init() {
	flag.DurationVar(&duration, "timeout", 10*time.Second, "use --timeout=<value> in seconds")
}

func main() {
	flag.Parse()
	args = flag.Args()
	if len(args) < 2 {
		fmt.Println("invalid usage of go-telnet\nuse command go-telnet [--timeout=x] <hostname/ipadress> <port>")
		return
	}
	sb := strings.Builder{}
	sb.WriteString(args[0]) // link/ip-address
	sb.WriteString(":")
	sb.WriteString(args[1]) // port
	client := NewTelnetClient(sb.String(), duration, os.Stdin, os.Stdout)
	err := client.Connect()
	// defer client.Close()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	client.Send()
	client.Receive()
	client.Close()
}
