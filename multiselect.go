package survey

import (
	"errors"
	"os"
	"github.com/djgilcrease/survey/core"
	"gopkg.in/AlecAivazis/survey.v1/terminal"
)

var DefaultMultiSelectQuestionTemplate = `
{{- if .ShowHelp }}{{- color "cyan"}}{{ HelpIcon }} {{ .DisplayHelp }}{{color "reset"}}{{"\n"}}{{end}}
{{- color "green+hb"}}{{ QuestionIcon }} {{color "reset"}}
{{- color "default+hb"}}{{ .DisplayMessage }}{{ .DisplayFilterMessage }}{{color "reset"}}
{{- if .ShowAnswer}} {{color "cyan"}} 
	{{- range $ix, $answer := .Answer}}
		{{- if ne $ix 0 }}{{- ", "}}{{- end}}
		{{- $answer.Display }}
	{{- end}}{{- color "reset"}}{{- "\n"}}
{{- else }}
  	{{- "  "}}{{- color "cyan"}}[Use arrows to move, type to filter{{- if and .DisplayHelp (not .ShowHelp)}}, {{ HelpInputRune }} for more help{{end}}]{{color "reset"}}
  	{{- "\n"}}
  	{{- range $ix, $option := .PageEntries}}
    	{{- if eq $ix $.SelectedIndex}}{{color "cyan"}}{{ SelectFocusIcon }}{{color "reset"}}{{else}} {{end}}
    	{{- if index $.Checked $option.Display}}{{color "green"}} {{ MarkedOptionIcon }} {{else}}{{color "default+hb"}} {{ UnmarkedOptionIcon }} {{end}}
    	{{- color "reset"}}
    	{{- " "}}{{$option.Display}}{{"\n"}}
  	{{- end}}
{{- end}}`

/*
MultiSelect is a prompt that presents a list of various options to the user
for them to select using the arrow keys and enter. Response type is a slice of strings.

	days := make(survey.Options, 0)
	prompt := survey.NewMultiSelect().SetMessage("What days do you prefer:").
			AddOption("Sunday", nil, false).
			AddOption("Monday", nil, false).
			AddOption("Tuesday", nil, false).
			AddOption("Wednesday", nil, false).
			AddOption("Thursday", nil, false).
			AddOption("Friday", nil, false).
			AddOption("Saturday", nil, false)
	survey.AskOne(prompt, &days, nil)
*/
type MultiSelect struct {
	*Select
	defaultValue       Options
	checked       map[string]bool
	tmpl string
}

/*
NewMultiSelect is a shortcut method to get a MultiSelect prompt
 */
func NewMultiSelect() *MultiSelect {
	return &MultiSelect{
		Select: NewSingleSelect(),
		defaultValue: make(Options, 0),
		tmpl: DefaultMultiSelectQuestionTemplate,
	}
}

/*
AddStringOption is a method to add an option to the selection and specify if it is the default value ot not
This returns a Selection interface to allow chaining of these method calls

	days := make(survey.Options, 0)
	prompt := survey.NewMultiSelect().SetMessage("What days do you prefer:").
			AddStringOption("Sunday", false).
			AddStringOption("Monday", false).
			AddStringOption("Tuesday", false).
			AddStringOption("Wednesday", false).
			AddStringOption("Thursday", false).
			AddStringOption("Friday", false).
			AddStringOption("Saturday", false)
	survey.AskOne(prompt, &days, nil)
 */
func (s *MultiSelect) AddStringOption(display string, defaultOption bool) Selection {
	opt := &Option{display, display}
	s.options = append(s.options, opt)
	if defaultOption {
		s.defaultValue = append(s.defaultValue, opt)
	}
	return s
}


