package survey

import (
	"fmt"
	"regexp"

	"github.com/djgilcrease/survey/core"
)

// Confirm is a regular text input that accept yes/no answers. Response type is a bool.
type Confirm struct {
	core.Renderer
	message string
	defaultValue bool
	help    string
	tmpl string
}

// data available to the templates when processing
type ConfirmTemplateData struct {
	*Confirm
	Answer   string
	ShowHelp bool
}

// Templates with Color formatting. See Documentation: https://github.com/mgutz/ansi#style-format
var DefaultConfirmQuestionTemplate = `
{{- if .ShowHelp }}{{- color "cyan"}}{{ HelpIcon }} {{ .DisplayHelp }}{{color "reset"}}{{"\n"}}{{end}}
{{- color "green+hb"}}{{ QuestionIcon }} {{color "reset"}}
{{- color "default+hb"}}{{ .DisplayMessage }} {{color "reset"}}
{{- if .Answer}}
  {{- color "cyan"}}{{.Answer}}{{color "reset"}}{{"\n"}}
{{- else }}
  {{- if and .DisplayHelp (not .ShowHelp)}}{{color "cyan"}}[{{ HelpInputRune }} for help]{{color "reset"}} {{end}}
  {{- color "white"}}{{ .DisplayDefault }} {{color "reset"}}
{{- end}}`

func NewConfirm() *Confirm {
	return &Confirm{
		tmpl: DefaultConfirmQuestionTemplate,
	}
}


/*
SetMessage is a method to set the prompt message for a selection
This returns a Selection interface to allow chaining of these method calls
 */
func (i *Confirm) SetTemplate(tmpl string) Defaulter {
	i.tmpl = tmpl
	return i
}

/*
SetMessage is a method to set the prompt message for a selection
This returns a Selection interface to allow chaining of these method calls
 */
func (i *Confirm) SetMessage(msg string) Defaulter {
	i.message = msg
	return i
}

/*
DisplayMessage is a method to set the prompt message for a selection
This returns a Selection interface to allow chaining of these method calls
 */
func (i *Confirm) DisplayMessage() string {
	return i.message
}

/*
SetHelp is a method to set the prompt help message for a selection
This returns a Selection interface to allow chaining of these method calls
 */
func (i *Confirm) SetHelp(help string) Defaulter {
	i.help = help
	return i
}

/*
SetHelp is a method to set the prompt help message for a selection
This returns a Selection interface to allow chaining of these method calls
 */
func (i *Confirm) DisplayHelp() string {
	return i.help
}

/*
SetMessage is a method to set the prompt message for a selection
This returns a Selection interface to allow chaining of these method calls
 */
func (i *Confirm) SetDefault(value interface{}) Defaulter {
	i.defaultValue = value.(bool)
	return i
}

/*
DisplayMessage is a method to set the prompt message for a selection
This returns a Selection interface to allow chaining of these method calls
 */
func (i *Confirm) DisplayDefault() string {
	if !i.defaultValue {
		return "(y/N)"
	}
	return "(Y/n)"
}

// the regex for answers
var (
	yesRx = regexp.MustCompile("^(?i:y(?:es)?)$")
	noRx  = regexp.MustCompile("^(?i:n(?:o)?)$")
)

func yesNo(t bool) string {
	if t {
		return "Yes"
	}
	return "No"
}

func (c *Confirm) getBool(showHelp bool) (bool, error) {
	rr := c.NewRuneReader()
	rr.SetTermMode()
	defer rr.RestoreTermMode()
	cursor := c.NewCursor()
	// start waiting for input
	for {
		line, err := rr.ReadLine(0)
		if err != nil {
			return false, err
		}
		// move back up a line to compensate for the \n echoed from terminal
		cursor.PreviousLine(1)
		val := string(line)

		// get the answer that matches the
		var answer bool
		switch {
		case yesRx.Match([]byte(val)):
			answer = true
		case noRx.Match([]byte(val)):
			answer = false
		case val == "":
			answer = c.defaultValue
		case val == string(core.HelpInputRune) && c.help != "":
			err := c.Render(
				DefaultConfirmQuestionTemplate,
				ConfirmTemplateData{Confirm: c, ShowHelp: true},
			)
			if err != nil {
				// use the default value and bubble up
				return c.defaultValue, err
			}
			showHelp = true
			continue
		default:
			// we didnt get a valid answer, so print error and prompt again
			if err := c.Error(fmt.Errorf("%q is not a valid answer, please try again.", val)); err != nil {
				return c.defaultValue, err
			}
			err := c.Render(
				DefaultConfirmQuestionTemplate,
				ConfirmTemplateData{Confirm: c, ShowHelp: showHelp},
			)
			if err != nil {
				// use the default value and bubble up
				return c.defaultValue, err
			}
			continue
		}
		return answer, nil
	}
	// should not get here
	return c.defaultValue, nil
}

/*
Prompt prompts the user with a simple text field and expects a reply followed
by a carriage return.

	likesPie := false
	prompt := &survey.Confirm{ Message: "What is your name?" }
	survey.AskOne(prompt, &likesPie, nil)
*/
func (c *Confirm) Prompt() (interface{}, error) {
	// render the question template
	err := c.Render(
		DefaultConfirmQuestionTemplate,
		ConfirmTemplateData{Confirm: c},
	)
	if err != nil {
		return "", err
	}

	// get input and return
	return c.getBool(false)
}

// Cleanup overwrite the line with the finalized formatted version
func (c *Confirm) Cleanup(val interface{}) error {
	// if the value was previously true
	ans := yesNo(val.(bool))
	// render the template
	return c.Render(
		DefaultConfirmQuestionTemplate,
		ConfirmTemplateData{Confirm: c, Answer: ans},
	)
}
