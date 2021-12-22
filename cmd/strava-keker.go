package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/google/uuid"
	"github.com/renbou/strava-keker/pkg/encoding"
	"github.com/renbou/strava-keker/pkg/fit"
	"github.com/renbou/strava-keker/pkg/strava"
	"github.com/spf13/cobra"
)

var (
	gpxFilePath         string
	desiredSpeed        float64
	activityType        string
	activityName        string
	activityDescription string
	startTime           int
	stravaEmail         string
	stravaPassword      string
)

var rootCmd = &cobra.Command{
	Use:   "strava-keker",
	Short: "strava-keker is a sophisticated strava activity faker",
	Long:  "a fake strava client built using the knowledge gained from reverse engineering the mobile app",
	Run:   run,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().StringVar(&activityType, "activity", "", "either run or ride")
	rootCmd.Flags().StringVar(&activityType, "name", "", "name of the activity")
	rootCmd.Flags().StringVar(&activityType, "description", "", "description of the activity")
	rootCmd.Flags().IntVar(&startTime, "start", 0, "start time in unix format")
	rootCmd.Flags().Float64Var(&desiredSpeed, "speed", 0.0, "desired average speed")
	rootCmd.Flags().StringVar(&gpxFilePath, "gpx", "", "path to gpx file with route")
	rootCmd.Flags().StringVar(&stravaEmail, "email", "", "your email login for strava")
	rootCmd.Flags().StringVar(&stravaPassword, "password", "", "your password for strava")

	rootCmd.MarkFlagRequired("activity")
	rootCmd.MarkFlagRequired("name")
	rootCmd.MarkFlagRequired("description")
	rootCmd.MarkFlagRequired("start")
	rootCmd.MarkFlagRequired("speed")
	rootCmd.MarkFlagRequired("gpx")
	rootCmd.MarkFlagRequired("email")
	rootCmd.MarkFlagRequired("password")
}

func run(cmd *cobra.Command, args []string) {
	// validate parameters
	if activityType != "run" && activityType != "ride" {
		log.Fatal("[!] activity must be either run or ride")
	}
	activityType = strings.Title(activityType)

	// create activity
	activity, err := strava.NewActivity(&strava.ActivityOptions{
		AppVersion:         1221988,
		MobileAppVersion:   "230.10 (1221988)",
		DeviceManufacturer: "Xiaomi",
		DeviceModel:        "Redmi Note 9 Pro",
		DeviceOsVersion:    "10",
		ActivityType:       strava.StravaActivityType(activityType),
		StartTime:          time.Unix(int64(startTime), 0),
		DesiredSpeed:       desiredSpeed,
	})
	if err != nil {
		log.Fatalf("[!] unable to create activity: %v", err)
	}
	log.Print("[+] strava activity created")

	// fill up the activity
	if err := readGpxIntoActivity(activity, gpxFilePath); err != nil {
		log.Fatalf("[!] error while converting gpx to strava activity: %v", err)
	}
	log.Print("[+] gpx successfully converted to strava activity")

	// construct the fit file for the activity
	fitFile, err := activity.BuildFitFile()
	if err != nil {
		log.Fatalf("[!] error while building the fit file for the activity: %v", err)
	}
	log.Print("[+] activity built into fit file")

	if err := sendActivity(fitFile); err != nil {
		log.Fatalf("[!] error while sending fit file: %v", err)
	}
	log.Print("[+] activity uploaded, check your strava!")
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
}

func sendActivity(fitFile *fit.FitFile) error {
	client := resty.New()

	// get an authorization token first, either from cache or a new one
	var token string
	if _, err := os.Stat(".token-cache"); err == nil {
		tokenBytes, err := os.ReadFile(".token-cache")
		if err != nil {
			return err
		}
		token = string(tokenBytes)
	} else {
		resp, err := client.R().
			SetBody(map[string]interface{}{
				"client_id":        2,
				"client_secret":    "3bf7cfbe375675dd9329e9de56d046b4f02a186f",
				"mobile_device_id": "",
				"email":            stravaEmail,
				"email_language":   "en_GB",
				"password":         stravaPassword,
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

		token = resp.Result().(*TokenResponse).AccessToken
		if err := os.WriteFile(".token-cache", []byte(token), os.ModePerm); err != nil {
			return err
		}
	}

	// now upload the activity
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
	if activityType == "Run" {
		workoutType = 3
	} else {
		workoutType = 12
	}
	metadataBytes, err := json.Marshal(map[string]interface{}{
		"activity_name":             "New activity",
		"activity_type":             activityType,
		"commute":                   false,
		"description":               "New description",
		"heartrate_opt_out":         false,
		"perceived_exertion":        4,
		"photo_ids":                 []string{},
		"prefer_perceived_exertion": false,
		"stats_visibility": []map[string]interface{}{
			{
				"type":       "calories",
				"visibility": "everyone",
			},
			{
				"type":       "speed",
				"visibility": "everyone",
			},
		},
		"visibility":   "everyone",
		"workout_type": workoutType,
	})
	if err != nil {
		return err
	}

	resp, err := client.R().
		SetFileReader("file", randomActivityId+"-activity.fit", activityBuffer).
		SetMultipartField("metadata", "", "application/json", bytes.NewBuffer(metadataBytes)).
		SetHeader("Authorization", "access_token "+token).
		Post("https://cdn-1.strava.com/api/v3/uploads/internal_fit?session_id=" + randomActivityId + "&hl=en")
	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusCreated {
		return fmt.Errorf("upload failed, server responded with code %d: %s",
			resp.StatusCode(), resp.String(),
		)
	}
	return nil
}

func readGpxIntoActivity(activity *strava.StravaActivity, gpxFilePath string) error {
	// read and unmarshal the actual file
	gpxFile, err := os.Open(gpxFilePath)
	if err != nil {
		return err
	}
	defer gpxFile.Close()

	bytes, err := ioutil.ReadAll(gpxFile)
	if err != nil {
		return err
	}

	gpx := new(GPX)
	if err := xml.Unmarshal(bytes, &gpx); err != nil {
		return err
	}

	// actually convert the gpx to a strava activity
	for _, trackPart := range gpx.Track.TrackSegment.TrackParts {
		lat, err := strconv.ParseFloat(trackPart.Lat, 64)
		if err != nil {
			return err
		}
		lon, err := strconv.ParseFloat(trackPart.Lon, 64)
		if err != nil {
			return err
		}

		if err := activity.AddRecord(strava.Record{
			Lat:      lat,
			Lon:      lon,
			Altitude: trackPart.Elevation,
		}); err != nil {
			return err
		}
	}

	return nil
}

// GPX file representation with only needed parts
type GPX struct {
	XMLName xml.Name `xml:"gpx"`
	Track   Track    `xml:"trk"`
}

// GPX trk
type Track struct {
	XMLName      xml.Name     `xml:"trk"`
	TrackSegment TrackSegment `xml:"trkseg"`
}

// GPX trkseg
type TrackSegment struct {
	XMLName    xml.Name    `xml:"trkseg"`
	TrackParts []TrackPart `xml:"trkpt"`
}

// GPX trkpt
type TrackPart struct {
	XMLName   xml.Name `xml:"trkpt"`
	Lat       string   `xml:"lat,attr"`
	Lon       string   `xml:"lon,attr"`
	Elevation float64  `xml:"ele"`
}

func main() {
	if err := Execute(); err != nil {
		fmt.Println(err)
	}
}
