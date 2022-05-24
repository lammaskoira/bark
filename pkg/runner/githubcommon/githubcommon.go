package githubcommon

import (
	"context"
	"net/http"
	"os"

	gh "github.com/google/go-github/v44/github"
	"golang.org/x/oauth2"
)

const (
	// nolint:gosec // This is a constant, not a secret
	GitHubTokenEnvKey = "GITHUB_TOKEN"
)

type GitHubCommon struct {
	client     *gh.Client
	httpClient *http.Client
}

func (gc *GitHubCommon) SetGHClient(client *gh.Client) {
	gc.client = client
}

func (gc *GitHubCommon) GetGHClient() *gh.Client {
	return gc.client
}

func (gc *GitHubCommon) SetupGitHubClient(ctx context.Context) {
	gc.SetGHClient(gh.NewClient(gc.BuildHTTPClient(ctx)))
}

func (gc *GitHubCommon) GetGitHubAccessToken() string {
	return os.Getenv(GitHubTokenEnvKey)
}

func (gc *GitHubCommon) BuildHTTPClient(ctx context.Context) *http.Client {
	if gc.httpClient != nil {
		return gc.httpClient
	}

	token := gc.GetGitHubAccessToken()
	if token == "" {
		return http.DefaultClient
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	return tc
}

func (gc *GitHubCommon) SetHTTPClient(client *http.Client) {
	gc.httpClient = client
}

func (gc *GitHubCommon) GetHTTPClient() *http.Client {
	return gc.httpClient
}
