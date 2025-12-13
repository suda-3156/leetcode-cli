package phases

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/suda-3156/leetcode-cli/internal/config"
	"github.com/suda-3156/leetcode-cli/internal/file"
	"github.com/suda-3156/leetcode-cli/internal/generator"
	"github.com/suda-3156/leetcode-cli/internal/tui/model"
)

type GenerationHandler struct{}

type fileGeneratedMsg struct {
	path string
}

func (h *GenerationHandler) Enter(m *model.Model) tea.Cmd {
	m.Loading = true
	return h.generateFile(m)
}

func (h *GenerationHandler) generateFile(m *model.Model) tea.Cmd {
	return func() tea.Msg {
		// Handle backup if needed
		if m.OverwriteChoice == config.OverwriteBackup && file.FileExists(m.OutPath) {
			backupPath, err := file.CreateBackup(m.OutPath)
			if err != nil {
				return errMsg{err: fmt.Errorf("create backup: %w", err)}
			}
			// Could optionally notify user about backup path
			_ = backupPath
		}

		// Generate file content
		date := m.Config.GetCurrentDate()
		content, err := generator.GenerateFileContent(
			date,
			m.QuestionDetail.QuestionFrontendID,
			m.QuestionDetail.Title,
			m.QuestionDetail.TitleSlug,
			m.SelectedLang.Lang,
			m.SelectedLang.LangSlug,
			m.SelectedLang.Code,
		)
		if err != nil {
			return errMsg{err: fmt.Errorf("generate file content: %w", err)}
		}

		// Save file
		err = file.Save(m.OutPath, content)
		if err != nil {
			return errMsg{err: fmt.Errorf("save file: %w", err)}
		}

		return fileGeneratedMsg{path: m.OutPath}
	}
}

func (h *GenerationHandler) Update(m *model.Model, msg tea.Msg) (tea.Cmd, *PhaseType) {
	switch msg := msg.(type) {
	case fileGeneratedMsg:
		m.Loading = false
		m.OutPath = msg.path
		next := DonePhase
		return nil, &next

	case errMsg:
		m.Err = msg.err
		return tea.Quit, nil
	}

	return nil, nil
}

func (h *GenerationHandler) View(m *model.Model) string {
	return fmt.Sprintf("%s Generating file...\n", m.Spinner.View())
}
