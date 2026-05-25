package app

import (
	"HeartOSC/heart"
	"HeartOSC/service"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
)

const (
	defaultOSCHost = "127.0.0.1"
	defaultOSCPort = 9000
)

type runnerStartedMsg struct{}

type runnerErrorMsg struct {
	err string
}

type reConnectMsg struct{}

type runnerTickMsg time.Time

type runnerModel struct {
	device        heart.Device
	host          string
	port          int
	rate          int32
	err           string
	reconnectErr  error
	started       bool
	connected     bool
	reConnecting  bool
	connectFailed bool
	srv           *service.Service
}

func newRunnerModel(d heart.Device, srv *service.Service) runnerModel {
	return runnerModel{
		device: d,
		host:   defaultOSCHost,
		port:   defaultOSCPort,
		srv:    srv,
	}
}

func (m *runnerModel) Init() tea.Cmd {
	return tea.Batch(startCmd(m.device, m.srv))
}

func (m *runnerModel) Active() tea.Cmd {
	return func() tea.Msg {
		return runnerTickMsg(time.Now())
	}
}

func (m *runnerModel) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "r":
			m.started = false
			m.err = ""
			return func() tea.Msg {
				err := reconnectDevice(m.device)
				if err != nil {
					return runnerErrorMsg{err: err.Error()}
				}
				return runnerStartedMsg{}
			}
		case "s":
			return func() tea.Msg { return openSettingsMsg{} }
		case "d":
			m.srv.Conf.Device = ""
			err := m.srv.Conf.Save()
			if err != nil {
				m.err = err.Error()
				return func() tea.Msg {
					return runnerErrorMsg{err: err.Error()}
				}
			}
			return func() tea.Msg {
				stopCmd(m.srv)
				return pickerRescanMsg{}
			}
		case "b":
			if m.started {
				return nil
			}
			return func() tea.Msg {
				return pickerRescanMsg{}
			}
		}
	case runnerStartedMsg:
		if m.srv.Conf.Device != m.device.Addr.String() {
			m.srv.Conf.Device = m.device.Addr.String()
			err := m.srv.Conf.Save()
			if err != nil {
				m.err = err.Error()
				return func() tea.Msg {
					return runnerErrorMsg{err: err.Error()}
				}
			}
		}
		m.started = true
		return runnerTickCmd()
	case runnerErrorMsg:
		m.err = msg.err
	case runnerTickMsg:
		if m.reConnecting {
			return nil
		}
		m.rate = heart.GetHeartRate()
		m.connected = heart.IsConnected()
		if !m.connected {
			m.reConnecting = true
			return tea.Tick(1*time.Second, func(_ time.Time) tea.Msg {
				return reConnectMsg{}
			})
		}
		return runnerTickCmd()
	case reConnectMsg:
		if !m.connectFailed {
			// update status
			m.connectFailed = false
			return func() tea.Msg {
				return reConnectMsg{}
			}
		}
		err := reconnectDevice(m.device)
		if err != nil {
			m.connectFailed = true
			m.reconnectErr = err
			return tea.Tick(5*time.Second, func(_ time.Time) tea.Msg {
				return reConnectMsg{}
			})
		}
		m.reConnecting = false
		m.connected = true
		return runnerTickCmd()
	case tea.QuitMsg:
		stopCmd(m.srv)
	}
	return nil
}

func (m *runnerModel) View() string {
	var b strings.Builder

	name := m.device.Name
	if name == "" {
		name = "(unnamed)"
	}

	b.WriteString(titleStyle.Render(fmt.Sprintf("接続: %s (%s)", name, m.device.Addr.String())))
	b.WriteString("\n\n")

	if m.err != "" {
		b.WriteString(errorStyle.Render(fmt.Sprintf("エラー: %s", m.err)))
		b.WriteString("\n\n")
		if m.started {
			b.WriteString(helpStyle.Render("終了: q / ctrl+c ・ 再接続: r  ・ 切断: d"))
		} else {
			b.WriteString(helpStyle.Render("終了: q / ctrl+c ・ 再接続: r  ・ 戻る: b"))
		}
		return b.String()
	}

	if !m.started {
		b.WriteString("接続して送信を開始しています...")
	} else {
		status := ""
		if m.connected {
			status = GreenStyle.Render("connected")
		} else {
			status = RedStyle.Render("disconnected")
		}
		b.WriteString(fmt.Sprintf("心拍数: %s bpm, 接続状況: %s", RedStyle.Render(strconv.Itoa(int(m.rate))), status))
		b.WriteString("\n")
		b.WriteString(helpStyle.Render(fmt.Sprintf("OSC送信先: %s:%d", m.host, m.port)))
		if m.reConnecting {
			b.WriteString("\n\n")
			if m.connectFailed {
				b.WriteString(errorStyle.Render("再接続に失敗しました、５秒後に再試行します..."))
				b.WriteString("\n")
				b.WriteString(errorStyle.Render("エラー: ", m.reconnectErr.Error()))
			} else {
				b.WriteString(errorStyle.Render("接続されていません、再接続中..."))
			}
		}

	}
	b.WriteString("\n\n")
	if m.started {
		if m.connected {
			b.WriteString(helpStyle.Render("終了: q / ctrl+c ・ 切断: d ・ 設定: s"))
		} else {
			b.WriteString(helpStyle.Render("終了: q / ctrl+c ・ 再接続: r ・ 切断: d ・ 設定: s"))
		}
	}

	return b.String()
}

func startCmd(d heart.Device, srv *service.Service) tea.Cmd {
	return func() tea.Msg {
		if err := heart.ConnectDevice(d.Addr); err != nil {
			msg := ""
			if errors.Is(err, heart.NotFoundServiceError) {
				msg = "HeartRateサービス見つかりませんでした、心拍計デバイスではないあるいはサポートされていません。"
			}
			if msg == "" {
				msg = fmt.Errorf("connect error: %s", err).Error()
			}
			return runnerErrorMsg{err: msg}
		}
		err := srv.StartSending()
		if err != nil {
			return runnerErrorMsg{
				err: fmt.Errorf("start sending error: %s", err).Error(),
			}
		}
		return runnerStartedMsg{}
	}
}

func stopCmd(srv *service.Service) {
	srv.StopSending()
	heart.DisconnectDevice()
}

func reconnectDevice(device heart.Device) error {
	heart.DisconnectDevice()
	err := heart.ConnectDevice(device.Addr)
	return err
}

func runnerTickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return runnerTickMsg(t)
	})
}
