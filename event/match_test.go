package event_test

import (
	"testing"

	"github.com/dwethmar/apostle/event"
)

func TestMatcherFunc_Match(t *testing.T) {
	t.Run("Match should return true if the event matches the matcher", func(t *testing.T) {
		matcher := event.MatcherFunc(func(e event.Event) bool {
			return e.Event() == "test"
		})
		if !matcher.Match(&MockEvent{event: "test"}) {
			t.Errorf("expected event to match")
		}
	})
}

func TestMatchAny(t *testing.T) {
	t.Run("MatchAny should return a MatcherFunc that matches any event in the given list", func(t *testing.T) {
		matcher := event.MatchAny("test1", "test2")
		if !matcher.Match(&MockEvent{event: "test1"}) {
			t.Errorf("expected event to match")
		}
		if !matcher.Match(&MockEvent{event: "test2"}) {
			t.Errorf("expected event to match")
		}
		if matcher.Match(&MockEvent{event: "test3"}) {
			t.Errorf("expected event to not match")
		}
	})
}
