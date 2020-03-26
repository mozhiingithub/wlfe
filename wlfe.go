package main

import (
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

func main() {

	var (
		ipv    bool   //是否采用ipv6
		p      string //端口号
		e      error  //错误变量
		ip     string //本机局域网ip地址
		url    string //服务开启后，文件或文件夹对应的url
		path   string //程序运行时所在目录
		imgDir string //生成二维码图片文件地址
	)

	//从环境变量获知是否采用ipv6
	if "true" == os.Getenv("WLFE_IPV6") {
		ipv = true
	} else {
		ipv = false
	}

	//获取本机IP
	ip, e = getIP(ipv)
	ifError(e)

	//获取端口号
	p = os.Getenv("WLFE_PORT")
	if "" == p {
		p = "8000"
	}

	//初始化文件/文件夹地址
	url = "http://" + ip + ":" + p
	if len(os.Args) > 1 { //文件名非空
		url += "/" + os.Args[1]
	}

	//按给定端口，初始化server
	server := &http.Server{
		Addr:    ":" + p,
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

//获取第一个非回环地址。参数ifIP6为是否返回ipv6
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

//判断错误是否非空，若非空，则输出错误内容并终止程序
func ifError(e error) {
	if nil != e {
		log.Println(e)
		os.Exit(0)
	}
}
