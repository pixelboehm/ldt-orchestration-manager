package wotmanager

import (
	"encoding/json"
	"io/ioutil"
	"strings"
)

type WoTDescription struct {
	Context             string `json:"@context"`
	ID                  string `json:"id"`
	Title               string `json:"title"`
	SecurityDefinitions struct {
		BasicSc struct {
			Scheme string `json:"scheme"`
			In     string `json:"in"`
		} `json:"basic_sc"`
	} `json:"securityDefinitions"`
	Security   []string `json:"security"`
	Properties struct {
		Status struct {
			Type  string `json:"type"`
			Forms []struct {
				Href string `json:"href"`
			} `json:"forms"`
		} `json:"status"`
		Address struct {
			Type  string `json:"type"`
			Value string `json:"value"`
		} `json:"address"`
	} `json:"properties"`
	Actions struct {
		On struct {
			Forms []struct {
				Href string `json:"href"`
			} `json:"forms"`
		} `json:"on"`
		Off struct {
			Forms []struct {
				Href string `json:"href"`
			} `json:"forms"`
		} `json:"off"`
	} `json:"actions"`
	Events struct {
		Overheating struct {
			Data struct {
				Type string `json:"type"`
			} `json:"data"`
			Forms []struct {
				Href        string `json:"href"`
				Subprotocol string `json:"subprotocol"`
			} `json:"forms"`
		} `json:"overheating"`
	} `json:"events"`
}

type WoTManager struct {
	description_raw string
}

func NewWoTmanager(ldt_path string) (*WoTManager, error) {
	var base string = "/usr/local/etc/orchestration-manager/"
	test := ldt_path[strings.LastIndex(ldt_path, ":")+1:]
	var replaced string
	if test == "latest" {
		replaced = strings.Replace(ldt_path, ":", "/", 1)
	} else {
		replaced = strings.Replace(ldt_path, ":", "/v", 1)
	}

	const location string = "/wotm/description.json"
	var path string = base + replaced + location
	buffer, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return &WoTManager{
		description_raw: string(buffer),
	}, nil
}

func (wotm *WoTManager) FetchWoTDescription() (WoTDescription, error) {
	var desc WoTDescription
	if err := json.Unmarshal([]byte(wotm.description_raw), &desc); err != nil {
		return WoTDescription{}, err
	}
	return desc, nil
}

func (wotm *WoTManager) getDeviceAddressFromDescription(desc WoTDescription) string {
	return desc.Properties.Address.Value
}
