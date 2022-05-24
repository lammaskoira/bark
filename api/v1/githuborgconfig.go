package v1

type GitHubOrgConfigDefinition struct {
	// Org is the organization that owns the repos
	Org string `json:"org" yaml:"org"`
}

const (
	OrgConfigInputKey = "orgconfig"
	AppsInOrgInputKey = "apps_in_org"
)
