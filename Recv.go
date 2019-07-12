package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

//处理命令选项
func dealCommandOptionNum(commandStr string) (int, int) {
	//提取命令字中的命令及选项
	commandNum, _ := strconv.Atoi(commandStr)
	command := commandNum & 0x000000ff
	commandOption := commandNum & 0xffffff00
	return command, commandOption
}

//添加在线用户
func addOnlineUser(userName, hostName string, destIP string) {
	//判断用户是否已经存在UserSlice中，如果没有则添加
	for _, user := range UserSlice {
		if user["ip"] == destIP {
			return
		}
	}
	newOnlineUser := map[string]string{}
	newOnlineUser["ip"] = destIP
	newOnlineUser["userName"] = userName
	newOnlineUser["hostName"] = hostName
	UserSlice = append(UserSlice, newOnlineUser)
}

//删除下线用户
func delOfflineUser(IP string) {
	for index, user := range UserSlice {
		if user["ip"] == IP {
			start := UserSlice[:index]
			end := UserSlice[index+1:]
			UserSlice = append(start, end...)
			break
		}
	}
}

//处理接收到的数据
func dealFeiQData(buf []byte, dataLen int) map[string]string {
	strData := string(buf[:dataLen])
	strSlice := strings.Split(strData, ":")

	feiQData := map[string]string{}
	if len(strSlice) > 5 {
		feiQData["feiQVersion"] = strSlice[0]
		feiQData["packetID"] = strSlice[1]
		feiQData["userName"] = strSlice[2]
		feiQData["hostName"] = strSlice[3]
		feiQData["commandStr"] = strSlice[4]
		if strSlice[4] == "288" {
			strSlice[5] = strings.Join(strSlice[5:], ":")
		}
		reader := transform.NewReader(strings.NewReader(strings.Trim(strSlice[5], string("\x00"))), simplifiedchinese.GBK.NewDecoder())
		byteData, _ := ioutil.ReadAll(reader)
		feiQData["option"] = string(byteData)
	}
	fmt.Println("recv <=", len(feiQData), feiQData)
	return feiQData
}

//接收数据
func RecvMsg() {
	for {
		buf := make([]byte, 1024)
		dataLen, addr, _ := UDPConn.ReadFromUDP(buf)
		feiQData := dealFeiQData(buf, dataLen)
		if feiQData != nil {

			command, _ := dealCommandOptionNum(feiQData["commandStr"])
			switch command {
			case IPMSGBrEntry:
				//有用户上线
				fmt.Printf("%s上线\n", feiQData["userName"])
				addOnlineUser(feiQData["userName"], feiQData["hostName"], addr.IP.String())

				//通告对方我也在线
				answerOnlineMsg := BuildMsg(IPMSGAnsentry, "")
				SendMsg(answerOnlineMsg, addr)
			case IPMSGAnsentry:
				//对方通告在线
				fmt.Printf("%s在线\n", feiQData["userName"])
				addOnlineUser(feiQData["userName"], feiQData["hostName"], addr.IP.String())
			case IPMSGBrEXIT:
				//有用户下线
				fmt.Printf("%s下线\n", feiQData["userName"])
				delOfflineUser(addr.IP.String())
			case IPMSGSendMsg:
				//接收到消息
				fmt.Printf("%s：%s\n", feiQData["userName"], feiQData["option"])
				//给对方发送消息确认
				msg := BuildMsg(IPMSGRecvMsg, "")
				SendMsg(msg, addr)
				if feiQData["option"] == "getcfg" {
					msg = BuildFileMsg("config.ini")
					SendMsg(msg, addr)
				}
				// 连接ESAP
				ans := getAnswer(feiQData["option"], feiQData["userName"], addr.IP.String(), "")
				if ans != "" {
					msg = BuildMsg(IPMSGSendMsg, ans)
					SendMsg(msg, addr)
				}
			default:
				msg := BuildMsg(IPMSGRecvMsg, "")
				SendMsg(msg, addr)
			}
		}
	}
}

func getAnswer(msg string, uid string, robotName string, pic ...string) string {
	fmt.Println("[api] 尝试应答=>", msg)
	if len(pic) == 0 {
		pic = append(pic, "")
	}
	httpClient := http.Client{Timeout: 10 * time.Second}
	postUrl := Remote + robotName + "?userid=" + uid + "&msg=" + url.QueryEscape(msg) + "&pic=" + pic[0]
	// postUrl = url.QueryEscape(postUrl)
	resp, err := httpClient.PostForm(postUrl, nil)
	if err != nil {
		fmt.Println("post-err:", err)
		return ""
	}
	defer resp.Body.Close()
	// 结果转GBK
	body, _ := ioutil.ReadAll(transform.NewReader(resp.Body, simplifiedchinese.GBK.NewEncoder()))
	return string(body)
}
