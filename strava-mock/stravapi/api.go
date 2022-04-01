// Copyright 2021 Artem Mikheev

package stravapi

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/renbou/jogmock/activities"
	"github.com/renbou/jogmock/fit-encoder/encoding"
	"github.com/renbou/jogmock/strava-mock/internal/stravafit"
)

var ErrUnauthorized = errors.New("strava api client unathorized")

type ApiClient struct {
	ApiConfig
	client *resty.Client
}

func NewClient(config *ApiConfig) (*ApiClient, error) {
	if config.internalAppVersion == 0 {
		return nil, errors.New("internal app version must not be zero")
	}
	if config.MobileAppVersion == "" {
		return nil, errors.New("full mobile app version must be set")
	}
	if config.DeviceManufacturer == "" {
		return nil, errors.New("device manufacturer must be set")
	}
	if config.DeviceModel == "" {
		return nil, errors.New("device model must be set")
	}

	switch config.DeviceOsVersion {
	case 11:
	case 10:
	case 9:
		break
	default:
		return nil, errors.New("device os version must be set to 9, 10 or 11")
	}

	restyClient := resty.New()
	restyClient.SetHeader("X-Strava-Trial-Status", "false")
	restyClient.SetHeader("X-Strava-Nav-Version", "2")
	restyClient.SetHeader("User-Agent", fmt.Sprintf("Strava/%.2f", config.publicAppVersion))
	// allow running through any proxy
	restyClient.SetTLSClientConfig(&tls.Config{
		InsecureSkipVerify: true,
	})

	client := &ApiClient{}
	client.ApiConfig = *config
	client.client = restyClient
	return client, nil
}

func (api *ApiClient) UploadActivity(activity *activities.Activity) error {
	if api.Token == "" {
		return ErrUnauthorized
	}

	a := stravafit.StravaActivity{
		AppVersion:         api.internalAppVersion,
		MobileAppVersion:   api.MobileAppVersion,
		DeviceManufacturer: api.DeviceManufacturer,
		DeviceModel:        api.DeviceModel,
		DeviceOsVersion:    strconv.Itoa(api.DeviceOsVersion),
		Activity:           activity,
	}
	fitFile, err := a.BuildFitFile()
	if err != nil {
		return err
	}

	activityBuffer := new(bytes.Buffer)
	encoder := encoding.NewEncoder(activityBuffer, encoding.BigEndian)
	if err := encoder.Encode(fitFile); err != nil {
		return fmt.Errorf("error while encoding fit file: %v", err)
	}

	randomActivityUUID, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	randomActivityId := randomActivityUUID.String()

	// and a second part for the metadata
	var workoutType int
	activityType := activity.Type()
	if activityType == activities.RunActivity {
		workoutType = 3
	} else if activityType == activities.RideActivity {
		workoutType = 12
	} else {
		panic(fmt.Sprintf("unknown activity type: %v", activityType))
	}
	metadataBuffer := new(bytes.Buffer)
	metadataEncoder := json.NewEncoder(metadataBuffer)
	err = metadataEncoder.Encode(map[string]interface{}{
		"activity_name":             activity.Name(),
		"activity_type":             stravafit.ActivityTypeToString(activityType),
		"commute":                   false,
		"description":               activity.Description(),
		"heartrate_opt_out":         false,
		"hide_from_home":            false,
		"perceived_exertion":        4,
		"photo_ids":                 []string{},
		"prefer_perceived_exertion": false,
		"selected_polyline_style":   "default",
		"stats_visibility": []map[string]interface{}{
			{
				"type":       "calories",
				"visibility": "everyone",
			},
			{
				"type":       "speed",
				"visibility": "everyone",
			},
			{
				"type":       "pace",
				"visibility": "everyone",
			},
		},
		"visibility":   "everyone",
		"workout_type": workoutType,
	})
	if err != nil {
		return err
	}

	resp, err := api.client.R().
		SetFileReader("file", randomActivityId+"-activity.fit", activityBuffer).
		SetMultipartField("metadata", "", "application/json", metadataBuffer).
		SetHeader("Authorization", "access_token "+api.Token).
		Post("https://cdn-1.strava.com/api/v3/uploads/internal_fit?session_id=" + randomActivityId + "&hl=en")
	if err != nil {
		return err
	}

	if resp.StatusCode() == http.StatusCreated {
		return nil
	} else if resp.StatusCode() == http.StatusUnauthorized {
		return ErrUnauthorized
	}
	return fmt.Errorf("upload failed, server responded with code %d: %s",
		resp.StatusCode(), resp.String(),
	)
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
}

func (api *ApiClient) Authorize(recaptchaToken string) error {
	resp, err := api.client.R().
		SetBody(map[string]interface{}{
			"client_id":        2,
			"client_secret":    "3bf7cfbe375675dd9329e9de56d046b4f02a186f",
			"mobile_device_id": "",
			"email":            api.Email,
			"email_language":   "en_GB",
			"password":         api.Password,
			"recaptcha":        recaptchaToken,
		}).
		SetResult(&TokenResponse{}).
		Post("https://cdn-1.strava.com/api/v3/oauth/internal/token?hl=en")
	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("authorization failed, server responded with code %d: %s",
			resp.StatusCode(), resp.String(),
		)
	}

	token := resp.Result().(*TokenResponse).AccessToken
	if token == "" {
		return fmt.Errorf("no token received in authorization response: %s", resp.String())
	}

	api.Token = token
	return nil
}

// Config extract this client's config.
// The returned config can be marshaled and reused to initialized
// the client again in the future.
func (api *ApiClient) Config() *ApiConfig {
	cfgCopy := api.ApiConfig
	return &cfgCopy
}
