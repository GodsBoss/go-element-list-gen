package elligen

import (
	"fmt"
)

// Generator is everything which generates a single element name.
type Generator interface {
	// Generate generates a single element name.
	Generate() string
}

// GeneratorFunc implements Generator by wrapping a function with the signature
// of Generator#Generate().
type GeneratorFunc func() string

// Generate calls f and returns its result.
func (f GeneratorFunc) Generate() string {
	return f()
}

// Generate attempts to generate count element names, avoiding duplicates.
// If the desired list could not be generated after the given maximum amount
// of attempts, an error is returned. A list containing names will still be
// generated.
// attempts determines the number of attempts on generating the full list.
// Negative values are equal to zero.
// Generate will make the minimum amount of calls to Generator#Generate().
func Generate(generator Generator, count int, attempts int) ([]string, error) {
	if count < 0 {
		return nil, negativeGenerateCountError(count)
	}
	remainingAttempts := attempts
	if remainingAttempts < 0 {
		remainingAttempts = 0
	}
	generated := map[string]struct{}{}
	for remainingAttempts >= 0 && len(generated) < count {
		remainingAttempts--
		missing := count - len(generated)
		for i := 0; i < missing; i++ {
			generated[generator.Generate()] = struct{}{}
		}
	}
	list := make([]string, 0, len(generated))
	for s := range generated {
		list = append(list, s)
	}
	var err error
	if len(list) != count {
		err = attemptsExhaustedError{
			attempts: attempts,
			count:    count,
			created:  len(list),
		}
	}
	return list, err
}

type negativeGenerateCountError int

func (err negativeGenerateCountError) negativeCount() int {
	return int(err)
}

func (err negativeGenerateCountError) Error() string {
	return fmt.Sprintf("count was %d, which is negative", int(err))
}

// IsNegativeGenerateCountError checks wether an error corresponds to a negative
// count given to Generate. If yes, also returns the count. If no, the count is
// not supposed to have any meaning.
func IsNegativeGenerateCountError(err error) (bool, int) {
	if castedErr, ok := err.(interface {
		negativeCount() int
	}); ok {
		return true, castedErr.negativeCount()
	}
	return false, 0
}

type attemptsExhaustedError struct {
	attempts int
	count    int
	created  int
}

func (err attemptsExhaustedError) Error() string {
	return fmt.Sprintf("created %d instead of %d names exhausting %d attempts", err.created, err.count, err.attempts)
}

// IsAttemptsExhaustedError checks wether exhausting attempts caused an error.
// attempts and count are the parameters given to Generate(), created is the
// number of names created.
func IsAttemptsExhaustedError(err error) (ok bool, attempts int, count int, created int) {
	if castedErr, ok := err.(attemptsExhaustedError); ok {
		return true, castedErr.attempts, castedErr.count, castedErr.created
	}
	return false, 0, 0, 0
}
