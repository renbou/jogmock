// Copyright 2021 Artem Mikheev

package stravapi

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func unmarshalConfig(a *assert.Assertions, raw string, valid bool, expected *ApiConfig) {
	cfg := ApiConfig{}
	err := yaml.Unmarshal([]byte(raw), &cfg)

	if valid {
		if a.NoError(err) {
			a.Equal(expected, &cfg)
		}
	} else {
		a.Error(err)
	}
}

func marshalConfig(a *assert.Assertions, cfg *ApiConfig, expected string) {
	actual, err := yaml.Marshal(cfg)
	if a.NoError(err) {
		a.Equal(expected, string(actual))
	}
}

func TestConfigUnmarshal(t *testing.T) {
	a := assert.New(t)

	// check that it works somehow
	unmarshalConfig(a, `
email: test@test.com
password: some-pass
app_version: 247.10 (1223782)
device_manufacturer: Honor
device_model: P20
device_os_version: 10`, true, &ApiConfig{
		rawApiConfig{
			Email:              "test@test.com",
			Password:           "some-pass",
			MobileAppVersion:   "247.10 (1223782)",
			DeviceManufacturer: "Honor",
			DeviceModel:        "P20",
			DeviceOsVersion:    10,
		},
		1223782,
	})

	// fields must be present and valid
	unmarshalConfig(a, `
password: some-pass
app_version: 247.10 (1223782)
device_manufacturer: Honor
device_model: P20
device_os_version: 10`, false, nil)

	unmarshalConfig(a, `
email: notanemail
password: some-pass
app_version: 247.10 (1223782)
device_manufacturer: Honor
device_model: P20
device_os_version: 10`, false, nil)

	unmarshalConfig(a, `
email: test@test.com
app_version: 247.10 (1223782)
device_manufacturer: Honor
device_model: P20
device_os_version: 10`, false, nil)

	unmarshalConfig(a, `
email: test@test.com
password: some-pass
device_manufacturer: Honor
device_model: P20
device_os_version: 10`, false, nil)

	unmarshalConfig(a, `
email: test@test.com
password: some-pass
app_version: 247.10 (1223782)
device_model: P20
device_os_version: 10`, false, nil)

	unmarshalConfig(a, `
email: test@test.com
password: some-pass
app_version: 247.10 (1223782)
device_manufacturer: Honor
device_os_version: 10`, false, nil)

	unmarshalConfig(a, `
email: test@test.com
password: some-pass
app_version: 247.10 (1223782)
device_manufacturer: Honor
device_model: P20`, false, nil)

	// token must be
}

func TestConfigMarshal(t *testing.T) {
	a := assert.New(t)

	marshalConfig(a,
		&ApiConfig{
			rawApiConfig{
				Email:              "test@test.com",
				Password:           "some-pass",
				Token:              "fake-token",
				MobileAppVersion:   "247.10 (1223782)",
				DeviceManufacturer: "Honor",
				DeviceModel:        "P20",
				DeviceOsVersion:    10,
			},
			uint32(1223782),
		}, `email: test@test.com
password: some-pass
token: fake-token
app_version: 247.10 (1223782)
device_manufacturer: Honor
device_model: P20
device_os_version: 10
`)
}
