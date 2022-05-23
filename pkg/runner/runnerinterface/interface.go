package runnerinterface

import (
	"context"
	"errors"

	apiv1 "github.com/lammaskoira/bark/api/v1"
)

var ErrEndOfTargets = errors.New("end of targets")

type TargetEval func(context.Context) (*apiv1.ContextualResult, error)

type Runner interface {
	Setup(context.Context) error
	Teardown(context.Context) error
	Next(ctx context.Context) (TargetEval, error)
}
