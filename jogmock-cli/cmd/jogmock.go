package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"gopkg.in/yaml.v2"

	bubblesCommon "github.com/mritd/bubbles/common"
	bubblesSelector "github.com/mritd/bubbles/selector"
	"github.com/renbou/jogmock/activities"
	bubblesAutoPrompt "github.com/renbou/jogmock/jogmock-cli/pkg/bubbles/autoprompt"
	bubblesPrompt "github.com/renbou/jogmock/jogmock-cli/pkg/bubbles/prompt"
	"github.com/renbou/jogmock/strava-mock/stravapi"
	"github.com/spf13/cobra"
)

type activityConfig struct {
	CommonSpeed     *activities.SpeedOptions `yaml:"common_speed_options"`
	RareSpeed       *activities.SpeedOptions `yaml:"rare_speed_options"`
	RareSpeedChance float64                  `yaml:"rare_speed_chance"`
	FadeDuration    int                      `yaml:"fade_duration"`
}

type UserConfig struct {
	StravaConfig       *stravapi.ApiConfig `yaml:"strava"`
	RunActivityConfig  *activityConfig     `yaml:"run_activity"`
	RideActivityConfig *activityConfig     `yaml:"ride_activity"`
}

// Arguments represents the possible commmand-line arguments
type Arguments struct {
	ConfigPath string
}

// LoadConfig reads and returns the config defined by args
func (args *Arguments) LoadConfig() (*UserConfig, error) {
	configFile, err := os.Open(args.ConfigPath)
	if err != nil {
		return nil, err
	}
	decoder := yaml.NewDecoder(configFile)

	config := new(UserConfig)
	if err := decoder.Decode(config); err != nil {
		return nil, err
	}

	if config.StravaConfig == nil {
		return nil, errors.New("currently only strava is supported so it must exist in the config")
	}
	return config, nil
}

type viewer interface {
	View() string
}

type modelStep struct {
	model viewer
	after func(value interface{})
}

func (s modelStep) Update(msg tea.Msg) (cmd tea.Cmd) {
	if msg == bubblesCommon.DONE {
		if s.after != nil {
			switch model := s.model.(type) {
			case interface{ Value() string }:
				s.after(model.Value())
			case interface{ Selected() interface{} }:
				s.after(model.Selected())
			default:
				s.after(nil)
			}
		}
		return nil
	} else {
		switch model := s.model.(type) {
		case *bubblesPrompt.Model:
			_, cmd = model.Update(msg)
		case *bubblesSelector.Model:
			_, cmd = model.Update(msg)
		case *bubblesAutoPrompt.Model:
			_, cmd = model.Update(msg)
		}
	}
	return
}

func (s modelStep) View() string {
	return s.model.View()
}

type simpleModel interface {
	Update(tea.Msg) tea.Cmd
	View() string
}

type ActivityModel struct {
	config      *UserConfig
	options     activities.ActivityOptions
	gpxFilePath string
	steps       []simpleModel
	index       int
}

