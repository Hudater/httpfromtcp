package main

import (
	"fmt"
	"net"

	"github.com/Hudater/httpfromtcp/internal/request"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		fmt.Printf("could not listen to connection: %v", err)
	}
	defer listener.Close()

	// wait for conn accept
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("failed to accept connection: %v", err)
		}
		fmt.Println("connection has been accepted")
		readRequest, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Printf("unable to get request from reader: %v", err)
		}
		fmt.Printf("Request line:\n - Method: %v \n - Target: %v \n - Version: %v\n", readRequest.RequestLine.Method, readRequest.RequestLine.RequestTarget, readRequest.RequestLine.HttpVersion)
	}
}
