package imo

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/bcpitutor/ostiki/appconfig"
	"github.com/bcpitutor/ostiki/utils"
)

type IMOListener struct {
	PC        *net.PacketConn
	IMOSender *IMOSender
	cfg       *appconfig.AppConfig
}

func NewIMOListener(cfg *appconfig.AppConfig, is *IMOSender) *IMOListener {
	listenAddr := fmt.Sprintf(":%d", cfg.IMONetPort)
	l := IMOListener{
		PC:        nil,
		IMOSender: is,
		cfg:       cfg,
	}

	if l.PC == nil {
		pc, _ := net.ListenPacket("udp", listenAddr)
		l.PC = &pc
	}

	return &l
}

func (l *IMOListener) Listen() {
	fmt.Printf("Starting to listen UDP on %d\n", l.cfg.IMONetPort)
	for {
		buf := make([]byte, 1024)
		pc := *(l.PC)

		n, addr, err := pc.ReadFrom(buf)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			continue
		}

		msg := string(buf[:n])
		msg = strings.TrimSpace(msg)

		l.checkMessage(addr, msg)
	}
}

func (l *IMOListener) checkMessage(addr net.Addr, msg string) {
	if len(msg) < 3 {
		log.Printf("Invalid message: %s\n", msg)
		return
	}

	remoteIP := strings.Split(addr.String(), ":")[0]
	myIP := utils.GetOutboundIP()

	msgType := msg[0:3]

	switch msgType {
	case "000":
		if remoteIP == myIP.String() {
			return
		}

		parts := strings.Split(addr.String(), ":")

		l.IMOSender.DirectMessage(
			fmt.Sprintf("%s:%d", parts[0], 8671),
			"001",
		)
	case "001":
		if remoteIP == myIP.String() {
			return
		}

		fmt.Printf("Join acknowledgement from %s\n", addr.String())
	case "100":
		analyze_message(msg)
	default:
		if remoteIP == myIP.String() {
			return
		}

		fmt.Printf("Unknown message: %s\n", msgType)
	}
}

func analyze_message(msg string) {
	rest := ""
	if len(msg) > 3 {
		rest = msg[4:]
	}

	fmt.Printf("Update message received for [%s]\n", rest)
	switch rest {
	case "domain":
		fmt.Printf("Domain update received\n")
	case "group":
		fmt.Printf("Group update received\n")
	case "ticket":
		fmt.Printf("Ticket update received\n")
	default:
		fmt.Printf("Unknown update received: %s\n", rest)
	}
}
