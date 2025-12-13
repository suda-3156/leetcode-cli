package init

// Model is the model for the overwrite confirmation prompt
type Model struct {
	path      string
	choice    int // 0: Yes, 1: No
	quitting  bool
	confirmed bool
}
