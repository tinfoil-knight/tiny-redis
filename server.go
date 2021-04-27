package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/tinfoil-knight/tiny-redis/commands"
	"github.com/tinfoil-knight/tiny-redis/resp"
	"github.com/tinfoil-knight/tiny-redis/store"
)

var (
	LF = []byte("\n")
	SP = []byte(" ")
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

func handleConn(kv *store.Store, c net.Conn) {
	defer c.Close()
	scanner := bufio.NewScanner(c)
	scanner.Split(untilEOF)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	byts := scanner.Bytes()
	if !(len(byts) == 0) {
		fmt.Printf("Recv: %+q\n", byts)
		var val interface{}
		if byts[0] == '*' {
			// resp
			val, _ = resp.Decode(byts)
		} else {
			// inline command format
			val = bytes.Split(bytes.TrimSuffix(byts, LF), SP)
		}
		fmt.Printf("Parsed: %s\n", val)
		_, ok := val.([]interface{})
		if !(ok) {
			v, err := commands.ExecuteCommand(kv, val)
			fmt.Printf("Send: %+q\n", resp.Encode(v))
			if err != nil {
				c.Write([]byte(resp.Encode(err)))
			} else {
				c.Write([]byte(resp.Encode(v)))
			}
		}
	}
}

func main() {
	port := flag.Int("p", 8001, "sets tcp port")
	flag.Parse()
	host := "[::]"
	address := fmt.Sprintf("%s:%d", host, *port)
	l, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Listening at: %s\n", l.Addr())
	defer l.Close()
	kv := store.New()
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(kv, conn)
	}
}
