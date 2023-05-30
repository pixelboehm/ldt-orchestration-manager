package wotmanager

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
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
		DeviceAddress struct {
			Type  string `json:"type"`
			Value string `json:"value"`
		} `json:"deviceAddress"`
		LdtAddress struct {
			Type  string `json:"type"`
			Value string `json:"value"`
		} `json:"ldtAddress"`
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

func NewWoTmanager(base string) (*WoTManager, error) {
	const location string = "/wotm/description.json"
	var path string = base + location
	path = strings.Replace(path, ":", "/", 1)
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

func (wotm *WoTManager) GetDeviceAddressFromDescription() string {
	desc, _ := wotm.FetchWoTDescription()
	return desc.Properties.DeviceAddress.Value
}

func (wotm *WoTManager) GetLdtPortFromDescription() int {
	desc, _ := wotm.FetchWoTDescription()
	ldt_address := desc.Properties.LdtAddress.Value
	res := ldt_address[strings.LastIndex(ldt_address, ":")+1:]
	port, _ := strconv.Atoi(res)
	return port
}
