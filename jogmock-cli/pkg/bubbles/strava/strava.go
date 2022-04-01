package strava

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	bubblesCommon "github.com/mritd/bubbles/common"
	"github.com/renbou/jogmock/activities"
	promptBubble "github.com/renbou/jogmock/jogmock-cli/pkg/bubbles/prompt"
	"github.com/renbou/jogmock/strava-mock/stravapi"
)

const (
	OkPrefix   = "[+]"
	ErrPrefix  = "[-]"
	ColorInfo  = "2"
	ColorWarn  = "#ed6c02"
	ColorError = "#d32f2f"
)

type Model struct {
	ActivityOptions *activities.ActivityOptions
	GpxFilePath     *string
	OutputPath      *string
	ApiConfig       *stravapi.ApiConfig

	apiClient *stravapi.ApiClient
	activity  *activities.Activity

	builtActivity     bool
	initializedClient bool
	uploading         bool
	uploaded          bool
	recaptchaPrompt   *promptBubble.Model
	err               error
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) buildActivity() error {
	// read and unmarshal the actual file
	gpxFile, err := os.Open(*m.GpxFilePath)
	if err != nil {
		return err
	}
	defer gpxFile.Close()

	b, err := ioutil.ReadAll(gpxFile)
	if err != nil {
		return err
	}

	return m.activity.BuildFromGPX(b)
}

type msg int

const (
	viewErrMsg msg = iota
	saveActivityMsg
	initApiClientMsg
	uploadActivityMsg
	uploadedActivityMsg
)

func viewErr() tea.Msg {
	return viewErrMsg
}

func saveActivity() tea.Msg {
	return saveActivityMsg
}

func initApiClient() tea.Msg {
	return initApiClientMsg
}

func uploadActivity() tea.Msg {
	return uploadActivityMsg
}

func uploadedActivity() tea.Msg {
	return uploadedActivityMsg
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.builtActivity && m.err == nil {
		m.activity, m.err = activities.NewActivity(m.ActivityOptions)
		if m.err != nil {
			return m, viewErr
		}

		m.err = m.buildActivity()
		if m.err != nil {
			return m, viewErr
		}

		m.builtActivity = true

		if *m.OutputPath != "" {
			return m, saveActivity
		}
		return m, initApiClient
	}

	switch msg {
	case viewErrMsg:
		return m, tea.Quit
	case saveActivityMsg:
		file, err := os.Create(*m.OutputPath)
		if err != nil {
			m.err = err
			return m, viewErr
		}

		encoder := json.NewEncoder(file)
		records := m.activity.Records()
		m.err = encoder.Encode(records)
		if m.err != nil {
			return m, viewErr
		}
		return m, bubblesCommon.Done
	case initApiClientMsg:
		m.apiClient, m.err = stravapi.NewClient(m.ApiConfig)
		if m.err != nil {
			return m, viewErr
		}
		m.initializedClient = true
		m.uploading = true
		return m, uploadActivity
	case uploadActivityMsg:
		m.err = m.apiClient.UploadActivity(m.activity)
		if m.err != nil {
			m.uploading = false
			if m.err == stravapi.ErrUnauthorized {
				m.err = nil
				m.recaptchaPrompt = &promptBubble.Model{
					Prompt:            bubblesCommon.FontColor("reCAPTCHA token (from mobile app): ", promptBubble.ColorPrompt),
					ValidateOkPrefix:  OkPrefix,
					ValidateErrPrefix: ErrPrefix,
				}
				m.recaptchaPrompt.Update(nil)
				return m, nil
			} else {
				return m, viewErr
			}
		}
		m.uploaded = true
		m.uploading = false
		return m, uploadedActivity
	case uploadedActivityMsg:
		return m, bubblesCommon.Done
	}

	if m.recaptchaPrompt != nil {
		_, cmd := m.recaptchaPrompt.Update(msg)
		if cmd != nil && cmd() == bubblesCommon.DONE {
			if err := m.apiClient.Authorize(m.recaptchaPrompt.Value()); err != nil {
				m.err = err
				return m, viewErr
			}
			m.uploading = true
			return m, uploadActivity
		}
		return m, cmd
	}
	return m, nil
}

func (m *Model) View() string {
	var lines []string
	if m.builtActivity {
		lines = append(lines,
			bubblesCommon.FontColor(OkPrefix+" Constructed activity from GPX", ColorInfo))
	}
	if m.initializedClient {
		lines = append(lines,
			bubblesCommon.FontColor(OkPrefix+" Initialized Strava API client", ColorInfo))
	}

	if m.recaptchaPrompt != nil {
		lines = append(lines, bubblesCommon.FontColor(ErrPrefix+" Api token expired or invalid, authorization required", ColorWarn))
		lines = append(lines, m.recaptchaPrompt.View())
	}

	if m.uploading {
		lines = append(lines,
			bubblesCommon.FontColor(OkPrefix+" Uploading activity to Strava...", ColorInfo))
	}
	if m.uploaded {
		lines = append(lines,
			bubblesCommon.FontColor(OkPrefix+" Uploaded activity, check your Strava!", ColorInfo))
	}

	if m.err != nil {
		lines = append(lines, bubblesCommon.FontColor(ErrPrefix+" Error: "+m.err.Error(), ColorError))
	}
	lines = append(lines, "")
	return strings.Join(lines, "\n")
}

func (m *Model) Value() interface{} {
	if m.apiClient != nil {
		return m.apiClient.Config()
	}
	return nil
}
