package githuborgconfig

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	gh "github.com/google/go-github/v44/github"
	"gopkg.in/yaml.v2"

	apiv1 "github.com/lammaskoira/bark/api/v1"
	"github.com/lammaskoira/bark/pkg/runner/githubcommon"
	rif "github.com/lammaskoira/bark/pkg/runner/runnerinterface"
	"github.com/lammaskoira/bark/pkg/util"
)

type gitHubOrgConfigRunner struct {
	ts *apiv1.TrickSet
	// We still create a stack here... but it's meant to only
	// hold one object
	org *util.Stack[string]
	githubcommon.GitHubCommon
	util.FileTracker
}

func NewGitHubOrgConfigRunner(ts *apiv1.TrickSet) (rif.Runner, error) {
	return &gitHubOrgConfigRunner{
		ts:  ts,
		org: util.NewStack[string](),
	}, nil
}

func (gr *gitHubOrgConfigRunner) Setup(ctx context.Context) error {
	gr.SetupFileTracker()
	gr.SetupGitHubClient(ctx)

	org := gr.ts.Context.GitHubOrgConfig.Org
	gr.org.Push(org)

	return nil
}

func (gr *gitHubOrgConfigRunner) Teardown(ctx context.Context) error {
	return gr.TearDownFileTracker(ctx)
}

func (gr *gitHubOrgConfigRunner) Next(ctx context.Context) (rif.TargetEval, error) {
	if gr.org.IsEmpty() {
		return nil, rif.ErrEndOfTargets
	}

	org := gr.org.Pop()
	// handles both the insertion of an empty value and the end of the stack
	if org == "" {
		return nil, rif.ErrEndOfTargets
	}

	orgInfo, _, orgErr := gr.GetGHClient().Organizations.Get(ctx, org)
	if orgErr != nil {
		return nil, fmt.Errorf("unable to get org info: %w", orgErr)
	}

	appsInstalled, _, appErr := gr.GetGHClient().Organizations.ListInstallations(ctx, org,
		&gh.ListOptions{})
	if appErr != nil {
		return nil, fmt.Errorf("unable to get apps installed in organization: %w", appErr)
	}

	return func(ctx context.Context) (*apiv1.ContextualResult, error) {
		input := map[string]interface{}{
			apiv1.OrgConfigInputKey: orgInfo,
			apiv1.AppsInOrgInputKey: appsInstalled,
		}
		targetDir, tderr := ioutil.TempDir("", "bark-github-org-config")
		if tderr != nil {
			return nil, fmt.Errorf("unable to create temp dir: %w", tderr)
		}
		gr.TrackFile(targetDir)

		inputfile, terr := ioutil.TempFile("", "bark-input")
		if terr != nil {
			return nil, fmt.Errorf("failed to create input file: %w", terr)
		}
		defer inputfile.Close()
		gr.TrackFile(inputfile.Name())

		if eerr := json.NewEncoder(inputfile).Encode(input); eerr != nil {
			return nil, fmt.Errorf("failed to encode input: %w", eerr)
		}

		tsfile, terr := ioutil.TempFile("", "bark-ts")
		if terr != nil {
			return nil, fmt.Errorf("failed to create temp trickset file: %w", terr)
		}
		defer tsfile.Close()
		gr.TrackFile(tsfile.Name())

		if eerr := yaml.NewEncoder(tsfile).Encode(gr.ts); eerr != nil {
			return nil, fmt.Errorf("failed to encode input: %w", eerr)
		}
		var outsb, errsb bytes.Buffer

		// this would be the environment
		cmd := os.Args[0]

		c := exec.Command(cmd, "bark",
			"-i", inputfile.Name(),
			"-t", tsfile.Name(),
			"-x", org,
			"-r", targetDir,
		)

		c.Stdout = &outsb
		c.Stderr = &errsb

		if err := c.Run(); err != nil {
			fmt.Printf("err: %s\n", err)
		}

		str := outsb.Bytes()

		cr := &apiv1.ContextualResult{}
		derr := json.Unmarshal(str, cr)
		if derr != nil {
			return nil, fmt.Errorf("failed to decode output: %w", derr)
		}

		return cr, nil
	}, nil
}