/*
AddOption is a method to add an option to the selection and specify if it is the default value ot not
This returns a Selection interface to allow chaining of these method calls

	days := make(survey.Options, 0)
	prompt := survey.NewMultiSelect().SetMessage("What days do you prefer:").
			AddOption("Sunday", nil, false).
			AddOption("Monday", nil, false).
			AddOption("Tuesday", nil, false).
			AddOption("Wednesday", nil, false).
			AddOption("Thursday", nil, false).
			AddOption("Friday", nil, false).
			AddOption("Saturday", nil, false)
	survey.AskOne(prompt, &days, nil)
 */
func (s *MultiSelect) AddOption(display string, value interface{}, defaultOption bool) Selection {
	if value == nil {
		value = display
	}

	opt := &Option{display, value}
	s.options = append(s.options, opt)
	if defaultOption {
		s.defaultValue = append(s.defaultValue, opt)
	}
	return s
}

/*
SetMessage is a method to set the prompt message for a selection
This returns a Selection interface to allow chaining of these method calls
 */
func (s *MultiSelect) SetTemplate(tmpl string) Selection {
	s.tmpl = tmpl
	return s
}

/*
SetMessage is a method to set the prompt message for a selection
This returns a Selection interface to allow chaining of these method calls
 */
func (s *MultiSelect) SetMessage(msg string) Selection {
	s.message = msg
	return s
}

/*
DisplayMessage is a method to set the prompt message for a selection
This returns a Selection interface to allow chaining of these method calls
 */
func (s *MultiSelect) DisplayMessage() string {
	return s.message
}

/*
SetHelp is a method to set the prompt help message for a selection
This returns a Selection interface to allow chaining of these method calls
 */
func (s *MultiSelect) SetHelp(help string) Selection {
	s.help = help
	return s
}

/*
SetHelp is a method to set the prompt help message for a selection
This returns a Selection interface to allow chaining of these method calls
 */
func (s *MultiSelect) DisplayHelp() string {
	return s.help
}

/*
SetFilterMessage is a method to set the prompt filter message for a selection
This returns a Selection interface to allow chaining of these method calls
 */
func (s *MultiSelect) SetFilterMessage(msg string) Selection {
	s.filterMessage = msg
	return s
}

/*
SetFilterMessage is a method to set the prompt filter message for a selection
This returns a Selection interface to allow chaining of these method calls
 */
func (s *MultiSelect) DisplayFilterMessage() string {
	return s.filterMessage
}

/*
SetVimMode is a method to turn on or off VimMode
This returns a Selection interface to allow chaining of these method calls
 */
func (s *MultiSelect) SetVimMode(vimMode bool) Selection {
	s.vimMode = vimMode
	return s
}

/*
SetVimMode is a method to turn on or off VimMode
This returns a Selection interface to allow chaining of these method calls
 */
func (s *MultiSelect) SetPageSize(pageSize int) Selection {
	s.pageSize = pageSize
	return s
}


// data available to the templates when processing
type MultiSelectTemplateData struct {
	*MultiSelect
	Answer        Options
	ShowAnswer    bool
	Checked       map[string]bool
	SelectedIndex int
	ShowHelp      bool
	PageEntries   Options
}

