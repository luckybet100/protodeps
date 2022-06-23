package scheme

import (
	"encoding/json"
	"github.com/luckybet100/protodeps/pkg/errors"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"strings"
)

type Parser interface {
	Parse() (*ProtoDepsConfig, error)
}

type decoder interface {
	Decode(v interface{}) error
}

type parserImpl struct {
	decoder decoder
}

func (parser *parserImpl) Parse() (*ProtoDepsConfig, error) {
	result := &ProtoDepsConfig{}
	if err := parser.decoder.Decode(result); err != nil {
		return nil, errors.ParsingError.Wrap(err, "failed to parse config")
	}
	return result, nil
}

func NewParserYAML(reader io.Reader) Parser {
	return &parserImpl{
		decoder: yaml.NewDecoder(reader),
	}
}

func NewParserJSON(reader io.Reader) Parser {
	return &parserImpl{
		decoder: json.NewDecoder(reader),
	}
}

func GetParser(file *os.File) (Parser, error) {
	if strings.HasSuffix(file.Name(), "json") {
		return NewParserJSON(file), nil
	}
	if strings.HasSuffix(file.Name(), "yml") || strings.HasSuffix(file.Name(), "yaml") {
		return NewParserYAML(file), nil
	}
	return nil, errors.ParsingError.New("unsupported file extension <%s>")
}
