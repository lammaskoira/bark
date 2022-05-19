package regolib

import (
	"github.com/open-policy-agent/opa/rego"

	"github.com/lammaskoira/bark/pkg/regolib/file"
)

func Library() []func(*rego.Rego) {
	return file.Library()
}
