package prompt

import (
	"strings"

	"github.com/manifoldco/promptui"
)

// PromptYesOrNo prompts yes/no
// TODO: how to test in code?
type promptConfirm struct {
	prompt *promptui.Prompt
	label  string
}

func (p *promptConfirm) run() bool {
	if p.prompt == nil {
		p.prompt = &promptui.Prompt{
			Label:     p.label,
			IsConfirm: true,
		}
	}
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

func PromptConfirm(label string) bool {
	prom := &promptConfirm{
		label: label,
	}
	return prom.run()
}
