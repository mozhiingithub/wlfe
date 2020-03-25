package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	ip, _ := getIP(false)
	fmt.Println(ip)
	p, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	http.Handle("/", http.FileServer(http.Dir(p)))
	http.ListenAndServe(":8000", nil)
}

//获取第一个非回环地址。参数ifIP4为是否返回ipv4。
func getIP(ifIP4 bool) (ip string, e error) {
	var addrs []net.Addr
	addrs, e = net.InterfaceAddrs()
	if nil != e {
		return
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if (ifIP4 && nil != ipnet.IP.To4()) || (!ifIP4 && nil == ipnet.IP.To4()) {
				ip = ipnet.IP.String()
				return
			}
		}
	}
	return
}
