package app

import (
	"HeartOSC/heart"
	"HeartOSC/service"

	tea "charm.land/bubbletea/v2"
)

type state int

const (
	stateStarting state = iota
	stateScanning
	stateSelecting
	stateRunning
	stateSettings
	stateError
)

type App struct {
	state      state
	start      startModel
	scan       scanModel
	picker     pickerModel
	runner     runnerModel
	settings   settingsModel
	error      ErrorModel
	srv        *service.Service
	windowSize tea.WindowSizeMsg
	chosen     *heart.Device
}

func NewModel() *App {
	return &App{
		state: stateStarting,
		start: newStartModel(),
	}
}

func (m *App) Chosen() *heart.Device {
	return m.chosen
}

func (m *App) Init() tea.Cmd {
	return m.start.Init()
}

type switchStateMsg int

func (m *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c":
			if m.srv != nil {
				m.srv.Close()
			}
			return m, tea.Quit
		case "q":
			if m.state != stateSettings {
				if m.srv != nil {
					m.srv.Close()
				}
				return m, tea.Quit
			}
		}
	case openSettingsMsg:
		m.settings = newSettingsModel(m.srv)
		m.state = stateSettings
		return m, m.settings.Init()
	case settingsDoneMsg:
		m.state = stateRunning
		return m, m.runner.Active()
	case startedMsg:
		m.state = stateScanning
		m.srv = msg.Srv
		m.scan = newScanModel(m.srv)
		return m, m.scan.Init()

	case scanFinishedMsg:
		_ = m.scan.Update(msg)
		m.picker = newPickerModel(msg.devices, m.srv, msg.err)
		m.state = stateSelecting
		return m, m.picker.Init()

	case pickerDoneMsg:
		d := msg.device
		m.chosen = &d
		m.runner = newRunnerModel(d, m.srv)
		m.state = stateRunning
		return m, m.runner.Init()

	case pickerRescanMsg:
		m.scan = newScanModel(m.srv)
		m.state = stateScanning
		return m, m.scan.Init()
	case error:
		m.error = newErrorModel(msg)
		m.state = stateError
		return m, nil
	case switchStateMsg:
		m.state = state(msg)
		return m, nil
	case tea.WindowSizeMsg:
		m.windowSize = msg
		return m, nil
	}

	var cmd tea.Cmd
	switch m.state {
	case stateStarting:
	case stateScanning:
		cmd = m.scan.Update(msg)
	case stateSelecting:
		cmd = m.picker.Update(msg)
	case stateRunning:
		cmd = m.runner.Update(msg)
	case stateSettings:
		cmd = m.settings.Update(msg)
	case stateError:
	default:
		panic("unhandled default case")
	}
	return m, cmd
}

const logo = "---------------------------------------\n═════════════════════════════════════\n █ █ █▀▀ █▀█ █▀▄ ▀█▀  ███  ███  ███\n ███ ██▄ █▀█ █▀▄  █   █ █  ▀▄▄  █  \n █ █ █▄▄ █ █ █ █  █   ███  ███  ███\n═════════════════════════════════════\n [>>>VRChat用心拍数送信プログラム<<<]\n\n---------------------------------------"

func (m *App) View() tea.View {
	var content string
	switch m.state {
	case stateStarting:
		content = m.start.View()
	case stateScanning:
		content = m.scan.View()
	case stateSelecting:
		content = m.picker.View()
	case stateRunning:
		content = m.runner.View()
	case stateSettings:
		content = m.settings.View()
	case stateError:
		content = m.error.View()
	}
	return tea.NewView(docStyle.Render(logo + "\n" + content))
}
