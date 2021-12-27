package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"gopkg.in/yaml.v2"

	"github.com/go-playground/validator/v10"
	bubblesCommon "github.com/mritd/bubbles/common"
	bubblesSelector "github.com/mritd/bubbles/selector"
	"github.com/renbou/strava-keker/pkg/encoding"
	"github.com/renbou/strava-keker/pkg/fit"
	bubblesPrompt "github.com/renbou/strava-keker/pkg/prompt"
	"github.com/renbou/strava-keker/pkg/strava"
	"github.com/spf13/cobra"
)

var (
	arguments Arguments
	model     = NewActivityModel()
	config    *UserConfig
	rootCmd   = &cobra.Command{
		Use:   "strava-keker",
		Short: "strava-keker is a sophisticated strava activity faker",
		Long:  "a fake strava client built using the knowledge gained from reverse engineering the mobile app",
		Run:   run,
	}

	infoStyle = lipgloss.
			NewStyle().
			Foreground(lipgloss.Color("2")).
			Inline(true)

	errorStyle = lipgloss.
			NewStyle().
			Foreground(lipgloss.Color("9")).
			Bold(true).
			Inline(true)
)

type UserConfig struct {
	Email                    string `validate:"required,email"`
	Password                 string `validate:"required"`
	internalMobileAppVersion uint32
	MobileAppVersion         string `yaml:"mobile_app_version" validate:"required"`
	DeviceManufacturer       string `yaml:"device_manufacturer" validate:"required"`
	DeviceModel              string `yaml:"device_model" validate:"required"`
	DeviceOsVersion          int    `yaml:"device_os_version" validate:"required"`
}

// Arguments represents the possible commmand-line arguments
type Arguments struct {
	ConfigPath string
	DontUpload bool
}

// GetConfig reads and returns the config defined by args
func (args *Arguments) GetConfig() (*UserConfig, error) {
	configFile, err := os.Open(args.ConfigPath)
	if err != nil {
		return nil, err
	}
	decoder := yaml.NewDecoder(configFile)

	var config UserConfig
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}
	validator := validator.New()
	if err := validator.Struct(config); err != nil {
		return nil, err
	}

	var publicMobileVer float64
	var internalMobileVer uint32
	if _, err := fmt.Sscanf(config.MobileAppVersion, "%f (%d)", &publicMobileVer, &internalMobileVer); err != nil {
		return nil, err
	}
	config.internalMobileAppVersion = internalMobileVer

	return &config, nil
}

type ActivityModel struct {
	gpxFilePath         string
	activityType        strava.StravaActivityType
	activityName        string
	activityDescription string
	startTime           time.Time
	desiredSpeed        float64
	selector            *bubblesSelector.Model
	prompt              *bubblesPrompt.Model
	inputIndex          int
	failedInput         bool
}

func pathExists(path string) error {
	_, path = userExpand(path)
	stat, err := os.Stat(path)
	if err != nil {
		return err
	}
	if stat.IsDir() {
		return errors.New("path is a directory")
	}
	return nil
}

func NewActivityModel() *ActivityModel {
	return &ActivityModel{
		inputIndex: 0,
		prompt: &bubblesPrompt.Model{
			Prompt:       bubblesCommon.FontColor("Input path to GPX file: ", bubblesPrompt.ColorPrompt),
			ValidateFunc: pathExists,
		},
	}
}

func strToTimestamp(val string) (time.Time, error) {
	t, err := time.Parse("02.01.2006 15:04:05", val)
	return t.Local(), err
}

func strIsTime(val string) error {
	if val == "" {
		return errors.New("input time as DD.MM.YYYY HH:MM:SS")
	}
	_, err := strToTimestamp(val)
	if err != nil {
		return errors.New("input time as DD.MM.YYYY HH:MM:SS, error: " + err.Error())
	}
	return err
}

func (m *ActivityModel) Init() tea.Cmd {
	return nil
}

// returns true if the path was expanded
func userExpand(path string) (bool, string) {
	if path == "" {
		return false, ""
	}
	usr, _ := user.Current()
	if path[0] == '~' {
		// expand the path
		path = filepath.Join(usr.HomeDir, path[1:])
		return true, path
	}
	return false, path
}

