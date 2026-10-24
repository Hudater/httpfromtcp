package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		fmt.Printf("could not resolve address: %v\n", err)
	}

	dialer, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Printf("could not create dialer: %v\n", err)
	}
	defer dialer.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println(">")
		readString, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("failed to read input: %v\n", err)
		}
		lineWriter, err := dialer.Write([]byte(readString))
		if err != nil {
			fmt.Printf("failed to line write: %v\n", err)
		}
		fmt.Printf("Number of bytes written: %v\n", lineWriter)
	}
}
