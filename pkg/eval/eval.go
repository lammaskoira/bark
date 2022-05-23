package eval

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/open-policy-agent/opa/rego"

	apiv1 "github.com/lammaskoira/bark/api/v1"
	barkerrors "github.com/lammaskoira/bark/pkg/errors"
	"github.com/lammaskoira/bark/pkg/regolib"
)

func EvaluateRules(
	ctx context.Context,
	target string,
	input any,
	rules []apiv1.RuleDefinition,
) (*apiv1.ContextualResult, error) {
	r := &apiv1.ContextualResult{
		Version: apiv1.Version,
		Target:  target,
	}
	for _, rule := range rules {
		if err := EvaluateOnePolicy(ctx, input, rule.InlinePolicy); err != nil {
			if errors.Is(err, barkerrors.ErrPolicyDenial) {
				// TODO: Generate ID
				r.AddResult(apiv1.Result{
					Name:   rule.Name,
					Status: apiv1.ResultStatusFail,
				})
				continue
			}
			r.AddResult(apiv1.Result{
				Name:         rule.Name,
				Status:       apiv1.ResultStatusError,
				ErrorMessage: err.Error(),
			})
			return r, err
		}
		r.AddResult(apiv1.Result{
			Name:   rule.Name,
			Status: apiv1.ResultStatusPass,
		})
	}

	return r, nil
}

func EvaluateOnePolicy(ctx context.Context, input any, policy string, strictBuiltin ...bool) error {
	args := []func(*rego.Rego){
		rego.Query("data.bark.allow"),
		rego.Module("bark.rego", policy),
		rego.Dump(os.Stderr),
		rego.EnablePrintStatements(true),
	}

	// This is useful for debugging our custom builtins.
	if len(strictBuiltin) > 0 && strictBuiltin[0] {
		args = append(args, rego.StrictBuiltinErrors(true))
	}

	r := rego.New(append(args, regolib.Library()...)...)
	pq, err := r.PrepareForEval(ctx)
	if err != nil {
		return fmt.Errorf("could not prepare Rego: %w - %s", barkerrors.ErrPolicyParseError, err)
	}

	rs, err := pq.Eval(ctx, rego.EvalInput(input))
	if err != nil || len(rs) == 0 {
		return fmt.Errorf("error evaluating policy. Might be wrong input or no results: %w - %s",
			barkerrors.ErrPolicyEvalError, err)
	}

	if rs.Allowed() {
		return nil
	}

	return barkerrors.ErrPolicyDenial
}
