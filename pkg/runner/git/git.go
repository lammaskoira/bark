package git

import (
	"context"

	apiv1 "github.com/lammaskoira/bark/api/v1"
	"github.com/lammaskoira/bark/pkg/runner/gitcommon"
	rif "github.com/lammaskoira/bark/pkg/runner/runnerinterface"
	"github.com/lammaskoira/bark/pkg/util"
)

func NewGitRunner(ts *apiv1.TrickSet) (rif.Runner, error) {
	return &gitRunner{
		ts:    ts,
		repos: util.NewStack[*apiv1.GitDefinition](),
	}, nil
}

type gitRunner struct {
	ts    *apiv1.TrickSet
	repos *util.Stack[*apiv1.GitDefinition]
	gitcommon.GitCommon
}

func (gr *gitRunner) Setup(ctx context.Context) error {
	gr.SetupCommon(ctx)
	gr.repos.Push(gr.ts.Context.Git)
	return nil
}

func (gr *gitRunner) Teardown(ctx context.Context) error {
	return gr.TearDownCommon(ctx)
}

func (gr *gitRunner) Next(ctx context.Context) (rif.TargetEval, error) {
	repoRef := gr.repos.Pop()
	// handles both the insertion of a nil value and the end of the stack
	if repoRef == nil {
		return nil, rif.ErrEndOfTargets
	}

	return func(ctx context.Context) (*apiv1.ContextualResult, error) {
		return gr.HandleGit(ctx, repoRef, gr.ts, map[string]interface{}{})
	}, nil
}
