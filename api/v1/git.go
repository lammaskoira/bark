package v1

type GitDefinition struct {
	URL    string `json:"url" yaml:"url"`
	Branch string `json:"branch" yaml:"branch"`

	// TODO(jaosorior): Add support for credentials
}
