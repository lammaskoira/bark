package git

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"gopkg.in/yaml.v2"

	apiv1 "github.com/lammaskoira/bark/api/v1"
	rif "github.com/lammaskoira/bark/pkg/runner/runnerinterface"
	"github.com/lammaskoira/bark/pkg/util"
)

const (
	noWorktree = true
)

func NewGitRunner(ts *apiv1.TrickSet) (rif.Runner, error) {
	return &gitRunner{
		ts:           ts,
		repos:        util.NewStack[*apiv1.GitDefinition](),
		filesToClean: util.NewStack[string](),
	}, nil
}

type gitRunner struct {
	ts           *apiv1.TrickSet
	repos        *util.Stack[*apiv1.GitDefinition]
	filesToClean *util.Stack[string]
}

func (gr *gitRunner) Setup(ctx context.Context) error {
	gr.repos.Push(gr.ts.Context.Git)
	return nil
}

func (gr *gitRunner) Teardown(ctx context.Context) error {
	for !gr.filesToClean.IsEmpty() {
		if err := os.RemoveAll(gr.filesToClean.Pop()); err != nil {
			return err
		}
	}
	return nil
}

func (gr *gitRunner) Next(ctx context.Context) (rif.TargetEval, error) {
	repoRef := gr.repos.Pop()
	if repoRef == nil {
		return nil, rif.ErrEndOfTargets
	}

	return func(ctx context.Context) error {
		targetDir, terr := ioutil.TempDir("", "bark-git-")
		if terr != nil {
			return fmt.Errorf("failed to create temp dir: %w", terr)
		}
		gr.filesToClean.Push(targetDir)

		_, err := git.PlainCloneContext(ctx, targetDir, noWorktree, &git.CloneOptions{
			URL:           repoRef.URL,
			ReferenceName: plumbing.NewBranchReferenceName(repoRef.Branch),
		})
		if err != nil {
			return fmt.Errorf("could not clone repo: %w", err)
		}

		// this would be the environment
		cmd := os.Args[0]

		inputfile, terr := ioutil.TempFile("", "bark-input")
		if terr != nil {
			return fmt.Errorf("failed to create input file: %w", terr)
		}
		defer inputfile.Close()
		gr.filesToClean.Push(inputfile.Name())

		if eerr := json.NewEncoder(inputfile).Encode(map[string]string{}); eerr != nil {
			return fmt.Errorf("failed to encode input: %w", eerr)
		}

		tsfile, terr := ioutil.TempFile("", "bark-ts")
		if terr != nil {
			return fmt.Errorf("failed to create temp trickset file: %w", terr)
		}
		defer tsfile.Close()
		gr.filesToClean.Push(tsfile.Name())

		if eerr := yaml.NewEncoder(tsfile).Encode(gr.ts); eerr != nil {
			return fmt.Errorf("failed to encode input: %w", eerr)
		}

		out, err := exec.Command(cmd, "bark",
			"-i", inputfile.Name(),
			"-t", tsfile.Name(),
			"-r", targetDir,
		).CombinedOutput()
		fmt.Printf("output: %s\n", out)
		fmt.Printf("err: %s\n", err)
		return nil
	}, nil
}
