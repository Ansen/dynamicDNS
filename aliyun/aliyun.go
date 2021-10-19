package aliyun

import (
	"fmt"
	"log"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
)

type aliDnsClient struct {
	client    *alidns.Client
	option    Options
	domain    string
	subDomain string
}

type Options struct {
	AccessKey       string `yaml:"access_key"`
	AccessKeySecret string `yaml:"access_key_secret"`
	RegionID        string `yaml:"region_id"`
	Domain          string `yaml:"domain"`
}

func NewAliDnsClient(option Options) *aliDnsClient {
	client, err := alidns.NewClientWithAccessKey(option.RegionID, option.AccessKey, option.AccessKeySecret)
	if err != nil {
		log.Fatal("unable create alidns client: ", err)
		return nil
	}
	domain := strings.Split(option.Domain, ".")
	mainDomain := fmt.Sprintf("%s.%s", domain[1], domain[2])

	return &aliDnsClient{
		client:    client,
		option:    option,
		domain:    mainDomain,
		subDomain: domain[0],
	}
}

func (a *aliDnsClient) getSubDomainRecord() ([]alidns.Record, error) {
	request := alidns.CreateDescribeSubDomainRecordsRequest()
	request.SubDomain = a.option.Domain
	request.DomainName = a.domain
	response, err := a.client.DescribeSubDomainRecords(request)
	if err != nil {
		log.Print("unable query subdomain records: ", err)
		return response.DomainRecords.Record, err
	}
	return response.DomainRecords.Record, nil
}

func (a *aliDnsClient) addSubDomainRecord(value string) error {
	request := alidns.CreateAddDomainRecordRequest()
	request.DomainName = a.domain
	request.RR = a.subDomain
	request.Type = "A"
	request.Value = value
	response, err := a.client.AddDomainRecord(request)
	if err != nil {
		log.Print("unable add domain record:", err)
		log.Print("response detail: ", response.String())
		return err
	}
	log.Printf("%s added, value: %s, requestID: %s", a.option.Domain, value, response.RequestId)
	return nil
}

func (a *aliDnsClient) updateSubDomainRecord(recordId, value string) error {
	request := alidns.CreateUpdateDomainRecordRequest()
	request.RecordId = recordId
	request.RR = a.subDomain
	request.Type = "A"
	request.Value = value
	response, err := a.client.UpdateDomainRecord(request)
	if err != nil {
		log.Print("unable add domain record:", err)
		return err
	}
	log.Printf("%s updated, value: %s, requestID: %s", a.option.Domain, value, response.RequestId)
	return nil
}

func (a *aliDnsClient) DynamicDNS(ip string) error {
	records, err := a.getSubDomainRecord()
	if err != nil {
		return err
	}
	switch len(records) {
	case 0:
		return a.addSubDomainRecord(ip)
	case 1:
		record := records[0]
		if record.Value == ip {
			log.Print("No change, Skip.")
			return nil
		}
		return a.updateSubDomainRecord(record.RecordId, ip)
	default:
		log.Printf("Does not support updating multiple records")
		return nil
	}
}
