// Package yamlpath, to get the YAMLPath for a given cursor position.
package yamlpath

import (
	"slices"

	"github.com/goccy/go-yaml"
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/tliron/commonlog"
	protocol "github.com/tliron/glsp/protocol_3_16"

	_ "github.com/tliron/commonlog/simple"
)

func init() {
	commonlog.SetMaxLevel(commonlog.Debug)
}

var logger = commonlog.GetLogger("yamlpath")

func hasPoint(node *sitter.Node, candidate *sitter.Point) bool {
	start := node.StartPoint()
	end := node.EndPoint()
	if candidate.Row < start.Row {
		return false
	}
	if candidate.Row == start.Row && candidate.Column < start.Column {
		return false
	}
	if candidate.Row == end.Row && candidate.Column > end.Column {
		return false
	}
	if candidate.Row > end.Row {
		return false
	}
	return true
}

var relevanNodeTypes = []string{
	"block_mapping_pair",
	"block_sequence_item",
}

func GetPath(tree *sitter.Tree, source []byte, position protocol.Position) *yaml.Path {

	point := sitter.Point{Row: position.Line, Column: position.Character}

	root := tree.RootNode()

	current := root
	current_new := true

	for current_new {
		current_new = false
		for i := uint32(0); i < current.ChildCount(); i++ {
			child := current.Child(int(i))
			if hasPoint(child, &point) {
				current = child
				current_new = true
			}
		}
	}

	logger.Debug("Found leaf node.", "node", current)

	var nodes []*sitter.Node

	for current != nil {
		if slices.Contains(relevanNodeTypes, current.Type()) {
			nodes = append(nodes, current)
		}
		current = current.Parent()
	}

	logger.Debug("Found nodes.", "len", len(nodes))

	builder := &yaml.PathBuilder{}
	builder = builder.Root()

	for i := len(nodes) - 1; i >= 0; i-- {
		node := nodes[i]
		switch node.Type() {
		case "block_mapping_pair": 
			builder.Child(extractBlockMappingPairKey(source, node))
		case "block_sequence_item":
			builder.Index(extractBlockSequenceItemIndex(node))
		}
	}

	return builder.Build()
}

func extractBlockMappingPairKey(source []byte, node *sitter.Node) string {
	logger.Debug("extractBlockMappingPairKey")
	return node.ChildByFieldName("key").Content(source)
}

func extractBlockSequenceItemIndex(node *sitter.Node) uint {
	logger.Debug("extractBlockSequenceItemIndex")
	parent := node.Parent()
	for i := uint32(0); i < parent.ChildCount(); i++ {
		if parent.Child(int(i)) == node {
			return uint(i)
		}
	}
	panic("Unreachable.")
}