func autoSuggestFile(path string) string {
	usr, _ := user.Current()
	wasExpanded, path := userExpand(path)

	// get directory and expand
	directory := filepath.Dir(path)
	file := filepath.Base(path)
	if file == "" {
		return ""
	}

	// expand by prefix
	dirListing, err := os.ReadDir(directory)
	if err != nil {
		return ""
	}
	for _, dirFile := range dirListing {
		if strings.HasPrefix(dirFile.Name(), file) {
			suggested := filepath.Join(directory, dirFile.Name())
			if dirFile.IsDir() {
				suggested += "/"
			}
			if wasExpanded {
				suggested = suggested[len(usr.HomeDir):]
				suggested = "~" + suggested
			}
			return suggested
		}
	}
	return ""
}

var updateViewMsg tea.Msg = "UPDATE_VIEW"

func updateView() tea.Msg {
	return updateViewMsg
}

func (m *ActivityModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// if entering tab and we are currently reading file, autosuggest it
	if key, ok := msg.(tea.KeyMsg); ok && key.Type == tea.KeyTab && m.inputIndex == 0 {
		if suggest := autoSuggestFile(m.prompt.Value()); suggest != "" {
			m.prompt.Err = m.prompt.ValidateFunc(suggest)
			m.prompt.Input.SetValue(suggest)
			m.prompt.Input.SetCursor(len(suggest))
			return m, nil
		}
	}
	// otherwise if simply reading values
	switch msg {
	case updateViewMsg:
		return m, tea.Quit
	case bubblesCommon.DONE:
		switch m.inputIndex {
		case 0:
			_, m.gpxFilePath = userExpand(m.prompt.Value())
			m.prompt = nil
			m.selector = &bubblesSelector.Model{
				Data: []interface{}{
					strava.StravaActivityTypeRun,
					strava.StravaActivityTypeRide,
				},
				HeaderFunc:     bubblesSelector.DefaultHeaderFuncWithAppend("Select activity type:"),
				SelectedFunc:   bubblesSelector.DefaultSelectedFuncWithIndex("[%d]"),
				UnSelectedFunc: bubblesSelector.DefaultUnSelectedFuncWithIndex(" %d."),
				FooterFunc: func(m bubblesSelector.Model, obj interface{}, gdIndex int) string {
					return ""
				},
			}
		case 1:
			m.activityType = m.selector.Selected().(strava.StravaActivityType)
			m.selector = nil
			m.prompt = &bubblesPrompt.Model{
				Prompt:       bubblesCommon.FontColor("Input activity name: ", bubblesPrompt.ColorPrompt),
				ValidateFunc: bubblesPrompt.VFNotBlank,
			}
		case 2:
			m.activityName = m.prompt.Value()
			m.prompt = &bubblesPrompt.Model{
				Prompt: bubblesCommon.FontColor("Input activity description: ", bubblesPrompt.ColorPrompt),
			}
		case 3:
			m.activityDescription = m.prompt.Value()
			m.prompt = &bubblesPrompt.Model{
				Prompt:       bubblesCommon.FontColor("Input activity start time (DD.MM.YYYY HH:MM:SS): ", bubblesPrompt.ColorPrompt),
				ValidateFunc: strIsTime,
			}
		case 4:
			m.startTime, _ = strToTimestamp(m.prompt.Value())
			m.prompt = &bubblesPrompt.Model{
				Prompt: bubblesCommon.FontColor("Input desired speed as float: ", bubblesPrompt.ColorPrompt),
				ValidateFunc: func(val string) error {
					_, err := strconv.ParseFloat(val, 64)
					return err
				},
			}
		case 5:
			m.desiredSpeed, _ = strconv.ParseFloat(m.prompt.Value(), 64)
		}
		m.inputIndex++

		// print view once after we are done
		if m.inputIndex == 6 {
			return m, updateView
		}
	}

	var cmd tea.Cmd
	if m.prompt != nil {
		_, cmd = m.prompt.Update(msg)
		if m.prompt.Canceled() {
			log.Println(errorStyle.Render("[!] Input required"))
			m.failedInput = true
			return m, tea.Quit
		}
	} else {
		_, cmd = m.selector.Update(msg)
		if m.selector.Canceled() {
			log.Println(errorStyle.Render("[!] Input required"))
			m.failedInput = true
			return m, tea.Quit
		}
	}
	return m, cmd
}

