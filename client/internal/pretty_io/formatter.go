package pretty_io

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Formatter struct {
	m *model
	p *tea.Program
}

func (f *Formatter) PrintMessage(msg string) {
	f.p.Send(newMsg{text: msg})
}

func (f *Formatter) GetInput() <-chan string {
	return f.m.input
}

func (f *Formatter) Run() error {
	_, err := f.p.Run()
	return err
}

func (f *Formatter) Stop() {
	f.p.Quit()
}

func NewFormatter() *Formatter {
	m := initialModel()
	p := tea.NewProgram(
		m,
		tea.WithAltScreen(),       // use the full size of the terminal in its "alternate screen buffer"
		tea.WithMouseCellMotion(), // turn on mouse support, so we can track the mouse wheel
	)

	return &Formatter{p: p, m: m}
}
