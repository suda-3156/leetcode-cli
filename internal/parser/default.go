package parser

type defaultParser struct{}

func newDefaultParser() (*defaultParser, error) {
	return &defaultParser{}, nil
}

func (p *defaultParser) ExtractTypeDefinition() (string, error) {
	return "", nil
}

func (p *defaultParser) GenerateImportStatement() (string, error) {
	return "", nil
}

func (p *defaultParser) Close() {}
