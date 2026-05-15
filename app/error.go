package app

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
)

type ErrorModel struct {
	err error
}

func newErrorModel(err error) ErrorModel {
	return ErrorModel{err: err}
}

func (m ErrorModel) Init() tea.Cmd {
	return nil
}

func (m ErrorModel) View() string {
	return fmt.Sprintf(
		"%s\n%s\n",
		errorStyle.Render(fmt.Errorf("エラー: %s", m.err).Error()),
		helpStyle.Render("終了: q / ctrl+c"),
	)
}

func (m ErrorModel) Update() (tea.Model, tea.Cmd) {
	return nil, nil
}
