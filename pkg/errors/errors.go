package errors

import "errors"

var (
	// ErrInvalidType is returned when the instantiated TrickSet
	// is not of the right type.
	ErrInvalidType = errors.New("invalid type")

	// ErrPolicyDenial is returned when a policy denies the query.
	ErrPolicyDenial = errors.New("denied by policy")

	// ErrPolicyParseError is returned when a policy cannot be parsed.
	ErrPolicyParseError = errors.New("policy parse error")

	// ErrPolicyEvalError is returned when a policy cannot be evaluated.
	ErrPolicyEvalError = errors.New("policy eval error")
)
