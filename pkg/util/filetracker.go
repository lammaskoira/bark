package util

import (
	"context"
	"os"
)

type FileTracker struct {
	filesToClean *Stack[string]
}

func (ft *FileTracker) SetupFileTracker() {
	ft.filesToClean = NewStack[string]()
}

func (ft *FileTracker) TrackFile(file string) {
	ft.filesToClean.Push(file)
}

func (ft *FileTracker) Len() int {
	return ft.filesToClean.Len()
}

func (ft *FileTracker) TearDownFileTracker(ctx context.Context) error {
	for !ft.filesToClean.IsEmpty() {
		if err := os.RemoveAll(ft.filesToClean.Pop()); err != nil {
			return err
		}
	}
	return nil
}
