package util

import (
	"crypto/md5"
	"encoding/hex"
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

func Md5(str string) string {
	h := md5.New()
	b := []byte(str)
	h.Write(b)
	return hex.EncodeToString(h.Sum(nil))
}

func Md5StringByNowTime() string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(time.Now().String()))
	return hex.EncodeToString(md5Ctx.Sum(nil))
}

func NowT() string {
	return time.Now().Format("2006-01-02T15:04:05-07:00")
}

func NowDate() string {
	return time.Now().Format("0102")
}

func NowN() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func NowDateName() string {
	return time.Now().Format("2006/01/02")
}

func NowDateTime() string {
	return time.Now().Format("01-02 15:04:05")
}
