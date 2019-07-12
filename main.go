package main

import (
	"fmt"
	"net"
	"net/http"
)

func init() {
	cfg, err := GetConfig("esap")
	if err != nil {
		fmt.Println("未找到配置文件")
	} else {
		Local = cfg["local"]
		Remote = cfg["remote"]
		Port = cfg["port"]
	}
}

func main() {
	//创建udp套接字
	addr, _ := net.ResolveUDPAddr("udp", "0.0.0.0:2425")
	UDPConn, _ = net.ListenUDP("udp", addr)
	defer UDPConn.Close()

	//循环接收数据
	go RecvMsg()

	//广播上线
	SendBroadcast(IPMSGBrEntry)

	//显示所有在线用户
	DisplayOnlineUser()

	http.HandleFunc("/p", func(w http.ResponseWriter, req *http.Request) {
		pi := req.FormValue("id")
		picfile := pi + ".jpg"
		http.ServeFile(w, req, picfile)
		return
	})
	http.ListenAndServe(fmt.Sprint(":", Port), nil)
}
