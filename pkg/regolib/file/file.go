package file

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/types"
)

func Library() []func(*rego.Rego) {
	return []func(*rego.Rego){
		Exists,
		ReadAll,
	}
}

var Exists = rego.Function1(
	&rego.Function{
		Name: "file.exists",
		Decl: types.NewFunction(types.Args(types.S), types.B),
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

var ReadAll = rego.Function1(
	&rego.Function{
		Name: "file.readall",
		Decl: types.NewFunction(types.Args(types.S), types.S),
	},
	func(bctx rego.BuiltinContext, op1 *ast.Term) (*ast.Term, error) {
		var path string
		if err := ast.As(op1.Value, &path); err != nil {
			return nil, err
		}

		cpath := filepath.Clean(path)
		f, err := os.Open(cpath)
		if err != nil {
			return nil, err
		}

		defer f.Close()

		all, rerr := io.ReadAll(f)
		if rerr != nil {
			return nil, rerr
		}

		allstr := ast.String(all)
		return ast.NewTerm(allstr), nil
	},
)
