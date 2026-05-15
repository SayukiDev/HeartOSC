package app

import (
	"HeartOSC/heart"
	"HeartOSC/service"
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
)

type pickerDoneMsg struct {
	device heart.Device
}

type pickerModel struct {
	devices  []heart.Device
	selected int
	srv      *service.Service
	err      error
}

type pickerRescanMsg struct{}

func newPickerModel(ds []heart.Device, srv *service.Service, err error) pickerModel {
	return pickerModel{
		devices: ds,
		err:     err,
		srv:     srv,
	}
}

func (m *pickerModel) Init() tea.Cmd {
	if m.srv.Conf.Device != "" {
		for _, d := range m.devices {
			if d.Addr.String() == m.srv.Conf.Device {
				return func() tea.Msg { return pickerDoneMsg{device: d} }
			}
		}
	}
	return nil
}

func (m *pickerModel) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "up", "k":
			if m.selected > 0 {
				m.selected--
			}
		case "down", "j":
			if m.selected < len(m.devices)-1 {
				m.selected++
			}
		case "enter":
			if len(m.devices) > 0 {
				d := m.devices[m.selected]
				return func() tea.Msg { return pickerDoneMsg{device: d} }
			}
		case "r":
			return func() tea.Msg {
				return pickerRescanMsg{}
			}
		}
	}
	return nil
}

func (m *pickerModel) View() string {
	var b strings.Builder

	if m.err != nil {
		b.WriteString(errorStyle.Render(fmt.Sprintf("スキャン失敗: %s", m.err)))
		b.WriteString("\n\n")
		b.WriteString(helpStyle.Render("r: 再スキャン ・ q: 終了"))
		return b.String()
	}
	if len(m.devices) == 0 {
		b.WriteString("デバイスが見つかりませんでした。")
		b.WriteString("\n\n")
		b.WriteString(helpStyle.Render("r: 再スキャン ・ q: 終了"))
		return b.String()
	}

	b.WriteString(titleStyle.Render(fmt.Sprintf("デバイスを選択 (%d 件)", len(m.devices))))
	b.WriteString("\n\n")
	for i, d := range m.devices {
		name := d.Name
		if name == "" {
			name = "(unnamed)"
		}
		line := fmt.Sprintf(" %s  %s  [RSSI %d]", d.Addr.String(), name, d.RSSI)
		if d.HaveHeartRate {
			line += " (可能な心拍計デバイス)"
			line = GreenStyle.Render(line)
		}
		if i == m.selected {
			b.WriteString(cursorStyle.Render("▶"))
			b.WriteString(selectedStyle.Render(line))
		} else {
			b.WriteString(" ")
			b.WriteString(line)
		}
		b.WriteString("\n")
	}
	b.WriteString("\n")
	b.WriteString(helpStyle.Render("Tips: \n一部の心拍計はスキャンするときに心拍計であるかどうかを識別できない場合があります\nその場合は名前で判断して接続してみてください。"))
	b.WriteString("\n\n")
	b.WriteString(helpStyle.Render("↑/↓: 移動 ・ enter: 決定 ・ r: 再スキャン ・ q: 終了"))
	return b.String()
}
