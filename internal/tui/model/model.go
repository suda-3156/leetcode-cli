package model

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/suda-3156/leetcode-cli/internal/api"
	"github.com/suda-3156/leetcode-cli/internal/config"
)

// Input holds the CLI flags and arguments passed to the TUI
type Input struct {
	ConfigPath string
	Keyword    string
	TitleSlug  string
	LangSlug   string
	OutPath    string
}

// Model holds the state of the TUI application
type Model struct {
	Err     error
	Spinner spinner.Model
	Loading bool

	Config *config.Config

	Input          *Input
	Questions      []api.Question
	SelectedQ      *api.Question
	QuestionDetail *api.QuestionDetail
	Languages      []api.CodeSnippet
	SelectedLang   *api.CodeSnippet
	Cursor         int
	TextInput      textinput.Model
	Client         *api.Client
	OutPath        string

	// Overwrite decision
	OverwriteChoice int // 0: overwrite, 1: backup, 2: return to path, 3: quit
}

func (m *Model) GetGeneratedPath() string {
	return m.OutPath
}

func (m *Model) GetError() error {
	return m.Err
}

func (m *Model) HasError() bool {
	return m.Err != nil
}
