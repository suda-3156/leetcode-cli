package parser

type Parser interface {
	ExtractTypeDefinition() (string, error)
	ExtractSolutionFuncName() (string, error)
	// NOTE: As for now, unused; Need to resolve conflict with existing imports of template code.
	// Only pythonParser implements this method experimentally.
	// GenerateImportStatement() (string, error)
	Close()
}

func NewParser(langSlug, codeSnippet string) (Parser, error) {
	switch langSlug {
	case "python3":
		return newPython3Parser(codeSnippet)
	case "golang":
		return newGoParser(codeSnippet)
	default:
		return newDefaultParser()
	}
}
