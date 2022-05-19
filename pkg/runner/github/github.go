package github

import (
	"context"

	apiv1 "github.com/lammaskoira/bark/api/v1"
	rif "github.com/lammaskoira/bark/pkg/runner/runnerinterface"
)

func NewGitRunner(ts *apiv1.TrickSet) (rif.Runner, error) {
	return &gitHubRunner{ts}, nil
}

type gitHubRunner struct {
	ts *apiv1.TrickSet
}

func (gr *gitHubRunner) Setup(ctx context.Context) error {
	return nil
}

func (gr *gitHubRunner) Teardown(ctx context.Context) error {
	return nil
}

func (gr *gitHubRunner) Next(ctx context.Context) (rif.TargetEval, error) {
	return nil, nil
}
