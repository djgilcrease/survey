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

func TestMultiSelectRender(t *testing.T) {
	prompt := NewMultiSelect()
	prompt.SetMessage("Pick your words:").
		AddOption("foo", nil, false).
		AddOption("bar", nil, true).
		AddOption("baz", nil, false).
		AddOption("buz", nil, true)

	helpfulPrompt := NewMultiSelect()
	helpfulPrompt.SetMessage("Pick your words:").
		AddOption("foo", nil, false).
		AddOption("bar", nil, true).
		AddOption("baz", nil, false).
		AddOption("buz", nil, true).SetHelp("This is helpful")

	tests := []struct {
		title    string
		prompt   *MultiSelect
		data     MultiSelectTemplateData
		expected string
	}{
		{
			"Test MultiSelect question output",
			prompt,
			MultiSelectTemplateData{
				SelectedIndex: 2,
				PageEntries:   prompt.options,
				Checked:       map[string]bool{"bar": true, "buz": true},
			},
			`? Pick your words:  [Use arrows to move, type to filter]
  ◯  foo
  ◉  bar
❯ ◯  baz
  ◉  buz
`,
		},
		{
			"Test MultiSelect answer output",
			prompt,
			MultiSelectTemplateData{
				Answer:     append(make(Options, 0), prompt.options[0], prompt.options[3]),
				ShowAnswer: true,
			},
			"? Pick your words: foo, buz\n",
		},
		{
			"Test MultiSelect question output with help hidden",
			helpfulPrompt,
			MultiSelectTemplateData{
				SelectedIndex: 2,
				PageEntries:   prompt.options,
				Checked:       map[string]bool{"bar": true, "buz": true},
			},
			`? Pick your words:  [Use arrows to move, type to filter, ? for more help]
  ◯  foo
  ◉  bar
❯ ◯  baz
  ◉  buz
`,
		},
		{
			"Test MultiSelect question output with help shown",
			helpfulPrompt,
			MultiSelectTemplateData{
				SelectedIndex: 2,
				PageEntries:   prompt.options,
				Checked:       map[string]bool{"bar": true, "buz": true},
				ShowHelp:      true,
			},
			`ⓘ This is helpful
? Pick your words:  [Use arrows to move, type to filter]
  ◯  foo
  ◉  bar
❯ ◯  baz
  ◉  buz
`,
		},
	}


	for _, test := range tests {
		r, w, err := os.Pipe()
		assert.Nil(t, err, test.title)

		test.prompt.WithStdio(terminal.Stdio{Out: w})
		test.data.MultiSelect = test.prompt
		err = test.prompt.Render(
			test.prompt.tmpl,
			test.data,
		)
		assert.Nil(t, err, test.title)

		w.Close()
		var buf bytes.Buffer
		io.Copy(&buf, r)

		assert.Contains(t, buf.String(), test.expected, test.title)
	}
}


func TestMultiSelectInterfaceValues(t *testing.T) {
	type value struct {
		Item string
		Other int
	}
	prompt := NewMultiSelect()
	prompt.SetMessage("Pick your words:").
		AddOption("foo", value{"foo", 0}, false).
		AddOption("bar", value{"bar", 5}, true).
		AddOption("baz", value{"baz", 100}, false).
		AddOption("buz", value{"buz", 999}, true)

	helpfulPrompt := NewMultiSelect()
	helpfulPrompt.SetMessage("Pick your words:").
		AddOption("foo", value{"foo", 0}, false).
		AddOption("bar", value{"bar", 5}, true).
		AddOption("baz", value{"baz", 100}, false).
		AddOption("buz", value{"buz", 999}, true).SetHelp("This is helpful")

	tests := []struct {
		title    string
		prompt   *MultiSelect
		data     MultiSelectTemplateData
		expected string
	}{
		{
			"Test MultiSelect question output",
			prompt,
			MultiSelectTemplateData{
				SelectedIndex: 2,
				PageEntries:   prompt.options,
				Checked:       map[string]bool{"bar": true, "buz": true},
			},
			`? Pick your words:  [Use arrows to move, type to filter]
  ◯  foo
  ◉  bar
❯ ◯  baz
  ◉  buz
`,
		},
		{
			"Test MultiSelect answer output",
			prompt,
			MultiSelectTemplateData{
				Answer:     append(make(Options, 0), prompt.options[0], prompt.options[3]),
				ShowAnswer: true,
			},
			"? Pick your words: foo, buz\n",
		},
		{
			"Test MultiSelect question output with help hidden",
			helpfulPrompt,
			MultiSelectTemplateData{
				SelectedIndex: 2,
				PageEntries:   prompt.options,
				Checked:       map[string]bool{"bar": true, "buz": true},
			},
			`? Pick your words:  [Use arrows to move, type to filter, ? for more help]
  ◯  foo
  ◉  bar
❯ ◯  baz
  ◉  buz
`,
		},
		{
			"Test MultiSelect question output with help shown",
			helpfulPrompt,
			MultiSelectTemplateData{
				SelectedIndex: 2,
				PageEntries:   prompt.options,
				Checked:       map[string]bool{"bar": true, "buz": true},
				ShowHelp:      true,
			},
			`ⓘ This is helpful
? Pick your words:  [Use arrows to move, type to filter]
  ◯  foo
  ◉  bar
❯ ◯  baz
  ◉  buz
`,
		},
	}

	for _, test := range tests {
		r, w, err := os.Pipe()
		assert.Nil(t, err, test.title)

		test.prompt.WithStdio(terminal.Stdio{Out: w})
		test.data.MultiSelect = test.prompt
		err = test.prompt.Render(
			test.prompt.tmpl,
			test.data,
		)
		assert.Nil(t, err, test.title)

		w.Close()
		var buf bytes.Buffer
		io.Copy(&buf, r)

		assert.Contains(t, buf.String(), test.expected, test.title)
	}
}