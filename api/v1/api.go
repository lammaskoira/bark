package v1

const (
	Version = "v1"
)

type TrickSet struct {
	// Version of the trickset format (should be v1)
	Version string `json:"version" yaml:"version"`
	// Context is the "context" of the trickset, or where
	// it will execute on.
	Context ContextDefinition `json:"context" yaml:"context"`
	// Rules is a list of rules that will be executed in order.
	Rules []RuleDefinition `json:"rules" yaml:"rules"`
}

type ValidContext string

const (
	GitContext    = ValidContext("git")
	GitHubContext = ValidContext("github")
)

type ContextDefinition struct {
	Provider ValidContext      `json:"provider" yaml:"provider"`
	Git      *GitDefinition    `json:"git,omitempty" yaml:"git,omitempty"`
	GitHub   *GitHubDefinition `json:"github,omitempty" yaml:"github,omitempty"`
}

type GitDefinition struct {
	URL    string `json:"url" yaml:"url"`
	Branch string `json:"branch" yaml:"branch"`

	// TODO(jaosorior): Add support for credentials
}

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

type GithubRepoConfig struct {
	Org string `json:"org" yaml:"org"`
}

type RuleDefinition struct {
	// * Base fields

	// Name of the rule
	Name string `json:"name" yaml:"name"`

	// Optional ID of the rule. If not specified, an
	// auto-generated ID will be used.
	ID string `json:"id,omitempty" yaml:"id,omitempty"`

	// * Policy execution fields

	// Include will import the rule from the remote location
	Include string `json:"include,omitempty" yaml:"include,omitempty"`

	// InlinePolicy specifies a raw rego policy inline
	InlinePolicy string `json:"inlinePolicy,omitempty" yaml:"inlinePolicy,omitempty"`

	RulesFrom *RulesFromDefinition `json:"rulesFrom,omitempty"`

	// * Logical Operator(s)

	// Or: If any of the conditions are true, the rule is true
	Or []*RuleDefinition `json:"or,omitempty" yaml:"or,omitempty"`
}

type RulesFromDefinition struct {
	// File to load the policies from
	File string `json:"file,omitempty" yaml:"file,omitempty"`
	// Trickset reference to load the policies from
	// e.g. <repository name>/<tricket file path>@<branch or version>
	TricksetRef string `json:"tricksetRef,omitempty" yaml:"tricksetRef,omitempty"`
}
