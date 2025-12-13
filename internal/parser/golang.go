package parser

import (
	"fmt"
	"strings"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_go "github.com/tree-sitter/tree-sitter-go/bindings/go"
)

const (
	GO_TYPEDEF_START_LINE_PREFIX = "/**\n * Definition"
)

type goParser struct {
	parser *tree_sitter.Parser
	tree   *tree_sitter.Tree
	src    []byte
	root   *tree_sitter.Node
}

func newGoParser(code string) (*goParser, error) {
	p := tree_sitter.NewParser()
	l := tree_sitter.NewLanguage(tree_sitter_go.Language())
	if err := p.SetLanguage(l); err != nil {
		return nil, fmt.Errorf("failed to set tree-sitter language: %w", err)
	}

	tree := p.Parse([]byte(code), nil)

	root := tree.RootNode()

	return &goParser{
		parser: p,
		tree:   tree,
		src:    []byte(code),
		root:   root,
	}, nil
}

// ExtractTypeDefinition extracts type definitions from top-level comments.
// The go parser recoginizes comment block `/** ... */` as a single comment node.
func (p *goParser) ExtractTypeDefinition() (string, error) {
	comments, err := p.extractTopLevelComments()
	if err != nil {
		return "", err
	}
	if len(comments) == 0 {
		return "", nil
	}

	var target string
	for _, comment := range comments {
		if strings.HasPrefix(comment, GO_TYPEDEF_START_LINE_PREFIX) {
			target = comment
			break
		}
	}

	if target == "" {
		return "", nil
	}

	lines := strings.Split(target, "\n")

	var typeDefLines []string
	for _, line := range lines[2 : len(lines)-1] { // Skip the first 2 lines (/**, * Definition...) and the last line (*/)
		typeDefLines = append(typeDefLines, line[3:])
	}

	typeDef := "\n" + strings.Join(typeDefLines, "\n") + "\n"
	return typeDef, nil
}

// extractTopLevelComments extracts comments that are direct children of the root node.
func (p *goParser) extractTopLevelComments() ([]string, error) {
	if p.root == nil {
		return nil, fmt.Errorf("root node is nil")
	}

	var comments []string
	for i := uint(0); i < p.root.ChildCount(); i++ { //nolint:intrange // ChildCount returns uint32
		child := p.root.Child(i)
		if child.Kind() == "comment" {
			comment := child.Utf8Text(p.src)
			if comment != "" {
				comments = append(comments, comment)
			}
		}
	}

	return comments, nil
}

// ExtractSolutionFuncName extracts the function name of the first direct child node with kind "function_declaration" of the root node.
func (p *goParser) ExtractSolutionFuncName() (string, error) {
	if p.root == nil {
		return "", fmt.Errorf("root node is nil")
	}

	for i := uint(0); i < p.root.ChildCount(); i++ { //nolint:intrange // ChildCount returns uint32
		child := p.root.Child(i)
		if child.Kind() == "function_declaration" {
			for j := uint(0); j < child.ChildCount(); j++ { //nolint:intrange // ChildCount returns uint32
				funcChild := child.Child(j)
				if funcChild.Kind() == "identifier" {
					funcName := funcChild.Utf8Text(p.src)
					return funcName, nil
				}
			}
		}
	}

	return "", fmt.Errorf("no function_declaration node found")
}

func (p *goParser) Close() {
	p.tree.Close()
	p.parser.Close()
}
