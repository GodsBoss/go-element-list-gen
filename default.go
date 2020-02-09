package elligen

import (
	"math/rand"
	"strings"
)

func DefaultGenerator(options ...DefaultGeneratorOption) (Generator, error) {
	generator := &defaultGenerator{
		r:          randomIntnFunc(rand.Intn),
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
	r randomIntn

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

func generateDefaultPart(r randomIntn, first, second []string) defaultGeneratorPart {
	return defaultGeneratorPart(first[r.Intn(len(first))] + second[r.Intn(len(second))])
}

type randomIntn interface {
	Intn(n int) int
}

type randomIntnFunc func(n int) int

func (f randomIntnFunc) Intn(n int) int {
	return f(n)
}
