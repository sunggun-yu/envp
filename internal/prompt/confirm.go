package prompt

import (
	"io"
	"strings"

	"github.com/manifoldco/promptui"
)

// PromptConfirm prompts yes/no
type PromptConfirm struct {
	prompt *promptui.Prompt
	label  string
}

// NewPromptConfirm create new PromptConfirm
func NewPromptConfirm(label string) PromptConfirm {
	return PromptConfirm{
		label: label,
		prompt: &promptui.Prompt{
			Label:     label,
			IsConfirm: true,
		},
	}
}

// run the prompt
func (p *PromptConfirm) run() bool {
	a, _ := p.prompt.Run()
	yeses := []string{"y", "yes", "yeah"}
	// only consider yeses as true. others are false.
	for _, y := range yeses {
		if strings.EqualFold(y, a) {
			return true
		}
	}
	return false
}

// SetIn sets the source for input data
// If newIn is nil, os.Stdin is used.
func (p *PromptConfirm) SetIn(in io.Reader) {
	p.prompt.Stdin = io.NopCloser(in)
}

// Prompt ask confirmation of y/N
func (p *PromptConfirm) Prompt() bool {
	return p.run()
}
