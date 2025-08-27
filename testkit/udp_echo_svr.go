package main

import (
	"fmt"
	"net"
	"os"
)

// running in s3.mojotv.cn
func main() {
	addr, err := net.ResolveUDPAddr("udp", ":8080")
	if err != nil {
		fmt.Println("Error resolving address:", err)
		os.Exit(1)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Error listening:", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("UDP echo server listening on :8080")

	buffer := make([]byte, 1024)
	for {
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading:", err)
			continue
		}

		fmt.Printf("Received %d bytes from %s: %s\n", n, clientAddr, string(buffer[:n]))

		_, err = conn.WriteToUDP(buffer[:n], clientAddr)
		if err != nil {
			fmt.Println("Error writing:", err)
		}
	}
}
