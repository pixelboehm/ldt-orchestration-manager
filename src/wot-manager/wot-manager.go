package wotmanager

import (
	"encoding/json"
	"io/ioutil"
	"strings"
)

type WoTDescription struct {
	Context             []any    `json:"@context,omitempty"`
	Type                []string `json:"@type,omitempty"`
	ID                  string   `json:"id,omitempty"`
	Name                string   `json:"name,omitempty"`
	SecurityDefinitions struct {
		BasicSc struct {
			Scheme string `json:"scheme,omitempty"`
			In     string `json:"in,omitempty"`
		} `json:"basic_sc,omitempty"`
	} `json:"securityDefinitions,omitempty"`
	Security   []string `json:"security,omitempty"`
	Properties struct {
		Status struct {
			Type       string `json:"@type,omitempty"`
			ReadOnly   bool   `json:"readOnly,omitempty"`
			WriteOnly  bool   `json:"writeOnly,omitempty"`
			Observable bool   `json:"observable,omitempty"`
			Type0      string `json:"type,omitempty"`
			Forms      []struct {
				Href          string `json:"href,omitempty"`
				ContentType   string `json:"contentType,omitempty"`
				HtvMethodName string `json:"htv:methodName,omitempty"`
				Op            string `json:"op,omitempty"`
			} `json:"forms,omitempty"`
		} `json:"status,omitempty"`
	} `json:"properties,omitempty"`
	Actions struct {
		Toggle struct {
			Type       string `json:"@type,omitempty"`
			Idempotent bool   `json:"idempotent,omitempty"`
			Safe       bool   `json:"safe,omitempty"`
			Forms      []struct {
				Href          string `json:"href,omitempty"`
				ContentType   string `json:"contentType,omitempty"`
				HtvMethodName string `json:"htv:methodName,omitempty"`
				Op            string `json:"op,omitempty"`
			} `json:"forms,omitempty"`
		} `json:"toggle,omitempty"`
	} `json:"actions,omitempty"`
	Events struct {
		Overheating struct {
			Data struct {
				Type string `json:"type,omitempty"`
			} `json:"data,omitempty"`
			Forms []struct {
				Href        string `json:"href,omitempty"`
				ContentType string `json:"contentType,omitempty"`
				Subprotocol string `json:"subprotocol,omitempty"`
				Op          string `json:"op,omitempty"`
			} `json:"forms,omitempty"`
		} `json:"overheating,omitempty"`
	} `json:"events,omitempty"`
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
