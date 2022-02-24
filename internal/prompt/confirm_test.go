package prompt

import (
	"fmt"
	"io"
	"strings"
	"testing"
)

func TestPromptConfirm(t *testing.T) {
	// var stdout bytes.Buffer

	var testCase = func(input string) bool {

		// PromptUI need
		s := fmt.Sprintf("%s\n", input)
		stringReader := strings.NewReader(s)
		stringReadCloser := io.NopCloser(stringReader)

		prom := NewPromptConfirm("label")
		prom.SetIn(stringReadCloser)
		// prom.SetOut(bufio.NewWriter(&stdout))
		return prom.Prompt()
	}

	t.Run("default should false", func(t *testing.T) {
		if testCase("some message") {
			t.Error("default should false")
		}
	})

	t.Run("y should true", func(t *testing.T) {
		if !testCase("y") {
			t.Error("should true")
		}
	})

	t.Run("Y should true", func(t *testing.T) {
		if !testCase("Y") {
			t.Error("should true")
		}
	})

	t.Run("Yes should true", func(t *testing.T) {
		if !testCase("Yes") {
			t.Error("should true")
		}
	})

	t.Run("n should false", func(t *testing.T) {
		if testCase("n") {
			t.Error("should false")
		}
	})

	t.Run("empty string should false", func(t *testing.T) {
		if testCase("") {
			t.Error("should false")
		}
	})
}
