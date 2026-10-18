package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		log.Fatalf("could not open file with error: %v", err)
	}
	defer file.Close()
	buffer := make([]byte, 8)
	for {
		bytesRead, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error reading buffer: %v", err)
		}
		fmt.Printf("read: %v\n", string(buffer[:bytesRead]))
	}
}
