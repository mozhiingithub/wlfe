package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	qrcode "github.com/skip2/go-qrcode"
)

//设置flag参数
var (
	f   = flag.String("f", "", "文件名。")
	p   = flag.String("p", "8000", "端口号。默认为8000。")
	ipv = flag.Bool("ipv6", false, "是否使用ipv6地址。默认为false。")
)

func main() {

	//解析flag
	flag.Parse()

	var (
		e      error
		ip     string
		url    string
		path   string
		imgDir string
	)

	//获取本机IP
	ip, e = getIP(*ipv)
	ifError(e)

	//初始化文件/文件夹地址
	url = "http://" + ip + ":" + *p
	if "" != *f { //文件名非空
		url += "/" + *f
	}

	//按给定端口，初始化server
	server := &http.Server{
		Addr:    ":" + *p,
		Handler: http.DefaultServeMux,
	}

	//获取当前目录
	path, e = filepath.Abs(filepath.Dir(os.Args[0]))
	ifError(e)

	//注册一个文件服务handle
	http.Handle("/", http.FileServer(http.Dir(path)))

	//生成二维码
	imgDir = strconv.Itoa(int(time.Now().UnixNano())) + ".png" //图片名
	qrcode.WriteFile(url, qrcode.Medium, 256, imgDir)

	//开启文件服务
	go server.ListenAndServe()

	//弹出二维码图片，并等待手动关闭图片
	cmd := exec.Command("eog", imgDir)
	cmd.Start()
	cmd.Wait()

	//关闭文件服务
	server.Close()

	//删除二维码图片
	e = os.Remove(imgDir)
	ifError(e)

}

//获取第一个非回环地址。参数ifIP6为是否返回ipv6。
func getIP(ifIP6 bool) (ip string, e error) {
	var addrs []net.Addr
	addrs, e = net.InterfaceAddrs()
	if nil != e {
		return
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if (ifIP6 && nil == ipnet.IP.To4()) || (!ifIP6 && nil != ipnet.IP.To4()) {
				ip = ipnet.IP.String()
				return
			}
		}
	}
	return
}

func ifError(e error) {
	if nil != e {
		log.Println(e)
		os.Exit(0)
	}
}
