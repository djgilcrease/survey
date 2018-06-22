package survey

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/djgilcrease/survey/core"
	"os"
	"gopkg.in/AlecAivazis/survey.v1/terminal"
	"bytes"
	"io"
)

func init() {
	// disable color output for all prompts to simplify testing
	core.DisableColor = true
}

func TestEditorRender(t *testing.T) {
	tests := []struct {
		title    string
		prompt   Editor
		data     EditorTemplateData
		expected string
	}{
		{
			"Test Editor question output without default",
			Editor{Message: "What is your favorite month:"},
			EditorTemplateData{},
			"? What is your favorite month: [Enter to launch editor] ",
		},
		{
			"Test Editor question output with default",
			Editor{Message: "What is your favorite month:", Default: "April"},
			EditorTemplateData{},
			"? What is your favorite month: (April) [Enter to launch editor] ",
		},
		{
			"Test Editor question output with HideDefault",
			Editor{Message: "What is your favorite month:", Default: "April", HideDefault: true},
			EditorTemplateData{},
			"? What is your favorite month: [Enter to launch editor] ",
		},
		{
			"Test Editor answer output",
			Editor{Message: "What is your favorite month:"},
			EditorTemplateData{Answer: "October", ShowAnswer: true},
			"? What is your favorite month: October\n",
		},
		{
			"Test Editor question output without default but with help hidden",
			Editor{Message: "What is your favorite month:", Help: "This is helpful"},
			EditorTemplateData{},
			"? What is your favorite month: [? for help] [Enter to launch editor] ",
		},
		{
			"Test Editor question output with default and with help hidden",
			Editor{Message: "What is your favorite month:", Default: "April", Help: "This is helpful"},
			EditorTemplateData{},
			"? What is your favorite month: [? for help] (April) [Enter to launch editor] ",
		},
		{
			"Test Editor question output without default but with help shown",
			Editor{Message: "What is your favorite month:", Help: "This is helpful"},
			EditorTemplateData{ShowHelp: true},
			`ⓘ This is helpful
? What is your favorite month: [Enter to launch editor] `,
		},
		{
			"Test Editor question output with default and with help shown",
			Editor{Message: "What is your favorite month:", Default: "April", Help: "This is helpful"},
			EditorTemplateData{ShowHelp: true},
			`ⓘ This is helpful
? What is your favorite month: (April) [Enter to launch editor] `,
		},
	}

	for _, test := range tests {
		r, w, err := os.Pipe()
		assert.Nil(t, err, test.title)
		test.prompt.WithStdio(terminal.Stdio{Out: w})

		test.data.Editor = test.prompt
		err = test.prompt.Render(
			EditorQuestionTemplate,
			test.data,
		)
		assert.Nil(t, err, test.title)

		w.Close()
		var buf bytes.Buffer
		io.Copy(&buf, r)

		assert.Contains(t, buf.String(), test.expected, test.title)
	}
}
