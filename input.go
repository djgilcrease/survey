package survey

import (
	"os"

	"github.com/djgilcrease/survey/core"
	"gopkg.in/AlecAivazis/survey.v1/terminal"
	"fmt"
)


// Templates with Color formatting. See Documentation: https://github.com/mgutz/ansi#style-format
var DefaultInputQuestionTemplate = `
{{- if .ShowHelp }}{{- color "cyan"}}{{ HelpIcon }} {{ .DisplayHelp }}{{color "reset"}}{{"\n"}}{{end}}
{{- color "green+hb"}}{{ QuestionIcon }} {{color "reset"}}
{{- color "default+hb"}}{{ .DisplayMessage }} {{color "reset"}}
{{- if .ShowAnswer}}
  {{- color "cyan"}}{{.Answer}}{{color "reset"}}{{"\n"}}
{{- else }}
  {{- if and .DisplayHelp (not .ShowHelp)}}{{color "cyan"}}[{{ HelpInputRune }} for help]{{color "reset"}} {{end}}
  {{- if .DisplayDefault}}{{color "white"}}({{.DisplayDefault}}) {{color "reset"}}{{end}}
{{- end}}`

/*
Input is a regular text input that prints each character the user types on the screen
and accepts the input with the enter key. Response type is a string.

	name := ""
	prompt := &survey.Input{ Message: "What is your name?" }
	survey.AskOne(prompt, &name, nil)
*/
type Input struct {
	core.Renderer
	message string
	defaultValue interface{}
	help    string
	tmpl string
}

func NewInput() *Input {
	return &Input{
		tmpl: DefaultInputQuestionTemplate,
	}
}

/*
SetMessage is a method to set the prompt message for a selection
This returns a Selection interface to allow chaining of these method calls
 */
func (i *Input) SetTemplate(tmpl string) Defaulter {
	i.tmpl = tmpl
	return i
}

/*
SetMessage is a method to set the prompt message for a selection
This returns a Selection interface to allow chaining of these method calls
 */
func (i *Input) SetMessage(msg string) Defaulter {
	i.message = msg
	return i
}

/*
DisplayMessage is a method to set the prompt message for a selection
This returns a Selection interface to allow chaining of these method calls
 */
func (i *Input) DisplayMessage() string {
	return i.message
}

/*
SetHelp is a method to set the prompt help message for a selection
This returns a Selection interface to allow chaining of these method calls
 */
func (i *Input) SetHelp(help string) Defaulter {
	i.help = help
	return i
}

/*
SetHelp is a method to set the prompt help message for a selection
This returns a Selection interface to allow chaining of these method calls
 */
func (i *Input) DisplayHelp() string {
	return i.help
}

/*
SetMessage is a method to set the prompt message for a selection
This returns a Selection interface to allow chaining of these method calls
 */
func (i *Input) SetDefault(value interface{}) Defaulter {
	i.defaultValue = value
	return i
}

/*
DisplayMessage is a method to set the prompt message for a selection
This returns a Selection interface to allow chaining of these method calls
 */
func (i *Input) DisplayDefault() string {
	if i.defaultValue == nil {
		return ""
	}
	return fmt.Sprintf("%v", i.defaultValue)
}

// data available to the templates when processing
type InputTemplateData struct {
	*Input
	Answer     string
	ShowAnswer bool
	ShowHelp   bool
}

func (i *Input) Prompt() (interface{}, error) {
	// render the template
	err := i.Render(
		i.tmpl,
		InputTemplateData{Input: i},
	)
	if err != nil {
		return "", err
	}

	// start reading runes from the standard in
	rr := terminal.NewRuneReader(os.Stdin)
	rr.SetTermMode()
	defer rr.RestoreTermMode()

	line := []rune{}
	// get the next line
	for {
		line, err = rr.ReadLine(0)
		if err != nil {
			return string(line), err
		}
		// terminal will echo the \n so we need to jump back up one row
		terminal.CursorPreviousLine(1)

		if string(line) == string(core.HelpInputRune) && i.help != "" {
			err = i.Render(
				i.tmpl,
				InputTemplateData{Input: i, ShowHelp: true},
			)
			if err != nil {
				return "", err
			}
			continue
		}
		break
	}

	// if the line is empty
	if line == nil || len(line) == 0 {
		// use the default value
		return i.defaultValue, err
	}

	// we're done
	return string(line), err
}

func (i *Input) Cleanup(val interface{}) error {
	return i.Render(
		i.tmpl,
		InputTemplateData{Input: i, Answer: val.(string), ShowAnswer: true},
	)
}
