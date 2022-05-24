package v1

type GitHubDefinition struct {
	// Org is the organization that owns the repos
	Org string `json:"org" yaml:"org"`

	// Branch to check for when evaluating the policy on repos
	Branch string `json:"branch" yaml:"branch"`

	// Override the default branch for a repo
	Overrides []GitHubDefinition `json:"overrides,omitempty" yaml:"overrides,omitempty"`

	// Exclude repos from the policy
	Exclude []string `json:"exclude,omitempty" yaml:"exclude,omitempty"`

	// TODO(jaosorior): Add support for credentials
}

const (
	// RepositoryMetadataInputKey is a key in the rego input will contain the metadata
	// available for the given repository.
	// The values come directly from the GitHub API [1]
	// [1]: https://docs.github.com/en/rest/repos/repos#get-a-repository
	RepositoryMetadataInputKey = "repometa"

	// VulnerabilityAlertsEnabledInputKey is a key in the rego input
	// that will contain a boolean that tells us whether vulnerability
	// alerts are enabled for a given repository or not.
	VulnerabilityAlertsEnabledInputKey = "vulnerability_alerts_enabled"
)
