package elligen

import (
	"fmt"
	"math/rand"
	"strings"
)

func DefaultGenerator(options ...DefaultGeneratorOption) (Generator, error) {
	generator := &defaultGenerator{
		r:          RandomIntnFunc(rand.Intn),
		consonants: consonants,
		vowels:     vowels,
		endings:    endings,
	}
	for i := range options {
		err := options[i].configure(generator)
		if err != nil {
			return nil, err
		}
	}
	return generator, nil
}

type DefaultGeneratorOption interface {
	configure(*defaultGenerator) error
}

type defaultGeneratorOptionFunc func(*defaultGenerator) error

func (f defaultGeneratorOptionFunc) configure(generator *defaultGenerator) error {
	return f(generator)
}

// WithOptions takes a bunch of options and combines them into a single option.
func WithOptions(options ...DefaultGeneratorOption) DefaultGeneratorOption {
	return defaultGeneratorOptionFunc(
		func(generator *defaultGenerator) error {
			for i := range options {
				if err := options[i].configure(generator); err != nil {
					return err
				}
			}
			return nil
		},
	)
}

func WithConsonants(consonants []string) DefaultGeneratorOption {
	return defaultGeneratorOptionFunc(
		func(generator *defaultGenerator) error {
			if len(consonants) == 0 {
				return fmt.Errorf("list of consonants is empty")
			}
			generator.consonants = consonants
			return nil
		},
	)
}

func WithVowels(vowels []string) DefaultGeneratorOption {
	return defaultGeneratorOptionFunc(
		func(generator *defaultGenerator) error {
			if len(vowels) == 0 {
				return fmt.Errorf("list of vowels is empty")
			}
			generator.vowels = vowels
			return nil
		},
	)
}

func WithEndings(endings []string) DefaultGeneratorOption {
	return defaultGeneratorOptionFunc(
		func(generator *defaultGenerator) error {
			if len(endings) == 0 {
				return fmt.Errorf("list of endings is empty")
			}
			generator.endings = endings
			return nil
		},
	)
}

func WithRandomIntn(r RandomIntn) DefaultGeneratorOption {
	return defaultGeneratorOptionFunc(
		func(generator *defaultGenerator) error {
			if r == nil {
				return fmt.Errorf("source of randomness must not be nil")
			}
			generator.r = r
			return nil
		},
	)
}

var consonants = []string{
	"b", "c", "d", "f", "g", "h", "j", "k", "l", "m", "n", "p", "q", "r", "s", "t", "v", "w", "x", "z",
}

var vowels = []string{
	"a", "e", "i", "o", "u", "y",
}

var endings = []string{
	"alt", "en", "ese", "ic", "ine", "ium", "on", "um", "us",
}

type defaultGenerator struct {
	r RandomIntn

	consonants []string
	vowels     []string
	endings    []string
}

type defaultGeneratorPart string

func (part defaultGeneratorPart) empty() bool {
	return len(part) == 0
}

func (part defaultGeneratorPart) endsWithOneOf(list []string) bool {
	for i := range list {
		if strings.HasSuffix(string(part), list[i]) {
			return true
		}
	}
	return false
}

func (generator *defaultGenerator) Generate() string {
	partCount := generator.r.Intn(3) + 2 // 2 to 4
	parts := make([]defaultGeneratorPart, partCount)
	if generator.r.Intn(2) == 0 {
		parts[0] = generateDefaultPart(generator.r, generator.consonants, generator.vowels)
	} else {
		parts[0] = generateDefaultPart(generator.r, generator.vowels, generator.consonants)
	}
	for i := 1; i < partCount; i++ {
		first, second := generator.consonants, generator.vowels
		if parts[i-1].endsWithOneOf(generator.consonants) && generator.r.Intn(2) == 0 {
			first, second = generator.vowels, generator.consonants
		}
		parts[i] = generateDefaultPart(generator.r, first, second)
	}
	builder := &strings.Builder{}
	for i := range parts {
		builder.WriteString(string(parts[i]))
	}
	if parts[partCount-1].endsWithOneOf(generator.vowels) {
		builder.WriteString(generator.consonants[generator.r.Intn(len(generator.consonants))])
	}
	builder.WriteString(generator.endings[generator.r.Intn(len(generator.endings))])
	return builder.String()
}

func generateDefaultPart(r RandomIntn, first, second []string) defaultGeneratorPart {
	return defaultGeneratorPart(first[r.Intn(len(first))] + second[r.Intn(len(second))])
}

// RandomIntn is the source of randomness used by the default generator.
type RandomIntn interface {
	Intn(n int) int
}

// RandomIntnFunc is a convenience wrapper so rand.Intn (and the like) can
// be used as RandomIntn easier.
type RandomIntnFunc func(n int) int

// Intn just calls f and returns the result.
func (f RandomIntnFunc) Intn(n int) int {
	return f(n)
}
