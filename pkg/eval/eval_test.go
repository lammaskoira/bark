package eval_test

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	barkerrors "github.com/lammaskoira/bark/pkg/errors"
	"github.com/lammaskoira/bark/pkg/eval"
	_ "github.com/lammaskoira/bark/tests"
)

func TestEval(t *testing.T) {
	t.Parallel()

	type args struct {
		rawinput   []byte
		policyPath string
	}
	tests := []struct {
		name        string
		args        args
		wantErr     bool
		errInstance error
	}{
		{
			name: "renovate is configured in this repo",
			args: args{
				rawinput:   nil,
				policyPath: "tests/data/rego/renovate.rego",
			},
		},
		{
			name: "unexistent file in repo should return an error",
			args: args{
				rawinput:   nil,
				policyPath: "tests/data/rego/unexistent-file.rego",
			},
			wantErr:     true,
			errInstance: barkerrors.ErrPolicyDenial,
		},
		{
			name: "invalid rego should return an error",
			args: args{
				rawinput:   nil,
				policyPath: "tests/data/rego/invalid.rego",
			},
			wantErr:     true,
			errInstance: barkerrors.ErrPolicyParseError,
		},
		{
			name: "rego with no results should output an error",
			args: args{
				rawinput:   nil,
				policyPath: "tests/data/rego/no-results-given.rego",
			},
			wantErr:     true,
			errInstance: barkerrors.ErrPolicyEvalError,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var input any
			if tt.args.rawinput == nil {
				input = nil
			} else {
				umerr := json.Unmarshal(tt.args.rawinput, &input)
				require.NoError(t, umerr, "unmarshalling raw input")
			}

			f, ferr := os.Open(tt.args.policyPath)
			require.NoError(t, ferr, "opening policy file")

			policy, rerr := io.ReadAll(f)
			require.NoError(t, rerr, "reading policy file")

			err := eval.EvaluateOnePolicy(context.TODO(), input, string(policy))
			if tt.wantErr {
				require.Error(t, err, "evaluation should errored out")
				require.ErrorIs(t, err, tt.errInstance)
			} else {
				require.NoError(t, err, "evaluation should have not errored out")
			}
		})
	}
}
