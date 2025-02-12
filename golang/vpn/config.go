package vpn

import (
	configuration "CFScanner/configuration"
	"CFScanner/utils"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var logger = []string{"debug",
	"info",
	"warning",
	"error",
	"none"}

func loglevel(level string) string {
	for _, log := range logger {
		if strings.ToLower(level) == log {
			return level
		}
	}
	return logger[4]
}

var XRayTemplate = `{
  "log": {
    "loglevel": "LOG"
  },
  "inbounds": [{
    "port": PORTPORT,
    "listen": "127.0.0.1",
    "tag": "socks-inbound",
    "protocol": "socks",
    "settings": {
      "auth": "noauth",
      "udp": false,
      "ip": "127.0.0.1"
    },
    "sniffing": {
      "enabled": true,
      "destOverride": ["http", "tls"]
    }
  }],
  "outbounds": [
    {
		"protocol": "vmess",
    "settings": {
      "vnext": [{
        "address": "IP.IP.IP.IP",
        "port": CFPORT,
        "users": [{"id": "IDID" }]
      }]
    },
		"streamSettings": {
        "network": "ws",
        "security": "tls",
        "wsSettings": {
            "headers": {
                "Host": "HOSTHOST"
            },
            "path": "ENDPOINTENDPOINT"
        },
        "tlsSettings": {
            "serverName": "RANDOMHOST",
            "allowInsecure": false
        }
    }
	}],
  "other": {}
}`

// XRayConfig create VPN configuration
func XRayConfig(IP string, testConfig *configuration.Configuration) string {
	localPortStr := strconv.Itoa(utils.GetFreePort())
	config := strings.ReplaceAll(XRayTemplate, "PORTPORT", localPortStr)
	config = strings.ReplaceAll(config, "LOG", loglevel(testConfig.LogLevel))
	config = strings.ReplaceAll(config, "IP.IP.IP.IP", IP)
	config = strings.ReplaceAll(config, "CFPORT", testConfig.Config.AddressPort)
	config = strings.ReplaceAll(config, "IDID", testConfig.Config.UserId)
	config = strings.ReplaceAll(config, "HOSTHOST", testConfig.Config.WsHeaderHost)
	config = strings.ReplaceAll(config, "ENDPOINTENDPOINT", testConfig.Config.WsHeaderPath)
	config = strings.ReplaceAll(config, "RANDOMHOST", testConfig.Config.Sni)

	configPath := fmt.Sprintf("%s/config-%s.json", configuration.DIR, IP)
	configFile, err := os.Create(configPath)
	if err != nil {
		log.Fatal(err)
	}
	defer func(configFile *os.File) {
		err := configFile.Close()
		if err != nil {

		}
	}(configFile)

	_, configErr := configFile.WriteString(config)
	if configErr != nil {
		return ""
	}

	return configPath
}
