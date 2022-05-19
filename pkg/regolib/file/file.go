package file

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/types"
)

func Library() []func(*rego.Rego) {
	return []func(*rego.Rego){
		Exists,
	}
}

var Exists = rego.Function1(
	&rego.Function{
		Name:    "file.exists",
		Decl:    types.NewFunction(types.Args(types.S), types.B),
		Memoize: true,
	},
	func(bctx rego.BuiltinContext, op1 *ast.Term) (*ast.Term, error) {
		var path string
		if err := ast.As(op1.Value, &path); err != nil {
			return nil, err
		}

		cpath := filepath.Clean(path)
		finfo, err := os.Stat(cpath)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return ast.BooleanTerm(false), nil
			}
			return nil, err
		}

		if finfo.IsDir() {
			return ast.BooleanTerm(false), nil
		}

		return ast.BooleanTerm(true), nil
	},
)
