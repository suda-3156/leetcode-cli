package main

import (
	"fmt"
	"os"
	"strings"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_go "github.com/tree-sitter/tree-sitter-go/bindings/go"
	tree_sitter_python "github.com/tree-sitter/tree-sitter-python/bindings/go"
)

func main() {
	args := os.Args

	if len(args) < 3 {
		fmt.Println("Usage: go run ./scripts/print_node.go <language> <filename>")
		return
	}

	langSlug := args[1]
	filename := args[2]
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	fmt.Printf("Content of %s:\n%s\n", filename, string(content))

	var (
		p = tree_sitter.NewParser()
		l *tree_sitter.Language
	)
	defer p.Close()
	switch langSlug {
	case "python", "python3":
		l = tree_sitter.NewLanguage(tree_sitter_python.Language())
	case "go", "golang":
		l = tree_sitter.NewLanguage(tree_sitter_go.Language())
	default:
		fmt.Printf("Unsupported language: %s\n", langSlug)
		return
	}

	if err := p.SetLanguage(l); err != nil {
		fmt.Printf("Failed to set tree-sitter language: %v\n", err)
		return
	}

	tree := p.Parse(content, nil)
	defer tree.Close()

	root := tree.RootNode()

	fmt.Println("\nDetailed Node Information:")
	printNode(root, content, 0)
}

func printNode(n *tree_sitter.Node, src []byte, indent int) {
	prefix := strings.Repeat("  ", indent)

	text := n.Utf8Text(src)
	text = strings.ReplaceAll(text, "\n", "\\n")
	if len(text) > 30 {
		text = text[:27] + "..."
	}
	fmt.Printf(
		"%sKind: %s, Start: %v, End: %v, Text: %s\n",
		prefix,
		n.Kind(),
		n.StartPosition(),
		n.EndPosition(),
		text,
	)
	var i uint
	for i = 0; i < n.ChildCount(); i++ { //nolint:intrange // ChildCount returns uint32
		child := n.Child(i)
		printNode(child, src, indent+1)
	}
}
