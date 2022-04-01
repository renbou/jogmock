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
	selectorBubble "github.com/mritd/bubbles/selector"
	"github.com/renbou/jogmock/activities"
	autoPromptBubble "github.com/renbou/jogmock/jogmock-cli/pkg/bubbles/autoprompt"
	promptBubble "github.com/renbou/jogmock/jogmock-cli/pkg/bubbles/prompt"
	stravaBubble "github.com/renbou/jogmock/jogmock-cli/pkg/bubbles/strava"
	"github.com/renbou/jogmock/strava-mock/stravapi"
	"github.com/spf13/cobra"
)

type activityConfig struct {
	CommonSpeed     *activities.SpeedOptions `yaml:"common_speed_options"`
	RareSpeed       *activities.SpeedOptions `yaml:"rare_speed_options"`
	RareSpeedChance float64                  `yaml:"rare_speed_chance"`
	FadeDuration    int                      `yaml:"fade_duration"`
	FadeFraction    float64                  `yaml:"fade_fraction"`
}

type UserConfig struct {
	StravaConfig       *stravapi.ApiConfig `yaml:"strava"`
	RunActivityConfig  *activityConfig     `yaml:"run_activity"`
	RideActivityConfig *activityConfig     `yaml:"ride_activity"`
}

// Arguments represents the possible commmand-line arguments
type Arguments struct {
	ConfigPath string
	OutputPath string
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

func (args *Arguments) SaveConfig(config *UserConfig) error {
	configFile, err := os.Create(args.ConfigPath)
	if err != nil {
		return err
	}

	encoder := yaml.NewEncoder(configFile)
	return encoder.Encode(config)
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
			case interface{ Value() interface{} }:
				s.after(model.Value())
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
		case *promptBubble.Model:
			_, cmd = model.Update(msg)
		case *selectorBubble.Model:
			_, cmd = model.Update(msg)
		case *autoPromptBubble.Model:
			_, cmd = model.Update(msg)
		case tea.Model:
			_, cmd = model.Update(msg)
		default:
			panic(fmt.Sprintf("unknown model: %v", s.model))
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
	args        *Arguments
	config      *UserConfig
	options     activities.ActivityOptions
	gpxFilePath string
	steps       []simpleModel
	index       int
}

func strToTimestamp(val string) (time.Time, error) {
	t, err := time.Parse("02.01.2006 15:04:05", val)
	if err != nil {
		return time.Time{}, err
	}

	// offset the time by local timezone offset before sending it as timestamp
	_, offset := time.Now().Zone()
	return t.Add(-time.Second * time.Duration(offset)), nil
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
	model := &ActivityModel{
		args:   &arguments,
		config: config,
	}
	model.steps = []simpleModel{
		modelStep{
			&autoPromptBubble.Model{
				Prompt:            bubblesCommon.FontColor("Path to GPX file: ", promptBubble.ColorPrompt),
				ValidateOkPrefix:  OkPrefix,
				ValidateErrPrefix: ErrPrefix,
			},
			func(value interface{}) {
				_, model.gpxFilePath = autoPromptBubble.UserExpand(value.(string))
			},
		},
		modelStep{
			&selectorBubble.Model{
				Data:           []interface{}{"Run", "Ride"},
				HeaderFunc:     selectorBubble.DefaultHeaderFuncWithAppend("Type:"),
				SelectedFunc:   selectorBubble.DefaultSelectedFuncWithIndex("[%d]"),
				UnSelectedFunc: selectorBubble.DefaultUnSelectedFuncWithIndex(" %d."),
				FooterFunc: func(m selectorBubble.Model, obj interface{}, gdIndex int) string {
					return ""
				},
				FinishedFunc: func(selected interface{}) string {
					return bubblesCommon.FontColor(OkPrefix+" Type: ", selectorBubble.ColorFinished) +
						fmt.Sprintln(selected)
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
					model.options.FadeFraction = activityCfg.FadeFraction
				}
			},
		},
		modelStep{
			&promptBubble.Model{
				Prompt:            bubblesCommon.FontColor("Name: ", promptBubble.ColorPrompt),
				ValidateFunc:      promptBubble.VFNotBlank,
				ValidateOkPrefix:  OkPrefix,
				ValidateErrPrefix: ErrPrefix,
			},
			func(value interface{}) {
				model.options.Name = value.(string)
			},
		},
		modelStep{
			&promptBubble.Model{
				Prompt:            bubblesCommon.FontColor("Description: ", promptBubble.ColorPrompt),
				ValidateOkPrefix:  OkPrefix,
				ValidateErrPrefix: ErrPrefix,
			},
			func(value interface{}) {
				model.options.Description = value.(string)
			},
		},
		modelStep{
			&promptBubble.Model{
				Prompt: bubblesCommon.FontColor("Start time (DD.MM.YYYY HH:MM:SS): ",
					promptBubble.ColorPrompt),
				ValidateFunc:      strIsTime,
				ValidateOkPrefix:  OkPrefix,
				ValidateErrPrefix: ErrPrefix,
			},
			func(value interface{}) {
				model.options.Start, _ = strToTimestamp(value.(string))
			},
		},
		modelStep{
			&promptBubble.Model{
				Prompt: bubblesCommon.FontColor("Desired speed (km/h) as float: ", promptBubble.ColorPrompt),
				ValidateFunc: func(val string) error {
					_, err := strconv.ParseFloat(val, 64)
					if err != nil {
						return errors.New("input speed as a float, error: " + err.Error())
					}
					return nil
				},
				ValidateOkPrefix:  OkPrefix,
				ValidateErrPrefix: ErrPrefix,
			},
			func(value interface{}) {
				model.options.DesiredSpeed, _ = strconv.ParseFloat(value.(string), 64)
			},
		},
		modelStep{
			&stravaBubble.Model{
				ActivityOptions: &model.options,
				ApiConfig:       model.config.StravaConfig,
				GpxFilePath:     &model.gpxFilePath,
				OutputPath:      &model.args.OutputPath,
			},
			func(value interface{}) {
				if value != nil {
					model.config.StravaConfig = value.(*stravapi.ApiConfig)
				}
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
		return m, tea.Quit
	case bubblesCommon.DONE:
		m.steps[m.index].Update(bubblesCommon.DONE)
		m.index++
		if m.index < len(m.steps) {
			// initialize the next prompt
			return m, m.steps[m.index].Update(nil)
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

	if err := arguments.SaveConfig(model.config); err != nil {
		fmt.Println(ErrPrefix+" Failed to save new config: "+err.Error(), ColorError)
		return
	}
}

func init() {
	rootCmd.Flags().StringVar(&arguments.ConfigPath, "config", "config.yml", "path to config file")
	rootCmd.Flags().StringVar(&arguments.OutputPath, "output", "",
		"path where to output the created route instead of uploading")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(ErrPrefix+" "+err.Error(), ColorError)
	}
}
