package parser

import (
	"fmt"
	"io"

	"gopkg.in/yaml.v3"

	apiv1 "github.com/lammaskoira/bark/api/v1"
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(r io.Reader) (*apiv1.TrickSet, error) {
	ts := &apiv1.TrickSet{}
	if derr := yaml.NewDecoder(r).Decode(ts); derr != nil {
		return nil, fmt.Errorf("could not decode Trickset: %w", derr)
	}
	return ts, nil
}