func strToTimestamp(val string) (time.Time, error) {
	t, err := time.Parse("02.01.2006 15:04:05", val)
	return t, err
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

const (
	OkPrefix   = "[+]"
	ErrPrefix  = "[-]"
	ColorInfo  = "2"
	ColorWarn  = "#ed6c02"
	ColorError = "#d32f2f"
)

func NewActivityModel(config *UserConfig) *ActivityModel {
	model := &ActivityModel{config: config}
	model.steps = []simpleModel{
		modelStep{
			&bubblesAutoPrompt.Model{
				Prompt:            bubblesCommon.FontColor("Path to GPX file: ", bubblesPrompt.ColorPrompt),
				ValidateOkPrefix:  OkPrefix,
				ValidateErrPrefix: ErrPrefix,
			},
			func(value interface{}) {
				_, model.gpxFilePath = bubblesAutoPrompt.UserExpand(value.(string))
			},
		},
		modelStep{
			&bubblesSelector.Model{
				Data:           []interface{}{"Run", "Ride"},
				HeaderFunc:     bubblesSelector.DefaultHeaderFuncWithAppend("Type:"),
				SelectedFunc:   bubblesSelector.DefaultSelectedFuncWithIndex("[%d]"),
				UnSelectedFunc: bubblesSelector.DefaultUnSelectedFuncWithIndex(" %d."),
				FooterFunc: func(m bubblesSelector.Model, obj interface{}, gdIndex int) string {
					return ""
				},
				FinishedFunc: func(selected interface{}) string {
					return bubblesCommon.FontColor(fmt.Sprintf(
						"%s Type: %v\n", OkPrefix, selected),
						bubblesSelector.ColorFinished)
				},
			},
			func(value interface{}) {
				text := value.(string)
				var activityCfg *activityConfig
				if text == "Run" {
					model.options.Type = activities.RunActivity
					activityCfg = model.config.RunActivityConfig
				} else {
					model.options.Type = activities.RideActivity
					activityCfg = model.config.RideActivityConfig
				}

				if activityCfg != nil {
					model.options.CommonSpeed = activityCfg.CommonSpeed
					model.options.RareSpeed = activityCfg.RareSpeed
					model.options.RareSpeedChance = activityCfg.RareSpeedChance
					model.options.FadeDuration = time.Duration(activityCfg.FadeDuration) * time.Second
				}
			},
		},
		modelStep{
			&bubblesPrompt.Model{
				Prompt:            bubblesCommon.FontColor("Name: ", bubblesPrompt.ColorPrompt),
				ValidateFunc:      bubblesPrompt.VFNotBlank,
				ValidateOkPrefix:  OkPrefix,
				ValidateErrPrefix: ErrPrefix,
			},
			func(value interface{}) {
				model.options.Name = value.(string)
			},
		},
		modelStep{
			&bubblesPrompt.Model{
				Prompt:            bubblesCommon.FontColor("Description: ", bubblesPrompt.ColorPrompt),
				ValidateOkPrefix:  OkPrefix,
				ValidateErrPrefix: ErrPrefix,
			},
			func(value interface{}) {
				model.options.Description = value.(string)
			},
		},
		modelStep{
			&bubblesPrompt.Model{
				Prompt: bubblesCommon.FontColor("Start time (DD.MM.YYYY HH:MM:SS): ",
					bubblesPrompt.ColorPrompt),
				ValidateFunc:      strIsTime,
				ValidateOkPrefix:  OkPrefix,
				ValidateErrPrefix: ErrPrefix,
			},
			func(value interface{}) {
				model.options.Start, _ = strToTimestamp(value.(string))
			},
		},
		modelStep{
			&bubblesPrompt.Model{
				Prompt: bubblesCommon.FontColor("Desired speed (km/h) as float: ", bubblesPrompt.ColorPrompt),
				ValidateFunc: func(val string) error {
					_, err := strconv.ParseFloat(val, 64)
					return err
				},
				ValidateOkPrefix:  OkPrefix,
				ValidateErrPrefix: ErrPrefix,
			},
			func(value interface{}) {
				model.options.DesiredSpeed, _ = strconv.ParseFloat(value.(string), 64)
			},
		},
	}
	return model
}

func (m *ActivityModel) Init() tea.Cmd {
	return nil
}

var updateViewMsg tea.Msg = "UPDATE_VIEW"

func updateView() tea.Msg {
	return updateViewMsg
}

func (m *ActivityModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// otherwise if simply reading values
	switch msg {
	case updateViewMsg:
		if m.index == len(m.steps) {
			return m, tea.Quit
		}
		return m, nil
	case bubblesCommon.DONE:
		m.steps[m.index].Update(bubblesCommon.DONE)
		m.index++
		if m.index < len(m.steps) {
			// initialize the next prompt
			m.steps[m.index].Update(nil)
		}
		return m, updateView
	}
	return m, m.steps[m.index].Update(msg)
}

func (m *ActivityModel) View() string {
	until := m.index
	if until == len(m.steps) {
		until -= 1
	}

	var view string
	for i := 0; i <= until; i++ {
		view += m.steps[i].View()
	}
	return view
}

// func readGpxIntoActivity(activity *strava.StravaActivity, gpxFilePath string) error {
// 	// read and unmarshal the actual file
// 	gpxFile, err := os.Open(gpxFilePath)
// 	if err != nil {
// 		return err
// 	}
// 	defer gpxFile.Close()

// 	bytes, err := ioutil.ReadAll(gpxFile)
// 	if err != nil {
// 		return err
// 	}

// 	gpx := new(GPX)
// 	if err := xml.Unmarshal(bytes, &gpx); err != nil {
// 		return err
// 	}

// 	// actually convert the gpx to a strava activity
// 	for _, trackPart := range gpx.Track.TrackSegment.TrackParts {
// 		lat, err := strconv.ParseFloat(trackPart.Lat, 64)
// 		if err != nil {
// 			return err
// 		}
// 		lon, err := strconv.ParseFloat(trackPart.Lon, 64)
// 		if err != nil {
// 			return err
// 		}

// 		if err := activity.AddRecord(&strava.Record{
// 			Lat:      lat,
// 			Lon:      lon,
// 			Altitude: trackPart.Elevation,
// 		}); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

var (
	rootCmd = &cobra.Command{
		Use:   "jogmock",
		Short: "jogmock is a sophisticated activity faker for Strava",
		Long:  "a fake Strava client built using the knowledge gained from reverse engineering the mobile app",
		Run:   run,
	}
	arguments Arguments
)

func run(cmd *cobra.Command, args []string) {
	config, err := arguments.LoadConfig()
	if err != nil {
		fmt.Println(bubblesCommon.FontColor(ErrPrefix+" Unable to load config: "+err.Error(), ColorError))
		return
	}

	model := NewActivityModel(config)
	prog := tea.NewProgram(model)
	if err := prog.Start(); err != nil {
		fmt.Println(ErrPrefix+" "+err.Error(), ColorError)
		return
	}

	fmt.Printf("%+v\n", model.options)

	// activity, err := activities.NewActivity(&activities.ActivityOptions{})

	// // create activity
	// activity, err := strava.NewActivity(&strava.ActivityOptions{
	// 	AppVersion:         config.internalMobileAppVersion,
	// 	MobileAppVersion:   config.MobileAppVersion,
	// 	DeviceManufacturer: config.DeviceManufacturer,
	// 	DeviceModel:        config.DeviceModel,
	// 	DeviceOsVersion:    strconv.Itoa(config.DeviceOsVersion),
	// 	ActivityType:       model.activityType,
	// 	StartTime:          model.startTime,
	// 	DesiredSpeed:       model.desiredSpeed,
	// })
	// if err != nil {
	// 	prettylog.Fatal("Unable to create activity: " + err.Error())
	// }
	// prettylog.Info("Activity created")
	// log.Println(infoStyle.Render(""))

	// // fill up the activity
	// if err := readGpxIntoActivity(activity, model.gpxFilePath); err != nil {
	// 	log.Fatal(errorStyle.Render("[!] Error while converting gpx to strava activity: ") + err.Error())
	// }
	// log.Println(infoStyle.Render("[+] GPX successfully converted to strava activity"))

	// // construct the fit file for the activity
	// fitFile, err := activity.BuildFitFile()
	// if err != nil {
	// 	log.Fatal(errorStyle.Render("[!] Error while building the fit file for the activity: " + err.Error()))
	// }
	// log.Println(infoStyle.Render("[+] Activity built into fit file"))

	// if err := sendActivity(fitFile); err != nil {
	// 	log.Fatal(errorStyle.Render("[!] Error while sending fit file: " + err.Error()))
	// }
	// log.Println(infoStyle.Render("[+] Activity uploaded, check your strava!"))
}

func init() {
	rootCmd.Flags().StringVar(&arguments.ConfigPath, "config", "config.yml", "path to config file")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(ErrPrefix+" "+err.Error(), ColorError)
	}
}
