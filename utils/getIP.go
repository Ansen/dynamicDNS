package utils

import (
	"dynamicDNS/config"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

// GetPublicIP 本机外网IP
func GetPublicIP() string {
	IPAddr := ""
	url := config.Conf.IPApi
	if url == "" {
		url = "http://ip.3322.net/"
	}

	client := http.Client{
		Timeout: 3 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		log.Println(err)
		return IPAddr
	}

	defer resp.Body.Close()

	s, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return IPAddr
	}

	IPAddr = string(s)
	IPAddr = strings.Trim(IPAddr, "\n")
	IPAddr = strings.Trim(IPAddr, " ")
	return IPAddr
}
