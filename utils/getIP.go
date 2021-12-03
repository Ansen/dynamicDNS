package utils

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// GetPublicIP 本机外网IP
func GetPublicIP(url string) string {
	IPAddr := ""
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

	return string(s)
}
