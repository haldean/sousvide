package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
)

var SockPort = flag.Int("sockport", 7897, "port used for streaming temperature data")

func StartSockServer() chan HistorySample {
	stream := make(chan HistorySample, 1)
	ready := make(chan bool)
	toSock := make(chan []byte, 1)
	go dispatchStream(stream, ready, toSock)
	go listenSock(ready, toSock)
	return stream
}

func dispatchStream(
		stream chan HistorySample, ready chan bool, toSock chan []byte) {
	sockReady := false
	for {
		select {
		case sample := <-stream:
			if (!sockReady) {
				continue
			}
			msg, err := json.Marshal(sample)
			if err != nil {
				log.Printf("warning: could not stream snapshot: %v", err)
				continue
			}
			toSock <- msg

		case sockReady = <-ready:
			// all good
		}
	}
}

func listenSock(ready chan bool, toSock chan []byte) {
	ss, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *SockPort))
	if err != nil {
		log.Fatalf("error: could not listen on %d: %v", *SockPort, err)
	}
	for {
		log.Printf("accepting new client on :%d", *SockPort)
		sock, err := ss.Accept()
		if err != nil {
			log.Printf("warning: could not accept on socket: %v", err)
		}
		log.Printf("sending snapshots to client %v", sock.RemoteAddr())
		ready <- true
		for {
			msg := <-toSock
			_, err := sock.Write(msg)
			if err != nil {
				log.Printf("warning: failed to write to socket: %v", err)
				break
			}
			_, err = sock.Write([]byte("\n\x00"))
		}
		ready <- false
		sock.Close()
	}
}
