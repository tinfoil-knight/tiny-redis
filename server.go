package main

import (
	"bufio"
	"fmt"
	"log"
	"net"

	"github.com/tinfoil-knight/tiny-redis/resp"
)

func untilEOF(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF {
		if len(data) == 0 {
			return 0, nil, nil
		}
		return 0, []byte{}, nil
	}
	return len(data), data, nil
}

func handleConn(c net.Conn) {
	defer c.Close()
	scanner := bufio.NewScanner(c)
	scanner.Split(untilEOF)
	scanner.Scan()
	err := scanner.Err()
	if err != nil {
		panic(err)
	}
	bytes := scanner.Bytes()
	fmt.Printf("%+q\n", bytes)
	val, _ := resp.Decode(bytes)
	fmt.Printf("Parsed: %v\n", val)
	c.Write([]byte("-ERR unknown command\r\n"))
}

func main() {
	listener, err := net.Listen("tcp", "localhost:8001")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		handleConn(conn)
	}
}
