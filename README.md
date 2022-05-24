# Lammaskoira / bark

`bark` is a utility to evaluate contexts against given policies.

A context may be a Git or a GitHub repository (more contexts may
come in the future).

In this repo the format to define policies is defined as well as
the `bark` program to evaluate them.

## Lammaskoira?

Lammaskoira is Finnish for "sheep dog"; a lot of the analogies
and concepts from this project come from sheep dog aspects.

## TrickSets

You'll teach a sheep dog tricks in order to guide lambs,
make them stay with the flock and get them safely to the
destination.

As such, a trickset is a file containing the policy we want
to evaluate, while also containing the context we want to
evaluate it against.

A sample looks as follows:

```yaml
---
version: v1
context:
  provider: github
  github:
    org: lammaskoira
rules:
  - name: Should have renovate configured
    inlinePolicy:  |
      package bark

      default allow := false

      allow {
          file.exists("./renovate.json")
      }

      allow {
          file.exists("./.github/renovate.json")
      }
  
  - name: Should have CodeQL analysis configured
    inlinePolicy:  |
      package bark

      default allow := false

      allow {
          some i
          workflowstr := file.readall("./.github/workflows/codeql-analysis.yml")
          workflow := yaml.unmarshal(workflowstr)
          steps := workflow.jobs.analyze.steps[i]
          contains(steps.uses, "github/codeql-action/analyze@")
      }
```

This trickset will evaluate the `lammaskoira` GitHub organization
and evaluate if the policy is fulfilled for all repositories there.

There are two rules in this sample:

* One that checks if the repository has a `renovate.json` file
  or a `.github/renovate.json` file.
* One that checks if the repository has a CodeQL analysis configured.

### Contexts

Currently, the following contexts are defined:

* `git`

* `github`

* `githubOrgConfig`


#### git

`git` is a context that is a single Git repository. It allows
for specifying the Git URL and branch to verify. A sample looks as
follows:

```yaml
context:
  provider: git
  git:
    url: https://github.com/lammaskoira/bark.git
    branch: main
```

One must always define the provider to be used, and specify the provider
configuration in the context.

#### github

`github` is a context allows for evaluating policies on GitHub repositories.
It allows for specifying the GitHub organization, which will verify all
repositories in that organization. A sample looks as follows:

```yaml
context:
  provider: github
  github:
    org: lammaskoira
```

While running the `bark` program, you can specify a GitHub token to
use for authentication. This is possible via the `GITHUB_TOKEN` environment
variable.

It's also possible to evaluate a policy against the repository metadata
retrived from the GitHub API. The current implementation adds the following keys
to the rego input:

* `repometa`: The repository metadata from the GitHub API.
  This comes from a [GET request to the GitHub API](https://docs.github.com/en/rest/repos/repos#get-a-repository).

* `vulnerability_alerts_enabled`: It's a boolean that indicates if the repository
  has the `vulnerability-alerts` feature enabled.

#### githubOrgConfig

`githubOrgConfig` is a context that allows for evaluating policies
on GitHub organizations. It allows for specifying the GitHub organization
to evaluate. A sample looks as follows:

```yaml
context:
  provider: githubOrgConfig
  githubOrgConfig:
    org: lammaskoira
```

While running the `bark` program, you can specify a GitHub token to
use for authentication. This is possible via the `GITHUB_TOKEN` environment
variable.

Policy evaluation relies entirely on the Organization information
retrieved from the GitHub API. The current implementation adds the following keys:

* `orgconfig`: The organization configuration from the GitHub API.
  This comes from a [GET request to the GitHub API](https://docs.github.com/en/rest/reference/orgs#get-an-organization).

### Policy language

The policy format is [rego](https://www.openpolicyagent.org/docs/latest/policy-language/)
which gives us a fairly versatile and powerful language to define
the policies.

`bark` runs [Open Policy Agent](https://www.openpolicyagent.org/docs/latest/)
to evaluate the policies.

### Policy assumptions

for the current `v1` version of the language, each policy must
use the `bark` package:

```rego
package bark
```

Each policy must also return a single boolean value:

```rego
default allow := false

allow {
    <your policy>
}
```

By default, the examples use the `allow` key.

### rego extensions

In order to allow the policies to be evaluated against the
contents of the repository, we need to define a couple of
extensions to the rego language. This is done by adding the
following builtin functions:

* `file.readall(path)`: reads the file at `path` and returns its contents.

* `file.exists(path)`: checks if the file at `path` exists.

More extensions will be added as needed.

## Building `bark`

In this repository do:

```bash
$ go build -o bark main.go
```

## Running `bark`

```bash
$ export GITHUB_TOKEN=<your token>
$ sudo -E ./bark -t trickset.yml
```

**Note:** the `sudo` is needed because `bark` will limit
OPA's access to the host machine to only the context. e.g.
`bark` clones the Git repository in a temporary directory,
changes the working directory towards the aforementioned
directory, `chroot`'s into that directory and then runs
the policy.