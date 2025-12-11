package tui

import (
	"fmt"
	"strings"

	"github.com/suda-3156/leetcode-cli/internal/tui/styles"
)

func (m *Model) View() string {
	// Show error if present
	if m.err != nil {
		return m.viewError()
	}

	switch m.phase {
	case InitialPhase:
		return m.viewInitialPhase()
	case DecideQuestionPhase:
		return m.viewDecideQuestionPhase()
	case DecideLanguagePhase:
		return m.viewDecideLanguagePhase()
	case DecidePathPhase:
		return m.viewDecidePathPhase()
	case GenerationPhase:
		return m.viewGenerationPhase()
	case DonePhase:
		return m.viewDonePhase()
	}
	return ""
}

func (m *Model) viewError() string {
	header := styles.StyleErrorHeader.Render("Error")
	return fmt.Sprintf("%s\n\n%v\n", header, m.err)
}

func (m *Model) viewInitialPhase() string {
	return fmt.Sprintf("%s Loading configuration...\n", m.spinner.View())
}

func (m *Model) viewDecideQuestionPhase() string {
	// Still loading questions
	if m.loading {
		return fmt.Sprintf("%s Searching for '%s'...\n", m.spinner.View(), m.keyword)
	}

	var sb strings.Builder
	sb.WriteString("Select a question:\n\n")

	for i, q := range m.questions {
		cursor := "  "
		if m.cursor == i {
			cursor = styles.StyleActive.Render("> ")
		}

		premium := ""
		if q.PaidOnly {
			premium = " [Premium]"
		}

		difficulty := strings.ToUpper(q.Difficulty)
		line := fmt.Sprintf("%s. %s (%s)%s",
			q.QuestionFrontendID,
			q.Title,
			difficulty,
			premium,
		)

		if m.cursor == i {
			line = styles.StyleActive.Render(line)
		}

		sb.WriteString(fmt.Sprintf("%s%s\n", cursor, line))
	}

	sb.WriteString("\n↑/k: up • ↓/j: down • enter: select • q/esc: quit\n")
	return sb.String()
}

func (m *Model) viewDecideLanguagePhase() string {
	// Still loading question detail
	if m.loading {
		slug := ""
		if m.selectedQ != nil {
			slug = m.selectedQ.TitleSlug
		} else if m.config != nil && m.config.TitleSlug != "" {
			slug = m.config.TitleSlug
		}
		return fmt.Sprintf("%s Fetching question details for '%s'...\n", m.spinner.View(), slug)
	}

	// Still waiting for language list to be set
	if m.loading {
		return fmt.Sprintf("%s Loading languages...\n", m.spinner.View())
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Select a language for %q:\n\n", m.questionDetail.Title))

	for i, lang := range m.languages {
		cursor := "  "
		line := lang.Lang

		if m.cursor == i {
			cursor = styles.StyleActive.Render("> ")
			line = styles.StyleActive.Render(line)
		}

		sb.WriteString(fmt.Sprintf("%s%s\n", cursor, line))
	}

	sb.WriteString("\n↑/k: up • ↓/j: down • enter: select • q/esc: quit\n")
	return sb.String()
}

func (m *Model) viewDecidePathPhase() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Output path for %q (%s):\n\n",
		m.questionDetail.Title,
		m.selectedLang.Lang,
	))
	sb.WriteString(m.textInput.View())
	sb.WriteString("\n\nenter: confirm • esc: quit\n")
	return sb.String()
}

func (m *Model) viewGenerationPhase() string {
	return fmt.Sprintf("%s Generating file...\n", m.spinner.View())
}

func (m *Model) viewDonePhase() string {
	return fmt.Sprintf("✓ File generated: %s\n", m.generatedPath)
}
