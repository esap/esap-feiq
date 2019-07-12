package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

//显示所有在线用户
func DisplayOnlineUser() {
	for index, user := range UserSlice {
		fmt.Printf("%d %s\n", index, user["userName"])
	}
	fmt.Println("userlist:", UserSlice)
}

//组装飞鸽传书的数据包
func BuildMsg(command int, optionData string) string {
	msg := FeiQVersion + ":" + strconv.FormatInt(time.Now().Unix(), 10) +
		":" + FeiQUserName + ":" + FeiQHostName + ":" + strconv.Itoa(command) + ":" + optionData
	return msg
}

//构建文件消息
func BuildFileMsg(fileName string) string {
	//文件序号:文件名:文件大小:修改时间:文件的属性
	fileInfo, err := os.Stat(fileName)
	if err != nil {
		panic("文件不存在")
	}

	fileNo := strconv.Itoa(0)
	fileSize := strconv.FormatInt(fileInfo.Size(), 16)
	fileCTime := strconv.FormatInt(fileInfo.ModTime().Unix(), 16)
	fileType := strconv.FormatInt(IPMSGFileRegular, 16)
	buildFileMsg := []string{fileNo, fileName, fileSize, fileCTime, fileType}
	optionStr := strings.Join(buildFileMsg, ":")
	fileStr := string("\x00") + optionStr + ":"
	commandNum := IPMSGSendMsg | IPMSGFileAttachOpt
	fmt.Println("filestr:", fileStr)
	return BuildMsg(commandNum, fileStr)
}

//发送消息
func SendMsg(msg string, DestIP *net.UDPAddr) {
	fmt.Println(msg, "=>", DestIP.IP)
	UDPConn.WriteToUDP([]byte(msg), DestIP)
}

//发送广播消息
func SendBroadcast(cmd int) {
	msg := BuildMsg(cmd, FeiQUserName)

	BroadCastIP := net.UDPAddr{
		IP:   net.IPv4(255, 255, 255, 255),
		Port: 2425,
	}
	SendMsg(msg, &BroadCastIP)
}
