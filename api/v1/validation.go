package v1

import "fmt"

func (ts *TrickSet) Validate() error {
	if ts.GetVersion() != Version {
		return fmt.Errorf("invalid version: %s", ts.GetVersion())
	}

	if err := ts.ValidateContext(); err != nil {
		return fmt.Errorf("could not validate context: %w", err)
	}

	if err := ts.ValidateRules(); err != nil {
		return fmt.Errorf("could not validate rules: %w", err)
	}

	return nil
}

func (ts *TrickSet) ValidateContext() error {
	switch ts.Context.Provider {
	case GitContext:
		if ts.Context.Git != nil {
			return nil
		}
	case GitHubContext:
		if ts.Context.GitHub != nil {
			return nil
		}
	default:
		return fmt.Errorf("invalid context provider: %s", ts.Context.Provider)
	}

	return fmt.Errorf("invalid context configuration for provider: %s",
		ts.Context.Provider)
}

func (ts *TrickSet) ValidateRules() error {
	for _, rule := range ts.Rules {
		if rule.Validate() != nil {
			return fmt.Errorf("could not validate rule: %w", rule.Validate())
		}
	}

	return nil
}

func (rule *RuleDefinition) Validate() error {
	if rule.Name == "" {
		return fmt.Errorf("rule must have a name")
	}

	policyDeclarations := 0
	if rule.Include != "" {
		policyDeclarations++
	}
	if rule.InlinePolicy != "" {
		policyDeclarations++
	}
	if rule.RulesFrom != nil {
		policyDeclarations++
	}
	if len(rule.Or) > 0 {
		policyDeclarations++
	}

	if policyDeclarations != 1 {
		return fmt.Errorf("rule must have exactly one policy declaration")
	}

	if len(rule.Or) > 0 {
		for _, orRule := range rule.Or {
			if err := orRule.Validate(); err != nil {
				return fmt.Errorf("could not validate 'or' rule: %w", err)
			}
		}
	}

	if rule.RulesFrom != nil {
		if (rule.RulesFrom.File == "" && rule.RulesFrom.TricksetRef == "") ||
			(rule.RulesFrom.File != "" && rule.RulesFrom.TricksetRef != "") {
			return fmt.Errorf("rulesFrom must have either a file or a tricksetRef")
		}
	}

	return nil
}
