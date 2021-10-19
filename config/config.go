package config

import (
	"dynamicDNS/aliyun"
	"flag"
	"io/ioutil"
	"log"
	"strings"

	"gopkg.in/yaml.v3"
)

type dnspod struct {
}

type config struct {
	Aliyun   aliyun.Options `yaml:"aliyun"`
	Dnspod   dnspod         `yaml:"dnspod"`
	Interval int            `yaml:"interval"`
	IPApi    string         `yaml:"ip_api"`
}

var Conf config

func LoadConfig() {
	confPath := flag.String("conf", "simple-conf.yaml", "path of config, default: simple-conf.yaml")
	flag.Parse()

	log.Print("load config: ", *confPath)
	bytes, err := ioutil.ReadFile(*confPath)
	if err != nil {
		log.Fatalf("load config file [%s] faild: %s", *confPath, err.Error())
	}

	err = yaml.Unmarshal(bytes, &Conf)
	if err != nil {
		log.Fatalf("parses [%s] faild: %s", *confPath, err.Error())
	}
	if Conf.Interval == 0 {
		Conf.Interval = 600
	}

	checkAliyun()

}

func checkAliyun() {
	if Conf.Aliyun.AccessKey == "" {
		log.Fatalf("access_key can not by empty")
	}
	if Conf.Aliyun.AccessKeySecret == "" {
		log.Fatalf("access_key_secret can not by empty")
	}
	if Conf.Aliyun.RegionID == "" {
		log.Fatalf("region_id can not by empty")
	}
	if Conf.Aliyun.Domain == "" {
		log.Fatalf("domain can not by empty")
	}
	if len(strings.Split(Conf.Aliyun.Domain, ".")) != 3 {
		log.Fatal("Invalid domain name: ", Conf.Aliyun.Domain)
	}
}
