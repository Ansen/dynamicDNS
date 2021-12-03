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
		ipv4 := ""
		if *config.Conf.IPV4 {
			response := utils.GetPublicIP(config.Conf.IPv4Api)
			ipv4 = utils.GetIPv4(response)
		}

		ipv6 := ""
		if *config.Conf.IPV6 {
			response := utils.GetPublicIP(config.Conf.IPv6Api)
			ipv6 = utils.GetIPv6(response)
		}

		if *config.Conf.IPV4 && ipv4 == "" {
			log.Print("unable get public ipv4")
			ticker(tick)
			continue
		}

		if *config.Conf.IPV6 && ipv6 == "" {
			log.Print("unable get public ipv6")
			ticker(tick)
			continue
		}

		if config.Conf.Aliyun != notConfig {
			log.Print("start handle ali yun dynamic dns...")
			ali := aliyun.NewAliDnsClient()
			err := ali.DynamicDNS(ipv4, ipv6)
			if err != nil {
				log.Print("unable update ali dns record: ", err)
			}
		}
		if config.Conf.Tencent != notConfig {
			log.Print("start handle tencent dynamic dns...")
			tencentClient := tencent.NewDnspodApi()
			err := tencentClient.DynamicDNS(ipv4, ipv6)
			if err != nil {
				log.Print("unable update dnspod recode: ", err)
			}
		}
		ticker(tick)
	}
}
