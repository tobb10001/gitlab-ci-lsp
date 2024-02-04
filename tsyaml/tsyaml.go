package tsyaml

import (
	"context"

	sitter "github.com/smacker/go-tree-sitter"
	yaml_language "github.com/smacker/go-tree-sitter/yaml"
)

var parser *sitter.Parser

func init() {
	parser = sitter.NewParser()
	parser.SetLanguage(yaml_language.GetLanguage())
}

func GetTree(source []byte) (*sitter.Tree, error) {
	return parser.ParseCtx(context.Background(), nil, source)
}

func UpdateTree(tree *sitter.Tree, source []byte) (*sitter.Tree, error) {
	return parser.ParseCtx(context.Background(), tree, source)
}
