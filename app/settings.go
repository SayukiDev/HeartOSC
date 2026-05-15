package app

import (
	"HeartOSC/service"
	"fmt"
	"strconv"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

type openSettingsMsg struct{}

type settingsDoneMsg struct {
	saved bool
}

const (
	fieldOSCHost = iota
	fieldOSCPort
	fieldParameter
	fieldEnableFilter
	fieldEnableSmoothing
	fieldCount
)

type settingsModel struct {
	srv          *service.Service
	inputs       [3]textinput.Model
	randomOffset bool
	smoothing    bool
	focus        int
	validateErr  string
	saveErr      error
}

func newSettingsModel(srv *service.Service) settingsModel {
	c := srv.Conf

	host := textinput.New()
	host.Placeholder = "127.0.0.1"
	host.CharLimit = 64
	host.SetWidth(30)
	host.SetValue(c.OSCHost)

	port := textinput.New()
	port.Placeholder = "9000"
	port.CharLimit = 5
	port.SetWidth(8)
	port.SetValue(strconv.Itoa(c.OscPort))

	param := textinput.New()
	param.Placeholder = "VRCOSC/Heartrate/Value"
	param.CharLimit = 128
	param.SetWidth(40)
	param.SetValue(c.Parameter)

	host.Focus()

	m := settingsModel{
		srv:          srv,
		randomOffset: c.EnableRandomOffset,
		smoothing:    c.EnableSmoothing,
		focus:        fieldOSCHost,
	}
	m.inputs[fieldOSCHost] = host
	m.inputs[fieldOSCPort] = port
	m.inputs[fieldParameter] = param
	return m
}

func (m *settingsModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *settingsModel) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "esc":
			return func() tea.Msg { return settingsDoneMsg{saved: false} }
		case "ctrl+s":
			return m.commit()
		case "tab", "down":
			m.focusNext()
			return nil
		case "shift+tab", "up":
			m.focusPrev()
			return nil
		case "enter":
			if m.focus == fieldEnableFilter || m.focus == fieldEnableSmoothing {
				return m.commit()
			}
			m.focusNext()
			return nil
		case "space":
			switch m.focus {
			case fieldEnableFilter:
				m.randomOffset = !m.randomOffset
				return nil
			case fieldEnableSmoothing:
				m.smoothing = !m.smoothing
				return nil
			}
		}
	}

	if m.focus >= fieldOSCHost && m.focus <= fieldParameter {
		var cmd tea.Cmd
		m.inputs[m.focus], cmd = m.inputs[m.focus].Update(msg)
		return cmd
	}
	return nil
}

func (m *settingsModel) focusNext() {
	if m.focus <= fieldParameter {
		m.inputs[m.focus].Blur()
	}
	m.focus = (m.focus + 1) % fieldCount
	if m.focus <= fieldParameter {
		m.inputs[m.focus].Focus()
	}
}

func (m *settingsModel) focusPrev() {
	if m.focus <= fieldParameter {
		m.inputs[m.focus].Blur()
	}
	m.focus = (m.focus - 1 + fieldCount) % fieldCount
	if m.focus <= fieldParameter {
		m.inputs[m.focus].Focus()
	}
}

func (m *settingsModel) commit() tea.Cmd {
	host := strings.TrimSpace(m.inputs[fieldOSCHost].Value())
	portStr := strings.TrimSpace(m.inputs[fieldOSCPort].Value())
	param := strings.TrimSpace(m.inputs[fieldParameter].Value())

	if host == "" {
		m.validateErr = "OSCHost は空にできません"
		return nil
	}
	port, err := strconv.Atoi(portStr)
	if err != nil || port < 1 || port > 65535 {
		m.validateErr = "OscPort は 1〜65535 の整数で指定してください"
		return nil
	}
	if param == "" {
		m.validateErr = "Parameter は空にできません"
		return nil
	}

	m.validateErr = ""
	m.srv.Conf.OSCHost = host
	m.srv.Conf.OscPort = port
	m.srv.Conf.Parameter = param
	m.srv.Conf.EnableRandomOffset = m.randomOffset
	m.srv.Conf.EnableSmoothing = m.smoothing
	if err := m.srv.Conf.Save(); err != nil {
		m.saveErr = err
		return nil
	}
	return func() tea.Msg { return settingsDoneMsg{saved: true} }
}

func (m settingsModel) View() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("設定 (変更は次回起動時に反映)"))
	b.WriteString("\n\n")

	labels := [fieldCount]string{
		"OSCHost        ",
		"OscPort        ",
		"Parameter      ",
		"EnableRandomOffset   ",
		"EnableSmoothing",
	}

	for i := 0; i < fieldCount; i++ {
		cursor := "  "
		label := labels[i]
		if m.focus == i {
			cursor = cursorStyle.Render("▶ ")
			label = selectedStyle.Render(label)
		}
		var value string
		switch i {
		case fieldEnableFilter:
			if m.randomOffset {
				value = "[x]"
			} else {
				value = "[ ]"
			}
		case fieldEnableSmoothing:
			if m.smoothing {
				value = "[x]"
			} else {
				value = "[ ]"
			}
		default:
			value = m.inputs[i].View()
		}
		b.WriteString(fmt.Sprintf("%s%s : %s\n", cursor, label, value))
	}

	b.WriteString("\n")
	if m.validateErr != "" {
		b.WriteString(errorStyle.Render(m.validateErr))
		b.WriteString("\n")
	}
	if m.saveErr != nil {
		b.WriteString(errorStyle.Render(fmt.Sprintf("保存失敗: %s", m.saveErr)))
		b.WriteString("\n")
	}
	b.WriteString("\n")
	b.WriteString(helpStyle.Render("tab/↑↓: 移動 ・ space: トグル) ・ ctrl+s: 保存 ・ esc: キャンセル"))
	return b.String()
}