func (m *ActivityModel) View() string {
	var view string
	if m.inputIndex > 0 {
		view += infoStyle.Render("[+] GPX file: ") + m.gpxFilePath + "\n"
	}
	if m.inputIndex > 1 {
		view += infoStyle.Render("[+] Activity type: ") + string(m.activityType) + "\n"
	}
	if m.inputIndex > 2 {
		view += infoStyle.Render("[+] Activity name: ") + m.activityName + "\n"
	}
	if m.inputIndex > 3 {
		view += infoStyle.Render("[+] Activity description: ") + m.activityDescription + "\n"
	}
	if m.inputIndex > 4 {
		view += infoStyle.Render("[+] Start time: ") + m.startTime.Format("02.01.2006 15:04:05") + "\n"
	}
	if m.inputIndex > 5 {
		view += infoStyle.Render("[+] Desired speed: ") + strconv.FormatFloat(m.desiredSpeed, 'f', 5, 64) + "\n"
	}
	if m.inputIndex < 6 {
		if m.prompt != nil {
			view += m.prompt.View()
		} else {
			view += m.selector.View()
		}
	}
	return view
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
				"email":            config.Email,
				"email_language":   "en_GB",
				"password":         config.Password,
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
	if model.activityType == strava.StravaActivityTypeRun {
		workoutType = 3
	} else {
		workoutType = 12
	}
	metadataBuffer := new(bytes.Buffer)
	metadataEncoder := json.NewEncoder(metadataBuffer)
	err = metadataEncoder.Encode(map[string]interface{}{
		"activity_name":             model.activityName,
		"activity_type":             model.activityType,
		"commute":                   false,
		"description":               model.activityDescription,
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

	if arguments.DontUpload {
		return nil
	}

	resp, err := client.R().
		SetFileReader("file", randomActivityId+"-activity.fit", activityBuffer).
		SetMultipartField("metadata", "", "application/json", metadataBuffer).
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

		if err := activity.AddRecord(&strava.Record{
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

func run(cmd *cobra.Command, args []string) {
	var err error

	// read config then read input params
	if config, err = arguments.GetConfig(); err != nil {
		log.Fatal(errorStyle.Render("[!] Unable to read config: " + err.Error()))
	}
	prog := tea.NewProgram(model)
	if err := prog.Start(); err != nil {
		log.Fatal(errorStyle.Render("[!] Error while getting input: " + err.Error()))
	}
	if model.failedInput {
		return
	}

	// create activity
	activity, err := strava.NewActivity(&strava.ActivityOptions{
		AppVersion:         config.internalMobileAppVersion,
		MobileAppVersion:   config.MobileAppVersion,
		DeviceManufacturer: config.DeviceManufacturer,
		DeviceModel:        config.DeviceModel,
		DeviceOsVersion:    strconv.Itoa(config.DeviceOsVersion),
		ActivityType:       model.activityType,
		StartTime:          model.startTime,
		DesiredSpeed:       model.desiredSpeed,
	})
	if err != nil {
		log.Fatal(errorStyle.Render("[!] Unable to create activity: " + err.Error()))
	}
	log.Println(infoStyle.Render("[+] Strava activity created"))

	// fill up the activity
	if err := readGpxIntoActivity(activity, model.gpxFilePath); err != nil {
		log.Fatal(errorStyle.Render("[!] Error while converting gpx to strava activity: ") + err.Error())
	}
	log.Println(infoStyle.Render("[+] GPX successfully converted to strava activity"))

	// construct the fit file for the activity
	fitFile, err := activity.BuildFitFile()
	if err != nil {
		log.Fatal(errorStyle.Render("[!] Error while building the fit file for the activity: " + err.Error()))
	}
	log.Println(infoStyle.Render("[+] Activity built into fit file"))

	if err := sendActivity(fitFile); err != nil {
		log.Fatal(errorStyle.Render("[!] Error while sending fit file: " + err.Error()))
	}
	log.Println(infoStyle.Render("[+] Activity uploaded, check your strava!"))
}

func init() {
	rootCmd.Flags().StringVar(&arguments.ConfigPath, "config", "config.yml", "path to config file")
	rootCmd.Flags().BoolVar(&arguments.DontUpload, "no-upload", false, "don't upload the generated activity (for debug)")
	log.SetFlags(0)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(errorStyle.Render("[!] Error running the app: " + err.Error()))
	}
}
