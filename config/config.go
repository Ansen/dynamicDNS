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
	IPApi    string `yaml:"ip_api"`
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