// OnChange is called on every keypress.
func (m *MultiSelect) OnChange(line []rune, pos int, key rune) (newLine []rune, newPos int, ok bool) {
	options := m.filterOptions()
	oldFilter := m.filter

	if key == terminal.KeyArrowUp || (m.vimMode && key == 'k') {
		// if we are at the top of the list
		if m.selectedIndex == 0 {
			// go to the bottom
			m.selectedIndex = len(options) - 1
		} else {
			// decrement the selected index
			m.selectedIndex--
		}
	} else if key == terminal.KeyArrowDown || (m.vimMode && key == 'j') {
		// if we are at the bottom of the list
		if m.selectedIndex == len(options)-1 {
			// start at the top
			m.selectedIndex = 0
		} else {
			// increment the selected index
			m.selectedIndex++
		}
		// if the user pressed down and there is room to move
	} else if key == terminal.KeySpace {
		if m.selectedIndex < len(options) {
			if old, ok := m.checked[options[m.selectedIndex].Display]; !ok {
				// otherwise just invert the current value
				m.checked[options[m.selectedIndex].Display] = true
			} else {
				// otherwise just invert the current value
				m.checked[options[m.selectedIndex].Display] = !old
			}
		}
		// only show the help message if we have one to show
	} else if key == core.HelpInputRune && m.help != "" {
		m.showingHelp = true
	} else if key == terminal.KeyEscape {
		m.vimMode = !m.vimMode
	} else if key == terminal.KeyDeleteWord || key == terminal.KeyDeleteLine {
		m.filter = ""
	} else if key == terminal.KeyDelete || key == terminal.KeyBackspace {
		if m.filter != "" {
			m.filter = m.filter[0 : len(m.filter)-1]
		}
	} else if key >= terminal.KeySpace {
		m.filter += string(key)
	}

	m.filterMessage = ""
	if m.filter != "" {
		m.filterMessage = " " + m.filter
	}
	if oldFilter != m.filter {
		// filter changed
		options = m.filterOptions()
		if len(options) > 0 && len(options) <= m.selectedIndex {
			m.selectedIndex = len(options) - 1
		}
	}
	// paginate the options

	// TODO if we have started filtering and were looking at the end of a list
	// and we have modified the filter then we should move the page back!
	opts, idx := m.Paginate(options)

	// render the options
	m.Render(
		m.tmpl,
		MultiSelectTemplateData{
			MultiSelect:   m,
			SelectedIndex: idx,
			Checked:       m.checked,
			ShowHelp:      m.showingHelp,
			PageEntries:   opts,
		},
	)

	// if we are not pressing ent
	return line, 0, true
}

func (m *MultiSelect) Prompt() (interface{}, error) {
	// compute the default state
	m.checked = make(map[string]bool)
	// if there is a default
	if len(m.defaultValue) > 0 {
		for _, dflt := range m.defaultValue {
			for _, opt := range m.options {
				// if the option correponds to the default
				if opt == dflt {
					// we found our initial value
					m.checked[opt.Display] = true
					// stop looking
					break
				}
			}
		}
	}

	// if there are no options to render
	if len(m.options) == 0 {
		// we failed
		return "", errors.New("please provide options to select from")
	}

	// paginate the options
	opts, idx := m.Paginate(m.options)

	// hide the cursor
	terminal.CursorHide()

	// show the cursor when we're done
	defer terminal.CursorShow()

	// ask the question
	err := m.Render(
		m.tmpl,
		MultiSelectTemplateData{
			MultiSelect:   m,
			SelectedIndex: idx,
			Checked:       m.checked,
			PageEntries:   opts,
		},
	)
	if err != nil {
		return "", err
	}

	rr := terminal.NewRuneReader(os.Stdin)
	rr.SetTermMode()
	defer rr.RestoreTermMode()

	// start waiting for input
	for {
		r, _, _ := rr.ReadRune()
		if r == '\r' || r == '\n' {
			break
		}
		if r == terminal.KeyInterrupt {
			return "", terminal.InterruptErr
		}
		if r == terminal.KeyEndTransmission {
			break
		}
		m.OnChange(nil, 0, r)
	}
	m.filter = ""
	m.filterMessage = ""

	answers := make(Options, 0)
	for _, option := range m.options {
		if val, ok := m.checked[option.Display]; ok && val {
			answers = append(answers, option)
		}
	}

	return answers, nil
}

// Cleanup removes the options section, and renders the ask like a normal question.
func (m *MultiSelect) Cleanup(val interface{}) error {
	// execute the output summary template with the answer
	return m.Render(
		m.tmpl,
		MultiSelectTemplateData{
			MultiSelect:   m,
			SelectedIndex: m.selectedIndex,
			Checked:       m.checked,
			Answer:        val.(Options),
			ShowAnswer:    true,
		},
	)
}
