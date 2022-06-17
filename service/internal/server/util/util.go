package util

import (
	uuid "github.com/satori/go.uuid"
	"net"
	"service/internal/server/crypto"
	"strings"
)

// GenUUID 生成uuid
func GenUUID() string {
	uuidFunc := uuid.NewV4()
	uuidStr := uuidFunc.String()
	uuidStr = strings.Replace(uuidStr, "-", "", -1)
	uuidByt := []rune(uuidStr)
	return string(uuidByt[8:24])
}

// GenClientId 对称加密IP和端口，当做 clientId
func GenClientId() string {
	raw := []byte(GetIntranetIp() + ":8081")
	str, err := crypto.Encrypt(raw, []byte("Adba723b7fe06819"))
	if err != nil {
		panic(err)
	}

	return str
}

func GenGroupKey(systemId, groupName string) string {
	return systemId + ":" + groupName
}

// 获取本机内网IP
func GetIntranetIp() string {
	addrs, _ := net.InterfaceAddrs()

	for _, addr := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}

		}
	}

	return ""
}
