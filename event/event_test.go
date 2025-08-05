package event_test

import (
	"testing"

	"github.com/dwethmar/apostle/event"
)

// MockEvent is a simple implementation of the Event interface for testing.
type MockEvent struct {
	event string
}

func (e *MockEvent) Event() string {
	return e.event
}

func TestNewBus(t *testing.T) {
	bus := event.NewBus(0)
	if bus == nil {
		t.Fatalf("NewBus() returned nil")
	}

	if len(bus.Subscriptions()) != 0 {
		t.Errorf("NewBus() handlers map should be empty initially, got %d", len(bus.Subscriptions()))
	}
}

func TestBus_Subscribe(t *testing.T) {
	bus := event.NewBus(0)
	handlerCalled := false

	bus.Subscribe(event.MatcherFunc(func(e event.Event) bool {
		return e.Event() == "testEvent"
	}), func(_ event.Event) error {
		handlerCalled = true
		return nil
	})

	if len(bus.Subscriptions()) != 1 {
		t.Errorf("Bus.Subscribe() should add a handler, got %d", len(bus.Subscriptions()))
	}

	// Trigger the event to check if the handler is called.
	if err := bus.Publish(&MockEvent{event: "testEvent"}); err != nil {
		t.Errorf("Bus.Publish() error = %v", err)
	}

	if !handlerCalled {
		t.Errorf("Handler was not called after Bus.Publish()")
	}
}

func TestBus_Publish(t *testing.T) {
	t.Run("publish an event", func(t *testing.T) {
		bus := event.NewBus(0)
		handlerCalled := false

		bus.Subscribe(event.MatcherFunc(func(e event.Event) bool {
			return e.Event() == "testEvent"
		}), func(_ event.Event) error {
			handlerCalled = true
			return nil
		})

		if err := bus.Publish(&MockEvent{event: "testEvent"}); err != nil {
			t.Errorf("Bus.Publish() error = %v", err)
		}

		if !handlerCalled {
			t.Errorf("Bus.Publish() did not trigger the handler")
		}

		// Test with an event that has no handlers
		handlerCalled = false
		if err := bus.Publish(&MockEvent{event: "nonExistentEvent"}); err != nil {
			t.Errorf("Bus.Publish() error = %v", err)
		}

		if handlerCalled {
			t.Errorf("Bus.Publish() should not trigger a handler for an event with no subscribers")
		}
	})
}

func TestBus_Unsubscribe(t *testing.T) {
	bus := event.NewBus(0)
	handlerCalled := false

	handler := func(_ event.Event) error {
		handlerCalled = true
		return nil
	}

	id := bus.Subscribe(event.MatcherFunc(func(e event.Event) bool {
		return e.Event() == "testEvent"
	}), handler)

	if err := bus.Publish(&MockEvent{event: "testEvent"}); err != nil {
		t.Errorf("Bus.Publish() error = %v", err)
	}

	if !handlerCalled {
		t.Errorf("Handler was not called after Bus.Publish()")
	}

	handlerCalled = false
	bus.Unsubscribe(id)
	if err := bus.Publish(&MockEvent{event: "testEvent"}); err != nil {
		t.Errorf("Bus.Publish() error = %v", err)
	}

	if handlerCalled {
		t.Errorf("Handler was called after Bus.Unsubscribe()")
	}
}

func TestBus_Subscriptions(t *testing.T) {
	t.Run("get all subscriptions", func(t *testing.T) {
		bus := event.NewBus(0)

		for range 10 {
			bus.Subscribe(event.MatcherFunc(func(e event.Event) bool {
				return e.Event() == "testEvent"
			}), func(_ event.Event) error { return nil })
		}

		subscriptions := bus.Subscriptions()

		if len(subscriptions) != 10 {
			t.Errorf("Bus.Subscriptions() should return all subscriptions, got %d", len(subscriptions))
		}
	})
}
