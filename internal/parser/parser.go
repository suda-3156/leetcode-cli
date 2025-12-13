package parser

type Parser interface {
	ExtractTypeDefinition() (string, error)
	// ExtractSolutionFuncName() (string, error) // TODO: Use parser to extract function name, not regex
	GenerateImportStatement() (string, error) // NOTE: As for now, unused; Need to resolve existing import statements
	Close()
}

func NewParser(langSlug, codeSnippet string) (Parser, error) {
	switch langSlug {
	case "python3":
		return newPython3Parser(codeSnippet)
	// case "go":
	// 	return newGoParser()
	default:
		return newDefaultParser()
	}
}
