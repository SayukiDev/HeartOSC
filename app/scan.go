package app

import (
	"HeartOSC/heart"
	"HeartOSC/service"
	"strings"
	"time"

	"charm.land/bubbles/v2/progress"
	tea "charm.land/bubbletea/v2"
)

const defaultScanDuration = 10 * time.Second

type scanFinishedMsg struct {
	devices []heart.Device
	err     error
}

type scanTickMsg time.Time

type scanModel struct {
	progress  progress.Model
	percent   float64
	startedAt time.Time
	duration  time.Duration
	finished  bool
}

func newScanModel(srv *service.Service) scanModel {
	p := progress.New(
		progress.WithDefaultBlend(),
		progress.WithoutPercentage(),
	)
	p.SetWidth(50)
	return scanModel{
		progress: p,
		duration: defaultScanDuration,
	}
}

func (m *scanModel) Init() tea.Cmd {
	m.startedAt = time.Now()
	return tea.Batch(m.scanCmd(), m.tickCmd())
}

func (m *scanModel) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		w := msg.Width - 8
		if w > 80 {
			w = 80
		}
		if w < 10 {
			w = 10
		}
		m.progress.SetWidth(w)

	case scanFinishedMsg:
		m.finished = true
		m.percent = 1

	case scanTickMsg:
		if m.finished {
			return nil
		}
		elapsed := time.Since(m.startedAt)
		p := float64(elapsed) / float64(m.duration)
		if p > 1 {
			p = 1
		}
		m.percent = p
		return m.tickCmd()
	}
	return nil
}

func (m *scanModel) View() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("Bluetoothデバイスをスキャン中..."))
	b.WriteString("\n\n")
	b.WriteString(m.progress.ViewAs(m.percent))
	b.WriteString("\n\n")
	b.WriteString(helpStyle.Render("中止: q / ctrl+c"))
	return b.String()
}

func (m *scanModel) scanCmd() tea.Cmd {
	duration := m.duration
	return func() tea.Msg {
		devices, err := heart.ScanDeviceWithTimeout(duration)
		return scanFinishedMsg{devices: devices, err: err}
	}
}

func (m scanModel) tickCmd() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return scanTickMsg(t)
	})
}
