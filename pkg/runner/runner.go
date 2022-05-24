package runner

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	apiv1 "github.com/lammaskoira/bark/api/v1"
	"github.com/lammaskoira/bark/pkg/runner/git"
	"github.com/lammaskoira/bark/pkg/runner/github"
	"github.com/lammaskoira/bark/pkg/runner/githuborgconfig"
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

	rep := &apiv1.Report{
		Version: apiv1.Version,
	}

	for {
		te, nerr := runner.Next(ctx)
		if nerr != nil {
			if errors.Is(nerr, rif.ErrEndOfTargets) {
				break
			}
			return fmt.Errorf("could not get next target: %w", nerr)
		}

		result, err := te(ctx)
		if err != nil {
			return fmt.Errorf("could not run target: %w", err)
		}
		rep.AddResult(result)
	}

	rep.GatherOverall()

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if encerr := enc.Encode(rep); encerr != nil {
		return fmt.Errorf("could not encode report: %w", encerr)
	}

	return runner.Teardown(context.TODO())
}

func GetContextualRunner(trickSet *apiv1.TrickSet) (rif.Runner, error) {
	switch trickSet.Context.Provider {
	case apiv1.GitContext:
		return git.NewGitRunner(trickSet)
	case apiv1.GitHubContext:
		return github.NewGitHubRunner(trickSet)
	case apiv1.GitHubOrgConfigContext:
		return githuborgconfig.NewGitHubOrgConfigRunner(trickSet)
	}
	return nil, fmt.Errorf("could not get runner for provider %s", trickSet.Context.Provider)
}
