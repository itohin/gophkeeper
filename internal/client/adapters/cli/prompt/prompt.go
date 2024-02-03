package prompt

import (
	"fmt"
	"github.com/manifoldco/promptui"
)

type Prompter interface {
	PromptGetInput(pc PromptContent, validate func(input string) error) (string, error)
	PromptGetSelect(pc PromptContent, items []SelectItem) (string, error)
}

type Prompt struct{}

func NewPrompt() *Prompt {
	return &Prompt{}
}

type PromptContent struct {
	Label string
	Mask  rune
}

type SelectItem struct {
	Label  string
	Action string
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

func (p *Prompt) PromptGetSelect(pc PromptContent, items []SelectItem) (string, error) {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "\u21E8 {{ .Label | cyan }}",
		Inactive: "  {{ .Label | cyan }}",
		Selected: "\u21E8 {{ .Label | red | cyan }}",
	}
	prompt := promptui.Select{
		Label:     pc.Label,
		Items:     items,
		Templates: templates,
	}

	i, s, err := prompt.Run()

	if err != nil {
		return s, fmt.Errorf("Prompt failed %v\n", err)
	}

	return items[i].Action, nil
}
