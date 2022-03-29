// Copyright 2021 Artem Mikheev

package stravapi

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type rawApiConfig struct {
	Email    string `yaml:"email" validate:"required,email"`
	Password string `yaml:"password" validate:"required"`
	Token    string `yaml:"token"`
	// full mobile application version identifier
	MobileAppVersion string `yaml:"app_version" validate:"required"`
	// company that produced the device, ex "Xiaomi" or Honor
	DeviceManufacturer string `yaml:"device_manufacturer" validate:"required"`
	// model of your device, ex "Redmi Note 9 Pro"
	DeviceModel string `yaml:"device_model" validate:"required"`
	// identifier of android os installed on the device, ex "10" or "11"
	DeviceOsVersion int `yaml:"device_os_version" validate:"required"`
}

// ApiConfig is a yaml-marshalable config
// for constructing a strava api client
type ApiConfig struct {
	rawApiConfig
	// internal identifier of the app version
	appVersion uint32
}

// UnmarshalYAML implements the yaml.Unmarshaler interface
func (cfg *ApiConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	rawCfg := &rawApiConfig{}
	if err := unmarshal(rawCfg); err != nil {
		return err
	}

	validator := validator.New()
	if err := validator.Struct(rawCfg); err != nil {
		return err
	}

	var publicMobileVer float64
	var internalMobileVer uint32
	if _, err := fmt.Sscanf(rawCfg.MobileAppVersion, "%f (%d)", &publicMobileVer, &internalMobileVer); err != nil {
		return err
	}

	cfg.rawApiConfig = *rawCfg
	cfg.appVersion = internalMobileVer
	return nil
}

// MarshalYAML implements the yaml.Marshaler interface
func (cfg *ApiConfig) MarshalYAML() (interface{}, error) {
	return cfg.rawApiConfig, nil
}
