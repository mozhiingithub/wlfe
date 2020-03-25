package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func main() {
	ip, _ := getIP(true)
	fmt.Println(ip)
	path, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	server := &http.Server{
		Addr:    ":8000",
		Handler: http.DefaultServeMux,
	}
	http.Handle("/", http.FileServer(http.Dir(path)))
	/*
		c := make(chan int)
		go func() {
			<-c
			srv.Close()
		}()
	*/
	go server.ListenAndServe()
	time.Sleep(10 * time.Second)
	fmt.Println(123)
	// c <- 1
	server.Close()
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
