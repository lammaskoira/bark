package runnerinterface

import (
	"context"
	"errors"
)

var ErrEndOfTargets = errors.New("end of targets")

type TargetEval func(context.Context) error

type Runner interface {
	Setup(context.Context) error
	Teardown(context.Context) error
	Next(ctx context.Context) (TargetEval, error)
}
