package tui

func (m *Model) View() string {
	// switch m.state {
	// case StateSearching:
	// 	if m.presetSlug != "" {
	// 		return fmt.Sprintf("Fetching question details for '%s'...\n", m.presetSlug)
	// 	}
	// 	return fmt.Sprintf("Searching for '%s'...\n", m.keyword)
	// case StateQuestionList:
	// 	return m.viewQuestionList()
	// case StateLanguageSelect:
	// 	return m.viewLanguageSelect()
	// case StatePathInput:
	// 	return m.viewPathInput()
	// case StateDone:
	// 	return fmt.Sprintf("✓ File generated: %s\n", m.generatedPath)
	// case StateError:
	// 	return fmt.Sprintf("Error: %v\n", m.err)
	// }
	return ""
}
