package util

import (
	"log"
	"net"
)

var localIp net.IP

func init() {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer conn.Close()

	localIp = conn.LocalAddr().(*net.UDPAddr).IP
}

func LocalAddr() net.IP {
	return localIp
}
