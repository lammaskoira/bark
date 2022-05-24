package input

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetInputFromFile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    any
		wantErr bool
	}{
		{
			name:  "can read valid json",
			input: `{"foo": "bar"}`,
			want: &map[string]any{
				"foo": "bar",
			},
		},
		{
			name: "can read valid json with whitespace",
			input: `{
	"foo": "bar"
			}`,
			want: &map[string]any{
				"foo": "bar",
			},
		},
		{
			name:    "reports an error if the file is not valid json",
			input:   `{ "foo": "bar" `,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			dirpath := t.TempDir()
			path := dirpath + "/input.json"

			writeerr := ioutil.WriteFile(path, []byte(tt.input), 0o600)
			require.NoError(t, writeerr, "failed to write input file")

			got, err := GetInputFromFile(path)
			if tt.wantErr {
				require.Error(t, err, "expected error")
			} else {
				require.NoError(t, err, "unexpected error")
				require.Equal(t, tt.want, got, "unexpected output")
			}
		})
	}
}

func TestGetInputFromFileFailsWithUnexistentFile(t *testing.T) {
	t.Parallel()

	_, err := GetInputFromFile("/does/not/exist")
	require.Error(t, err, "expected error")
}

func TestGetInputFromReaderFailsWithNilValue(t *testing.T) {
	t.Parallel()

	_, err := GetInputFromReader(nil)
	require.Error(t, err, "expected error")
}
