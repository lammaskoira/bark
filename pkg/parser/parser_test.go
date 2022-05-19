package parser_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	apiv1 "github.com/lammaskoira/bark/api/v1"
	"github.com/lammaskoira/bark/pkg/parser"
	_ "github.com/lammaskoira/bark/tests"
)

func TestSimpleInlineTrickSet(t *testing.T) {
	t.Parallel()

	f, oerr := os.Open("tests/data/v1/tricksets/simple_inline.yaml")
	require.NoError(t, oerr)
	defer f.Close()

	ts, perr := parser.NewParser().Parse(f)
	require.NoError(t, perr, "Parse() should not return an error.")

	require.Equal(t, apiv1.Version, ts.Version, "The version of the parser should match the version of the object.")
	require.NotNil(t, ts.Context.Git, "The context should have a Git context.")
	require.Nil(t, ts.Context.GitHub, "The context should not have a GitHub context.")

	require.Equal(t, "https://github.com/lammaskoira/lammaskoira.git", ts.Context.Git.URL,
		"The Git context should have the right URL.")
	require.Equal(t, "main", ts.Context.Git.Branch, "The Git context should have the right branch.")

	require.Len(t, ts.Rules, 1, "There should be one rule.")
	require.NotEmpty(t, ts.Rules[0].Name, "The rule should have a name.")
	require.NotEmpty(t, ts.Rules[0].InlinePolicy, "The rule should have an inline policy.")
}

func TestInvalidTrickset(t *testing.T) {
	t.Parallel()

	f, oerr := os.Open("tests/data/v1/tricksets/invalid.yaml")
	require.NoError(t, oerr)
	defer f.Close()

	rawts, perr := parser.NewParser().Parse(f)
	require.Nil(t, rawts, "Parse() should return nil.")
	require.Error(t, perr, "Parse() should return an error.")
}
