package main

import (
	"dynamicDNS/aliyun"
	"dynamicDNS/config"
	"dynamicDNS/tencent"
	"dynamicDNS/utils"
	"log"
	"time"
)

func init() {
	config.LoadConfig()
}

func ticker(tick *time.Ticker) {
	log.Printf("Retry in %d seconds ", config.Conf.Interval)
	select {
	case <-tick.C:
	}
}

func main() {

	notConfig := config.Option{}
	tick := time.NewTicker(time.Duration(config.Conf.Interval) * time.Second)

	for {
		publicIP := utils.GetPublicIP()
		if publicIP == "" {
			log.Print("unable get public ip")
			ticker(tick)
			continue
		}

		if config.Conf.Aliyun != notConfig {
			log.Print("start handle ali yun dynamic dns...")
			ali := aliyun.NewAliDnsClient()
			err := ali.DynamicDNS(publicIP)
			if err != nil {
				log.Print("unable update ali dns record: ", err)
			}
		}
		if config.Conf.Tencent != notConfig {
			log.Print("start handle tencent dynamic dns...")
			tencentClient := tencent.NewDnspodApi()
			err := tencentClient.DynamicDNS(publicIP)
			if err != nil {
				log.Print("unable update dnspod recode: ", err)
			}
		}
		ticker(tick)
	}
}
