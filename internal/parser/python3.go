package parser

import (
	"fmt"
	"strings"

	"github.com/suda-3156/leetcode-cli/pkg/set"
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_python "github.com/tree-sitter/tree-sitter-python/bindings/go"
)

const (
	TYPE_DEF_START_LINE_PREFIX = "# Definition"
)

var typingTypes = set.New(
	"Optional",
	"List",
	"Dict",
	"Set",
	"Tuple",
	"Union",
	"Any",
	"Callable",
	"Iterator",
	"Iterable",
	"Sequence",
	"Mapping",
	"Type",
	"TypeVar",
	"Generic",
	"Protocol",
	"Literal",
	"Final",
	"ClassVar",
	"Annotated",
	"TypedDict",
	"NamedTuple",
)

type python3Parser struct {
	parser *tree_sitter.Parser
	tree   *tree_sitter.Tree
	src    []byte
	root   *tree_sitter.Node
}

func newPython3Parser(code string) (*python3Parser, error) {
	p := tree_sitter.NewParser()
	l := tree_sitter.NewLanguage(tree_sitter_python.Language())
	if err := p.SetLanguage(l); err != nil {
		return nil, fmt.Errorf("failed to set tree-sitter language: %w", err)
	}

	tree := p.Parse([]byte(code), nil)

	root := tree.RootNode()

	return &python3Parser{
		parser: p,
		tree:   tree,
		src:    []byte(code),
		root:   root,
	}, nil
}

// extractTopLevelComments extracts comments that are direct children of the root node.
func (p *python3Parser) extractTopLevelComments() []string {
	if p.root == nil {
		return nil
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

	return comments
}

// ExtractTypeDefinition extracts type definitions from top-level comments.
func (p *python3Parser) ExtractTypeDefinition() (string, error) {
	comments := p.extractTopLevelComments()
	if len(comments) == 0 || !strings.HasPrefix(comments[0], TYPE_DEF_START_LINE_PREFIX) {
		return "", nil
	}

	var typeDefLines []string
	for _, comment := range comments[1:] {
		typeDefLines = append(typeDefLines, comment[2:])
	}

	typeDef := strings.Join(typeDefLines, "\n")
	return typeDef, nil
}

// extractTypeNames extracts type names from a type node.
// This function detects:
//   - the identifier node which is a direct child of type nodes
//   - the identifier nodes which is a direct child of generic_type nodes, it should be a grandchild of type nodes
//   - the identifier nodes which is a direct child of subscript nodes
func (p *python3Parser) extractTypeNames(n *tree_sitter.Node, src []byte, result *set.Set[string]) {
	if n == nil {
		return
	}

	if n.Kind() == "type" || n.Kind() == "subscript" {
		for i := uint(0); i < n.ChildCount(); i++ { //nolint:intrange // ChildCount returns uint32
			child := n.Child(i)
			// Direct identifier of type or subscript node
			if child.Kind() == "identifier" {
				result.Add(child.Utf8Text(src))
			}

			// Generic type node -> identifier grandchild
			if child.Kind() == "generic_type" {
				for j := uint(0); j < child.ChildCount(); j++ { //nolint:intrange // ChildCount returns uint32
					grandChild := child.Child(j)
					if grandChild.Kind() == "identifier" {
						result.Add(grandChild.Utf8Text(src))
						break
					}
				}
			}
		}
	}

	// recursively search children
	for i := uint(0); i < n.ChildCount(); i++ { //nolint:intrange // ChildCount returns uint32
		child := n.Child(i)
		p.extractTypeNames(child, src, result)
	}
}

func (p *python3Parser) listTypingTypes() set.Set[string] {
	types := set.New[string]()
	p.extractTypeNames(p.root, p.src, &types)

	result := set.Union(types, typingTypes)
	return result
}

// GenerateImportStatement generates import statements for used typing types.
func (p *python3Parser) GenerateImportStatement() (string, error) {
	usedTypes := p.listTypingTypes()
	if usedTypes.Len() == 0 {
		return "", nil
	}

	var typeList []string
	for t := range usedTypes {
		typeList = append(typeList, t)
	}
	importStmt := "from typing import " + strings.Join(typeList, ", ")
	return importStmt, nil
}

// Close releases resources held by the parser.
func (p *python3Parser) Close() {
	p.tree.Close()
	p.parser.Close()
}
