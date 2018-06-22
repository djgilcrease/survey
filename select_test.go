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

func TestSelectRender(t *testing.T) {
	prompt := NewSingleSelect()
	prompt.SetMessage("Pick your word:").
		AddStringOption("foo", false).
		AddStringOption("bar", false).
		AddStringOption("baz", true).
		AddStringOption("buz", false)

	helpfulPrompt := NewSingleSelect()
	helpfulPrompt.SetMessage("Pick your word:").
		AddStringOption("foo", false).
		AddStringOption("bar", false).
		AddStringOption("baz", true).
		AddStringOption("buz", false).SetHelp("This is helpful")

	tests := []struct {
		title    string
		prompt   *Select
		data     SelectTemplateData
		expected string
	}{
		{
			"Test Select question output",
			prompt,
			SelectTemplateData{SelectedIndex: 2, PageEntries: prompt.options},
			`? Pick your word:  [Use arrows to move, type to filter]
  foo
  bar
❯ baz
  buz
`,
		},
		{
			"Test Select answer output",
			prompt,
			SelectTemplateData{Answer: prompt.options[3], ShowAnswer: true, PageEntries: prompt.options},
			"? Pick your word: buz\n",
		},
		{
			"Test Select question output with help hidden",
			helpfulPrompt,
			SelectTemplateData{SelectedIndex: 2, PageEntries: prompt.options},
			`? Pick your word:  [Use arrows to move, type to filter, ? for more help]
  foo
  bar
❯ baz
  buz
`,
		},
		{
			"Test Select question output with help shown",
			helpfulPrompt,
			SelectTemplateData{SelectedIndex: 2, ShowHelp: true, PageEntries: prompt.options},
			`ⓘ This is helpful
? Pick your word:  [Use arrows to move, type to filter]
  foo
  bar
❯ baz
  buz
`,
		},
	}

	for _, test := range tests {
		r, w, err := os.Pipe()
		assert.Nil(t, err, test.title)

		test.prompt.WithStdio(terminal.Stdio{Out: w})
		test.data.Select = test.prompt
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

func TestSelectInterfaceValues(t *testing.T) {
	type value struct {
		Item string
		Other int
	}
	prompt := NewSingleSelect()
	prompt.SetMessage("Pick your word:").
		AddOption("foo", value{"foo", 0}, false).
		AddOption("bar", value{"bar", 5}, false).
		AddOption("baz", value{"baz", 100}, true).
		AddOption("buz", value{"buz", 999}, false)

	helpfulPrompt := NewSingleSelect()
	helpfulPrompt.SetMessage("Pick your word:").
		AddOption("foo", value{"foo", 0}, false).
		AddOption("bar", value{"bar", 5}, false).
		AddOption("baz", value{"baz", 100}, true).
		AddOption("buz", value{"buz", 999}, false).SetHelp("This is helpful")

	tests := []struct {
		title    string
		prompt   *Select
		data     SelectTemplateData
		expected string
	}{
		{
			"Test Select question output",
			prompt,
			SelectTemplateData{SelectedIndex: 2, PageEntries: prompt.options},
			`? Pick your word:  [Use arrows to move, type to filter]
  foo
  bar
❯ baz
  buz
`,
		},
		{
			"Test Select answer output",
			prompt,
			SelectTemplateData{Answer: prompt.options[3], ShowAnswer: true, PageEntries: prompt.options},
			"? Pick your word: buz\n",
		},
		{
			"Test Select question output with help hidden",
			helpfulPrompt,
			SelectTemplateData{SelectedIndex: 2, PageEntries: prompt.options},
			`? Pick your word:  [Use arrows to move, type to filter, ? for more help]
  foo
  bar
❯ baz
  buz
`,
		},
		{
			"Test Select question output with help shown",
			helpfulPrompt,
			SelectTemplateData{SelectedIndex: 2, ShowHelp: true, PageEntries: prompt.options},
			`ⓘ This is helpful
? Pick your word:  [Use arrows to move, type to filter]
  foo
  bar
❯ baz
  buz
`,
		},
	}

	for _, test := range tests {
		r, w, err := os.Pipe()
		assert.Nil(t, err, test.title)

		test.prompt.WithStdio(terminal.Stdio{Out: w})
		test.data.Select = test.prompt
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

func TestSelectionPagination_tooFew(t *testing.T) {
	prompt := NewSingleSelect()
	prompt.SetMessage("Pick your word:").
		AddStringOption("choice1", false).
		AddStringOption("choice2", false).
		AddStringOption("choice3", false).
		SetPageSize(4)
	// the current selection
	prompt.selectedIndex = 3

	// compute the page info
	page, idx := prompt.Paginate(prompt.options)

	// make sure we see the full list of options
	assert.Equal(t, prompt.options, page)
	// with the second index highlighted (no change)
	assert.Equal(t, 3, idx)
}

func TestSelectionPagination_firstHalf(t *testing.T) {
	prompt := NewSingleSelect()
	prompt.SetMessage("Pick your word:").
		AddStringOption("choice1", false).
		AddStringOption("choice2", false).
		AddStringOption("choice3", false).
		AddStringOption("choice4", false).
		AddStringOption("choice5", false).
		AddStringOption("choice6", false).
		SetPageSize(4)
	// the current selection
	prompt.selectedIndex = 2

	// compute the page info
	page, idx := prompt.Paginate(prompt.options)

	// we should see the first three options
	assert.Equal(t, prompt.options[0:4], page)
	// with the second index highlighted
	assert.Equal(t, 2, idx)
}

func TestSelectionPagination_middle(t *testing.T) {
	prompt := NewSingleSelect()
	prompt.SetMessage("Pick your word:").
		AddStringOption("choice1", false).
		AddStringOption("choice2", false).
		AddStringOption("choice3", false).
		AddStringOption("choice4", false).
		AddStringOption("choice5", false).
		SetPageSize(2)
	// the current selection
	prompt.selectedIndex = 3

	// compute the page info
	page, idx := prompt.Paginate(prompt.options)

	// we should see the first three options
	assert.Equal(t, prompt.options[2:4], page)
	// with the second index highlighted
	assert.Equal(t, 1, idx)
}

func TestSelectionPagination_lastHalf(t *testing.T) {
	prompt := NewSingleSelect()
	prompt.SetMessage("Pick your word:").
		AddStringOption("choice1", false).
		AddStringOption("choice2", false).
		AddStringOption("choice3", false).
		AddStringOption("choice4", false).
		AddStringOption("choice5", false).
		SetPageSize(3)
	// the current selection
	prompt.selectedIndex = 4

	// compute the page info
	page, idx := prompt.Paginate(prompt.options)

	// we should see the last three options
	assert.Equal(t, prompt.options[2:5], page)
	// we should be at the bottom of the list
	assert.Equal(t, 2, idx)
}
