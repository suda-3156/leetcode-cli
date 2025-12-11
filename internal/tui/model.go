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

type Input struct {
	ConfigPath string
	Keyword    string
	TitleSlug  string
	LangSlug   string
	OutPath    string
}

type Model struct {
	err     error
	spinner spinner.Model
	phase   Phase
	loading bool

	config *config.Config

	input          Input
	questions      []api.Question
	selectedQ      *api.Question
	questionDetail *api.QuestionDetail
	languages      []api.CodeSnippet
	selectedLang   *api.CodeSnippet
	cursor         int
	textInput      textinput.Model
	client         *api.Client
	outPath        string
}

func New(input Input) Model {
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
		input:     input,
		textInput: ti,
		client:    api.NewClient(),
	}

	return m
}

func (m *Model) GetGeneratedPath() string {
	return m.outPath
}

func (m *Model) GetError() error {
	return m.err
}

func (m *Model) HasError() bool {
	return m.err != nil
}
