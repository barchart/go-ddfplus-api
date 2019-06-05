// Go ddfplus API
//
// Copyright 2019 Barchart.com, Inc. All rights reserved.
//
// This Source Code Form is subject to the terms of the GNU license
// available at https://github.com/barchart/go-ddfplus-api/blob/master/LICENSE.
package ddf

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const UserSettingsURL = "http://www.ddfplus.com/json/usersettings/"

type UserSettings struct {
	Login struct {
		Username    string `json:"username"`
		Status      bool   `json:"status"`
		Credentials bool   `json:"credentials"`
	} `json:"login"`
	Service struct {
		Id         string `json:"id"`
		MaxSymbols int    `json:"maxsymbols"`
	} `json:"service"`
	Exchanges []string `json:"exchanges"`
	Servers   struct {
		Stream     []string `json:"stream"`
		WSS        []string `json:"wss"`
		Historical []string `json:"historical"`
		Extras     []string `json:"extras"`
		News       []string `json:"news"`
	} `json:"servers"`
}

func GetUserSettings(credentials *Credentials) (UserSettings, error) {
	var (
		settings UserSettings
	)

	resp, err := http.Get(UserSettingsURL + "?username=" + credentials.Username + "&password=" + credentials.Password)
	if err == nil {
		bytes, _ := ioutil.ReadAll(resp.Body)

		tmp := new(struct {
			Settings UserSettings `json:"usersettings"`
		})

		err := json.Unmarshal(bytes, &tmp)
		if err != nil {
			return settings, err
		}

		settings = tmp.Settings
		resp.Body.Close()
	} else {
		return settings, err
	}

	return settings, nil
}
