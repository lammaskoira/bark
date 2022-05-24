package util_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/lammaskoira/bark/pkg/util"
)

func TestFileTracker(t *testing.T) {
	t.Parallel()

	ft := &util.FileTracker{}
	ft.SetupFileTracker()

	tdir := t.TempDir()

	for i := 0; i < 10; i++ {
		fpath := filepath.Join(tdir, fmt.Sprintf("testfile-%d", i))
		_, err := os.Create(fpath)
		require.NoError(t, err, "Create() should not return an error.")
		ft.TrackFile(fpath)
	}
	require.Equal(t, 10, ft.Len(), "There should be three files tracked.")

	err := ft.TearDownFileTracker(context.Background())
	require.NoError(t, err)
}

func TestFileTrackerSucceedsTeardownWithNoTrackedFiles(t *testing.T) {
	t.Parallel()

	ft := &util.FileTracker{}
	ft.SetupFileTracker()

	err := ft.TearDownFileTracker(context.Background())
	require.NoError(t, err)
}

// TearDown with unexistent files succeeds as the desired state is
// for the file to be deleted.
func TestFileTrackerSuceedsTearDownWithUnexistentFiles(t *testing.T) {
	t.Parallel()

	ft := &util.FileTracker{}
	ft.SetupFileTracker()

	fpath := t.TempDir() + "/testfile"
	ft.TrackFile(fpath)

	err := ft.TearDownFileTracker(context.Background())
	require.NoError(t, err, "TearDownFileTracker() should not return an error.")
}
