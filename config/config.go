package config

import (
	"flag"
	"io/ioutil"
	"log"
	"strings"

	"gopkg.in/yaml.v3"
)

type config struct {
	Aliyun   Option `yaml:"aliyun"`
	Tencent  Option `yaml:"tencent"`
	Interval int    `yaml:"interval"`
	IPv4Api  string `yaml:"ipv4_api"`
	IPv6Api  string `yaml:"ipv6_api"`
	IPV4     *bool  `yaml:"ipv_4"`
	IPV6     *bool  `yaml:"ipv_6"`
}

type Option struct {
	AccessKey       string `yaml:"access_key"`
	AccessKeySecret string `yaml:"access_key_secret"`
	Domain          string `yaml:"domain"`
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
		log.Print("set default interval: 600")
		Conf.Interval = 600
	}

	if Conf.IPv4Api == "" && Conf.IPv6Api == "" {
		Conf.IPv4Api = "http://ip.3322.net/"
	}

	if *Conf.IPV4 == false && *Conf.IPV6 == false {
		*Conf.IPV4 = true
		log.Println("set default ipv4")
	}

	if Conf.IPv6Api == "" && *Conf.IPV6 {
		Conf.IPv6Api = Conf.IPv4Api
	}

	notConfig := Option{}
	if Conf.Aliyun != notConfig {
		checkOption(&Conf.Aliyun)
	}

	if Conf.Tencent != notConfig {
		checkOption(&Conf.Tencent)
	}

}

func checkOption(option *Option) {
	if option.AccessKey == "" {
		log.Fatalf("access_key can not by empty")
	}
	if option.AccessKeySecret == "" {
		log.Fatalf("access_key_secret can not by empty")
	}
	if option.Domain == "" {
		log.Fatalf("domain can not by empty")
	}
	if len(strings.Split(option.Domain, ".")) != 3 {
		log.Fatal("Invalid domain name: ", Conf.Aliyun.Domain)
	}
}
