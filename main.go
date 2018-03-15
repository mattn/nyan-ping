package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func main() {
	var sleep time.Duration
	var timeout time.Duration
	var ipv6 bool

	flag.DurationVar(&sleep, "s", 1*time.Millisecond, "sleep")
	flag.DurationVar(&timeout, "t", 100*time.Millisecond, "timeout")
	flag.BoolVar(&ipv6, "6", false, "ip6")
	flag.Parse()

	proto := "ip4"
	host := "pong4.kooshin.net"
	if ipv6 {
		proto = "ip6"
		host = "pong6.kooshin.net"
	}

	ip, err := net.ResolveIPAddr(proto, host)
	if err != nil {
		log.Fatalf("ResolveIPAddr: %v", err)
	}

	c, err := icmp.ListenPacket(proto+":icmp", "0.0.0.0")
	if err != nil {
		log.Fatalf("ListenPacket: %v", err)
	}
	defer c.Close()

	for i := 1; i <= 1400; i++ {
		wm := icmp.Message{
			Type: ipv4.ICMPTypeEcho,
			Code: 0,
			Body: &icmp.Echo{
				ID: os.Getpid() & 0xffff, Seq: i,
				Data: []byte("HELLO-R-U-THERE"),
			},
		}
		wb, err := wm.Marshal(nil)
		if err != nil {
			log.Fatalf("Marshal: %v", err)
		}
		if _, err := c.WriteTo(wb, &net.IPAddr{IP: ip.IP}); err != nil {
			log.Fatalf("WriteTo: %v", err)
		}

		c.SetReadDeadline(time.Now().Add(timeout))
		rb := make([]byte, 1500)
		n, _, err := c.ReadFrom(rb)
		if err != nil {
			fmt.Print("U")
		} else {
			rm, err := icmp.ParseMessage(ipv4.ICMPTypeEcho.Protocol(), rb[:n])
			if err == nil && rm.Type == ipv4.ICMPTypeEchoReply {
				fmt.Print("!")
			} else {
				fmt.Print("U")
			}
		}
		if i%70 == 0 {
			fmt.Println()
		}
		time.Sleep(sleep)
	}
}
