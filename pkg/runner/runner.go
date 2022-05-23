package runner

import (
	"context"
	"errors"
	"fmt"

	apiv1 "github.com/lammaskoira/bark/api/v1"
	"github.com/lammaskoira/bark/pkg/runner/git"
	"github.com/lammaskoira/bark/pkg/runner/github"
	rif "github.com/lammaskoira/bark/pkg/runner/runnerinterface"
)

func Run(ctx context.Context, ts *apiv1.TrickSet) error {
	if err := ts.Validate(); err != nil {
		return fmt.Errorf("could not validate TrickSet: %w", err)
	}

	runner, err := GetContextualRunner(ts)
	if err != nil {
		return fmt.Errorf("could not get runner: %w", err)
	}

	if err := runner.Setup(ctx); err != nil {
		return fmt.Errorf("could not setup context: %w", err)
	}

	for {
		te, nerr := runner.Next(ctx)
		if nerr != nil {
			if errors.Is(nerr, rif.ErrEndOfTargets) {
				break
			}
			return fmt.Errorf("could not get next target: %w", nerr)
		}

		if err := te(ctx); err != nil {
			return fmt.Errorf("could not run target: %w", err)
		}
	}

	return runner.Teardown(context.TODO())
}

func GetContextualRunner(trickSet *apiv1.TrickSet) (rif.Runner, error) {
	switch trickSet.Context.Provider {
	case apiv1.GitContext:
		return git.NewGitRunner(trickSet)
	case apiv1.GitHubContext:
		return github.NewGitHubRunner(trickSet)
	}
	return nil, fmt.Errorf("could not get runner for provider %s", trickSet.Context.Provider)
}
