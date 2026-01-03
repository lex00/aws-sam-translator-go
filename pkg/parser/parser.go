// Package parser provides YAML/JSON template parsing with intrinsic function detection.
package parser

import (
	"encoding/json"

	"github.com/lex00/aws-sam-translator-go/pkg/types"
	"gopkg.in/yaml.v3"
)

// Parser handles SAM/CloudFormation template parsing.
type Parser struct{}

// New creates a new Parser instance.
func New() *Parser {
	return &Parser{}
}

// ParseYAML parses a YAML template.
func (p *Parser) ParseYAML(data []byte) (*types.Template, error) {
	var template types.Template
	if err := yaml.Unmarshal(data, &template); err != nil {
		return nil, err
	}
	return &template, nil
}

// ParseJSON parses a JSON template.
func (p *Parser) ParseJSON(data []byte) (*types.Template, error) {
	var template types.Template
	if err := json.Unmarshal(data, &template); err != nil {
		return nil, err
	}
	return &template, nil
}
