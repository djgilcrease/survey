package survey


// Prompt is the primary interface for the objects that can take user input
// and return a response.
type Prompt interface {
	Prompt() (interface{}, error)
	Cleanup(interface{}) error
	Error(error) error
	DisplayMessage() string
	DisplayHelp() string
}

// Selection is the interface for a prompt that will present a section list to the user
type Selection interface {
	Prompt
	SetMessage(msg string) Selection
	SetHelp(help string) Selection
	SetTemplate(tmpl string) Selection
	AddStringOption(display string, defaultOption bool) Selection
	AddOption(display string, value interface{}, defaultOption bool) Selection
	SetFilterMessage(msg string) Selection
	DisplayFilterMessage() string
	SetVimMode(vimMode bool) Selection
	SetPageSize(pageSize int) Selection
	Paginate(choices Options) (Options, int)
	OnChange(line []rune, pos int, key rune) (newLine []rune, newPos int, ok bool)
}

//
type Defaulter interface {
	Prompt
	SetMessage(msg string) Defaulter
	SetHelp(help string) Defaulter
	SetTemplate(tmpl string) Defaulter
	SetDefault(value interface{}) Defaulter
	DisplayDefault() string
}