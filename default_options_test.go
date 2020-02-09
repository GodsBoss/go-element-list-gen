package elligen_test

import (
	"testing"

	elligen "github.com/GodsBoss/go-element-list-gen"
)

func TestInvalidDefaultGeneratorOptions(t *testing.T) {
	invalidOptions := map[string]elligen.DefaultGeneratorOption{
		"WithConsonants":             elligen.WithConsonants(make([]string, 0)),
		"WithVowels":                 elligen.WithVowels(make([]string, 0)),
		"WithEndings":                elligen.WithEndings(make([]string, 0)),
		"WithRandomIntn":             elligen.WithRandomIntn(nil),
		"WithMinimumParts":           elligen.WithMinimumParts(0),
		"WithMaximumAdditionalParts": elligen.WithMaximumAdditionalParts(-5),
		"WithOptions":                elligen.WithOptions(elligen.WithRandomIntn(nil)),
	}

	for name := range invalidOptions {
		t.Run(
			name,
			func(option elligen.DefaultGeneratorOption) func(*testing.T) {
				return func(*testing.T) {
					generator, err := elligen.DefaultGenerator(option)
					if generator != nil {
						t.Errorf("expected creating default generator to fail, but got %+v", generator)
					}
					if err == nil {
						t.Errorf("expected error")
					}
				}
			}(invalidOptions[name]),
		)
	}
}
