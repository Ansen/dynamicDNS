package aliyun

import (
	"dynamicDNS/config"
	"dynamicDNS/utils"
	"fmt"
	"log"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
)

type aliDnsClient struct {
	client    *alidns.Client
	domain    string
	subDomain string
}

func NewAliDnsClient() *aliDnsClient {
	client, err := alidns.NewClientWithAccessKey("cn-hangzhou", config.Conf.Aliyun.AccessKey, config.Conf.Aliyun.AccessKeySecret)
	if err != nil {
		log.Fatal("unable create aliyun dns client: ", err)
		return nil
	}
	domain := strings.Split(config.Conf.Aliyun.Domain, ".")
	mainDomain := fmt.Sprintf("%s.%s", domain[1], domain[2])

	return &aliDnsClient{
		client:    client,
		domain:    mainDomain,
		subDomain: domain[0],
	}
}

func (a *aliDnsClient) getSubDomainRecord() ([]alidns.Record, error) {
	request := alidns.CreateDescribeSubDomainRecordsRequest()
	request.SubDomain = config.Conf.Aliyun.Domain
	request.DomainName = a.domain
	response, err := a.client.DescribeSubDomainRecords(request)
	if err != nil {
		log.Print("unable query subdomain records: ", err)
		return response.DomainRecords.Record, err
	}
	return response.DomainRecords.Record, nil
}

func (a *aliDnsClient) addSubDomainRecord(recordType, value string) error {
	request := alidns.CreateAddDomainRecordRequest()
	request.DomainName = a.domain
	request.RR = a.subDomain
	request.Type = recordType
	request.Value = value
	response, err := a.client.AddDomainRecord(request)
	if err != nil {
		log.Print("unable add domain record:", err)
		log.Print("response detail: ", response.String())
		return err
	}
	log.Printf("%s added, value: %s, requestID: %s", config.Conf.Aliyun.Domain, value, response.RequestId)
	return nil
}

func (a *aliDnsClient) updateSubDomainRecord(recordId, recordType, value string) error {
	request := alidns.CreateUpdateDomainRecordRequest()
	request.RecordId = recordId
	request.RR = a.subDomain
	request.Type = recordType
	request.Value = value
	response, err := a.client.UpdateDomainRecord(request)
	if err != nil {
		log.Print("unable add domain record:", err)
		return err
	}
	log.Printf("%s updated, value: %s, requestID: %s", config.Conf.Aliyun.Domain, value, response.RequestId)
	return nil
}

func (a *aliDnsClient) DynamicDNS(ipv4, ipv6 string) error {
	if ipv4 == "" && ipv6 == "" {
		return utils.PublicIPEmpty
	}
	records, err := a.getSubDomainRecord()
	if err != nil {
		return err
	}
	switch len(records) {
	case 0:
		if ipv4 != "" {
			return a.addSubDomainRecord("A", ipv4)
		} else {
			return a.addSubDomainRecord("AAAA", ipv6)
		}
	case 1:
		record := records[0]
		if record.Type == "A" && record.Value != ipv4 {
			return a.updateSubDomainRecord(record.RecordId, record.Type, ipv4)
		} else if record.Type == "AAAA" && record.Value != ipv6 {
			return a.updateSubDomainRecord(record.RecordId, record.Type, ipv6)
		} else {
			log.Print("no need to update")
		}
	case 2:
		var err error
		for _, r := range records {
			if r.Type == "A" && r.Value != ipv4 {
				err = a.updateSubDomainRecord(r.RecordId, r.Type, ipv4)
			}
			if r.Type == "AAAA" && r.Value != ipv6 {
				err = a.updateSubDomainRecord(r.RecordId, r.Type, ipv6)
			}
			if err != nil {
				return err
			}
		}
	default:
		return utils.UnSupportMultiRecord
	}
	return nil
}
