package yamlpath_test

import (
	"gitlab-ci-lsp/yamlpath"

	"testing"

	"github.com/lithammer/dedent"
	"github.com/stretchr/testify/require"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func TestGetPath(t *testing.T) {

	for _, tt := range []struct {
		name     string
		source   []byte
		position protocol.Position
		expected string
	}{
		{
			name: "local_include",
			source: []byte(dedent.Dedent(`
				include:
				  - local: 'directory/.gitlab-ci.yml'
			`)),
			position: protocol.Position{Line: 2, Character: 6},
			expected: "$.include[0].local",
		},
		{
			name: "local_include",
			source: []byte(dedent.Dedent(`
				include:
				  - local: 'directory/.gitlab-ci.yml'
			`)),
			position: protocol.Position{Line: 2, Character: 17},
			expected: "$.include[0].local",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := yamlpath.GetPath(tt.source, tt.position)
			require.NoError(t, err)
			require.Equal(t, tt.expected, *actual)
		})
	}
}
