/*
This initializes the current working directory
to be the root of the repository.
*/
package tests

import (
	"fmt"
	"os"
	"path"
	"runtime"
)

//nolint:gochecknoinits // This helps initialize the tests package.
func init() {
	//nolint:dogsled // We really don't need the rest of the values
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..")

	if err := os.Chdir(dir); err != nil {
		panic(err)
	}

	fmt.Printf("Current working directory: %s\n", dir)
}
