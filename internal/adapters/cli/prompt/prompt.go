package prompt

import (
	"fmt"
	"github.com/manifoldco/promptui"
)

type Prompter interface {
	PromptGetInput(pc PromptContent, validate func(input string) error) (string, error)
	PromptGetSelect(pc PromptContent, items []string) (string, error)
}

type Prompt struct{}

func NewPrompt() *Prompt {
	return &Prompt{}
}

type PromptContent struct {
	Label string
	Mask  rune
}

func (p *Prompt) PromptGetInput(pc PromptContent, validate func(input string) error) (string, error) {
	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }}",
		Valid:   "{{ . | green }}",
		Invalid: "{{ . | red }}",
		Success: "{{ . | bold }}",
	}

	prompt := promptui.Prompt{
		Label:     pc.Label,
		Templates: templates,
		Validate:  validate,
		Mask:      pc.Mask,
	}

	result, err := prompt.Run()
	if err != nil {
		return result, fmt.Errorf("prompt failed %v\n", err)
	}
	return result, nil
}

func (p *Prompt) PromptGetSelect(pc PromptContent, items []string) (string, error) {
	prompt := promptui.Select{
		Label: pc.Label,
		Items: items,
	}

	_, result, err := prompt.Run()

	if err != nil {
		return result, fmt.Errorf("Prompt failed %v\n", err)
	}

	return result, nil
}
