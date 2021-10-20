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

func (d *dnspodApi) addSubDomainRecord(value string) (record, error) {
	log.Print("start add record ...")
	resp, err := http.PostForm(recordCreate, url.Values{
		"login_token": {d.loginToken},
		"format":      {d.format},
		"domain":      {d.domain},
		"sub_domain":  {d.subDomain},
		"record_type": {"A"},
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

func (d *dnspodApi) updateSubDomainRecord(id, value string) (modifyRecord, error) {
	log.Print("start update record ...")
	resp, err := http.PostForm(recordModify, url.Values{
		"login_token": {d.loginToken},
		"format":      {d.format},
		"domain":      {d.domain},
		"sub_domain":  {d.subDomain},
		"record_id":   {id},
		"record_type": {"A"},
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

func (d *dnspodApi) DynamicDNS(ip string) error {
	if ip == "" {
		return utils.PublicIPEmpty
	}

	records, err := d.getSubdomainRecord()
	if err != nil {
		return err
	}
	switch len(records) {
	case 0:
		added, err := d.addSubDomainRecord(ip)
		if err != nil {
			return err
		}
		log.Printf("added: %s => %s", added.Name, ip)
	case 1:
		if ip == records[0].Value {
			log.Print(utils.NoChangeSkip.Error())
			return nil
		}
		modify, err := d.updateSubDomainRecord(records[0].ID, ip)
		if err != nil {
			return err
		}
		log.Printf("updated: %s => %s", modify.Name, modify.Value)
	default:
		return utils.UnSupportMultiRecord
	}

	return nil
}