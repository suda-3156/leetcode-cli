package tui

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/suda-3156/leetcode-cli/internal/api"
	"github.com/suda-3156/leetcode-cli/internal/config"
	"github.com/suda-3156/leetcode-cli/internal/tui/styles"
)

type Phase int

const (
	InitialPhase Phase = iota
	DecideQuestionPhase
	DecideLanguagePhase
	DecidePathPhase
	GenerationPhase
	DonePhase
)

type Model struct {
	err     error
	spinner spinner.Model
	phase   Phase
	loading bool

	config *config.Config

	keyword        string
	questions      []api.Question
	selectedQ      *api.Question
	questionDetail *api.QuestionDetail
	languages      []api.CodeSnippet
	selectedLang   *api.CodeSnippet
	cursor         int
	textInput      textinput.Model
	client         *api.Client

	generatedPath string
}

func New(keyword string, config *config.Config) Model {
	ti := textinput.New()
	ti.Placeholder = "Enter output path..."
	ti.CharLimit = 256
	ti.Width = 60

	spinner := spinner.New(
		spinner.WithSpinner(spinner.Dot),
		spinner.WithStyle(styles.StyleSpinner),
	)

	m := Model{
		spinner:   spinner,
		phase:     InitialPhase,
		loading:   true,
		config:    config,
		keyword:   keyword,
		textInput: ti,
		client:    api.NewClient(),
	}

	return m
}

func (m *Model) GetGeneratedPath() string {
	return m.generatedPath
}

func (m *Model) GetError() error {
	return m.err
}

func (m *Model) HasError() bool {
	return m.err != nil
}
