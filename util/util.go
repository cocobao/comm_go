package util

import (
	"fmt"
	"net"
	"runtime"
	"time"

	"github.com/cocobao/log"
)

func PrintStack() {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	fmt.Printf("[%s] - %s", time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05"), buf[:n])
}

//获取本机ip地址
func GetLocalIPAddr() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Warn("get local ip fail", err)
		return ""
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}