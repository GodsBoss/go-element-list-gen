package elligen_test

import (
	"fmt"
	"testing"

	"github.com/GodsBoss/elligen"
)

func TestGenerateNegativeCount(t *testing.T) {
	f := elligen.GeneratorFunc(
		func() string {
			return ""
		},
	)
	list, err := elligen.Generate(f, -7, 1000)
	if list != nil {
		t.Errorf("expected generated element names to be nil, but got %+v", list)
	}
	isNegativeCountErr, count := elligen.IsNegativeGenerateCountError(err)
	if !isNegativeCountErr {
		t.Errorf("expected error to be negative count error")
	}
	if count != -7 {
		t.Errorf("expected negative count extracted from error to be -7, but got %d", count)
	}
}

func TestGenerateCountIsZero(t *testing.T) {
	f := elligen.GeneratorFunc(
		func() string {
			return ""
		},
	)
	list, err := elligen.Generate(f, 0, 1000)
	if list == nil {
		t.Errorf("expected generated list not to be nil")
	}
	if len(list) != 0 {
		t.Errorf("expected empty list, not %v", list)
	}
	assertNoErr(t, err)
}

func TestGenerateSucceeds(t *testing.T) {
	current := 'a'
	f := elligen.GeneratorFunc(
		func() string {
			ret := string(current)
			current++
			return ret
		},
	)
	list, err := elligen.Generate(f, 5, 1)
	assertNoErr(t, err)
	if len(list) != 5 {
		t.Errorf("expected list with 5 elements, got %v", list)
	}
	assertListContains(t, list, "a")
	assertListContains(t, list, "b")
	assertListContains(t, list, "c")
	assertListContains(t, list, "d")
	assertListContains(t, list, "e")
}

func TestGenerateSucceedsWithNegativeAttempts(t *testing.T) {
	current := 'a'
	f := elligen.GeneratorFunc(
		func() string {
			ret := string(current)
			current++
			return ret
		},
	)
	list, err := elligen.Generate(f, 5, -2)
	assertNoErr(t, err)
	if len(list) != 5 {
		t.Errorf("expected list with 5 elements, got %v", list)
	}
}

func TestGenerateCreatesPartialList(t *testing.T) {
	values := []string{"foo", "bar", "baz"}
	current := 0
	f := elligen.GeneratorFunc(
		func() string {
			ret := values[current]
			current = (current + 1) % len(values)
			return ret
		},
	)
	list, err := elligen.Generate(f, 5, 100)
	if len(list) != 3 {
		t.Errorf("expected list with 3 elements, got %v", list)
	}
	assertListContains(t, list, "foo")
	assertListContains(t, list, "bar")
	assertListContains(t, list, "baz")
	isExhaustedAttemptsErr, attempts, count, created := elligen.IsAttemptsExhaustedError(err)
	if !isExhaustedAttemptsErr {
		t.Errorf("expected attempts exhaustion error, not %v", err)
	}
	if attempts != 100 {
		t.Errorf("expected attempts extracted from error to be 100, got %d", attempts)
	}
	if count != 5 {
		t.Errorf("expected count extracted from error to be 5, got %d", count)
	}
	if created != 3 {
		t.Errorf("expected number of created names to be 3, got %d", created)
	}
}

func TestIsNotNegativeGenerateCountError(t *testing.T) {
	err := fmt.Errorf("just a random error")
	isNegativeCountErr, _ := elligen.IsNegativeGenerateCountError(err)
	if isNegativeCountErr {
		t.Errorf("expected %v not to be a negative count error", err)
	}
}

func TestIsNotAttemptsExhaustedError(t *testing.T) {
	err := fmt.Errorf("just a random error")
	isAttemptsExhaustedErr, _, _, _ := elligen.IsAttemptsExhaustedError(err)
	if isAttemptsExhaustedErr {
		t.Errorf("expected %v not to be an attempts exhaustion error", err)
	}
}
