package survey

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/djgilcrease/survey/core"
	"gopkg.in/AlecAivazis/survey.v1/terminal"
)

func init() {
	// disable color output for all prompts to simplify testing
	core.DisableColor = true
}

func TestInputRender(t *testing.T) {

	tests := []struct {
		title    string
		prompt   Prompt
		data     InputTemplateData
		expected string
	}{
		{
			"Test Input question output without default",
			NewInput().SetMessage("What is your favorite month:"),
			InputTemplateData{},
			"? What is your favorite month: ",
		},
		{
			"Test Input question output with default",
			NewInput().SetMessage("What is your favorite month:").SetDefault("April"),
			InputTemplateData{},
			"? What is your favorite month: (April) ",
		},
		{
			"Test Input answer output",
			NewInput().SetMessage("What is your favorite month:"),
			InputTemplateData{Answer: "October", ShowAnswer: true},
			"? What is your favorite month: October\n",
		},
		{
			"Test Input question output without default but with help hidden",
			NewInput().SetMessage("What is your favorite month:").SetHelp("This is helpful"),
			InputTemplateData{},
			"? What is your favorite month: [? for help] ",
		},
		{
			"Test Input question output with default and with help hidden",
			NewInput().SetMessage("What is your favorite month:").SetDefault("April").SetHelp("This is helpful"),
			InputTemplateData{},
			"? What is your favorite month: [? for help] (April) ",
		},
		{
			"Test Input question output without default but with help shown",
			NewInput().SetMessage("What is your favorite month:").SetHelp("This is helpful"),
			InputTemplateData{ShowHelp: true},
			`ⓘ This is helpful
? What is your favorite month: `,
		},
		{
			"Test Input question output with default and with help shown",
			NewInput().SetMessage("What is your favorite month:").SetDefault("April").SetHelp("This is helpful"),
			InputTemplateData{ShowHelp: true},
			`ⓘ This is helpful
? What is your favorite month: (April) `,
		},
	}

	outputBuffer := bytes.NewBufferString("")
	terminal.Stdout = outputBuffer

	for _, test := range tests {
		outputBuffer.Reset()
		test.data.Input = test.prompt.(*Input)
		err := test.data.Render(
			test.data.tmpl,
			test.data,
		)
		assert.Nil(t, err, test.title)
		assert.Equal(t, test.expected, outputBuffer.String(), test.title)
	}
}
