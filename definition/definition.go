package definition

import (
	"gitlab-ci-lsp/files"
	"path"
	"strings"

	"github.com/tliron/commonlog"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

var log = commonlog.GetLogger("definition")

func TextDocumentDefinition(context *glsp.Context, params *protocol.DefinitionParams) (any, error) {
	log.Noticef("textDocumentDefinition: %v", *params)

	// Find target.
	file, err := files.Get(params.TextDocument.URI)
	if err != nil {
		return nil, err
	}

	yaml_path := file.GetYamlPath(params.Position)

	// this will have to become more sophisticated later
	// but let's just assume that the path has the pattern
	// $.include[<index>].local
	// path.
	var included string
	err = file.ReadYamlPath(yaml_path, &included)
	if err != nil {
		return nil, err
	}

	requestor_dir := path.Dir(strings.TrimPrefix(params.TextDocument.URI, "file://"))
	target_absolute := path.Join(requestor_dir, included)

	location := protocol.Location{
		URI: "file://" + target_absolute,
		Range: protocol.Range{
			Start: protocol.Position{Line: 0, Character: 0},
			End: protocol.Position{Line: 0, Character: 0},
		},
	}

	return location, nil	
}
