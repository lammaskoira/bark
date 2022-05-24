package github

import (
	"context"
	"fmt"

	gh "github.com/google/go-github/v44/github"

	apiv1 "github.com/lammaskoira/bark/api/v1"
	"github.com/lammaskoira/bark/pkg/runner/gitcommon"
	"github.com/lammaskoira/bark/pkg/runner/githubcommon"
	rif "github.com/lammaskoira/bark/pkg/runner/runnerinterface"
	"github.com/lammaskoira/bark/pkg/util"
)

func NewGitHubRunner(ts *apiv1.TrickSet) (rif.Runner, error) {
	return &gitHubRunner{
		ts:    ts,
		repos: util.NewStack[string](),
	}, nil
}

type gitHubRunner struct {
	ts    *apiv1.TrickSet
	repos *util.Stack[string]
	githubcommon.GitHubCommon
	gitcommon.GitCommon
}

func (gr *gitHubRunner) Setup(ctx context.Context) error {
	gr.SetupCommon(ctx)
	gr.SetupGitHubClient(ctx)

	// Populate access token internally.
	gr.SetAccessToken(gr.GetGitHubAccessToken())

	org := gr.ts.Context.GitHub.Org

	repos, _, repoErr := gr.GetGHClient().Repositories.ListByOrg(ctx, org, &gh.RepositoryListByOrgOptions{})
	if repoErr != nil {
		return fmt.Errorf("unable to setup GitHub runner: %w", repoErr)
	}

	for _, repo := range repos {
		gr.repos.Push(*repo.Name)
	}

	return nil
}

func (gr *gitHubRunner) Teardown(ctx context.Context) error {
	return gr.TearDownCommon(ctx)
}

func (gr *gitHubRunner) Next(ctx context.Context) (rif.TargetEval, error) {
	if gr.repos.IsEmpty() {
		return nil, rif.ErrEndOfTargets
	}

	repoName := gr.repos.Pop()
	// handles both the insertion of a nil value and the end of the stack
	if repoName == "" {
		return nil, rif.ErrEndOfTargets
	}

	repoInfo, _, repoErr := gr.GetGHClient().Repositories.Get(ctx,
		gr.ts.Context.GitHub.Org, repoName)
	if repoErr != nil {
		return nil, fmt.Errorf("unable to get repo info: %w", repoErr)
	}

	vulnalerts, _, vulnErr := gr.GetGHClient().Repositories.GetVulnerabilityAlerts(ctx,
		gr.ts.Context.GitHub.Org, repoName)

	if vulnErr != nil {
		return nil, fmt.Errorf("unable to get repo vulnerability alerts: %w", vulnErr)
	}

	grepo := githubRepoToRepoRef(repoInfo)

	return func(ctx context.Context) (*apiv1.ContextualResult, error) {
		input := map[string]interface{}{
			apiv1.RepositoryMetadataInputKey: repoInfo,
			"vulnerability_alerts_enabled":   vulnalerts,
		}
		return gr.HandleGit(ctx, grepo, gr.ts, input)
	}, nil
}

func githubRepoToRepoRef(repo *gh.Repository) *apiv1.GitDefinition {
	return &apiv1.GitDefinition{
		URL:    repo.GetSVNURL(),
		Branch: repo.GetDefaultBranch(),
	}
}
