package survey

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/djgilcrease/survey/core"
	"gopkg.in/AlecAivazis/survey.v1/terminal"
	"os"
	"io"
)

func init() {
	// disable color output for all prompts to simplify testing
	core.DisableColor = true
}

func TestConfirmRender(t *testing.T) {

	tests := []struct {
		title    string
		prompt   Prompt
		data     ConfirmTemplateData
		expected string
	}{
		{
			"Test Confirm question output with default true",
			NewConfirm().SetMessage("Is pizza your favorite food?").SetDefault(true),
			ConfirmTemplateData{},
			`? Is pizza your favorite food? (Y/n) `,
		},
		{
			"Test Confirm question output with default false",
			NewConfirm().SetMessage("Is pizza your favorite food?").SetDefault(false),
			ConfirmTemplateData{},
			`? Is pizza your favorite food? (y/N) `,
		},
		{
			"Test Confirm answer output",
			NewConfirm().SetMessage("Is pizza your favorite food?"),
			ConfirmTemplateData{Answer: "Yes"},
			"? Is pizza your favorite food? Yes\n",
		},
		{
			"Test Confirm with help but help message is hidden",
			NewConfirm().SetMessage("Is pizza your favorite food?").SetHelp("This is helpful"),
			ConfirmTemplateData{},
			"? Is pizza your favorite food? [? for help] (y/N) ",
		},
		{
			"Test Confirm help output with help message shown",
			NewConfirm().SetMessage("Is pizza your favorite food?").SetHelp("This is helpful"),
			ConfirmTemplateData{ShowHelp: true},
			`ⓘ This is helpful
? Is pizza your favorite food? (y/N) `,
		},
	}


	for _, test := range tests {
		r, w, err := os.Pipe()
		assert.Nil(t, err, test.title)
		test.data.Confirm = test.prompt.(*Confirm)
		test.data.WithStdio(terminal.Stdio{Out: w})

		err = test.data.Render(
			DefaultConfirmQuestionTemplate,
			test.data,
		)
		assert.Nil(t, err, test.title)

		w.Close()
		var buf bytes.Buffer
		io.Copy(&buf, r)

		assert.Contains(t, buf.String(), test.expected, test.title)
	}
}
