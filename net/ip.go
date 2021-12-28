package net

import (
	"errors"
	"net"
)

//获取本机内网ip
func GetLocalIp() string {
	addr, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}
	for _, add := range addr {
		if ipnet, ok := add.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "127.0.0.1"
}

//获取本地ip

func GetIntranetIp() (string, error) {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		return "", err
	}

	for _, address := range addrs {

		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}

		}
	}

	return "", errors.New("ip empty")
}
