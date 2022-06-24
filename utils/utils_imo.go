package utils

import (
	"net"
)

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "255.255.255.255:53")
	if err != nil {
		return nil
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

// func GetBroadcastAddress() string {
// 	outboundIp := GetOutboundIP()
// 	net := &net.IPNet{
// 		IP:   outboundIp,
// 		Mask: outboundIp.DefaultMask(),
// 	}
// 	a, err := lastAddr(net)
// 	if err != nil {
// 		return ""
// 	}
// 	return a.String()
// }

// func lastAddr(n *net.IPNet) (net.IP, error) { // works when the n is a prefix, otherwise...
// 	if n.IP.To4() == nil {
// 		return net.IP{}, errors.New("does not support IPv6 addresses.")
// 	}
// 	ip := make(net.IP, len(n.IP.To4()))
// 	binary.BigEndian.PutUint32(ip, binary.BigEndian.Uint32(n.IP.To4())|^binary.BigEndian.Uint32(net.IP(n.Mask).To4()))
// 	return ip, nil
// }
