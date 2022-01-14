package tencent

import (
	"dynamicDNS/config"
	"dynamicDNS/utils"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type dnspodApi struct {
	loginToken string
	format     string
	domain     string
	subDomain  string
}

type status struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type record struct {
	ID      string `json:"id,omitempty"`
	Name    string `json:"name"`
	Type    string `json:"type,omitempty"`
	Value   string `json:"value,omitempty"`
	Enabled string `json:"enabled,omitempty"`
}

type modifyRecord struct {
	ID int `json:"id"`
	record
}

func NewDnspodApi() *dnspodApi {
	domain := strings.Split(config.Conf.Tencent.Domain, ".")
	mainDomain := fmt.Sprintf("%s.%s", domain[1], domain[2])
	return &dnspodApi{
		loginToken: fmt.Sprintf("%s,%s", config.Conf.Tencent.AccessKey, config.Conf.Tencent.AccessKeySecret),
		format:     "json",
		domain:     mainDomain,
		subDomain:  domain[0],
	}
}

func (d *dnspodApi) handleResponseBody(body *io.ReadCloser, responseObject interface{}) error {
	bytes, err := io.ReadAll(*body)
	if err != nil {
		log.Print("unable read body: ", err)
		return err
	}
	if err := json.Unmarshal(bytes, responseObject); err != nil {
		log.Print("unable unmarshal body: ", err)
		log.Print("bytes: ", string(bytes))
		return err
	}
	return nil

}

func (d *dnspodApi) getSubdomainRecord() ([]record, error) {
	resp, err := http.PostForm(recordList, url.Values{
		"login_token": {d.loginToken},
		"format":      {d.format},
		"domain":      {d.domain},
		"sub_domain":  {d.subDomain},
	})
	if err != nil {
		log.Print("request record list failed: ", err)
		return nil, err
	}
	defer resp.Body.Close()

	content := struct {
		Status  status   `json:"status"`
		Records []record `json:"records"`
	}{}
	if err := d.handleResponseBody(&resp.Body, &content); err != nil {
		return nil, err
	}
	if content.Status.Code != "1" {
		log.Printf("code: %s, message: %s", content.Status.Code, content.Status.Message)
	} else {
		log.Print(content.Records)
	}
	return content.Records, nil
}

func (d *dnspodApi) addSubDomainRecord(recordType, value string) (record, error) {
	log.Print("start add record ...")
	resp, err := http.PostForm(recordCreate, url.Values{
		"login_token": {d.loginToken},
		"format":      {d.format},
		"domain":      {d.domain},
		"sub_domain":  {d.subDomain},
		"record_type": {recordType},
		"value":       {value},
		"record_line": {"默认"},
	})
	if err != nil {
		log.Print("request record create failed: ", err)
		return record{}, err
	}
	defer resp.Body.Close()

	content := struct {
		Status status `json:"status"`
		Record record `json:"record"`
	}{}
	if err := d.handleResponseBody(&resp.Body, &content); err != nil {
		return record{}, err
	}
	if content.Status.Code != "1" {
		log.Printf("code: %s, message: %s", content.Status.Code, content.Status.Message)
		return record{}, errors.New(content.Status.Message)
	}
	log.Print(content.Record)
	return content.Record, nil
}

func (d *dnspodApi) updateSubDomainRecord(id, recordType, value string) (modifyRecord, error) {
	log.Print("start update record ...")
	resp, err := http.PostForm(recordModify, url.Values{
		"login_token": {d.loginToken},
		"format":      {d.format},
		"domain":      {d.domain},
		"sub_domain":  {d.subDomain},
		"record_id":   {id},
		"record_type": {recordType},
		"value":       {value},
		"record_line": {"默认"},
	})
	if err != nil {
		log.Print("request record create failed: ", err)
		return modifyRecord{}, err
	}
	defer resp.Body.Close()
	content := struct {
		Status status       `json:"status"`
		Record modifyRecord `json:"record"`
	}{}
	if err := d.handleResponseBody(&resp.Body, &content); err != nil {
		return modifyRecord{}, err
	}
	if content.Status.Code != "1" {
		log.Printf("code: %s, message: %s", content.Status.Code, content.Status.Message)
	} else {
		log.Print(content.Record)
	}
	return content.Record, nil
}

func (d *dnspodApi) DynamicDNS(ipv4, ipv6 string) error {
	if ipv4 == "" && ipv6 == "" {
		return utils.PublicIPEmpty
	}

	records, err := d.getSubdomainRecord()
	if err != nil {
		return err
	}
	for _, record := range records {
		log.Printf("remote record Type: %s, Value: %s, Enabled: %s", record.Type, record.Value, record.Enabled)
	}
	switch len(records) {
	case 0:
		if ipv4 != "" {
			if _, err := d.addSubDomainRecord("A", ipv4); err != nil {
				return err
			}
		}
		if ipv6 != "" {
			if _, err := d.addSubDomainRecord("AAAA", ipv6); err != nil {
				return err
			}
		}
	case 1:
		if records[0].Type == "A" && ipv4 != records[0].Value {
			if _, err := d.updateSubDomainRecord(records[0].ID, "A", ipv4); err != nil {
				return err
			}
		} else if records[0].Type == "AAAA" && ipv6 != records[0].Value {
			if _, err := d.updateSubDomainRecord(records[0].ID, "AAAA", ipv6); err != nil {
				return err
			}
		} else {
			log.Print("remote ip is same as local ip, skip update")
		}
	default:
		return utils.UnSupportMultiRecord
	}
	return nil
}
