package prompt

import (
	"fmt"
	"github.com/itohin/gophkeeper/internal/infrastructure/logger"
	"github.com/manifoldco/promptui"
)

type Prompt struct {
	log logger.Logger
}

func NewPrompt(logger logger.Logger) *Prompt {
	return &Prompt{log: logger}
}

type PromptContent struct {
	ErrorMsg string
	Label    string
	Mask     rune
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
	index := -1
	var result string
	var err error

	for index < 0 {
		prompt := promptui.SelectWithAdd{
			Label: pc.Label,
			Items: items,
		}

		index, result, err = prompt.Run()

		if index == -1 {
			items = append(items, result)
		}
	}

	if err != nil {
		return result, fmt.Errorf("Prompt failed %v\n", err)
	}

	return result, nil
}
