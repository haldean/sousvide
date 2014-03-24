package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

var bcastEnabled = flag.Bool("broadcast", true, "enable broadcast")
var bcastPort = flag.Int("broadcast_port", 9786, "broadcast port")

var SV_BCAST_MSG = []byte(fmt.Sprintf(
	"NOTIFY * HTTP/1.1\r\nHost: 255.255.255.255:%d\r\nX-263A-Capabilities: sousvide\r\n", *bcastPort))
var SV_BCAST_FIELD_STREAM = "X-263A-Stream: %v:%d"

func StartBroadcast() {
	if (!*bcastEnabled) {
		log.Printf("broadcast disabled")
		return
	}
	log.Printf("broadcast enabled, starting broadcast loop")

	sock, err := net.DialUDP("udp4", nil, &net.UDPAddr{
		IP: net.IPv4bcast, Port: *bcastPort,
	})
	if err != nil {
		log.Panicf("error creating broadcast port: %v", err)
	}
	defer sock.Close()

	var addr string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Panicf("error getting interface addresses: %v", err)
	}
	for _, a := range addrs {
		ipAddr := a.String()
		ipAddr = ipAddr[0:strings.Index(ipAddr, "/")]
		ip := net.ParseIP(ipAddr)
		if ip != nil && !ip.IsLoopback() && ip.To4() != nil {
			log.Printf("broadcasting address %v", ipAddr)
			addr = ipAddr
		}
	}

	field := []byte(fmt.Sprintf(SV_BCAST_FIELD_STREAM, addr, *SockPort))
	msg := make([]byte, len(SV_BCAST_MSG) + len(field))
	copy(msg, SV_BCAST_MSG)
	for i := 0; i < len(field); i++ {
		msg[i + len(SV_BCAST_MSG)] = field[i]
	}

	msg_len := len(msg)
	tick := time.Tick(5 * time.Second)
	for _ = range tick {
		n, err := sock.Write(msg)
		if err != nil {
			log.Printf("warning, could not send broadcast: %v", err)
			continue
		}
		if n != msg_len {
			log.Printf("warning, sent %d bytes of %d in broadcast", n, msg_len)
			continue
		}
	}
}
