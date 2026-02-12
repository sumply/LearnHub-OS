package get_quizzes

type OptionScope string

const (
	ScopeEmpty     OptionScope = ""
	ScopeOwner     OptionScope = "owner"
	ScopeAvailable OptionScope = "available"
)

func (q OptionScope) Name() string {
	return "scope"
}

func (q OptionScope) Variants() []string {
	return []string{
		string(ScopeOwner),
		string(ScopeAvailable),
	}
}

func (o OptionScope) Validate() bool {
	switch o {
	case ScopeOwner, ScopeAvailable, ScopeEmpty:
		return true
	default:
		return false
	}
}

type OptionInclude string

const (
	IncludeEmpty       OptionInclude = ""
	IncludeLastAttempt OptionInclude = "last_attempt"
)

func (o OptionInclude) Validate() bool {
	switch o {
	case IncludeEmpty, IncludeLastAttempt:
		return true
	default:
		return false
	}
}

func (q OptionInclude) Name() string {
	return "include"
}

func (q OptionInclude) Variants() []string {
	return []string{
		string(IncludeLastAttempt),
	}
}

type Options struct {
	Scope   OptionScope
	Include OptionInclude
}
