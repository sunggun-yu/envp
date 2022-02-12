package prompt

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/manifoldco/promptui"
)

func TestPromptConfirm(t *testing.T) {

	var testCase = func(input string) bool {

		// PromptUI need
		s := fmt.Sprintf("%s\n", input)
		stringReader := strings.NewReader(s)
		stringReadCloser := io.NopCloser(stringReader)
		p := &promptui.Prompt{
			Label:     "label",
			IsConfirm: true,
			IsVimMode: true,
			Stdin:     stringReadCloser,
		}
		prom := promptConfirm{
			prompt: p,
		}
		return prom.run()
	}

	t.Run("default should false", func(t *testing.T) {
		if PromptConfirm("some message") {
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
