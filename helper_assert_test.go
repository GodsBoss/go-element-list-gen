package elligen_test

import (
	"testing"
)

func assertNoErr(t *testing.T, err error) {
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func assertListContains(t *testing.T, list []string, needle string) {
	for i := range list {
		if list[i] == needle {
			return
		}
	}
	t.Errorf("%s not found in %v", needle, list)
}
