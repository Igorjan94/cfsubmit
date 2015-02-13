package cfsubmit

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"gopkg.in/xmlpath.v2"
)

const SettingsFileName = "cfsubmit_settings.json"

var (
	ErrUnexpectedFormat = errors.New("Unexpected page format. Please check connection and/or X-User setting")
)

type CFSettings struct {
	XUser        string            `json:"X-User"`
	CSRF         string            `json:"CSRF-Token"`
	Handle       string            `json:"Handle"`
	CFDomain     string            `json:"CF-Domain"`
	CheckResults bool              `json:"Check-Results"`
	ExtId        map[string]string `json:"Ext-ID"`
	CFIdCodes    map[string]string `json:"CF-id-codes"`
}

func ReadSettings() (*CFSettings, error) {
	jsonData, err := os.Open(SettingsFileName)
	defer jsonData.Close()
	if err != nil {
		return nil, err
	}

	var settings CFSettings
	if err := json.NewDecoder(jsonData).Decode(&settings); err != nil {
		return nil, err
	}
	return &settings, nil
}

func WriteSettings(settings *CFSettings) error {
	jsonData, err := os.Create(SettingsFileName)
	defer jsonData.Close()
	if err != nil {
		return err
	}

	b, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}

	if _, err := jsonData.Write(b); err != nil {
		return err
	}
	return nil
}

func UpdateSettings() error {
	settings, err := ReadSettings()

	req, err := http.NewRequest("GET", "http://codeforces.com/problemset/submit", nil)
	req.AddCookie(&http.Cookie{Name: "X-User", Value: settings.XUser})

	resp, err := http.DefaultClient.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	node, err := xmlpath.ParseHTML(resp.Body)
	if err != nil {
		return err
	}

	settings.CFIdCodes = make(map[string]string)

	path := xmlpath.MustCompile(".//select[@name='programTypeId']/option")
	codeAttr := xmlpath.MustCompile("./@value")

	for option := path.Iter(node); option.Next(); {
		code, ok := codeAttr.String(option.Node())
		if !ok {
			return ErrUnexpectedFormat
		}
		settings.CFIdCodes[code] = option.Node().String()
	}

	return WriteSettings(settings)
}
