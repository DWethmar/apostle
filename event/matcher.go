package event

import "slices"

// Matcher is an interface that defines a method to match events.
type Matcher interface {
	// Match returns true if the given event matches the matcher.
	Match(e Event) bool
}

// MatcherFunc is a function that implements the Matcher interface.
type MatcherFunc func(e Event) bool

// Match calls the function itself.
func (f MatcherFunc) Match(e Event) bool {
	return f(e)
}

func MatchAny(t ...string) MatcherFunc {
	return func(e Event) bool { return slices.Contains(t, e.Event()) }
}
