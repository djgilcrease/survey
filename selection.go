package survey

// PageSize is the default maximum number of items to show in select/multiselect prompts
var DefaultPageSize = 7

// Option is a struct to selerate the display value from the actual value of a selection
type Option struct {
	Display string
	Value interface{}
}

// String to define the default output of Option
func (o *Option) String() string {
	return o.Display
}

// Options alias for []*Option
type Options = []*Option

func OptionsValues(o Options) (values []interface{}) {
	values = make([]interface{}, len(o))
	for i, v := range o {
		values[i] = v.Value
	}

	return
}