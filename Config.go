package main

import (
	"errors"
	"net"

	"github.com/larspensjo/config"
)

const (
	FeiQPort    = "2425"
	FeiQVersion = "1"
	// FeiQVersion  = "1_lbt6_0#128#305A3A52C610#0#0#0#4001#9"
	FeiQUserName = "robot"
	FeiQHostName = "ESAP"
	Broadcast    = "255.255.255.255"

	IPMSGBrEntry       = 0x00000001 //上线提醒消息命令
	IPMSGBrEXIT        = 0x00000002 //下线提醒消息命令
	IPMSGSendMsg       = 0x00000020 //表示发送消息
	IPMSGAnsentry      = 0x00000003 //对方也在线
	IPMSGRecvMsg       = 0x00000021 //确认收到消息
	IPMSGFileAttachOpt = 0x00200000
	IPMSGFileRegular   = 0x00000001 //普通文件
)

var (
	//保存在线用户列表
	UserSlice []map[string]string

	UDPConn *net.UDPConn = nil

	Port   = "19090"
	Local  = "192.168.99.10"
	Remote = "http://192.168.99.10:9090/robot/"
)

func GetConfig(sec string) (map[string]string, error) {
	targetConfig := make(map[string]string)
	cfg, err := config.ReadDefault("config.ini")
	if err != nil {
		return targetConfig, err
	}
	sections := cfg.Sections()
	if len(sections) == 0 {
		return targetConfig, errors.New("no " + sec + " config")
	}
	for _, section := range sections {
		if section != sec {
			continue
		}
		sectionData, _ := cfg.SectionOptions(section)
		for _, key := range sectionData {
			value, err := cfg.String(section, key)
			if err == nil {
				targetConfig[key] = value
			}
		}
		break
	}
	return targetConfig, nil
}
