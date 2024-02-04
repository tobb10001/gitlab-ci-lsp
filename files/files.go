package files

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/goccy/go-yaml"
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/tliron/commonlog"
	protocol "github.com/tliron/glsp/protocol_3_16"

	"gitlab-ci-lsp/tsyaml"
	"gitlab-ci-lsp/yamlpath"
)

var logger = commonlog.GetLogger("files")
var filestore map[string]*File = make(map[string]*File)


type File struct {
	uri    string
	local  string
	source []byte
	tree   *sitter.Tree
}

func (f *File) GetYamlPath(position protocol.Position) *yaml.Path {
	return yamlpath.GetPath(f.tree, f.source, position)
}

func (f *File) ReadYamlPath(path *yaml.Path, v interface{}) error {
	return path.Read(bytes.NewReader(f.source), v)
}

func Get(uri string) (*File, error) {
	logger.Debug("Getting file.", "uri", uri)

	file, ok := filestore[uri]
	if ok {
		logger.Debug("Cache hit.", "local", file.local)
		return file, nil
	}

	logger.Debug("Cache miss.")

	file, err := new(uri)
	if err != nil {
		return nil, err
	}

	filestore[uri] = file

	logger.Debug("Successfully initialized file.", "local", file.local)

	return file, nil
}

func new(uri string) (*File, error) {
	var local string
	if strings.HasPrefix(uri, "file://") {
		local = strings.TrimPrefix(uri, "file://")
		logger.Debug("Local file requrested.", "local", local)
	} else {
		return nil, fmt.Errorf("URI '%s' can't be interpreted.", uri)
	}

	source, err := read(local)
	if err != nil {
		return nil, err
	}

	logger.Debug("Successfully read file.")

	tree, err := tsyaml.GetTree(source)
	if err != nil {
		return nil, err
	}

	logger.Debug("Successfully built tree.")

	file := File{
		uri: uri,
		local: local,
		source: source,
		tree: tree,
	}
	return &file, nil;
}

func read(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	bs := make([]byte, stat.Size())
	_, err = bufio.NewReader(file).Read(bs)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return bs, nil
}
