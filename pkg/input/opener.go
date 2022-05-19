package input

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func GetInputFromFile(path string) (any, error) {
	f, oerr := os.Open(path)
	if oerr != nil {
		return nil, fmt.Errorf("failed to open input file: %w", oerr)
	}

	defer f.Close()

	return GetInputFromReader(f)
}

func GetInputFromReader(r io.Reader) (any, error) {
	output := &map[string]any{}
	d := json.NewDecoder(r)

	if err := d.Decode(output); err != nil {
		return nil, fmt.Errorf("failed to decode input: %w", err)
	}

	return output, nil
}
