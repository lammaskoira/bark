package gitcommon

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"gopkg.in/yaml.v2"

	apiv1 "github.com/lammaskoira/bark/api/v1"
	"github.com/lammaskoira/bark/pkg/util"
)

const (
	withWorktree = false
)

type GitCommon struct {
	filesToClean *util.Stack[string]
	accessToken  string
}

func (gc *GitCommon) SetAccessToken(tok string) {
	gc.accessToken = tok
}

func (gc *GitCommon) GetAccessToken() string {
	return gc.accessToken
}

func (gc *GitCommon) SetupCommon(ctx context.Context) {
	gc.filesToClean = util.NewStack[string]()
}

func (gc *GitCommon) TrackFile(file string) {
	gc.filesToClean.Push(file)
}

func (gc *GitCommon) TearDownCommon(ctx context.Context) error {
	for !gc.filesToClean.IsEmpty() {
		if err := os.RemoveAll(gc.filesToClean.Pop()); err != nil {
			return err
		}
	}
	return nil
}

func (gc *GitCommon) HandleGit(
	ctx context.Context,
	gd *apiv1.GitDefinition,
	ts *apiv1.TrickSet,
	input any,
) (*apiv1.ContextualResult, error) {
	targetDir, terr := ioutil.TempDir("", "bark-git-")
	if terr != nil {
		return nil, fmt.Errorf("failed to create temp dir: %w", terr)
	}
	gc.TrackFile(targetDir)

	opts := &git.CloneOptions{
		URL:           gd.URL,
		ReferenceName: plumbing.NewBranchReferenceName(gd.Branch),
	}
	if gc.GetAccessToken() != "" {
		opts.Auth = &http.BasicAuth{
			// the Username can be anything but it can't be empty
			Username: "JAORMX",
			Password: gc.GetAccessToken(),
		}
	}
	_, err := git.PlainCloneContext(ctx, targetDir, withWorktree, opts)
	if err != nil {
		return nil, fmt.Errorf("could not clone repo: %w", err)
	}

	// this would be the environment
	cmd := os.Args[0]

	inputfile, terr := ioutil.TempFile("", "bark-input")
	if terr != nil {
		return nil, fmt.Errorf("failed to create input file: %w", terr)
	}
	defer inputfile.Close()
	gc.TrackFile(inputfile.Name())

	if eerr := json.NewEncoder(inputfile).Encode(input); eerr != nil {
		return nil, fmt.Errorf("failed to encode input: %w", eerr)
	}

	tsfile, terr := ioutil.TempFile("", "bark-ts")
	if terr != nil {
		return nil, fmt.Errorf("failed to create temp trickset file: %w", terr)
	}
	defer tsfile.Close()
	gc.TrackFile(tsfile.Name())

	if eerr := yaml.NewEncoder(tsfile).Encode(ts); eerr != nil {
		return nil, fmt.Errorf("failed to encode input: %w", eerr)
	}

	var outsb, errsb bytes.Buffer

	c := exec.Command(cmd, "bark",
		"-i", inputfile.Name(),
		"-t", tsfile.Name(),
		"-x", getTargetString(gd.URL, gd.Branch),
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
}

func getTargetString(url, branch string) string {
	return fmt.Sprintf("%s@%s", url, branch)
}
