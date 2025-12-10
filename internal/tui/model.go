package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/suda-3156/leetcode-cli/internal/api"
)

// State is an enumeration of the TUI states
type State int

const (
	StateSearching State = iota
	StateQuestionList
	StateLanguageSelect
	StatePathInput
	StateDone
	StateError
)

// Model is the Bubble Tea model
type Model struct {
	state          State
	keyword        string
	questions      []api.Question
	selectedQ      *api.Question
	questionDetail *api.QuestionDetail
	languages      []api.CodeSnippet
	selectedLang   *api.CodeSnippet
	outputPath     string
	cursor         int
	textInput      textinput.Model
	err            error
	client         *api.Client

	// Preset values set by flags
	presetSlug string
	presetLang string
	presetPath string

	// Results
	generatedPath string
}

func NewModel(keyword, slug, lang, path string) Model {
	ti := textinput.New()
	ti.Placeholder = "Enter output path..."
	ti.CharLimit = 256
	ti.Width = 60

	m := Model{
		state:      StateSearching,
		keyword:    keyword,
		client:     api.NewClient(),
		textInput:  ti,
		presetSlug: slug,
		presetLang: lang,
		presetPath: path,
	}

	return m
}

func (m Model) GetGeneratedPath() string {
	return m.generatedPath
}

func (m Model) GetError() error {
	return m.err
}

func (m Model) IsDone() bool {
	return m.state == StateDone
}

func (m Model) HasError() bool {
	return m.state == StateError
}
