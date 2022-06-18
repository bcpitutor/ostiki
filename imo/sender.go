package imo

import (
	"fmt"
	"log"
	"net"

	"github.com/bcpitutor/ostiki/appconfig"
	"github.com/bcpitutor/ostiki/utils"
)

type IMOSender struct {
	broadcastAddress string
	broadcastPort    int
	myIP             net.IP
}

func NewIMOSender(cfg *appconfig.AppConfig) *IMOSender {
	ims := IMOSender{
		broadcastAddress: utils.GetBroadcastAddress(),
		broadcastPort:    cfg.IMONetPort,
		myIP:             utils.GetOutboundIP(),
	}
	return &ims
}

func (ims *IMOSender) BroadcastMessage(msg string) {
	addr := fmt.Sprintf("%s:%d", utils.GetBroadcastAddress(), ims.broadcastPort)
	conn, _ := net.Dial("udp", addr)

	_, err := conn.Write([]byte(msg))
	if err != nil {
		log.Printf("BroadcastMessage.Error: %s\n", err)
	}
	defer conn.Close()
}

func (ims *IMOSender) DirectMessage(addr string, msg string) {
	conn, _ := net.Dial("udp", addr)

	_, err := conn.Write([]byte(msg))
	if err != nil {
		log.Printf("DirectMessage.Error: %s\n", err)
	}
	defer conn.Close()
}

func (ims *IMOSender) JoinCluster() {
	ims.BroadcastMessage("000")
}
