package github

import (
	"context"
	"fmt"
	"net/http"
	"os"

	gh "github.com/google/go-github/v44/github"
	"golang.org/x/oauth2"

	apiv1 "github.com/lammaskoira/bark/api/v1"
	"github.com/lammaskoira/bark/pkg/runner/gitcommon"
	rif "github.com/lammaskoira/bark/pkg/runner/runnerinterface"
	"github.com/lammaskoira/bark/pkg/util"
)

const (
	// nolint:gosec // This is a constant, not a secret
	GitHubTokenEnvKey = "GITHUB_TOKEN"
)

func NewGitHubRunner(ts *apiv1.TrickSet) (rif.Runner, error) {
	return &gitHubRunner{
		ts:    ts,
		repos: util.NewStack[*gh.Repository](),
	}, nil
}

type gitHubRunner struct {
	ts         *apiv1.TrickSet
	repos      *util.Stack[*gh.Repository]
	client     *gh.Client
	httpClient *http.Client
	gitcommon.GitCommon
}

func (gr *gitHubRunner) Setup(ctx context.Context) error {
	gr.SetupCommon(ctx)
	gr.client = gh.NewClient(gr.buildHTTPClient(ctx))

	org := gr.ts.Context.GitHub.Org

	repos, _, repoErr := gr.client.Repositories.ListByOrg(ctx, org, &gh.RepositoryListByOrgOptions{})
	if repoErr != nil {
		return fmt.Errorf("unable to setup GitHub runner: %w", repoErr)
	}

	for _, repo := range repos {
		gr.repos.Push(repo)
	}

	return nil
}

func (gr *gitHubRunner) Teardown(ctx context.Context) error {
	return gr.TearDownCommon(ctx)
}

func (gr *gitHubRunner) Next(ctx context.Context) (rif.TargetEval, error) {
	repoInfo := gr.repos.Pop()
	// handles both the insertion of a nil value and the end of the stack
	if repoInfo == nil {
		return nil, rif.ErrEndOfTargets
	}

	grepo := githubRepoToRepoRef(repoInfo)

	return func(ctx context.Context) (*apiv1.ContextualResult, error) {
		// TODO(jaosorior): Add GitHub info as input
		return gr.HandleGit(ctx, grepo, gr.ts, map[string]interface{}{})
	}, nil
}

// allows for dependency injection.
// nolint:unused // TODO(jaosorior): write unit tests
func (gr *gitHubRunner) setHTTPClient(c *http.Client) {
	gr.httpClient = c
}

func (gr *gitHubRunner) buildHTTPClient(ctx context.Context) *http.Client {
	if gr.httpClient != nil {
		return gr.httpClient
	}

	if token := os.Getenv(GitHubTokenEnvKey); token != "" {
		gr.SetAccessToken(token)
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: gr.GetAccessToken()},
		)
		return oauth2.NewClient(ctx, ts)
	}

	return nil
}

func githubRepoToRepoRef(repo *gh.Repository) *apiv1.GitDefinition {
	return &apiv1.GitDefinition{
		URL:    repo.GetSVNURL(),
		Branch: repo.GetDefaultBranch(),
	}
}
