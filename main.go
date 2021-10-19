package main

import (
	"dynamicDNS/aliyun"
	"dynamicDNS/config"
	"dynamicDNS/utils"
	"log"
	"time"
)

func init() {
	config.LoadConfig()
}

func main() {
	tick := time.NewTicker(time.Duration(config.Conf.Interval) * time.Second)
	for {
		publicIP := utils.GetPublicIP()
		if publicIP == "" {
			log.Print("unable get public ip")
			continue
		}
		ali := aliyun.NewAliDnsClient(config.Conf.Aliyun)
		err := ali.DynamicDNS(publicIP)
		if err != nil {
			log.Print("unable update record: ", err)
		}
		log.Printf("Retry in %d seconds ", config.Conf.Interval)
		select {
		case <-tick.C:

		}
	}
}
